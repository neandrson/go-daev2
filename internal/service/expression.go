package service

import (
	"container/list"
	"strconv"
	"strings"

	"github.com/neandrson/go-daev2/pkg/rpn"
)

const (
	StatusError     = "Error"
	StatusDone      = "Done"
	StatusInProcess = "In process"
)

const (
	TokenTypeNumber = iota
	TokenTypeOperation
	TokenTypeTask
)

type Token interface {
	Type() int
}

type NumToken struct {
	Value float64
}

type OpToken struct {
	Value string
}

type TaskToken struct {
	ID int64
}

type Expression struct {
	*list.List
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
	Source string `json:"source"`
}

// Структура для ответа по запросу на endpoint expressions/{id}
type ExpressionUnit struct {
	Expr Expression `json:"expression"`
}

// Структура для ответа по запросу на endpoint expressions
type ExpressionList struct {
	Exprs []Expression `json:"expressions"`
}

// структура связывающая узел списка, в который нужно положить
// результат вычисления, с ID выражения, которое хранит это список
type ExprElement struct {
	ID  string
	Ptr *list.Element
}

func (num NumToken) Type() int {
	return TokenTypeNumber
}

func (num OpToken) Type() int {
	return TokenTypeOperation
}

func (num TaskToken) Type() int {
	return TokenTypeTask
}

func NewExpression(id, expr string) (*Expression, error) {
	rpn, err := rpn.NewRPN(expr)
	if err != nil {
		expression := Expression{
			List:   list.New(),
			ID:     id,
			Status: StatusError,
			Result: "",
			Source: expr,
		}
		return &expression, err
	}

	if len(rpn) == 1 {
		expression := Expression{
			List:   list.New(),
			ID:     id,
			Status: StatusDone,
			Result: rpn[0],
			Source: expr,
		}
		return &expression, nil
	}

	expression := Expression{
		List:   list.New(),
		ID:     id,
		Status: StatusInProcess,
		Result: rpn[len(rpn)],
		Source: expr,
	}
	for _, val := range rpn {
		if strings.Contains("-+*/", val) {
			expression.PushBack(OpToken{val})
		} else {
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			expression.PushBack(NumToken{num})
		}
	}
	return &expression, nil
}
