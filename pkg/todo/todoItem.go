package todo

type ItemStatus string

const (
	NotStarted ItemStatus = "not started"
	Started    ItemStatus = "started"
	Completed  ItemStatus = "completed"
)

type Item struct {
	ID          string
	Description string
	Status      ItemStatus
}
