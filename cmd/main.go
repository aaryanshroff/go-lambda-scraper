package main

import (
	"github.com/aaryanshroff/go-lambda-scraper/pkg/scraper"
)

func main() {
	scraper.HandleRequest()
	// lambda.Start(scraper.HandleRequest)
}
