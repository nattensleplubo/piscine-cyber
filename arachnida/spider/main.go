package main

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	// "path/filepath"
	"spider/downloader"
	"time"
)

func printOptions(url string, recursive bool, depth int, folder_path string) {
	fmt.Printf("URL: %s | recursive: %v | depth: %d | path: %s\n",
		url, recursive, depth, folder_path)
}

func main() {
	recursive := flag.Bool("r", false, "recursively download images")
	depth := flag.Int("l", 5, "max depth for recursive download")
	path := flag.String("p", "./data", "path to save downloaded files")
	current_depth := 0
	spinr := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spider [-r] [-l N] [-p PATH] URL")
		os.Exit(1)
	}

	err := os.MkdirAll(*path, os.ModePerm)
	if err != nil {
		os.Exit(1)
	}

	*path += "/"

	url := flag.Arg(0)

	printOptions(url, *recursive, *depth, *path)
	doc, _ := downloader.GetHtmlFromUrl(url)
	fmt.Println("")
	fmt.Println("                    █▄            ")
	fmt.Println("              ▀▀    ██       ▄    ")
	fmt.Println("  ▄██▀█ ████▄ ██ ▄████ ▄█▀█▄ ████▄")
	fmt.Println("  ▀███▄ ██ ██ ██ ██ ██ ██▄█▀ ██   ")
	fmt.Println(" █▄▄██▀▄████▀▄██▄█▀███▄▀█▄▄▄▄█▀   ")
	fmt.Println("        ██                        ")
	spinr.Prefix = "        ▀   "
	spinr.Suffix = " crawling\n"
	spinr.Color("red", "bold")
	spinr.Start()
	image_links := downloader.ExtractLinks(doc, *depth, *recursive, current_depth)
	spinr.Stop()
	fmt.Printf("Done\n")
	for _, img := range image_links {
		filename, err := downloader.GetFilenameFromUrl(img)
		fmt.Println("filename: ", filename)
		if err == nil && downloader.CheckFileExtension(filename) != false {
			namepath := *path + filename
			downloader.DownloadImageFromUrl(img, namepath)
		}
	}
}
