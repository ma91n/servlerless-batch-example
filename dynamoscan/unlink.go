package unlink

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

var db = dynamodb.New(session.Must(session.NewSession(aws.NewConfig().WithRegion("ap-northeast-1"))))

type Resp struct {
	StartKey map[string]*dynamodb.AttributeValue
	Payload  []Ncu
}

type Report struct {
	// 何かしらの実行結果
}

func main() {
	ctx := context.Background()
	parallel := 4
	if err := Do(ctx, parallel); err != nil {
		log.Fatal(err)
	}
}

func Do(ctx context.Context, parallel int) error {
	reports := make([]*Report, 0, parallel)
	m := sync.Mutex{}

	eg := errgroup.Group{}
	for i := 0; i < parallel; i++ {
		i := i
		eg.Go(func() error {
			err := ScanAndLogic(ctx, int64(parallel), int64(i))
			if err != nil {
				return err
			}

			m.Lock()
			reports = append(reports, &Report{}) // なにか実行結果を返したいとき
			m.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("checkUnreach failed: %w", err)
	}

	fmt.Printf("reporte: %+v", reports)
	return nil
}

func ScanAndLogic(ctx context.Context, total, seg int64) error {
	var startKey map[string]*dynamodb.AttributeValue // 初回はnilでOK
	for {
		resp, startKey, err := ScanSegment(ctx, total, seg, startKey)
		if err != nil {
			return fmt.Errorf("ScanSegment: %w", err)
		}

		// TODO respに対して何かしらのビジネスロジック
		fmt.Printf("resp: %+v\n", resp)

		if len(startKey) == 0 {
			break // 続きが無いということなので終了
		}
	}
	return nil
}

func ScanSegment(ctx context.Context, total, seg int64, startKey map[string]*dynamodb.AttributeValue) ([]Resp, map[string]*dynamodb.AttributeValue, error) {
	out, err := db.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName:         aws.String("<DynamoDB Scan Table>"),
		TotalSegments:     aws.Int64(total), // セグメントへの分割数
		Segment:           aws.Int64(seg),   // 処理番号（0,1,2,3を指定）
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
