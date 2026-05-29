package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/net/html"
)

func makeGetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Since we can't do request if we don't have a user agent, we have to set one
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("bad status: %d (%s)", response.StatusCode, response.Status)
	}

	return response, nil
}

// Returns a []string of the link of the images found in a given Node of html
func ExtractLinks(doc *html.Node, depth int, recursive bool, current_depth int) []string {
	var images []string

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
					if recursive && current_depth < depth {
						if isHTMLPage(attr.Val) {
							doc, _ := GetHtmlFromUrl(attr.Val)
							if doc != nil {
								new_images := ExtractLinks(doc, depth, recursive, current_depth+1)
								images = append(images, new_images...)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	images = clearDoubles(images)
	return images
}

func clearDoubles(images []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range images {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func GetHtmlFromUrl(url string) (*html.Node, error) {
	response, err := makeGetRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, nil
}

func DownloadImageFromUrl(url string, filename string) {
	response, err := makeGetRequest(url)
	if err != nil {
		return
		log.Fatal(err)
	}
	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		return
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file.
	// This supports huge files
	_, e := io.Copy(file, response.Body)
	if e != nil {
		return
		log.Fatal(err)
	}
}

func GetFilenameFromUrl(image_url string) (string, error) {
	parsed_url, err := url.Parse(image_url)
	if err != nil {
		return "", err
	}
	filename := path.Base(parsed_url.Path)
	if filename == "." || filename == "/" || filename == "" {
		return "", err
	}
	return filename, nil
}
