package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

// TODO: only CreateItem has full context handling

const (
	errCantLoadListFmt = "TODO List Store Actor: can't load list: %w"
	errCantSaveListFmt = "TODO List Store Actor: can't save list: %w"
	errItemNotFoundFmt = "TODO List Store Actor: item not found ID: %s"
)

const (
	ActRead   int = 1
	ActCreate int = 2
	ActUpdate int = 3
	ActDelete int = 4
)

type actionRequest struct {
	ctx        context.Context
	resultChan chan ListActionResult

	actionType int // ActRead etc

	description string
	status      todo.ItemStatus
	ID          string
}

var actionChannel = make(chan actionRequest)

type ListActionResult struct {
	list      todo.List
	createdID string

	err        error
	isApiError bool
}

func ReadItems(ctx context.Context) ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType: ActRead,
	}
	actionChannel <- request
	return <-resultChan
}

func CreateItem(ctx context.Context, description string, status todo.ItemStatus) ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx:         ctx,
		resultChan:  resultChan,
		actionType:  ActCreate,
		description: description,
		status:      status,
	}

	select {
	case actionChannel <- request:
		select {
		case result := <-resultChan:
			return result
		case <-ctx.Done():
			return ListActionResult{err: fmt.Errorf("CreateItem: context Done: %w", ctx.Err())}
		}

	case <-ctx.Done():
		return ListActionResult{err: fmt.Errorf("CreateItem: context Done: %w", ctx.Err())}

		// this will cause an error if the channel is blocked with another request
		//default:
		//	return ListActionResult{err: fmt.Errorf("CreateItem: Action channel full")}
	}
}

func UpdateItem(ctx context.Context, ID string, description string, status todo.ItemStatus) ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType:  ActUpdate,
		ID:          ID,
		description: description, status: status,
	}
	actionChannel <- request
	return <-resultChan
}

// DeleteItem returns success for an item not found scenario
func DeleteItem(ctx context.Context, ID string) ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType: ActDelete,
		ID:         ID,
	}
	actionChannel <- request
	return <-resultChan
}

func StartTodolistStoreActor(logger *slog.Logger) {
	// TODO: no graceful shutdown mechanism

	go func() {
		logger.Info("StartTodolistStoreActor: enter")
		defer logger.Info("StartTodolistStoreActor: leave")

		for {
			action := <-actionChannel
			result := processAction(action)
			action.resultChan <- result
		}
	}()
}

func processAction(action actionRequest) ListActionResult {
	ctx := action.ctx

	var createdId string
	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		return ListActionResult{err: fmt.Errorf(errCantLoadListFmt, err)}
	}

	switch action.actionType {
	case ActRead:
		{
		}

	case ActCreate:
		if action.description == "" {
			action.description = "new-item-" + time.Now().Format(time.RFC3339)
		}

		createdId = uuid.New().String()
		todoList.Items = append(todoList.Items, todo.Item{ID: createdId, Description: action.description, Status: action.status})

		err = common.SaveTodoList(ctx, todoList)
		if err != nil {
			return ListActionResult{err: fmt.Errorf(errCantSaveListFmt, err)}
		}

	case ActUpdate:
		var toUpdate *todo.Item
		for i, item := range todoList.Items {
			if item.ID == action.ID {
				toUpdate = &todoList.Items[i]
				break
			}
		}

		if toUpdate != nil {
			toUpdate.Status = action.status
			if action.description != "" {
				toUpdate.Description = action.description
			}

			err = common.SaveTodoList(ctx, todoList)
			if err != nil {
				return ListActionResult{err: fmt.Errorf(errCantSaveListFmt, err)}
			}
		} else {
			return ListActionResult{err: fmt.Errorf(errItemNotFoundFmt, action.ID)}
		}

	case ActDelete:
		for i, item := range todoList.Items {
			if item.ID == action.ID {
				todoList.Items = append(todoList.Items[:i], todoList.Items[i+1:]...)

				err = common.SaveTodoList(ctx, todoList)
				if err != nil {
					return ListActionResult{err: fmt.Errorf(errCantSaveListFmt, err)}
				}
			}
		}
	}

	return ListActionResult{list: todoList, createdID: createdId}
}
