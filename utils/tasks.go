package utils

import (
	"errors"
	"log"

	"google.golang.org/api/tasks/v1"
)

func CreateTask(srv *tasks.Service, task *tasks.Task, tasklistID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Insert(tasklistID, task).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func GetTasks(srv *tasks.Service, id string) ([]*tasks.Task, error) {
	r, err := srv.Tasks.List(id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve tasks. %v", err)
	}
	if len(r.Items) == 0 {
		return nil, errors.New("No Tasks found")
	}
	return r.Items, nil
}

func GetTaskInfo(srv *tasks.Service, id string, taskID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Get(id, taskID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve tasks. %v", err)
	}
	return r, nil
}
