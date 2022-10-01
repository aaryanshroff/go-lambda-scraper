package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gocolly/colly"
)

type Item struct {
	Title string
	URL   string
	City  string
	Price int
}

func HandleRequest() {

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// On every element which has class "result-row" call callback
	c.OnHTML(".result-row", func(e *colly.HTMLElement) {
		price, err := strconv.Atoi(normalizePrice(e.ChildText(".result-price")))
		if err != nil {
			log.Println("Error parsing price:", err)
			return
		}
		item := Item{
			Title: e.ChildText(".result-title"),
			URL:   e.ChildAttr("a.result-title", "href"),
			City:  "Kitchener",
			Price: price,
		}

		// Convert the item to AttributeValues.
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatalf("Got error marshalling new listing item: %s", err)
		}

		// Create item in table Scraper History
		tableName := "ScraperHistory"

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}

		fmt.Println("Succesffuly added '" + item.URL + "' to table" + tableName)

	})

	c.Visit("https://kitchener.craigslist.org/search/apa")
}

// normalizePrice removes the dollar sign and commas from a price string.
func normalizePrice(price string) string {
	return strings.ReplaceAll(strings.ReplaceAll(price, "$", ""), ",", "")
}
