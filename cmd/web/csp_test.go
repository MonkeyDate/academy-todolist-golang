package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

const (
	numConcurrentCreates = 10
	itemDescriptionFmt   = "concurrent-item-%d"
)

func TestCspSupportsConcurrentCreates(t *testing.T) {
	t.Parallel()

	logger := slog.Default()
	ctx := common.SetLogger(context.Background(), logger)
	StartTodolistStoreActor(logger)

	var wg sync.WaitGroup

	for i := 0; i < numConcurrentCreates; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			desc := fmt.Sprintf(itemDescriptionFmt, i)
			result := CreateItem(ctx, desc, todo.NotStarted)
			if result.err != nil {
				t.Errorf("CreateItem failed: %v", result.err)
				return
			}

			found := false
			for _, item := range result.list.Items {
				if item.Description == desc {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Created item not found in returned list")
			}
		}()
	}

	wg.Wait()
}
