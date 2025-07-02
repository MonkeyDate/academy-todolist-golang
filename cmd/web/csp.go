package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"context"
	"fmt"
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

func StartTodolistStoreActor() {
	go func() {
		for {
			select {
			case action := <-actionChannel:
				ctx := action.ctx

				switch action.actionType {
				case ActCreate:
					{
						if action.description == "" {
							action.description = "new-item-" + time.Now().Format(time.RFC3339)
						}

						todoList, err := common.LoadTodoList(ctx)
						if err != nil {
							action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't load list: %w", err)}
							return
						}

						todoList.Items = append(todoList.Items, todo.Item{Description: action.description, Status: action.status})

						err = common.SaveTodoList(ctx, todoList)
						if err != nil {
							action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: can't save list: %w", err)}
							return
						}

						action.resultChan <- ListActionResult{list: todoList}
						return
					}

				default:
					action.resultChan <- ListActionResult{err: fmt.Errorf("TODO List Store Actor: unsupported action : %d", action.actionType)}
					return
				}
			}
		}
	}()
}
