package utils

import (
	"errors"
	"log"

	"google.golang.org/api/tasks/v1"
)

func GetTaskLists(srv *tasks.Service) ([]*tasks.TaskList, error) {
	r, err := srv.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	if len(r.Items) == 0 {
		return nil, errors.New("No Tasklist found")
	}
	return r.Items, nil
}
