package utils

import (
	"sort"

	"google.golang.org/api/tasks/v1"
)

func Sort(tasks []*tasks.Task, sortBy string) {
	switch sortBy {
	case "due":
		sort.SliceStable(tasks, func(i, j int) bool {
			if tasks[i].Due == "" {
				return false
			}

			if tasks[j].Due == "" {
				return true
			}

			return tasks[i].Due < tasks[j].Due
		})
	case "title":
		sort.SliceStable(tasks, func(i, j int) bool {
			return tasks[i].Title < tasks[j].Title
		})
	case "position":
	default:
		sort.SliceStable(tasks, func(i, j int) bool {
			return tasks[i].Position < tasks[j].Position && tasks[i].Parent == ""
		})
	}
}
