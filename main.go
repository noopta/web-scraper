package main

import (
    "context"
    "bufio"
    "sync"
    "fmt"
    "io/ioutil"
    "net/http"
    "golang.org/x/net/html"
    "strings"
    "os"
    // "reflect"
    "encoding/json"	
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

var homeTeam string
var awayTeam string
var gameDate string
var ticketQuantity string
var sectionVal string
var rowVal string

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

func visitPage(inputLink string) StubHubItem{

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
    // row := "11"
    // section := "332"

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
    
    for i = 0; i < len(tempList); i++ {
        fmt.Println(tempList[i].Section)
        if(tempList[i].Section == sectionVal && tempList[i].Row == rowVal) {
            fmt.Println("Found ticket")
            fmt.Println()
            return tempList[i]
        }

        if err != nil {
            fmt.Println(err)
        }
    }

    return StubHubItem{}
}

func callGPT() {
    fmt.Println("calling Chat GPT")
    fmt.Println()
    
	client := openai.NewClient("")
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

    getVividSeatsTickets()

    return 
    fmt.Println("Enter the home team name: ")
    scanner := bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        homeTeam = scanner.Text()
    }

    fmt.Println("Enter the away team name: ")

    scanner = bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        awayTeam = scanner.Text()
    }

    fmt.Println("Enter the game date: ")

    scanner = bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        gameDate = scanner.Text()
    }

    fmt.Println("Enter the section number: ")
    scanner = bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        sectionVal = scanner.Text()
    }

    fmt.Println("Enter the row number: ")
    scanner = bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        rowVal = scanner.Text()
    }

    fmt.Println("Enter the ticket quantity: ")
    scanner = bufio.NewScanner(os.Stdin)

    if scanner.Scan() {
        ticketQuantity = scanner.Text()
    }

    fmt.Println(homeTeam + " " + awayTeam + " " + gameDate + " " + sectionVal + " " + rowVal + " " + ticketQuantity)

    go callGPT()

    url := "https://www.stubhub.ca/chicago-bulls-tickets/performer/2863/"
    fmt.Printf("HTML code of %s ...\n", url)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("err 1")
        panic(err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("err 2")
        panic(err)
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("err 3")
        panic(err)
    }
    // fmt.Printf("%s\n", html)
    ioutil.WriteFile("ticketText.txt", html, 0644)

    htmlString, err := ioutil.ReadFile("ticketText.txt")

    if err != nil {
        fmt.Println("err 4")
        fmt.Println(err)
    }

    data := parse( string(htmlString[:]))

    i := 1
    isFound := false

    for i = 0; i < len(data); i++ {
        if strings.Contains(data[i], homeTeam) && strings.Contains(data[i], awayTeam) && strings.Contains(data[i], "@graph"){
            fmt.Println("found the link")
            data[i] = strings.ReplaceAll(data[i], " ", "")
            fmt.Println(data[i])
            ioutil.WriteFile("data.txt", []byte(data[i]), 0644)
            isFound = true
            break
        } else {
            continue
        }
    }

    if(!isFound) {
        fmt.Println("RETURNING")
        return
    }

    // OPEN AI API KEY = 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

    fileText, err := os.ReadFile("data.txt")

    if err != nil {
        fmt.Println("err 5")
        fmt.Println(err)
    }

    convertedString := string(fileText)

    allData := TicketGraph{}

    err = json.Unmarshal([]byte(convertedString), &allData)
    if err != nil {
        // panic
        fmt.Println("err 6")
        fmt.Println(err)
    }

    temp := allData.AllTickets

    i = 0

    urlsToVisit := []string{}
    isParking := false

    for i < len(temp) {
        // take date as input

        if isParking {
            if strings.Contains(temp[i].URL, gameDate) && strings.Contains(temp[i].URL, "parking-passes"){
                urlsToVisit = append(urlsToVisit, "https://www.stubhub.ca" + temp[i].URL + "?quantity=" + ticketQuantity)
            }
        } else {
            if strings.Contains(temp[i].URL, gameDate) && !strings.Contains(temp[i].URL, "parking-passes"){
                urlsToVisit = append(urlsToVisit, "https://www.stubhub.ca" + temp[i].URL + "?quantity=" + ticketQuantity)
            }
        }

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

    var wg sync.WaitGroup

    wg.Add(len(urlsToVisit))
    tempItem := StubHubItem{}
    for i = 0; i < len(urlsToVisit); i++ {
         
        go func(i int) {
            defer wg.Done()
            tempItem = visitPage(urlsToVisit[i])
        }(i)
    }
    wg.Wait()

    fmt.Println("Visited all pages")

    fmt.Println(tempItem.Section + " " + tempItem.Row + " " + tempItem.PriceWithFees)
}

func getVividSeatsTickets() {
    url := "https://www.vividseats.com/chicago-bulls-tickets--sports-nba-basketball/performer/161"
    fmt.Printf("HTML code of %s ...\n", url)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("err 1")
        panic(err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36")
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("err 2")
        panic(err)
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("err 3")
        panic(err)
    }
    // fmt.Printf("%s\n", html)
    ioutil.WriteFile("vsTicketText.txt", html, 0644)

    // htmlString, err := ioutil.ReadFile("vsTicketText.txt")

    // if err != nil {
    //     fmt.Println("err 4")
    //     fmt.Println(err)
    // }

    // data := parse( string(htmlString[:]))
}
