package main

import (
    "io"
    "net/http"
    "html/template"
    "fmt"
    "io/ioutil"
    "os"
    "math/rand"
    "log"
    "time"
    "strconv"
    "regexp"
)

// TODO break out list of files into const array
// total log count
var counter int

// check error function 
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// index function 
func index(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("/templates/index.html")
    t.Execute(w, "/templates/index.html")
}

// take in file, write to logs and Stdout and Stderr
func process_file(file string) {
    // file contents as string
    content, err := ioutil.ReadFile(file)
    check(err)

    // increment log counter
    counter += 1

    // convert log content to string and append counter
    str := strconv.Itoa(counter) + ": " + string(content)

    f, err := os.OpenFile("/var/log/reference-logging.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
    check(err)

    defer f.Close()

    log.SetOutput(io.MultiWriter(f, os.Stdout, os.Stderr))
    log.Println(str)

}

// loop through files in directory, return slice of all file names
func loop(filepath string) (files []string) {
    dirname := filepath
    d, err := os.Open(dirname)
    check(err)

    defer d.Close()

    fi, err := d.Readdir(-1)
    check(err)

    slice := []string{}
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            file := string(fi.Name())
            slice = append(slice, file)
        }
    }
    return slice
}

// currently only returns go stack trace
func stack_traces(w http.ResponseWriter, r *http.Request) {
    // get array of file names loop
    files := loop("/models")

    stacks := []string{}
    // iterate over file names  
    for _, file := range files {
        // process file include "stack" keyword
        r, _ := regexp.Compile("stack")
        if r.FindStringSubmatch(file) != nil {
            stacks = append(stacks, file)
        }        
    }

    // process random file
    process_file("/models/" + stacks[rand.Intn(len(stacks))])

    http.Redirect(w, r, "/", http.StatusFound)
}

// prints log with random level to Stdout Stderr and reference-logging
func levels(w http.ResponseWriter, r *http.Request) {
    // get array of file names loop
    files := loop("/models")

    levels := []string{}
    // iterate over file names  
    for _, file := range files {
        // process file include "level" keyword
        r, _ := regexp.Compile("level")
        if r.FindStringSubmatch(file) != nil {
            levels = append(levels, file)
        }        
    }

    // process random file
    process_file("/models/" + levels[rand.Intn(len(levels))])

    http.Redirect(w, r, "/", http.StatusFound)
}

// send batch logs all at once (currently only stack traces)
func batch(w http.ResponseWriter, r *http.Request) {
    str_num := r.URL.Query().Get("num")
    // convert num string to int
    num, err := strconv.Atoi(str_num)
    check(err)

    files := loop("/models")
    fmt.Println(files)

    // send num number of logs
    for i := 0; i < num; i++ {
        // process one random file from all files in models
        process_file("/models/" + files[rand.Intn(len(files))])
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

// call function every x millis
func doEvery(d time.Duration, f func(time.Time)) {
    for x := range time.Tick(d) {
        f(x)
    }
}

// reoccuring message
func random_message(t time.Time) {
    // get array of file names loop
    files := loop("/models")

    // process random file
    process_file("/models/" + files[rand.Intn(len(files))])
}

// take in file, write to logs and Stdout and Stderr
func random_message_simple(t time.Time) {
    i := 10
    for i > 0 {

        content := "[ERROR]: PANIC This parachute is a napsack!"

        // increment log counter
        counter += 1

        // convert log content to string and append counter
        str := strconv.Itoa(counter) + ": " + string(content)

        log.Println(str)
        i--

    }

}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/stack_traces", stack_traces)
    http.HandleFunc("/levels", levels)
    http.HandleFunc("/batch", batch)
    // go doEvery(1000*time.Millisecond, random_message)
    // go doEvery(time.Millisecond, random_message_simple)
    http.ListenAndServe(":8080", nil)
}