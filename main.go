package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/
// Different Games in navbar Home/Games/{Game}
// Detect what game users

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
var help bool
var install string

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
	Success int    `json:"success"`
	Cache   string `json:"cache"`
	Content []Mod  `json:"content"`
	Total   int    `json:"total"`
	PerPage int    `json:"perpage"`
}

var games = map[string]string{
	"Payday 2":                            `C:\Program Files (x86)\Steam\SteamApps\common\PAYDAY 2\mods`,
	"Noita":                               `C:\Program Files (x86)\Steam\SteamApps\common\Noita\mods`,
	"Enter the Gungeon":                   `C:\Program Files (x86)\Steam\steamapps\common\Enter the Gungeon\Mods`,
	"Payday: The Heist":                   `C:\Program Files (x86)\Steam\SteamApps\common\PAYDAY The Heist\mods`,
	"Final Fantasy XV":                    ".",
	"Stolen Realm":                        ".",
	"RAID: World War II":                  ".",
	"Aurora":                              ".",
	"Zuma":                                ".",
	"Luxor":                               ".",
	"VRChat":                              `C:\Program Files (x86)\Steam\SteamApps\common\VRChat\mods`,
	"Left 4 Dead 2":                       `C:\Program Files (x86)\Steam\SteamApps\common\Left 4 Ded 2\left4dead2\addons`,
	"Hitman 3":                            ".",
	"Monster Sanctuary":                   `C:\Program Files (x86)\Steam\SteamApps\common\Monster Sanctuary\BapInEx\plugins`,
	"Fallout 4":                           `C:\Program Files (x86)\Steam\SteamApps\common\Fallout 4\Data`,
	"Teardown":                            `C:\Program Files (x86)\Steam\SteamApps\common\Teardown\data`,
	"Black Mesa":                          ".",
	"Yakuza Kiwami 2":                     ".",
	"Hotline Miami 2: Wrong Number":       `C:\Program Files (x86)\Steam\SteamApps\common\Hotline Miami 2\mods`,
	"Friday Night Funkin'":                ".",
	"Hotdogs, Horseshoes & Hand Grenades": ".",
	"Yakuza Kiwami 1":                     ".",
	"100% Orange Juice":                   `C:\Program Files (x86)\Steam\SteamApps\common\100% Orange Juice\mods`,
	"Hyperdimension Neptunia Re;Birth2":   ".",
	"Non-games / Plugins":                 ".",
	"Yakuza 0":                            ".",
	"One Step From Eden":                  ".",
	"OVERKILL's The Walking Dead":         ".",
	"The Elder Scrolls V: Skyrim - Legendary Edition": `C:\Program Files (x86)\Steam\SteamApps\common\Skyrim\Data`,
	"SCP: Containment Breach":                         ".",
	"Fallout: New Vegas":                              `C:\Program Files (x86)\Steam\SteamApps\common\Fallout New Vegas\Data`,
	"OneShot":                                         ".",
	"SteamVR":                                         `C:\Program Files (x86)\Steam\SteamApps\common\SteamVR\bin\win64`,
	"Criminal Girls: Invite Only":                     ".",
	"Gal*Gun: Double Peace":                           `C:\Program Files (x86)\Steam\SteamApps\common\GalGun Double Peace`,
	"Warhammer: End Times - Vermintide":               `C:\Program Files (x86)\Steam\SteamApps\common\Warhammer End Times Vermintide\binaries\mods`,
	"Tales of Berseria":                               ".",
	"Team Fortress 2":                                 `C:\Program Files (x86)\Steam\SteamApps\common\Team Fortress 2\tf\custom`,
	"Hyperdimension Neptunia Re;Birth3":               ".",
	"Hyperdimension Neptunia Re;Birth1":               ".",
	"Metal Gear Solid V: The Phantom Pain":            ".",
	"Skyrim Special Edition":                          ".",
	"Forspoken":                                       ".",
}

var modResponseObject Response

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	parseCliArgs(c)
}

