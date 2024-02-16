package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"github.com/gocolly/colly"
)


// takes a collector, instantiates callback functions for them
func PrintScraper(col *colly.Collector) {

	col.OnRequest(func(r *colly.Request) {

		fmt.Println("Visiting: ", r.URL)
	})


	col.OnError(func(r *colly.Response, err error) {

		fmt.Println("Something went wrong: ", err)
	})


	col.OnResponse(func(r *colly.Response) {

		fmt.Println("Page visited: ", r.Request.URL)
	})


	col.OnHTML("a", func(e *colly.HTMLElement) {

		log.Println("%v", e.Attr("href"))
	})


	col.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " Scraped!")

	})
}


// fields of data we want to collect
type WebProduct struct {
	image, name, price string
}

//creates everything we want to do with the file into one nice big function :)
func CSVFILE(wp []WebProduct, col *colly.Collector) {

	col.OnHTML("li.product", func(e *colly.HTMLElement) {

		webProd := WebProduct{}
		fmt.Println("Gathering HTML Elements: ", e)
		webProd.image = e.ChildAttr("img", "src")
		webProd.price = e.ChildText(".price")
		webProd.name = e.ChildText("h2")

		//adds all the scraped data to the list of products
		wp = append(wp, webProd)
	})


	file, err := os.Create("Scraper.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()


	//file writer initialization
	write := csv.NewWriter(file)


	//defining CSV headers
	headers := []string{

		"image",
		"name",
		"price",
	}


	//write column headers
	write.Write(headers)


	//for loop to add each product to CSV output file
	for _, web := range wp {

		record := []string{
			web.image,
			web.name,
			web.price,
		}
		write.Write(record)
	}

	
	defer write.Flush()
}

func main() {
	fmt.Println("Hello world!")

	c := colly.NewCollector()
	var webProduct []WebProduct

	c.Visit("https://scrapeme.live/shop/")

	PrintScraper(c)

	CSVFILE(webProduct, c)
}
