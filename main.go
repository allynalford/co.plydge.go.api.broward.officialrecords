package main

import (
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chromedp/chromedp"
)

var (
	_baseURL = "https://officialrecords.broward.org/AcclaimWeb/search/SearchTypeParcel"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	parcelID, ok := request.QueryStringParameters["SN"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Parcel ID", "1", parcelID)
	}

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		//log.Fatal(err)
		return GenerateErrorResponse(err.Error(), "2", "")
	}

	// run task list
	var html string

	err = c.Run(ctxt, submit(_baseURL, parcelID, &html))
	if err != nil {
		//log.Fatal(err)
		return GenerateErrorResponse(err.Error(), "3", "")
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		//log.Fatal(err)
		return GenerateErrorResponse(err.Error(), "4", "")
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		//log.Fatal(err)
		return GenerateErrorResponse(err.Error(), "5", "")
	}

	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)

	or := ParseRecordTable(doc)

	//fmt.Println(Marshal(or)) // Links:FooBarBazTEXT I WANT

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(Marshal(or)),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil

}

// ParseRecordTable Parse the records table of official records
func ParseRecordTable(doc *goquery.Document) model.OfficialRecords {

	//The Parent Struct returned
	or := model.OfficialRecords{}

	//Loop the Table Rows
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		//Init the Struct
		p := model.Parcel{}
		//Set the values by grabing each Table Data's contents
		p.ParcelID = tr.Find("td:nth-child(3)").Contents().Text()
		p.FirstDirectName = tr.Find("td:nth-child(4)").Contents().Text()
		p.FirstInDirectName = tr.Find("td:nth-child(5)").Contents().Text()
		p.BookType = tr.Find("td:nth-child(6)").Contents().Text()
		p.BookPage = tr.Find("td:nth-child(7)").Contents().Text()
		p.DateRecorded = tr.Find("td:nth-child(8)").Contents().Text()
		p.DocType = tr.Find("td:nth-child(9)").Contents().Text()
		p.InstrumentNumber = tr.Find("td:nth-child(10)").Contents().Text()
		p.Legal = tr.Find("td:nth-child(11)").Contents().Text()

		//Append to the array
		or.Records = append(or.Records, p)
	})

	return or
}

//submit page and forms
func submit(urlstr string, parcelid string, html *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitReady(`/html/body/div[2]/div/div/div/div[2]/form`),
		chromedp.Submit(`#btnButton`),
		chromedp.WaitVisible(`//input[@name="ParcelId"]`),
		chromedp.SendKeys(`//input[@name="ParcelId"]`, parcelid),
		chromedp.Click(`//*[@id="btnSearch"]`),
		chromedp.WaitReady(`//*[@id="0"]`),
		chromedp.OuterHTML(`//*[@id="RsltsGrid"]/div[4]/table`, html),
	}
}

// Marshal Convert BCPA	to string
func Marshal(or interface{}) string {
	b, err := json.Marshal(or)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "0"
	}
	//fmt.Println(string(b))

	return string(b)
}

// GenerateErrorResponse function to create base error message with events.APIGatewayProxyResponse
func GenerateErrorResponse(m string, c string, o string) (events.APIGatewayProxyResponse, error) {
	ge, err := model.GenerateErrorString(m, c, o)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(ge),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
