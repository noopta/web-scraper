package main

import (
	// "fmt"
	// "time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// func main() {
//     // Run Chrome browser
//     service, err := selenium.NewChromeDriverService("/../usr/bin/chromedriver", 4444)
//     if err != nil {
//         panic(err)
//     }
//     defer service.Stop()

//     caps := selenium.Capabilities{}
//     caps.AddChrome(chrome.Capabilities{Args: []string{
//         "--no-sandbox",
//         "window-size=1920x1080",
//         "--disable-dev-shm-usage",
//         "disable-gpu",
//         "--user-data-dir=~/.config/google-chrome",
//         // "--headless",  // comment out this line to see the browser
//     }})

//     driver, err := selenium.NewRemote(caps, "")
//     if err != nil {
//         panic(err)
//     }

// 	// Navigate to the website
// 	err = driver.Get("https://boards.greenhouse.io/cloudflare/jobs/1443863?gh_jid=1443863&utm_source=Simplify#appm")
// 	if err != nil {
// 		fmt.Println("Failed to navigate to website:", err)
// 		return
// 	}

// 	// Find the text field and enter text
// 	elem, err := driver.FindElement(selenium.ByCSSSelector, "input[type='text']")
// 	if err != nil {
// 		fmt.Println("Failed to find text field:", err)
// 		return
// 	}
// 	err = elem.SendKeys("Hello, world!")
// 	if err != nil {
// 		fmt.Println("Failed to enter text:", err)
// 		return
// 	}

// 	// Wait for a few seconds to see the results
// 	time.Sleep(5 * time.Second)
// }
func main() {
    // Run Chrome browser
    service, err := selenium.NewChromeDriverService("/../usr/bin/chromedriver", 4444)
    if err != nil {
        panic(err)
    }
    defer service.Stop()

    caps := selenium.Capabilities{}
    caps.AddChrome(chrome.Capabilities{Args: []string{
        "window-size=1920x1080",
        "--user-data-dir=/home/ubuntu/.config/google-chrome/",
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "disable-gpu",
        // "--headless",  // comment out this line to see the browser
    }})

    driver, err := selenium.NewRemote(caps, "")
    if err != nil {
        panic(err)
    }

    driver.Get("https://www.google.com")
}
