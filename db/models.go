package db

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	Login     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // bcrypt hash
	CreatedAt time.Time
}

type Expression struct {
	ID         string `gorm:"primarykey"` // uuid
	Expr       string
	Status     string // pending, done
	Result     *float64
	RootTaskID string
	UserID     uint `gorm:"index"`
	CreatedAt  time.Time
}

type Task struct {
	ID            string `gorm:"primarykey"`
	ExpressionID  string `gorm:"index"`
	Operator      string
	Arg1          *float64
	Arg2          *float64
	DepTask1      string
	DepTask2      string
	OperationTime int
	Status        string // pending, running, done
	Result        *float64
	CreatedAt     time.Time
}
