package rpn

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/neandrson/go-daev2/pkg/stack"
)

const (
	emptyToken = iota
	wrongToken
	numberToken
	operatorToken
	leftBracketToken
	rightBracketToken
)

func NewRPN(input string) ([]string, error) {
	input = strings.ReplaceAll(input, "+", " + ")
	input = strings.ReplaceAll(input, "-", " - ")
	input = strings.ReplaceAll(input, "*", " * ")
	input = strings.ReplaceAll(input, "/", " / ")
	input = strings.ReplaceAll(input, "(", " ( ")
	input = strings.ReplaceAll(input, ")", " ) ")

	tokens := strings.Fields(input)

	rpn := make([]string, 0, len(tokens))
	stack := stack.NewStack[string]()
	prevToken := emptyToken

	for _, token := range tokens {
		curToken := emptyToken

		if isOperator(token) {
			if isUnaryOperator(token, prevToken) {
				rpn = append(rpn, "0")
				stack.Push(token)
				continue
			}
			for !stack.Empty() && isOperator(stack.Top()) {
				op := stack.Pop()
				if operatorPriority(op) <= operatorPriority(token) {
					rpn = append(rpn, op)
				} else {
					stack.Push(op)
					break
				}
			}
			stack.Push(token)
			curToken = operatorToken
		} else if token == "(" {
			stack.Push(token)
			curToken = leftBracketToken
		} else if token == ")" {
			for !stack.Empty() && stack.Top() != "(" {
				rpn = append(rpn, stack.Pop())
			}
			if stack.Empty() {
				return nil, fmt.Errorf("error")
			}
			stack.Pop()
			curToken = rightBracketToken
		} else {
			// токен является числом
			_, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, fmt.Errorf("incorrect token: '%s'", token)
			}
			rpn = append(rpn, token)
			curToken = numberToken
		}

		if !checkTokens(prevToken, curToken) {
			return nil, fmt.Errorf("incorrect sequence near token: '%s'", token)
		}
		prevToken = curToken
	}

	for !stack.Empty() {
		token := stack.Pop()
		if token == "(" {
			return nil, fmt.Errorf("unpaired brackets")
		}
		rpn = append(rpn, token)
	}
	if prevToken != numberToken && prevToken != rightBracketToken {
		return nil, fmt.Errorf("incorrect sequence near last token")
	}

	return rpn, nil
}

func operatorPriority(op string) int {
	if op == "*" || op == "/" {
		return 1
	} else if op == "-" || op == "+" {
		return 2
	}
	return -1
}

func isOperator(op string) bool {
	if op == "+" || op == "-" || op == "*" || op == "/" {
		return true
	}
	return false
}

func isUnaryOperator(op string, prevToken int) bool {
	if (op == "-" || op == "+") &&
		(prevToken == emptyToken || prevToken == leftBracketToken ||
			prevToken == operatorToken) {
		return true
	}
	return false
}

func checkTokens(prev, cur int) bool {
	switch cur {
	case numberToken:
		return prev == emptyToken || prev == operatorToken ||
			prev == leftBracketToken
	case leftBracketToken:
		return prev == emptyToken || prev == operatorToken ||
			prev == leftBracketToken
	case rightBracketToken:
		return prev == numberToken || prev == rightBracketToken
	case operatorToken:
		return prev == numberToken || prev == rightBracketToken
	default:
		return false
	}
}
