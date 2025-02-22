package rpn_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Vojan-Najov/daec/pkg/rpn"
)

func TestRpn(t *testing.T) {
	input := [...]string{
		"1 + 2",
		"2*3 + 4*5",
		"2 * -3",
		"2 * (2 + -4)",
	}
	expectedArr := [...][]string{
		[]string{"1", "2", "+"},
		[]string{"2", "3", "*", "4", "5", "*", "+"},
		[]string{"2", "0", "3", "-", "*"},
		[]string{"2", "2", "0", "4", "-", "+", "*"},
	}
	expectedErrors := [...]error{
		nil,
		nil,
		nil,
		nil,
	}
	for i, expected := range expectedArr {
		expectedErr := expectedErrors[i]
		actual, err := rpn.NewRPN(input[i])
		fmt.Println(actual)
		if expectedErr == nil && err != nil ||
			expectedErr != nil && err == nil {
			t.Errorf("expected %v, actual %v", expectedErr, err)
		}
		if !slices.Equal(actual, expected) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}
