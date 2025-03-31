package scripts

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gocolly/colly"
)

const (
	REPO_URL = "https://github.com/nostr-protocol/nips"
	RAW_FILES_URL = "https://raw.githubusercontent.com/nostr-protocol/nips/refs/heads/master/"
)

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

func getNips() []string {
	collector := colly.NewCollector()
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error:", e)
	})

	nips := []string{}

	collector.OnHTML(".react-directory-filename-cell", func(e *colly.HTMLElement) {
		if isNip(e.Text) {
			nips = append(nips, e.Text)
		}
	})

	collector.Visit(REPO_URL)

	return nips
}

func getTagsFromText(text string) []string {
	tags := []string{}

	if idx := regexp.MustCompile(`(?m)^##`).FindStringIndex(text); idx != nil {
		text = text[:idx[0]]
	}

	re := regexp.MustCompile("(?:`\\w+`\\s*)+\n")
	line_match := re.FindString(text)
	re_words := regexp.MustCompile("`\\w+`")
	matches := re_words.FindAllString(line_match, -1)

	tags = append(tags, matches...)

	return tags
}

func GetCurrentNipsTags() map[string][]string {
	nip_texts := map[string][]string{}
	var tmp_text string

	for _, nip := range getNips() {
		res, err := http.Get(RAW_FILES_URL + nip)
		if err != nil {
			fmt.Println("Error fetching file:", err)
			continue
		}
		content, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}
		tmp_text = string(content)
		nip_texts[nip] = getTagsFromText(tmp_text)
	}

	return nip_texts
}
