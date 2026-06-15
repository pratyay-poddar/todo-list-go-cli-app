package todo

import (
	"fmt"
	"time"
)

///
///

type Item struct {
	Title       string
	Status      bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// ? creation of a custom type mean where we will pass this we can pass any named variable
// ? which have Item obj in there every index
type Items []Item

///

// add a task
func (items *Items) Add(title string) error {

	for i := 0; i < len(*items); i++ {
		if (*items)[i].Title == title {
			return fmt.Errorf("Task with same name exists")
		}
	}

	item := Item{
		Title:     title,
		Status:    false,
		CreatedAt: time.Now(),
	}
	*items = append(*items, item)

	return nil

}

///

// edit a task
func (items *Items) Edit(prevTitle string, title string) {
	for i := 0; i < len(*items); i++ {

		if (*items)[i].Title == prevTitle {
			(*items)[i].Title = title
			return
		}
	}
}

///

// delete a task
func (items *Items) Delete(title string) {
	for i := 0; i < len(*items); i++ {
		if (*items)[i].Title == title {
			*items = append((*items)[:i], (*items)[i+1:]...)
			return
		}
	}
}

// mark a task as complete
func (items *Items) Complete(title string) {
	for i := 0; i < len(*items); i++ {
		if (*items)[i].Title == title {
			(*items)[i].Status = true
			now := time.Now()
			(*items)[i].CompletedAt = &now
			return
		}
	}
}
