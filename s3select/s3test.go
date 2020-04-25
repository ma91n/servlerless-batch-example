package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
)

func main() {
	if err := ExecuteQuery(); err != nil {
		log.Fatal(err)
	}
}

var svc = s3.New(session.Must(session.NewSession(aws.NewConfig().WithRegion("ap-northeast-1"))))

type Resp struct {
	ItemNo  int64 `json:"tem_no"`
	ItemName string `json:"item_name"`
}

func ExecuteQuery() error {
	resp, err := svc.SelectObjectContent(&s3.SelectObjectContentInput{
		Bucket:         aws.String("<Your S3 Bucket>"),
		Key:            aws.String("<S3 Key Name>.csv"),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     aws.String("SELECT cast(item_no integer), item_name FROM s3object s WHERE cast(item_no integer) % 4 = 0"), // 4分割のうちから0~3を指定
		InputSerialization: &s3.InputSerialization{
			CompressionType: aws.String("NONE"),
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			JSON: &s3.JSONOutput{},
		},
	})
	if err != nil {
		return fmt.Errorf("svc.SelectObjectContent: %w", err)
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		switch v := event.(type) {
		case *s3.RecordsEvent:
			r := bufio.NewReader(bytes.NewReader(v.Payload))
			for {
				line, _, err := r.ReadLine()
				if err == io.EOF {
					break
				} else if err != nil {
					return fmt.Errorf("readLine :%w", err)
				}
				var resp Resp
				if err := json.Unmarshal(line, &resp); err != nil {
					return err
				}
				// TODO 何かしらのビジネスロジック
				fmt.Printf("%+v\n", resp)
			}
		}

	}

	if err := resp.EventStream.Err(); err != nil {
		return err
	}
	return nil
}
