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
	recursive := flag.Bool("r", true, "recursively download images")
	depth := flag.Int("l", 2, "max depth for recursive download")
	path := flag.String("p", "./data/", "path to save downloaded files")
	current_depth := 0

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spider [-r] [-l N] [-p PATH] URL")
		os.Exit(1)
	}

	url := flag.Arg(0)

	printOptions(url, *recursive, *depth, *path)
	// downloader.DownloadImageFromUrl(url, "testing.jpg")
	doc, _ := downloader.GetHtmlFromUrl("https://www.42.fr")
	image_links := downloader.ExtractLinks(doc, *depth, *recursive, current_depth)
	fmt.Println("\n\n[ALL LINKS] : \n", image_links)
}

//go run main.go https://upload.wikimedia.org/wikipedia/commons/thumb/d/d0/Assemblage_chien.jpg/330px-Assemblage_chien.jpg
