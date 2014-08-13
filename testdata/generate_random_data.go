package main

import (
    "io/ioutil"
    "os"
    "bufio"
    "math/rand"
    "os/exec"
    "flag"
    "fmt"
)

func gen_random_text() string {
    // Cue crude way of generating random data
    file, err := os.Open("/usr/share/dict/words")
    if err != nil { panic(err) }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    i := 0
    output := ""
    for scanner.Scan() && i < 20000 {
        if rand.Intn(100) > 95 {
            output += " " + scanner.Text()
            i++
        }
    }
    return output
}

func generate_random_file(ch chan bool) {
    output := gen_random_text()
    uuid, err := exec.Command("uuidgen").Output()
    if err != nil { panic(err) }
    uuid = uuid[:len(uuid)-1] // Cut off the newline char
    filename := "testdata_" + string(uuid) + ".tmp"
    err = ioutil.WriteFile(filename, []byte(output), 0644)
    if err != nil { panic(err) }
    ch <- true
}

func main() {
    clean := flag.Bool("clean", false, "Remove testdata files")
    nr_of_files := flag.Int("files", 10, "Number of files to generate")
    flag.Parse()
    if *clean {
        fmt.Println("Cleaning up!")
        err := exec.Command("rm", "testdata.").Run()
        if err != nil { panic(err) }
        os.Exit(0)
    }
    finished := make(chan bool)
    for i := 0; i < *nr_of_files; i++ {
        go generate_random_file(finished)
    }
    for i :=0; i < *nr_of_files; i++ {
        <-finished
    }
}
