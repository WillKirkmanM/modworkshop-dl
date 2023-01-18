package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
var writer = uilive.New()

func main() {

	// modURL := user input
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	visitWebsitesAndDownload(c)

	// WolfHud
	//c.Visit("https://modworkshop.net/mod/15901") // Doesn't Have Valid Download Link (Error)

	// OneShot Mod
	//c.Visit("https://modworkshop.net/mod/40265")
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
	modArray := loadModsFromText()

	for i := 0; i < len(modArray); i++ {
		title, downloadID := getModInformation(c, modArray[i])
		resp := downloadMod(title, downloadID)
		writer.Stop()
		//unzipSource(resp.Filename, destination)
		unzipFile(resp.Filename)
		os.Remove(resp.Filename)
	}
	fmt.Println("Done! The Mods Have Been Downloaded and Installed!")
}

func loadModsFromText() []string {
	file, err := os.Open("modlist.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	modArray := []string{}
	for scanner.Scan() {
		modArray = append(modArray, scanner.Text())
	}

	if len(modArray) == 0 {
		log.Fatal("There are no mods specified in modlist.txt!")
	}
	return modArray
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
		println("Zip FIle")
		err := archiver.DefaultZip.Extract(file, destination, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "tar":
		println("tar ball")
		err := archiver.DefaultTar.Extract(file, destination, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "rar":
		println("rar File")
		err := archiver.DefaultRar.Extract(file, destination, destination)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}

// Thanks: https://gosamples.dev/unzip-file/
/*
func unzipSource(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
			return err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
			return err
	}

	for _, f := range reader.File {
			err := unzipFile(f, destination)
			if err != nil {
					return err
			}
	}

	return nil
}

/*
func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
					return err
			}
			return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
			return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
			return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return err
	}
	return nil
}
*/