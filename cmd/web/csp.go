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
	err  error
	list todo.List
}

func BeginReadItems(ctx context.Context) chan ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType: ActRead,
	}
	actionChannel <- request

	return request.resultChan
}

func BeginCreateItem(ctx context.Context, description string, status todo.ItemStatus) chan ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType:  ActCreate,
		description: description, status: status,
	}
	actionChannel <- request

	return request.resultChan
}

func BeginUpdateItem(ctx context.Context, ID string, description string, status todo.ItemStatus) chan ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType:  ActUpdate,
		ID:          ID,
		description: description, status: status,
	}
	actionChannel <- request

	return request.resultChan
}

func BeginDeleteItem(ctx context.Context, ID string) chan ListActionResult {
	resultChan := make(chan ListActionResult, 1)
	request := actionRequest{
		ctx: ctx, resultChan: resultChan,
		actionType: ActDelete,
		ID:         ID,
	}
	actionChannel <- request

	return request.resultChan
}

func StartTodolistStoreActor(logger *slog.Logger) {
	go func() {
		logger.Info("StartTodolistStoreActor: enter")
		defer logger.Info("StartTodolistStoreActor: leave")

		for {
			processAction(<-actionChannel)
		}
	}()
}

func processAction(action actionRequest) {
	ctx := action.ctx

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't load list: %w", err)}
		return
	}

	switch action.actionType {
	case ActRead:
		{
		}

	case ActCreate:
		{
			if action.description == "" {
				action.description = "new-item-" + time.Now().Format(time.RFC3339)
			}

			id := uuid.New().String()
			todoList.Items = append(todoList.Items, todo.Item{ID: id, Description: action.description, Status: action.status})

			err = common.SaveTodoList(ctx, todoList)
			if err != nil {
				action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't save list: %w", err)}
				return
			}
		}

	case ActUpdate:
		{
			var toUpdate *todo.Item
			for i, item := range todoList.Items {
				if item.ID == action.ID {
					toUpdate = &todoList.Items[i]
				}
			}

			if toUpdate != nil {
				toUpdate.Status = action.status
				if action.description != "" {
					toUpdate.Description = action.description
				}

				err = common.SaveTodoList(ctx, todoList)
				if err != nil {
					action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't save list: %w", err)}
					return
				}
			}
		}

	case ActDelete:
		{
			for i, item := range todoList.Items {
				if item.ID == action.ID {
					todoList.Items = append(todoList.Items[:i], todoList.Items[i+1:]...)

					err = common.SaveTodoList(ctx, todoList)
					if err != nil {
						action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't save list: %w", err)}
						return
					}

					return
				}
			}
		}
	}

	action.resultChan <- ListActionResult{list: todoList}
}
