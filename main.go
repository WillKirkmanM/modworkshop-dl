package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/
// Different Games in navbar Home/Games/{Game}
// Detect what game users

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
var modsDirectory string = "."
var assetsDirectory string = "."
var writer = uilive.New()

// CLI Argument Variables
var file string
var search string
var install string
var help bool 

type Mod struct {
	Did                  int         `json:"did"`
	Name                 string      `json:"name"`
	SuspendedStatus      int         `json:"suspended_status"`
	FileStatus           int         `json:"file_status"`
	Hidden               int         `json:"hidden"`
	Thumbnail            string      `json:"thumbnail"`
	SubmitterUID         int         `json:"submitter_uid"`
	Submitter            string      `json:"submitter"`
	CollaboratorsNr      int         `json:"collaborators_nr"`
	ShortDescription     string      `json:"short_description"`
	Game                 string      `json:"game"`
	GameShort            string      `json:"game_short"`
	Gid                  int         `json:"gid"`
	Category             string      `json:"category"`
	Cid                  int         `json:"cid"`
	Views                interface{} `json:"views"`
	Downloads            interface{} `json:"downloads"`
	Likes                int         `json:"likes"`
	Date                 string      `json:"date"`
	PubDate              string      `json:"pub_date"`
	Timeago              string      `json:"timeago"`
	TimeagoPub           string      `json:"timeago_pub"`
	DateTimestamp        int         `json:"date_timestamp"`
	PubDateTimestamp     int         `json:"pub_date_timestamp"`
	IsNsfw               int         `json:"is_nsfw"`
	UsesDefaultThumbnail interface{} `json:"uses_default_thumbnail"`
} 

type Response struct {
	Success       int       `json:"success"`
	Cache         string    `json:"cache"`
	Content       []Mod     `json:"content"`
	Total         int       `json:"total"`
	PerPage       int       `json:"perpage"`
}

var modResponseObject Response

func main() {
	beforeChecks()

	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	parseCliArgs(c)
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
			log.Fatalf("There was an error while running the mods you specified. Please Look in the %s file for any formatting errors. %s\n", file,  err)
		}		
	return title, downloadID
}

func downloadFromFile(c *colly.Collector) {

	modsArray, assetsArray := parseText(file)

	if len(modsArray) > 0 {
		fmt.Println("Downloading Mods!")
		for i := 0; i < len(modsArray); i++ {
			title, downloadID := getModInformation(c, modsArray[i])
			resp := downloadFile(title, downloadID, modsDirectory)
			unzipFile(resp.Filename, modsDirectory)
			os.Remove(resp.Filename)
		}
	}
	
	if len(assetsArray) > 0 {
		fmt.Println("Downloading Assets!")
		for i := 0; i < len(assetsArray); i++ {
			title, downloadID := getModInformation(c, assetsArray[i])
			resp := downloadFile(title, downloadID, assetsDirectory)
			unzipFile(resp.Filename, assetsDirectory)
			os.Remove(resp.Filename)
		}
	}
	fmt.Println("Done! The Mods and/or Assets Have Been Downloaded and Installed!")
	}

func parseText(filePath string) (modsArray []string, assetsArray []string) {
	file, err := os.Open(filePath)
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

		if text == "Assets" {
			matchedMods = false
			matchedAssets = true
			continue
		}

		if matchedMods {
			modsArray = append(modsArray, text)				
		}

		if matchedAssets {
			assetsArray = append(assetsArray, text)
		}
	}
	return modsArray, assetsArray
}

func downloadFile(title string, downloadID string, destination string) (resp *grab.Response) {
	downloadLink := apiURL + downloadID + "/download?"

	resp, err := grab.Get(destination, downloadLink)
	if err != nil {
		log.Fatal(err)
	}

	progress := writer.Newline()
	writer.Start()
	defer writer.Stop()

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	fmt.Fprintf(writer, "Downloading: %s\n", title)
	Downloading:
		for {
			select {
			case <-t.C:
				fmt.Println("Downloading")
				fmt.Fprintf(progress, "Downloaded %v / %v (%.2f%%)\n", resp.BytesComplete(), resp.Size(), 100*resp.Progress())
			case <-resp.Done:
				fmt.Fprintf(progress, "The Download has Complete! Took %v\n", resp.Duration())
				break Downloading
			}
		}
	return resp
}

func unzipFile(file string, destination string) {
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

func beforeChecks() {
	switch runtime.GOOS {
	case "windows":
		modsDirectory = `C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods`
		assetsDirectory = `C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods\assets`
	}

	if _, err := os.Stat(modsDirectory); os.IsNotExist(err) {
		os.Mkdir(modsDirectory, os.ModeDir)
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}

	if _, err := os.Stat(assetsDirectory); os.IsNotExist(err) {
		os.Mkdir(assetsDirectory, os.ModeDir)
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}
}

func parseCliArgs(c *colly.Collector) {
	flag.StringVar(&file, "file", "", "The text file containing the mods.")
	flag.StringVar(&file, "f", "", "The text file containing the mods.")
	flag.StringVar(&search, "search", "", "The Mod To Search.")
	flag.StringVar(&search, "S", "", "The Mod To Search")
	flag.StringVar(&install, "I", "", "The Mod To Install.")
	flag.BoolVar(&help, "h", false, "View all Commands.")
	flag.BoolVar(&help, "help", false, "View all Commands.")
	flag.Parse()

	if file != "" {
		downloadFromFile(c)
	}

	if search != "" {
		searchForMod(search)
	} 

	if help == true {
		fmt.Printf(
		`
Modworkshop-dl allows for installing mods with ease.

usage: modworkshop-dl [<command>] [<argument>]

The following commands are available:
search, S			The mod to search 				[-S WolfHud]
file, f				The text file containing the mods		[-f modlist.txt]
install, I			The mod to install				[-I  1]
		`)
	} 
}

func searchForMod(query string) {
	res, err := http.Get(fmt.Sprintf("https://modworkshop.net/mws/api/modsapi.php?count_total=1&query=%s&func=mods&page=1", query))
	if err != nil {
			log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(responseData, &modResponseObject)

	// Print Stuff to Out

	for i := 0; i < 10; i++ {
			fmt.Printf(
		`
		[%v] %s 
			🧍 %s
			📍 %s / %s
			❤ %v 😋 %v 👁️  %v		🕑 %s
		
		`, i + 1, 
		modResponseObject.Content[i].Name, 
		modResponseObject.Content[i].Submitter,
	
		modResponseObject.Content[i].Game,
		modResponseObject.Content[i].Category,
	
		modResponseObject.Content[i].Likes,
		modResponseObject.Content[i].Downloads,
		modResponseObject.Content[i].Views,
	
		modResponseObject.Content[i].Timeago,
	)
	}
}

func choose(index int) {
	// I'm thinking take in array then choose from that!
}