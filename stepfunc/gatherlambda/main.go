package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type InEvent struct {
	Results []TaskResult `json:"task_results"`
}

type TaskResult struct {
	Count int64 `json:"count"`
}

type OutEvent struct {
	TotalCount int64 `json:"total_count"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(e InEvent) (*OutEvent, error) {
	totalCount := int64(0)
	for _, v := range e.Results {
		totalCount += v.Count
	}

	log.Printf("total: %d", totalCount)
	return &OutEvent{
		TotalCount: totalCount,
	}, nil
}
