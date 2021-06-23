package api

import (
	"errors"
	"log"

	"google.golang.org/api/tasks/v1"
)

type TaskList []tasks.TaskList

func (e TaskList) Len() int {
	return len(e)
}

func (e TaskList) Less(i, j int) bool {
	return e[i].Title < e[j].Title
}

func (e TaskList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func GetTaskLists(srv *tasks.Service) ([]tasks.TaskList, error) {
	r, err := srv.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	var list []tasks.TaskList

	if len(r.Items) == 0 {
		return nil, errors.New("no Tasklist found")
	}

	for _, item := range r.Items {
		list = append(list, *item)
	}

	return list, nil
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
