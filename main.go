package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/
// Different Games in navbar Home/Games/{Game}
// Detect what game users 

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gocolly/colly"
	"github.com/gosuri/uilive"
	"github.com/mholt/archiver/v3"
)

var baseURL string = "modworkshop.net"
var apiURL string = "https://modworkshop.net/api/files/"
var destination string = "."
var assetsDirectory string = "."
var writer = uilive.New()

func main() {
	getModDirectory()
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)
	visitWebsitesAndDownload(c)
}

func getModDirectory() {
	switch runtime.GOOS {
	case "windows":
		destination = `C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods`
		assetsDirectory =`C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods\Assets`
	}
}

func getModInformation(c *colly.Collector, mod string) (title string, downloadID string) {
		c.OnHTML("div.flex-grow-1.p-3.d-flex.flex-column.data", func(r *colly.HTMLElement) {
			title = r.ChildText("[id=title]")
			downloadButtonText := r.ChildAttr("[id=download-button]", "href")
			if !strings.Contains(downloadButtonText, "download") {
				log.Fatal("The Mod You have Requested is Invalid / Does not Have a Download Link. Check the GitHub Page of the Mod!")
			}
			downloadID = strings.Split(downloadButtonText, "/download/")[1]
	})
		err := c.Visit(mod)
		if err != nil {
			log.Fatal("There was an error while running the mods you specified. Please Look in the modlist.txt file for any errors.\n", err)
		}		
	return title, downloadID
}

func visitWebsitesAndDownload(c *colly.Collector) {

	modsArray, assetsArray := parseText()

	for i := 0; i < len(modsArray); i++ {
		title, downloadID := getModInformation(c, modsArray[i])
		resp := downloadMod(title, downloadID)
		writer.Stop()
		//unzipSource(resp.Filename, destination)
		unzipFile(resp.Filename)
		os.Remove(resp.Filename)
	}
	fmt.Println("Done! The Mods Have Been Downloaded and Installed!")
}

func parseText(filePath string) (modsArray []string, assetsArray []string) {
	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	matchedMods := false
	matchedAssets := false

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		if text == "Mods" {
			matchedMods = true
			continue
		}

		if matchedMods {
			modsArray = append(modsArray, text)				
		}

		if text == "Assets" {
			matchedMods = false
			matchedAssets = true
			continue
		}

		if matchedAssets {
			assetsArray = append(assetsArray, text)
		}
	}

	if len(modsArray) == 0 {
		log.Fatal("There are no mods specified in modlist.txt!")
	}
	return modsArray, assetsArray
}

func downloadMod(title string, downloadID string) (resp *grab.Response) {
	downloadLink := apiURL + downloadID + "/download?"

	resp, err := grab.Get(destination, downloadLink)
	if err != nil {
		log.Fatal(err)
	}

	writer.Start()

	t := time.NewTicker(5 * time.Millisecond)
	defer t.Stop()

	fmt.Fprintf(writer, "Downloading: %s\n", title)
	Downloading:
		for {
			select {
			case <-t.C:
				fmt.Fprintf(writer, "Downloaded %v / %v (%.2f%%)\n", resp.BytesComplete(), resp.Size(), 100*resp.Progress())
			case <-resp.Done:
				fmt.Fprintf(writer, "The Download has Complete! Took %v\n", resp.Duration())
				break Downloading
			}
		}
	return resp
}

func unzipFile(file string) {
	switch file[len(file)-3:] {
	case "zip":
		err := archiver.DefaultZip.Unarchive(file, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "tar":
		err := archiver.DefaultTar.Unarchive(file, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "rar":
		err := archiver.DefaultRar.Unarchive(file, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}
