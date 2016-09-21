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

func stack_traces(w http.ResponseWriter, r *http.Request, language string) {
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
    // use log

    http.Redirect(w, r, "/", http.StatusFound)
}

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

    message := levels_array[rand.Intn(len(levels_array))]


    f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    check(err)
    defer f.Close()

    log.SetOutput(io.MultiWriter(f, os.Stdout))

    // print log with level and message
    // log.Println("[" + message[0] + "] " + message[1])
    http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/stack_traces", stack_traces)
    http.HandleFunc("/levels", levels)
    http.ListenAndServe(":8080", nil)
}