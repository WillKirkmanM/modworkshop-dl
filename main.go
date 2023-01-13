package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/

import (
	"fmt"
	"os"
	"time"

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
		modSize := r.ChildAttr("[id=download-button]", "style.font-size=13px")

		if downloadLink == "https://modworkshop.net" {
			downloadLink = "No Download Link Found, Check the GitHub Page of that mod."
			os.Exit(0)
		}

		fmt.Println("Title:", title)
		fmt.Println("Download Link:", downloadLink)
		fmt.Println("ModSize:", modSize)

		fmt.Println("Would you like to download this mod? (Y/n)")

		var answer string
		fmt.Scan(&answer)

		if answer == "Y" || answer == "y" {
			downloadFile("./temp", downloadLink)
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

func downloadFile(filePath string, url string) (err error) {
	client := grab.NewClient()
	req, _ := grab.NewRequest(filePath, url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			if err := resp.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
				os.Exit(1)
			}
			break Loop
		}
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	return nil
}