package main

import (
	"github.com/aaryanshroff/go-lambda-scraper/pkg/scraper"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(scraper.HandleRequest)
}
