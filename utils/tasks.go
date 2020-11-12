package utils

import (
	"errors"
	"log"

	"google.golang.org/api/tasks/v1"
)

//CreateTask used to create tasks
func CreateTask(srv *tasks.Service, task *tasks.Task, tasklistID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Insert(tasklistID, task).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

//GetTasks used to retreive tasks
func GetTasks(srv *tasks.Service, id string, showCompleted bool) ([]*tasks.Task, error) {
	r, err := srv.Tasks.List(id).ShowHidden(showCompleted).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve tasks. %v", err)
	}
	if len(r.Items) == 0 {
		return nil, errors.New("No Tasks found")
	}
	return r.Items, nil
}

//GetTaskInfo to get more info about a task
func GetTaskInfo(srv *tasks.Service, id string, taskID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Get(id, taskID).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

//UpdateTask used to update task data
func UpdateTask(srv *tasks.Service, t *tasks.Task, tListID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Patch(tListID, t.Id, t).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

//DeleteTask used to delete a task
func DeleteTask(srv *tasks.Service, id string, tid string) error {
	err := srv.Tasks.Delete(tid, id).Do()
	return err
}
