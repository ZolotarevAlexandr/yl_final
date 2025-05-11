package orchestrator

import (
	"errors"

	"github.com/ZolotarevAlexandr/yl_final/calculator/calculator"
)

// Expression represents an expression submitted by the user.
type Expression struct {
	ID         string   `json:"id"`
	Expr       string   `json:"expression"`
	Status     string   `json:"status"` // "pending" or "done"
	Result     *float64 `json:"result,omitempty"`
	RootTaskID string   `json:"-"`
}

// Node represents a node in the expression tree.
type Node struct {
	IsLiteral bool
	Value     float64 // if IsLiteral is true
	Operator  string  // if node represents an operation
	Left      *Node
	Right     *Node
	TaskID    string
}

// getOperationTime returns the operation execution time using environment variables.
func getOperationTime(op string) int {
	switch op {
	case "+":
		return AdditionTimeMs
	case "-":
		return SubtractionTimeMs
	case "*":
		return MultiplicationTimeMs
	case "/":
		return DivisionTimeMs
	default:
		return 1000
	}
}

// buildExpressionTree builds an expression tree from tokens in Reverse Polish Notation.
func buildExpressionTree(tokens []calculator.Token) (*Node, error) {
	var stack []*Node
	for _, token := range tokens {
		if token.IsOperand {
			val, _ := token.GetOperand()
			stack = append(stack, &Node{IsLiteral: true, Value: val})
		} else if token.IsOperator {
			if len(stack) < 2 {
				return nil, errors.New("not enough operands")
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			node := &Node{
				IsLiteral: false,
				Operator:  token.Value.(string),
				Left:      left,
				Right:     right,
			}
			stack = append(stack, node)
		} else {
			return nil, errors.New("unexpected token")
		}
	}
	if len(stack) != 1 {
		return nil, errors.New("invalid expression")
	}
	return stack[0], nil
}
