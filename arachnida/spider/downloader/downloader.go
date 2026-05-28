package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func makeGetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// Since we can't do request if we don't have a user agent, we have to set one
	req.Header.Set("User-Agent", "Mozilla/5.0")

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", response.Status)
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Bad status: %s", response.Status)
	}

	return response, nil
}

func ExtractLinks(doc *html.Node) ([]string, []string) {
	var images []string
	var sublinks []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			// range creates an iterator named `attr`
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					images = append(images, attr.Val)
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					sublinks = append(sublinks, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	fmt.Println("images", images)
	fmt.Println("sublinks", sublinks)
	return images, sublinks
}

func GetHtmlFromUrl(url string) *html.Node {
	response, _ := makeGetRequest(url)
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func DownloadImageFromUrl(url string, filename string) {
	response, _ := makeGetRequest(url)
	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file.
	// This supports huge files
	written, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success with file size of ", written)
}
