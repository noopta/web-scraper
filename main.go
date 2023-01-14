package main

import (
    "fmt"
	"net/http"
    "github.com/gocolly/colly"
)

type userAgentRoundTripper struct {
    userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
    r.Header.Set("User-Agent", rt.userAgent)
    return http.DefaultTransport.RoundTrip(r)
}

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains("stubhub.ca"),
    )

	
    // Add custom headers
    c.WithTransport(&http.Transport{
        // set custom User-Agent header
        RoundTripper: &userAgentRoundTripper{
            userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
        },
    })

    // Find and visit all links
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        c.Visit(e.Request.AbsoluteURL(link))
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    c.Visit("https://www.stubhub.ca/toronto-raptors-tickets/performer/7549/")
}
