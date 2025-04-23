package nipTags

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	// URL of the GitHub repository where the NIPs documentation is stored
	REPO_URL = "https://github.com/nostr-protocol/nips"
	// URL prefix for fetching the raw markdown files of the NIPs
	RAW_FILES_URL = "https://raw.githubusercontent.com/nostr-protocol/nips/refs/heads/master/"
)

// Takes in a file name and checks if it is (probably) a NIP file
func isNip(lookup string) bool {
	switch lookup {
	case
		"BREAKING.md",
		"README.md",
		"CONTRIBUTING.md",
		"LICENSE.md",
		"CHANGELOG.md":
		return false
	}
	return true
}

// TODO : check if no internet / nothing retrieved
// Fetches the list of files in the github repo and extracts the NIPs
func getNips() ([]string, error) {
	nips := []string{}

	collector := colly.NewCollector()
	collector.OnHTML(".react-directory-filename-cell", func(e *colly.HTMLElement) {
		if isNip(e.Text) {
			nip := e.Text[:len(e.Text)-len(".md")]
			nips = append(nips, nip)
		}
	})

	err := collector.Visit(REPO_URL)

	return nips, err
}

// Reads the raw markdown of a NIP file and returns the extracted tags
func getTagsFromText(text string) ([]string, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	tags := []string{}

	// Finds the line with the tags (i.e. line with only tags)
	re := regexp.MustCompile("\n(?:`[a-zA-Z]+`\\s*)+\n")
	line_matches := re.FindAllString(text, -1)

	if len(line_matches) > 1 {
		fmt.Println(line_matches)
		return nil, fmt.Errorf("multiple lines with possible tags found")
	}

	if len (line_matches) == 0 {
		return nil, nil
	}

	// Extracts the tags from the line
	re_words := regexp.MustCompile("[a-zA-Z]+")
	matches := re_words.FindAllString(line_matches[0], -1)

	tags = append(tags, matches...)

	return tags, nil
}

// Extracts the tags from the NIPs in the GitHub repository
func GetNipsTags() (map[string][]string, error) {
	nip_texts := map[string][]string{}
	var tmp_text string
	nips, err := getNips()

	if err != nil {
		return nil, fmt.Errorf("error fetching NIPs: %w", err)
	}

	log.Info("extracting the tags of ", len(nips), " NIPs...")
	for _, nip := range nips {
		nip_texts[nip] = []string{}
		res, _ := http.Get(RAW_FILES_URL + nip + ".md")
		if res.StatusCode != 200 {
			fmt.Println("error fetching file", nip, ":", res.Status)
		}

		content, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println("error reading response body for", nip, ":", err)
			continue
		}
		tmp_text = string(content)
		tags, err := getTagsFromText(tmp_text)
		if err != nil {
			fmt.Println("error extracting tags for", nip, ":", err)
			continue
		}
		nip_texts[nip] = tags
	}

	return nip_texts, nil
}
