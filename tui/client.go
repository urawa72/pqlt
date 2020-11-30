package tui

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Client *dynamodb.DynamoDB

// unused
func NewClient() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

    svc := dynamodb.New(
		sess,
		aws.NewConfig(),
	)

	Client = svc
}
