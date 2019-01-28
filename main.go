package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var questions_regex = regexp.MustCompile(`questions/\d+/`)

func isQuestion(link string) bool {
	return questions_regex.Find([]byte(link)) != nil
}

func parseEmbededReferenceLinks(doc *goquery.Document) {
	doc.Find(".answercell .post-text").First().Find("a").Each(func(i int, s *goquery.Selection) {
		href, success := s.Attr("href")
		if success == true {
			href = " (" + href + ")"
			s.AppendHtml(href)
		}
	})
}

func performRequest(url string) string {
	doc, err := goquery.NewDocument(url + "?answertab=votes")
	if err != nil {
		log.Fatal(err)
	}
	parseEmbededReferenceLinks(doc)
	return doc.Find(".answercell .post-text").First().Text()

}

func performSearch(query string) ([]string, error) {
	searchURL := "http://www.google.com/search?q=site:stackoverflow.com/questions%20" + url.QueryEscape(query)
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
	if err != nil {
		log.Fatal("error reading document", err)
	}
	var links []string
	doc.Find("h3.r a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists {
			u, err := url.Parse(str)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := url.ParseQuery(u.RawQuery)
			link := m["q"][0]
			if isQuestion(link) {
				links = append(links, link)
			}

		}
	})

	if len(links) == 0 {
		return nil, errors.New("search failed")
	}

	return links, nil
}

func printAnswer(answer, url string) {
	fmt.Println(answer)
	fmt.Println("Url:", url)

}

func main() {
	query := strings.Join(os.Args[1:], " ")
	questions, err := performSearch(query)
	if err != nil {
		log.Fatal("sorry, i couldn't find what you're looking for :(")
	}
	fmt.Println("Hello World")
	answer := performRequest(questions[0])
	printAnswer(answer, questions[0])
}
