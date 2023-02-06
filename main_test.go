package main

import (
	"testing"

	"github.com/gocolly/colly"
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

	_, err := downloadFile("Test File", "14602" /* 4kb | Low Violence Mode */, tempDir)
	if err != nil {
		t.Errorf("TestDownloadFile:  %s", err)
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