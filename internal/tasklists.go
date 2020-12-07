package internal

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

func UpdateTaskList(srv *tasks.Service, tl *tasks.TaskList) (*tasks.TaskList, error) {
	r, err := srv.Tasklists.Patch(tl.Id, tl).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func DeleteTaskList(srv *tasks.Service, tID string) error {
	err := srv.Tasklists.Delete(tID).Do()
	return err
}
