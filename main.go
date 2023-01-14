package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/

import (
	"fmt"
	"os"

	"strings"
	//"time"

	"github.com/gocolly/colly"
)

// Thanks (https://twin.sh/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go)
var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"

var baseURL string = "modworkshop.net"
var apiURL string = "https://modworkshop.net/api/files/"

func main() {
	// Start of File! (If you hadn't noticed.)
	// modURL := user input
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	c.OnHTML("div.flex-grow-1.p-3.d-flex.flex-column.data", func(r *colly.HTMLElement) {
		title := r.ChildText("[id=title]")
		downloadButtonText := r.ChildAttr("[id=download-button]", "href")

		downloadID := strings.Split(downloadButtonText, "/download/")[1]
		fmt.Println(downloadID)

		fmt.Printf("%sFound Mod: %s\n", Green, title)
		fmt.Printf("Size N/A\n")
		// TODO: Get Size of Mod

		// TODO: Do this if no text file parameter was set (in order to prevent annoyance)
		fmt.Println("Would you like to download this mod? (Y/n)")

		var answer string
		fmt.Scan(&answer)

		if answer == "Y" || answer == "y" {
			fmt.Println(White)
			os.Exit(0)
		} else {
			fmt.Println(White)
			fmt.Println("Alright, Exiting Process...")
			os.Exit(0)
		}
	})
	// WolfHud
	//c.Visit("https://modworkshop.net/mod/15901")

	// OneShot Mod
	c.Visit("https://modworkshop.net/mod/40265")
}
