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
		downloadLink := "https://" + baseURL + downloadButtonText

		downloadID := strings.Split(downloadButtonText, "/download/")[1]
		fmt.Println(downloadID)

		fmt.Println(downloadButtonText)

		fmt.Println("Title:", title)
		fmt.Println("Download Link:", downloadLink)

		fmt.Println("Would you like to download this mod? (Y/n)")

		var answer string
		fmt.Scan(&answer)

		if answer == "Y" || answer == "y" {
			os.Exit(0)
		} else {
			fmt.Println("Alright, Exiting Process...")
			os.Exit(0)
		}
	})
	// WolfHud
	//c.Visit("https://modworkshop.net/mod/15901")

	// OneShot Mod
	c.Visit("https://modworkshop.net/mod/40265")
}
