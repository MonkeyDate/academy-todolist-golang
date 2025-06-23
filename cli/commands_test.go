package cli

import (
	"academy-todo/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func setup(t *testing.T) {
//	// setup work
//	t.Cleanup(func() {
//		// teardown code here
//	})
//}

func Test_AddItemToListCommand(t *testing.T) {
	t.Run("should add item to end of already populated list", func(t *testing.T) {
		commandLine := []string{"add", "-d=third_item_started", "-started"}
		startingList := setupPrepopulatedStartingList()

		list, err := addItemToListCommand(startingList, commandLine[1:])

		expectedItem := models.TodoItem{Status: models.Started, Description: "third_item_started"}
		assert.Nil(t, err, "should not return error")
		assert.Equal(t, len(list), 3, "list should have 3 items")
		assert.ObjectsAreEqualValues(expectedItem, list[2])
	})

	t.Run("should add item with default values if values are not specified", func(t *testing.T) {
		commandLine := []string{"add"}
		startingList := setupPrepopulatedStartingList()

		list, _ := addItemToListCommand(startingList, commandLine[1:])

		expectedItem := models.TodoItem{Status: models.NotStarted, Description: "new-item"}
		assert.ObjectsAreEqualValues(expectedItem, list[2])
	})
}

func setupPrepopulatedStartingList() (startingList []models.TodoItem) {
	startingList = make([]models.TodoItem, 2)
	startingList[0] = models.TodoItem{Status: models.NotStarted, Description: "item-1 not started"}
	startingList[1] = models.TodoItem{Status: models.Started, Description: "item-2 started"}

	return
}
