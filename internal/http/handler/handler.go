// В этом пакетк содержится код обработчиков http запросов.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/neandrson/go-daev2/internal/result"
	"github.com/neandrson/go-daev2/internal/service"
	"github.com/neandrson/go-daev2/internal/task"
)

// тип Decorator служит для добавления middleware к обработчикам
type Decorator func(http.Handler) http.Handler

// объект для обработки запросов
type calcStates struct {
	CalcService   *service.CalcService
	ClientGetTask *Client
}

type Client struct {
	http.Client
	Hostname string
	Port     int
}

func NewHandler(ctx context.Context, calcService *service.CalcService) (http.Handler, error) {
	serveMux := http.NewServeMux()

	calcState := calcStates{
		CalcService: calcService,
	}

	serveMux.HandleFunc("POST /api/v1/calculate", calcState.calculate)
	serveMux.HandleFunc("GET /api/v1/expressions", calcState.listAll)
	serveMux.HandleFunc("GET /api/v1/expressions/{id}", calcState.listByID)
	serveMux.HandleFunc("GET /internal/task", calcState.sendTask)
	serveMux.HandleFunc("POST /internal/task", calcState.receiveResult)

	return serveMux, nil
}

// функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

// Добавление вычисления арифметического выражения
func (cs *calcStates) calculate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !slices.Contains(r.Header["Content-Type"], "application/json") {
		http.Error(w, "Incorrect header", http.StatusUnprocessableEntity)
		return
	}

	type Expression struct {
		Id         string `json:"id"`
		Expression string `json:"expression"`
	}
	var expr Expression
	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expr.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	if err = cs.CalcService.AddExpression(expr.Id, expr.Expression); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println(&expr)
}

func (cs *calcStates) listAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	lst := cs.CalcService.ListAll()
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(&lst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cs *calcStates) listByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := r.PathValue("id")
	expr, err := cs.CalcService.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cs *calcStates) sendTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newTask := cs.CalcService.GetTask() //cs.CalcService.GetTask()
	if newTask == nil {
		http.Error(w, "no tasks", http.StatusNotFound)
		return
	}

	answer := struct {
		Task *task.Task `json:"task"`
	}{
		Task: newTask,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(&answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cs *calcStates) GetTask() *task.Task {
	url := fmt.Sprintf("http://%s:%d/internal/task", cs.ClientGetTask.Hostname, cs.ClientGetTask.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cs.ClientGetTask.Do(req.WithContext(ctx))
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
		time.Sleep(500)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	answer := struct {
		Task task.Task `json:"task"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
		return nil
	}

	return &answer.Task
}

func (cs *calcStates) receiveResult(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var res result.Result
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, err := strconv.ParseFloat(res.Value, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	if err = cs.CalcService.PutResult(res.ID, value); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
