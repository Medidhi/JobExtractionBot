package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {

	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML(".jobsearch-SerpJobCard", func(e *colly.HTMLElement) {

		// Print link
		fmt.Printf("Link found: %q %q \n", e.Attr("h2"), e.Attr("summary"), e.Text)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.indeed.com/jobs?q=software%20internship&l=california")

}
