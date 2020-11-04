package main

import "fmt"

type task struct {
	ID       string `json:"ID"`
	TaskType string `json:"TaskType"`
	Args     string `json:"Args"`
	Status   string `json:"Status"`
	Output   string `json:"Output"`
}

type tasks []*task

func createTask(taskType string, args string) *task {
	var t task

	t.ID = getNewID()
	t.TaskType = taskType
	t.Args = args
	t.Status = "Created"
	t.Output = ""

	return &t
}

func (ts *tasks) addNewTask(taskType string, args string) *task {
	newTask := createTask(taskType, args)
	*ts = append(*ts, newTask)

	return newTask
}

func (ts *tasks) appendTask(t *task) {
	*ts = append(*ts, t)
}

func (ts tasks) filterID(id string) *task {
	for i := range ts {
		if ts[i].ID == id {
			return ts[i]
		}
	}
	return nil
}

func (ts tasks) filterStatus(status string) tasks {
	var filteredTasks tasks
	for _, t := range ts {
		if t.Status == status {
			filteredTasks.appendTask(t)
		}
	}
	return filteredTasks
}

func (ts *tasks) setStatus(status string) {
	for _, t := range *ts {
		t.Status = status
	}
}

func (t *task) print() {
	fmt.Printf("ID: %v\n", t.ID)
	fmt.Printf("TaskType: %v\n", t.TaskType)
	fmt.Printf("Args: %v\n", t.Args)
	fmt.Printf("Status: %v\n", t.Status)
	fmt.Println()
}

func (ts *tasks) print() {
	fmt.Println("Tasks:")
	debugPrintBanner()

	if ts != nil {
		for _, t := range *ts {
			t.print()
		}
	} else {
		fmt.Println("None")
	}
	fmt.Println()
}