func getModInformation(c *colly.Collector, modLink string) (title string, downloadID string) {
	c.OnHTML("div.flex-grow-1.p-3.d-flex.flex-column.data", func(r *colly.HTMLElement) {
		title = r.ChildText("[id=title]")
		downloadButtonText := r.ChildAttr("[id=download-button]", "href")
		if !strings.Contains(downloadButtonText, "download") {
			log.Fatal("The Mod You have Requested is Invalid / Does not Have a Download Link. Check the GitHub Page of the Mod!")
		}
		downloadID = strings.Split(downloadButtonText, "/download/")[1]
	})

	err := c.Visit(modLink)
	if err != nil {
		log.Fatalf("There was an error while running the mods you specified. Please Look in the %s file for any formatting errors. %s\n", file, err)
	}
	return title, downloadID
}

func downloadFromFile(c *colly.Collector) error {
	modsArray, assetsArray, err := parseText(file)
	if err != nil {
		log.Fatal(err)
	}

	modsDirectory = games["Payday 2"]
	assetsDirectory = games["Payday 2"] + `\assets`
	fmt.Println(assetsDirectory)

	endStr := "Done! The "

	if len(modsArray) > 0 {
		endStr += "Mods "
		fmt.Println("Downloading Mods!")
		for i := 0; i < len(modsArray); i++ {
			title, downloadID := getModInformation(c, modsArray[i])
			resp, err := downloadFile(title, downloadID, modsDirectory)
			if err != nil {
				return nil
			}
			unzipFile(resp.Filename, modsDirectory)
			os.Remove(resp.Filename)
		}
	}

	if len(modsArray) > 0 && len(assetsArray) > 0 {
		endStr += "and "
	}

	if len(assetsArray) > 0 {
		endStr += "Assets "
		fmt.Println("Downloading Assets!")
		for i := 0; i < len(assetsArray); i++ {
			title, downloadID := getModInformation(c, assetsArray[i])
			resp, err := downloadFile(title, downloadID, assetsDirectory)
			if err != nil {
				return nil
			}
			unzipFile(resp.Filename, assetsDirectory)
			os.Remove(resp.Filename)
		}
	}
	fmt.Println(endStr, "Have Been Downloaded and Installed!")
	return nil
}

func parseText(filePath string) (modsArray []string, assetsArray []string, error error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("There was an error while parsing the %s. Please check for any possible errors or contact the Developer on GitHub! %s", file.Name(), err)
		return nil, nil, errors.New("Error While Parsing File")
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
	return modsArray, assetsArray, nil
}

