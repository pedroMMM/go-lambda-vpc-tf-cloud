package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(ctx context.Context) (string, error) {
	log.Println("hi")
	return "output", nil
}

func main() {
	lambda.Start(LambdaHandler)
}
