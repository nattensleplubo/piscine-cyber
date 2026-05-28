package main

import (
	"flag"
	"fmt"
	"os"
	"spider/downloader"
)

func printOptions(url string, recursive bool, depth int, folder_path string) {
	fmt.Printf("URL: %s | recursive: %v | depth: %d | path: %s\n",
		url, recursive, depth, folder_path)
}

func main() {
	recursive := flag.Bool("r", false, "recursively download images")
	depth := flag.Int("l", 5, "max depth for recursive download")
	path := flag.String("p", "./data/", "path to save downloaded files")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spider [-r] [-l N] [-p PATH] URL")
		os.Exit(1)
	}

	url := flag.Arg(0)

	printOptions(url, *recursive, *depth, *path)
	// downloader.DownloadImageFromUrl(url, "testing.jpg")
	doc := downloader.GetHtmlFromUrl("https://www.42.fr")
	downloader.ExtractLinks(doc)
}
