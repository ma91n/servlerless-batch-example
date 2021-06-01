package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var db = dynamodb.New(session.Must(session.NewSession()))
var kc = kinesis.New(session.Must(session.NewSession()))

type Job struct {
	Total int64 `json:"total"`
	Seg   int64 `json:"seg"`
}

func main() {
	lambda.Start(Handle)
}

func Handle(ctx context.Context) error {

	resp, err := db.ScanWithContext(ctx, &dynamodb.ScanInput{
		Select:    aws.String(dynamodb.SelectCount),
		TableName: aws.String("TestTable"),
	})
	if err != nil {
		return err
	}

	total := 1
	if int(*resp.Count) > 1000 {
		total = int(*resp.Count) / 1000
	}

	for i := 0; i < total; i++ {
		log.Printf("i=%d\n", i)
		job := Job{
			Total: int64(total),
			Seg:   int64(i),
		}
		b, err := json.Marshal(job)
		if err != nil {
			return err
		}

		if _, err = kc.PutRecordWithContext(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String("scatter"),
			PartitionKey: aws.String(fmt.Sprintf("partitionKey-%d", i)),
			Data:         b,
		}); err != nil {
			return err
		}

	}
	log.Println("done")
	return nil
}
