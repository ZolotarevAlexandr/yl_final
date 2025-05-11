package orchestrator

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ZolotarevAlexandr/yl_final/calculator/calculator"
	"github.com/ZolotarevAlexandr/yl_final/db"
)

// CreateExpression
//  1. parse & build the expression tree
//  2. insert an Expression row
//  3. walk the tree and insert Task rows
//  4. update the Expression.RootTaskID
func CreateExpression(exprStr string, userID uint) (string, error) {
	// tokenize & RPN
	tokens, err := calculator.Tokenize(exprStr)
	if err != nil {
		return "", err
	}
	rpn, err := calculator.ShuntingYard(tokens)
	if err != nil {
		return "", err
	}
	rootNode, err := buildExpressionTree(rpn)
	if err != nil {
		return "", err
	}

	exprID := uuid.NewString()
	expr := &db.Expression{
		ID:     exprID,
		Expr:   exprStr,
		Status: "pending",
		UserID: userID,
	}

	tx := db.DB.Begin()
	if err := tx.Create(expr).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// recursively persist tasks; returns the final/root task ID
	rootTaskID, err := persistNode(tx, exprID, rootNode)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// update the expression with the root task
	if err := tx.Model(expr).
		Update("root_task_id", rootTaskID).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	return exprID, tx.Commit().Error
}

// persistNode writes one operation-node as a Task row,
// recurses into children, wires up dependencies.
// If node.IsLiteral, returns "" (no task).
func persistNode(tx *gorm.DB, exprID string, node *Node) (taskID string, _ error) {
	if node.IsLiteral {
		return "", nil
	}
	// left
	dep1 := ""
	if !node.Left.IsLiteral {
		dep1, _ = persistNode(tx, exprID, node.Left)
	}
	// right
	dep2 := ""
	if !node.Right.IsLiteral {
		dep2, _ = persistNode(tx, exprID, node.Right)
	}
	// literal operands?
	var a1, a2 *float64
	if node.Left.IsLiteral {
		a1 = &node.Left.Value
	}
	if node.Right.IsLiteral {
		a2 = &node.Right.Value
	}

	t := &db.Task{
		ID:            uuid.NewString(),
		ExpressionID:  exprID,
		Operator:      node.Operator,
		Arg1:          a1,
		Arg2:          a2,
		DepTask1:      dep1,
		DepTask2:      dep2,
		OperationTime: getOperationTime(node.Operator),
		Status:        "pending",
	}
	if err := tx.Create(t).Error; err != nil {
		return "", err
	}
	return t.ID, nil
}
