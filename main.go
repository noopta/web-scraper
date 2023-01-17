package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "golang.org/x/net/html"
    "strings"
    "os"
    // "reflect"
    "encoding/json"	
)

type TicketData struct {
    Name string  `json:"name"`
    URL string  `json:"url"`
    StartDate string  `json:"startDate"`
    DoorTime string `json:"doorTime"`
}

type TicketGraph struct {
    Context string `json:"@context"`
    AllTickets []TicketData `json:"@graph"`
}

type StubHubTicket struct {
    SectionNumber int 
    RowNumber int
    Price int
    TicketLink string
}

type StubHubTicketGrid struct {
    AppName string `json:"appName"`
    Grid ItemObject `json:"grid"`
}

type ItemObject struct {
    Items []StubHubItem `json:"items"`
}

type StubHubItem struct {
    ID int `json:"id"`
    Section string `json:"section"`
    SectionID string `json:"sectionId"`
    SectionMapName string `json:"sectionMapName"`
    SectionType int `json:"sectionType"`
    Row string `json:"row"`
    SeatFrom string `json:"seatFrom"`
    SeatTo string `json:"@seatTo"`
    AvailableTickets int `json:"availableTickets"`
    AvailableQuantities int `json:"availableQuantities"`
    TicketClass string `json:"ticketClass"`
    TicketClassName string `json:"ticketClassName"`
    BestSellingInSectionMessage SellingMessage `json:"bestSellingInSectionMessage"`
    Price string `json:"price"`
    PriceWithFees string `json:"priceWithFees"`
    ListingCurrencyCode string `json:"listingCurrencyCode"`
    QualityRank int `json:"qualityRank"`
}

type SellingMessage struct {
    Message string `json:"message"`
    Qualifier string `json:"qualifier"`
    Disclaimer string `json:"disclaimer"`
    HasValue bool `json:"hasValue"`
    FeatureTrackingKey string `json:"featureTrackingKey"`
}

func parse(text string) (data []string) {

    tkn := html.NewTokenizer(strings.NewReader(text))

    var vals []string

    var isLi bool

    for {

        tt := tkn.Next()

        switch {

        case tt == html.ErrorToken:
            return vals

        case tt == html.StartTagToken:

            t := tkn.Token()
            isLi = t.Data == "script"

        case tt == html.TextToken:

            t := tkn.Token()

            if isLi {
                vals = append(vals, t.Data)
            }

            isLi = false
        }
    }
}

func visitPage(inputLink string) {
    url := inputLink
    fmt.Printf("HTML code of %s ...\n", url)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // fmt.Println(html)
    ioutil.WriteFile("StubHubEvent.txt", html, 0644)

    htmlString, err := ioutil.ReadFile("StubHubEvent.txt")

    if err != nil {
        fmt.Println(err)
    }

    data := parse(string(htmlString[:]))

    i := 0

    for i < len(data) {
        if strings.Contains(data[i], "sectionMapName") {
            data[i] = strings.ReplaceAll(data[i], " ", "")
            ioutil.WriteFile("moreData.txt", []byte(data[i]), 0644)
        }
        i++
    } 

    fileText, err := os.ReadFile("moreData.txt")

    if err != nil {
        fmt.Println(err)
    }

    convertedString := string(fileText)

    allData := StubHubTicketGrid{}

    err = json.Unmarshal([]byte(convertedString), &allData)
    if err != nil {
        // panic
        fmt.Println(err)
    }

    fmt.Println(allData)
}

func main() {
    url := "https://www.stubhub.ca/toronto-raptors-tickets/performer/7549/"
    // url := "stubhub.ca/milwaukee-bucks-milwaukee-tickets-1-17-2023/event/150337076/?quantity=1"
    fmt.Printf("HTML code of %s ...\n", url)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    // fmt.Printf("%s\n", html)
    ioutil.WriteFile("ticketText.txt", html, 0644)

    htmlString, err := ioutil.ReadFile("ticketText.txt")

    if err != nil {
        fmt.Println(err)
    }

    data := parse( string(htmlString[:]))

    i := 1

    for i < len(data) {
        if strings.Contains(data[i], "at Toronto Raptors") {
            data[i] = strings.ReplaceAll(data[i], " ", "")
            ioutil.WriteFile("data.txt", []byte(data[i]), 0644)
        }
        i++
    }

    fileText, err := os.ReadFile("data.txt")

    if err != nil {
        fmt.Println(err)
    }

    convertedString := string(fileText)

    allData := TicketGraph{}

    err = json.Unmarshal([]byte(convertedString), &allData)
    if err != nil {
        // panic
        fmt.Println(err)
    }

    temp := allData.AllTickets

    i = 0

    for i < len(temp) {
       temp[i].URL = "https://www.stubhub.ca" + temp[i].URL + "?quantity=1"
        i++
    }

    // visit each URL with a ticket quantity of 1 
    // get the price, and section number 
    // write the data to a text file similar to how we already did 
    // store the price, section number, and row I guess to a struct for 
    // all the StubHub tickets for this game

    // we want to create a list of tickets 
    // each ticket has the following information
    // section number, price, row, link to buy

    visitPage(temp[0].URL)
}
