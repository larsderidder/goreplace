package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const debug = false

func main() {
	args := os.Args[1:]
	if debug {
		fmt.Println(os.Args)
	}
	if len(args) < 3 {
		fmt.Println("Usage: goreplace TOFIND TOREPLACE TARGETS...")
		os.Exit(0)
	}
	toFind := args[0]
	toReplace := args[1]
	targets := args[2:]
	if debug {
		fmt.Println(toFind)
		fmt.Println(toReplace)
		fmt.Println(targets)
	}
	changed := goReplace(toFind, toReplace, targets)

	fmt.Printf("Done! Changed %d file(s).\n", changed)
}

func goReplace(toFind string, toReplace string, targets []string) int {
	replaced := make(chan bool)

	// Find all files with exact matches and goreplace
	routines, unmatched := replaceByExactMatch(toFind, toReplace, targets, replaced)

	// Find all files that match leftover target and goreplace
	routines += replaceByPatternMatch(toFind, toReplace, unmatched, replaced)

	// Wait for all goroutines to report back.
	changed := 0
	for i := 0; i < routines; i++ {
		if <-replaced {
			changed += 1
		}
	}
	return changed
}

func replaceByExactMatch(toFind string, toReplace string, filenames []string, replaced chan bool) (int, []string) {
	routines := 0
	unmatched := make([]string, 0)
	for _, filename := range filenames {
		fileInfo, err := os.Stat(filename)
		if err == nil && !fileInfo.IsDir() {
			routines += 1
			go replaceFile(filename, toFind, toReplace, replaced)
		} else if err != nil && !os.IsNotExist(err) {
			panic(err)
		} else {
			unmatched = append(unmatched, filename)
		}
	}
	return routines, unmatched
}

func replaceByPatternMatch(toFind string, toReplace string, patterns []string, replaced chan bool) int {
	routines := 0
	files, _ := ioutil.ReadDir(".")
	for _, file := range files {
		if !file.IsDir() {
			for _, pattern := range patterns {
				match, _ := regexp.Match(pattern, []byte(file.Name()))
				if match {
					routines += 1
					go replaceFile(file.Name(), toFind, toReplace, replaced)
				}
			}
		}
	}
	return routines
}

func replaceFile(filename, toFind, toReplace string, ch chan bool) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			ch <- false
			return
		}
		panic(err)
	}
	filecontents := string(file)
	if strings.Contains(filecontents, toFind) {
		filecontents = strings.Replace(filecontents, toFind, toReplace, -1)
		err = ioutil.WriteFile(filename, []byte(filecontents), 0644)
		if err != nil {
			panic(err)
		}
		ch <- true
	} else {
		ch <- false
	}
}
