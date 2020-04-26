package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type InEvent struct {
	Total int64 `json:"total"`
	Seg   int64 `json:"seg"`
}

type OutEvent struct {
	Count int64 `json:"count"`
}

type Resp struct {
	// Any fields
}

var (
	db    = dynamodb.New(session.Must(session.NewSession(aws.NewConfig().WithRegion("ap-northeast-1"))))
	table = os.Getenv("DYNAMO_TABLE")
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(e InEvent) (*OutEvent, error) {
	log.Printf("InEvent: %+v", e)
	total, seg := e.Total, e.Seg

	var startKey map[string]*dynamodb.AttributeValue
	count := 0
	for {
		resp, sk, err := ScanSegment(context.Background(), total, seg, startKey)
		if err != nil {
			log.Printf("scan error: %v", err)
			return nil, fmt.Errorf("ScanSegment: %w", err)
		}
		count += len(resp)
		log.Printf("count: %v", count)

		startKey = sk
		if len(startKey) == 0 {
			break
		}
	}

	log.Printf("total count: %v", count)
	return &OutEvent{
		Count: int64(count),
	}, nil
}

func ScanSegment(ctx context.Context, total, seg int64, startKey map[string]*dynamodb.AttributeValue) ([]Resp, map[string]*dynamodb.AttributeValue, error) {
	out, err := db.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName:         aws.String(table),
		TotalSegments:     aws.Int64(total),
		Segment:           aws.Int64(seg),
		ExclusiveStartKey: startKey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("db.ScanWithContext: %w", err)
	}

	var resp []Resp
	if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &resp); err != nil {
		return nil, nil, fmt.Errorf("dynamodbattribute.UnmarshalListOfMaps: %w", err)
	}
	return resp, out.LastEvaluatedKey, nil
}
