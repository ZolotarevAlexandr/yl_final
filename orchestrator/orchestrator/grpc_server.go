package orchestrator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/ZolotarevAlexandr/yl_final/db"
	pb "github.com/ZolotarevAlexandr/yl_final/grpc"
)

// taskServer implements pb.TaskServiceServer
type taskServer struct {
	pb.UnimplementedTaskServiceServer
}

// GetTask finds one task in status="pending" whose dependencies are done,
// marks it "running", and returns its payload.
// If no ready task exists → NotFound.
func (s *taskServer) GetTask(ctx context.Context, _ *emptypb.Empty) (*pb.TaskResponse, error) {
	// wrap in a transaction to avoid races
	tx := db.DB.Begin()
	defer func() {
		// if we didn't commit, rollback
		_ = tx.Rollback()
	}()

	var tasks []db.Task
	if err := tx.Where("status = ?", "pending").Find(&tasks).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "db lookup failed: %v", err)
	}

	for _, t := range tasks {
		ready := true

		// check dep1
		if t.DepTask1 != "" {
			var d1 db.Task
			if err := tx.First(&d1, "id = ?", t.DepTask1).Error; err != nil || d1.Status != "done" {
				ready = false
			}
		}
		// check dep2
		if ready && t.DepTask2 != "" {
			var d2 db.Task
			if err := tx.First(&d2, "id = ?", t.DepTask2).Error; err != nil || d2.Status != "done" {
				ready = false
			}
		}

		if !ready {
			continue
		}

		// mark this task running
		if err := tx.Model(&db.Task{}).
			Where("id = ?", t.ID).
			Update("status", "running").Error; err != nil {
			return nil, status.Errorf(codes.Internal, "could not mark running: %v", err)
		}
		if err := tx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "tx commit failed: %v", err)
		}

		// return the task
		return &pb.TaskResponse{
			Id:            t.ID,
			Arg1:          *t.Arg1,
			Arg2:          *t.Arg2,
			Operation:     t.Operator,
			OperationTime: int32(t.OperationTime),
		}, nil
	}

	// nothing ready
	return nil, status.Error(codes.NotFound, "no task available")
}

// SubmitTask marks a running task done, records its result,
// and if it’s the root task, also marks the expression done.
func (s *taskServer) SubmitTask(ctx context.Context, in *pb.TaskResult) (*pb.ResultAck, error) {
	tx := db.DB.Begin()
	defer func() { _ = tx.Rollback() }()

	var task db.Task
	if err := tx.First(&task, "id = ?", in.Id).Error; err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	if task.Status != "running" {
		return nil, status.Errorf(codes.FailedPrecondition, "task not in running state: %s", task.Status)
	}

	// update the task
	if err := tx.Model(&db.Task{}).
		Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"status": "done",
			"result": in.Result,
		}).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "could not update task: %v", err)
	}

	// if this is the root task, update the expression
	var expr db.Expression
	if err := tx.First(&expr, "root_task_id = ?", task.ID).Error; err == nil {
		_ = tx.Model(&db.Expression{}).
			Where("id = ?", expr.ID).
			Updates(map[string]interface{}{
				"status": "done",
				"result": in.Result,
			}).Error
	}
	// Note: if expression not found, we ignore silently-could be a sub‐expression.

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "tx commit failed: %v", err)
	}

	return &pb.ResultAck{Status: "ok"}, nil
}
