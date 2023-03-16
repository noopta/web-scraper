package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "net/http"
    "golang.org/x/net/html"
    "strings"
    "os"
    // "reflect"
    "encoding/json"	
    "strconv"
    openai "github.com/sashabaranov/go-openai"
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
    SectionID int `json:"sectionId"`
    SectionMapName string `json:"sectionMapName"`
    SectionType int `json:"sectionType"`
    Row string `json:"row"`
    SeatFrom string `json:"seatFrom"`
    SeatTo string `json:"@seatTo"`
    AvailableTickets int `json:"availableTickets"`
    AvailableQuantities []int `json:"availableQuantities"`
    TicketClass int `json:"ticketClass"`
    TicketClassName string `json:"ticketClassName"`
    BestSellingInSectionMessage SellingMessage `json:"bestSellingInSectionMessage"`
    RawPrice float64 `json:"rawPrice"`
    RawTicketPrice float64 `json:"rawTicketPrice"`
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

type AllStubHubData struct {
    Events []StubHubTicketGrid
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

func visitPage(inputLink string) []StubHubItem{
    mostProfitable := []float64{0.0, 0.0}
    currentMinDistance := 10000.0

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
    row := "15"
    section := "325"

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

    tempList := allData.Grid.Items

    i = 0
    // var sectionValue string
    // var rowValue string

    for i < len(tempList) {
        convertedString, err := strconv.Atoi(tempList[i].PriceWithFees[1:])
        
        if float64(convertedString) - tempList[i].RawPrice < currentMinDistance {
            currentMinDistance = float64(convertedString) - tempList[i].RawPrice
            mostProfitable[0] = tempList[i].RawPrice
            mostProfitable[1] = float64(convertedString)
            // sectionValue = tempList[i].Section
            // rowValue = tempList[i].Row
        }

        fmt.Println(tempList[i].Section + " " + tempList[i].Row)
        if(tempList[i].Section == section && tempList[i].Row == row) {
            fmt.Println("Found ticket")
            fmt.Println(tempList[i])
        }

        if err != nil {
            fmt.Println(err)
        }
        i++
    }

    return tempList
}

func callGPT() {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Describe the view and quality of the seats at section 323 and row 15 for Chicago Bulls home games from 2021",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func main() {

    callGPT()
    return

    // url := "https://www.stubhub.ca/toronto-raptors-tickets/performer/7549/"
    url := "https://www.stubhub.ca/chicago-bulls-tickets/performer/2863/"
    //   url := "https://www.stubhub.ca/chicago-bulls-tickets/"
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
        if strings.Contains(data[i], "at Chicago Bulls") && strings.Contains(data[i], "Minnesota Timberwolves") {
            data[i] = strings.ReplaceAll(data[i], " ", "")
            ioutil.WriteFile("data.txt", []byte(data[i]), 0644)
        }
        i++
    }

    // OPEN AI API KEY = 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

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

    i = 0
    // https://www.stubhub.ca/chicago-bulls-chicago-tickets-3-15-2023/event/150341877/?quantity=1&listingId=6143848653&listingQty=
    // we get a list of tickets with the rows so just find the match with the same row and section 

    for i < len(temp) {
        visitPage(temp[i].URL)

        // Section string `json:"section"`
        // SectionID int `json:"sectionId"`
        // SectionMapName string `json:"sectionMapName"`
        // SectionType int `json:"sectionType"`
        // Row string `json:"row"`
        i++
    }
}
