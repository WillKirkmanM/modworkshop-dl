package main

/*
	Goal of Project:
	- Be able to download files from modworkshop (https://modworkshop.net/)
*/

import (
	"fmt"
	"os"
	"strings"
	"io"
	"archive/zip"
	"path/filepath"
	"net/http"
	"net/url"

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
		if !strings.Contains(downloadButtonText, "download") {
			panic("The Mod You have Requested is Invalid / Does not Have a Download Link. Check the GitHub Page of the Mod!")
		}

		downloadID := strings.Split(downloadButtonText, "/download/")[1]

		fmt.Printf("%sFound Mod%s: %s%s%s\n", Green, White, Yellow, title, White)
		fmt.Printf("Size N/A\n")
		// TODO: Get Size of Mod

		// TODO: Do this if no text file parameter was set (in order to prevent annoyance)
		fmt.Println("Would you like to download this mod? (Y/n)")

		var answer string
		fmt.Scan(&answer)

		if answer == "Y" || answer == "y" {
			downloadMod(downloadID)
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

// Thanks: https://www.sohamkamani.com/golang/exec-shell-command/
// Dont want to use Dependency but Thank: https://golangdocs.com/golang-download-files
func downloadMod(downloadID string) {
	downloadLink := apiURL + downloadID + "/download?"
	fmt.Printf("downloadLink: %v\n", downloadLink)

	//unzipSource(resp.Filename, ".")
	//os.Remove(resp.Filename)
}

func download(link string, destination string) {
	fileURL, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	client := http.Client {
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	io.Copy(file, resp.Body)

 	defer file.Close()
}

// Thanks: https://gosamples.dev/unzip-file/
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
	os.Remove(source)
	return nil
}

func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination) + string(os.PathSeparator)) {
			return fmt.Errorf("Invalid file path: %s", filePath)
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

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, f.Mode())
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
