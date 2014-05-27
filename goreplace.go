package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "os"
    "io/ioutil"
    "strings"
    "regexp"
)

const debug = false

func main() {
    args := os.Args[1:]
    if debug {
        fmt.Println(os.Args)
    }
    if len(args) < 3 {
        panic("Usage: goreplace TOFIND TOREPLACE PATTERN...")
    }
    toFind := args[0]
    toReplace := args[1]
    patterns := args[2:]
    if debug {
        fmt.Println(toFind)
        fmt.Println(toReplace)
        fmt.Println(patterns)
    }
    changed := goReplace(toFind, toReplace, patterns)

    fmt.Printf("Done! Changed %d file(s).\n", changed)
}

func goReplace(toFind string, toReplace string, patterns []string) int {
    replaced := make(chan bool)
    routines := 0
    for _, filename := range patterns {
        fileInfo, err := os.Stat(filename)
        if err == nil && ! fileInfo.IsDir() {
            routines += 1
            go replaceFile(filename, toFind, toReplace, replaced)
        } else if err != nil && ! os.IsNotExist(err) {
            panic(err)
        }
    }
    files, _ := ioutil.ReadDir(".")
    for _, file := range files {
        if ! file.IsDir() {
            for _, pattern := range patterns {
                match, _ := regexp.Match(pattern, []byte(file.Name()))
                if match {
                    routines += 1
                    go replaceFile(file.Name(), toFind, toReplace, replaced)
                }
            }
        }
    }
    changed := 0
    for i := 0; i < routines; i++ {
        if <-replaced {
            changed += 1
        }
    }
    return changed
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

func catFile(filename string) {
    cmd := exec.Command("cat", filename)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Run()
    fmt.Println(out.String())
}
