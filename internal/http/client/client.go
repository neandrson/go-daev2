package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "os"
	"time"

	"github.com/neandrson/go-daev2/internal/result"
	"github.com/neandrson/go-daev2/internal/task"
)

type Client struct {
	http.Client
	Hostname string
	Port     int
}

func (client *Client) GetTask() *task.Task {
	url := fmt.Sprintf("http://%s:%d/internal/task", client.Hostname, client.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	return &answer.Task
}

func (client *Client) SendResult(result result.Result) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(result)
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/internal/task", client.Hostname, client.Port)
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return
	}

	//fmt.Println(req)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return
	}
	//fmt.Println(resp)
	defer resp.Body.Close()
}
