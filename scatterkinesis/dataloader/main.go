package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
)

func main() {

	var db = dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	for i := 0; i < 1000*1000; i++ {
		av, _ := dynamodbattribute.MarshalMap(map[string]interface{}{"id": strconv.Itoa(i), "body": "example"})
		if _, err := db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String("TestTable"),
			Item:     av ,
		}); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("done")
}
