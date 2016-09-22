package main

import (
    "io"
    "net/http"
    "html/template"
    "fmt"
    "io/ioutil"
    "os"
    "math/rand"
    // "bytes"
    "log"
    // "strings"
    "time"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, "templates/index.html")
}

// currently only returns go stack trace
func stack_traces(w http.ResponseWriter, r *http.Request) {
    // file contents as string
    trace, err := ioutil.ReadFile("models/go_stack_trace")
    check(err)
    // convert trace to string
    str := string(trace)

    f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0600)
    check(err)

    defer f.Close()

    if _, err = f.WriteString(str + "\n"); err != nil {
        panic(err)
    }

    // print stack traces to stdout & stderr
    fmt.Printf(str)
    fmt.Fprintln(os.Stderr, "hello world")

    http.Redirect(w, r, "/", http.StatusFound)
}

// prints log with random level to Stdout Stderr and reference-logging
func levels(w http.ResponseWriter, r *http.Request) {
    // array of levels and level messages
    var levels_array = [6][2]string{ 
                {"Fatal", "We're going doooowwwwnnnnn!!!!!!"}, 
                {"Panic", "This parachute is a napsack!"}, 
                {"Error", "Negatory...does not compute."}, 
                {"Warn", "Hey buddy - think again!"},
                {"Debug", "Dude. Get to work."},
                {"Trace", "Happy hunting."},
            }

    // picks random level and message
    message := levels_array[rand.Intn(len(levels_array))]

    // opens reference-logging file
    f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    check(err)

    defer f.Close()

    // sets output to print to file, stdout, stderr
    log.SetOutput(io.MultiWriter(f, os.Stdout, os.Stderr))

    // print log with level and message
    log.Println("[" + message[0] + "] " + message[1])

    http.Redirect(w, r, "/", http.StatusFound)
}

// send batch logs all at once
func batch(w http.ResponseWriter, r *http.Request) {
    str_num := r.URL.Query().Get("num")
    // convert num string to int
    num, err := strconv.Atoi(str_num)
    check(err)

    // send num number of logs
    for i := 0; i < num; i++ {
        // file contents as string
        trace, err := ioutil.ReadFile("models/levels")
        check(err)
        // convert trace to string
        str := string(trace)

        f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0600)
        check(err)

        defer f.Close()

        if _, err = f.WriteString(str + "\n"); err != nil {
            panic(err)
        }

        // print stack traces to stdout & stderr
        fmt.Printf(str)
        fmt.Fprintln(os.Stderr, "hello world")

        http.Redirect(w, r, "/", http.StatusFound)
    }

    fmt.Println(num)
    http.Redirect(w, r, "/", http.StatusFound)
}

// call function every x millis
func doEvery(d time.Duration, f func(time.Time)) {
    for x := range time.Tick(d) {
        f(x)
    }
}

// reoccuring message
func ahoy(t time.Time) {
    log.Println("Ahoy.")
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/stack_traces", stack_traces)
    http.HandleFunc("/levels", levels)
    http.HandleFunc("/batch", batch)
    // TODO this blocks the rest of the website
    // doEvery(1000*time.Millisecond, ahoy)
    http.ListenAndServe(":8080", nil)
}