package main

import "github.com/aws/aws-lambda-go/lambda"

type InEvent struct {
	Total int `json:"total"`
}

type OutEvent struct {
	TaskDefinitions []TaskDefinition `json:"task_definitions"`
}

type TaskDefinition struct {
	Total int64 `json:"total"`
	Seg   int64 `json:"seg"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(e InEvent) (*OutEvent, error) {
	total := 4 // default
	if e.Total != 0 {
		total = e.Total
	}

	defs := make([]TaskDefinition, 0, total)
	for i := 0; i < total; i++ {
		defs = append(defs, TaskDefinition{
			Total: int64(total),
			Seg:   int64(i),
		})
	}
	return &OutEvent{
		TaskDefinitions: defs,
	}, nil
}
