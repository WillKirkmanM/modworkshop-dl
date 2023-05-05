package main

import (
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"path/filepath"
	"testing"
	"os"
	"fmt"
	"bytes"
	"log"

	"github.com/gocolly/colly"
	"github.com/mholt/archiver/v3"
)

func TestHelpArguments(t *testing.T) {
}

// Testing Utility Functions
func TestEnsureDir(t *testing.T) {
	tempDir := t.TempDir()
	err := ensureDir(tempDir + "./test")
	if err != nil {
		t.Errorf("TestEnsureDir:  %s", err)
	}
}

// Testing Download Functions

func TestDownloadFile(t *testing.T) {
	tempDir := t.TempDir()

	resp, err := downloadFile("Test File", "14602" /* 4kb | Low Violence Mode */, tempDir)
	if err != nil {
		t.Fatalf("downloadFile failed: %s", err)
	}

	if resp == nil {
		t.Fatal("downloadFile failed: expected a valid response")
	}

	if !doesExist(filepath.Join(tempDir, resp.Filename)) {
		t.Errorf("downloadFile failed: expected file '%s' to exist in temp directory, but it doesn't", resp.Filename)
	}
}


func TestDownloadFromFile(t *testing.T) {
}

func TestDownloadModFromID(t *testing.T) {
	tempDir := t.TempDir()
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	err := downloadModFromID(14602 /* 4kb | Low Violence Mode */ , c, tempDir)
	if err != nil {
		t.Errorf("TestDownloadModFromID:  %s", err)
	}
}

func TestDownloadModFromLink(t *testing.T) {
	tempDir := t.TempDir()
	c := colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	err := downloadModFromLink("https://modworkshop.net/mod/14602", c, tempDir)
	if err != nil {
		t.Errorf("TestDownloadModFromLink:  %s", err)
	}
}

func TestDownloadModFromIndex(t *testing.T) {
}

func TestGetModInformation(t *testing.T) {
	// Initialize a new collector
	c := colly.NewCollector()

	// Set up a mock server to serve test HTML content
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Return a mock HTML response with a valid mod title and download link
			fmt.Fprint(w, `<div class="flex-grow-1 p-3 d-flex flex-column data">
					<h1 id="title">Test Mod</h1>
					<a id="download-button" href="/download/test-mod">Download</a>
			</div>`)
	}))
	defer mockServer.Close()

	// Call the function and check the results
	expectedTitle := "Test Mod"
	expectedDownloadID := "test-mod"
	actualTitle, actualDownloadID := getModInformation(c, mockServer.URL)
	if actualTitle != expectedTitle {
			t.Errorf("Expected title %q, but got %q", expectedTitle, actualTitle)
	}
	if actualDownloadID != expectedDownloadID {
			t.Errorf("Expected download ID %q, but got %q", expectedDownloadID, actualDownloadID)
	}
}
func TestDoesExist(t *testing.T) {
	// Create a temporary file
	file, err := ioutil.TempFile("", "test")
	if err != nil {
	t.Fatal(err)
	}
	// Close the file to release the handle
	defer os.Remove(file.Name())
	// Test for an existing file

	exists := doesExist(file.Name())
	if !exists {
		t.Errorf("doesExist(%s) returned false, expected true", file.Name())
	}

	// Test for a non-existing file
	exists = doesExist("non-existing-file")
	if exists {
		t.Errorf("doesExist(non-existing-file) returned true, expected false")
	}
}



