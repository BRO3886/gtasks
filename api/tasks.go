package api

import (
	"errors"
	"fmt"

	"google.golang.org/api/tasks/v1"
)

// CreateTask used to create tasks
func CreateTask(srv *tasks.Service, task *tasks.Task, tasklistID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Insert(tasklistID, task).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTasks used to retreive tasks
func GetTasks(srv *tasks.Service, id string, includeCompleted bool) ([]*tasks.Task, error) {
	var allTasks []*tasks.Task
	pageToken := ""

	for {
		call := srv.Tasks.List(id).ShowHidden(includeCompleted).MaxResults(100)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		r, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve tasks: %v", err)
		}

		allTasks = append(allTasks, r.Items...)

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	if len(allTasks) == 0 {
		return nil, errors.New("no Tasks found")
	}

	if includeCompleted {
		return allTasks, nil
	}

	var list []*tasks.Task
	for _, task := range allTasks {
		if task.Status != "completed" {
			list = append(list, task)
		}
	}
	return list, nil
}

func MakeMap(taskList []*tasks.Task) map[string]tasks.Task {
	m := make(map[string]tasks.Task)
	for _, t := range taskList {
		m[t.Id] = *t
	}
	return m
}

// GetTaskInfo to get more info about a task
func GetTaskInfo(srv *tasks.Service, id string, taskID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Get(id, taskID).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateTask used to update task data
func UpdateTask(srv *tasks.Service, t *tasks.Task, tListID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Patch(tListID, t.Id, t).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteTask used to delete a task
func DeleteTask(srv *tasks.Service, id string, tid string) error {
	err := srv.Tasks.Delete(tid, id).Do()
	return err
}
