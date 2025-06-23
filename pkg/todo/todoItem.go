package todo

type ItemStatus string

const (
	NotStarted ItemStatus = "not started"
	Started    ItemStatus = "started"
	Completed  ItemStatus = "completed"
)

type Item struct {
	Description string
	Status      ItemStatus
}
