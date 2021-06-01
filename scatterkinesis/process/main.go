package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Job struct {
	Total int64 `json:"total"`
	Seg   int64 `json:"seg"`
}

func main() {
	lambda.Start(Handle)
}

func Handle(ctx context.Context, e events.KinesisEvent) error {
	for _, r := range e.Records {
		var job Job

		if err := json.Unmarshal(r.Kinesis.Data, &job); err != nil {
			return err
		}

		count, err := executeBizLogic(ctx, job.Total, job.Seg)
		if err != nil {
			return err
		}

		log.Printf("count: %v", count)
	}
	return nil
}

var db = dynamodb.New(session.Must(session.NewSession()))

func executeBizLogic(ctx context.Context, total, seg int64) (int64, error) {
	out, err := db.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName:     aws.String("TestTable"),
		TotalSegments: aws.Int64(total),
		Segment:       aws.Int64(seg),
		Select:        aws.String(dynamodb.SelectCount),
	})
	if err != nil {
		return 0, fmt.Errorf("db.ScanWithContext: %w", err)
	}

	return *out.Count, nil
}
