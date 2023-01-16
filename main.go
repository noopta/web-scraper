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

func main() {
    url := "https://www.stubhub.ca/toronto-raptors-tickets/performer/7549/"
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

    // fmt.Println(data)

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

    fmt.Println(len(allData.AllTickets))

    temp := allData.AllTickets

    i = 0

    for i < len(temp) {
       temp[i].URL = "stubhub.ca" + temp[i].URL + "?quantity=1"
       fmt.Println(temp[i].URL)
        i++
    }

    // visit each URL with a ticket quantity of 1 
    // get the price, and section number 
    // write the data to a text file similar to how we already did 
    // store the price, section number, and row I guess to a struct for 
    // all the StubHub tickets for this game
}