func downloadFile(title string, downloadID string, destination string) (resp *grab.Response, error error) {

	downloadLink := apiURL + downloadID + "/download?"

	resp, err := grab.Get(destination, downloadLink)
	if err != nil {
		log.Fatalf("There was an error while downloading %s. %s", resp.Filename, err)
		return nil, errors.New("Error While Downloading The File")
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
	return resp, nil
}

func unzipFile(file string, destination string) {
	switch file[len(file)-3:] {
	case "zip":
		err := archiver.DefaultZip.Unarchive(file, destination)
		if err != nil {
			log.Fatalf("There was an error while unzipping the Zip File %s", err)
		}
		break
	case "tar":
		err := archiver.DefaultTar.Unarchive(file, destination)
		if err != nil {
			log.Fatalf("There was an error while unzipping the Tar File %s", err)
		}
		break
	case "rar":
		err := archiver.DefaultRar.Unarchive(file, destination)
		if err != nil {
			log.Fatalf("There was an error while unzipping the Rar File %s", err)
		}
		break
	}
}

func parseCliArgs(c *colly.Collector) {
	flag.StringVar(&file, "file", "", "The text file containing the mods.")
	flag.StringVar(&file, "f", "", "The text file containing the mods.")
	flag.StringVar(&search, "search", "", "The Mod To Search.")
	flag.StringVar(&search, "S", "", "The Mod To Search")
	flag.BoolVar(&help, "h", false, "View all Commands.")
	flag.BoolVar(&help, "help", false, "View all Commands.")
	flag.StringVar(&install, "I", "", "The Mod To Install")
	flag.StringVar(&install, "install", "", "The Mod To Install")
	flag.Parse()

	helpMsg :=
		`
Modworkshop-dl allows for installing mods with ease.

usage: modworkshop-dl [<command>] [<argument>]

The following commands are available:
search, S			The mod to search 				[-S <Name>]
file, f				The text file containing the mods		[-f <File>]
install, I			The Link / ModID To Be Installed		[-I <Link / ModID>]	
		`

	if len(os.Args) > 1 {
		if os.Args[1] == "" {
			fmt.Println("none")
		}
		if os.Args[1] == "-S" {
			search = strings.Join(os.Args[2:], " ")
			searchForMod(search, c)
		}

		if file != "" {
			downloadFromFile(c)
		}

		if install != "" {
			installMod(install, c)
		}

		if help == true {
			fmt.Printf(helpMsg)
		}
	} else {
		fmt.Println(helpMsg)
	}
}

func searchForMod(query string, c *colly.Collector) {
	res, err := http.Get(fmt.Sprintf("https://modworkshop.net/mws/api/modsapi.php?count_total=1&query=%s&func=mods&page=1", query))
	if err != nil {
		log.Fatalf("There was an error while fetching the API! Try Again in 5 Minutes or Contact the Developer on GitHub! %s", err)
	}

	responseData, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(responseData, &modResponseObject)

	// Print Stuff to Out

	if len(modResponseObject.Content) > 0 {
		for i := 0; i < len(modResponseObject.Content); i++ {
			if i == 10 {
				break
			}
			fmt.Printf(
				`
		[%v] %s 
			🧍 %s
			📍 %s / %s
			❤ %v 😋 %v 👁️ %v		🕑 %s
		
		`, i+1,
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
		fmt.Println("\n\nWhich mod would you like to install?")
		var modChoice int
		fmt.Scanln(&modChoice)

		downloadModFromIndex(modChoice, c)

	} else {
		fmt.Println("No mods found")
	}
}

func downloadModFromID(ID int, c *colly.Collector, destination string)  error {
	iID := strconv.Itoa(ID)

	link := "https://modworkshop.net/mod/" + iID
	err := downloadModFromLink(link, c, destination)
	if err != nil {
		return err
	}
	return nil
}

func downloadModFromLink(link string, c *colly.Collector, destination string) error {
	title, downloadID := getModInformation(c, link)

	resp, err := downloadFile(title, downloadID, destination)
	if err != nil {
		return err
	}
	unzipFile(resp.Filename, destination)
	os.Remove(resp.Filename)
	return nil
}

func downloadModFromIndex(index int, c *colly.Collector) {
	index = index - 1

	if index > modResponseObject.Total {
		index = index + 1
		log.Fatalf("The index %d is out the bounds! Select a Mod Provided", index)
	}

	title := modResponseObject.Content[index].Name
	downloadLink := "https://modworkshop.net/mod/" + strconv.Itoa(modResponseObject.Content[index].Did)

	title, downloadID := getModInformation(c, downloadLink)

	gameName := modResponseObject.Content[index].Game
	gameDir := games[gameName]

	ensureDir(gameDir)

	resp, err := downloadFile(title, downloadID, gameDir)
	if err != nil {
		return
	}
	unzipFile(resp.Filename, gameDir)
	os.Remove(resp.Filename)
}

func installMod(mod string, c *colly.Collector) {
	if strings.Contains(mod, "http") {
		downloadModFromLink(mod, c, modsDirectory)
		return
	}

	if len(mod) > 3 && len(mod) < 8 {
		iMod, err := strconv.Atoi(mod)
		if err != nil {
			log.Fatalf("There was a conversion error! Please Contact the Developer on GitHub %s", err)
		}
		downloadModFromID(iMod, c, modsDirectory)
		return
	}
}

func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModeSticky|os.ModePerm)
	} else {
		if err != nil {
			log.Fatalf("There was an error ensuring that %s exists! %s\n", dirPath, err)
			return errors.New("The Directory Could Not Be Created!")
		}
		return nil
	}

	if strings.Contains(dirPath, "PAYDAY 2") {
		ensureDir(`C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods`)
		ensureDir(`C:\Program Files (x86)\Steam\steamapps\common\PAYDAY 2\mods\assets`)
	}
	ensureDir(dirPath)
	return nil
}
