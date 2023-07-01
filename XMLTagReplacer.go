/*
XMLTagReplacer.go
Author: Sergio Lima
Creation date: 24 June 2023

XMLTagReplacer is a command-line utility written in Go.
It replaces the contents of a specific XML tag in all .xml files
in the current directory and its subdirectories.

The command-line arguments to XMLTagReplacer are the name of the
tag to replace and the new value for that tag.

For example:
`XMLTagReplacer mytag false` replaces the contents of the `mytag` XML tags with `false`.

Use `XMLTagReplacer --help` or `XMLTagReplacer -h` to display help information.

Use `XMLTagReplacer --version` or `XMLTagReplacer -v` to display the version number of XMLTagReplacer.

To compile this program for Linux, use the following command:
    go build -o XMLTagReplacer

To compile this program for Windows, use the following command:
    GOOS=windows GOARCH=amd64 go build -o XMLTagReplacer.exe

Please note that this program does not currently support XML namespaces.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const version = "1.1"

const helpDescription = "XMLTagReplacer is a utility for replacing the contents of " +
	"a specific XML tag in all .xml files in the current directory and its subdirectories.\n\n" +
	"Usage:\n" +
	"XMLTagReplacer [tag] [newvalue]: Replace the contents of [tag] with [newvalue] " +
	"in all .xml files in the current directory and its subdirectories.\n" +
	"XMLTagReplacer --version or -v: Print the version number of XMLTagReplacer.\n" +
	"XMLTagReplacer --help or -h: Print this help message.\n\n" +
	"Example:\n" +
	"XMLTagReplacer mytag false: Replace the contents of 'mytag' with 'false' in all .xml files in the current directory and its subdirectories.\n"

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func processFile(path string, tag string, newValue string) bool {
	data, err := ioutil.ReadFile(path)
	check(err)

	content := string(data)

	re := regexp.MustCompile(fmt.Sprintf("<%s>.*</%s>", tag, tag))
	newTag := fmt.Sprintf("<%s>%s</%s>", tag, newValue, tag)

	modifiedContent := re.ReplaceAllString(content, newTag)

	if content != modifiedContent {
		err = ioutil.WriteFile(path, []byte(modifiedContent), 0)
		check(err)
		return true
	}

	return false
}

func printHelp() string {
	return helpDescription
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: XMLTagReplacer [tag] [newvalue] or XMLTagReplacer --version or XMLTagReplacer -v or XMLTagReplacer --help or XMLTagReplacer -h")
		os.Exit(1)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("XMLTagReplacer version %s\n", version)
		os.Exit(0)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println(printHelp())
		os.Exit(0)
	}

	tag := os.Args[1]
	newValue := os.Args[2]

	totalCounter := 0
	modifiedCounter := 0

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		check(err)

		if filepath.Ext(path) == ".xml" {
			totalCounter++
			if processFile(path, tag, newValue) {
				modifiedCounter++
			}
		}

		return nil
	})

	check(err)

	fmt.Printf("Number of XML files read: %d\n", totalCounter)
	fmt.Printf("Number of XML files where '<%s>' tag content was replaced: %d\n", tag, modifiedCounter)
}
