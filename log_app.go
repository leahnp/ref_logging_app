package main

import (
    // "io"
    "net/http"
    "html/template"
    "fmt"
    "io/ioutil"
    "os"
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

func stack_traces(w http.ResponseWriter, r *http.Request) {
    // file contents as string
    trace, err := ioutil.ReadFile("models/go_stack_trace")
    check(err)
    // convert trace to string
    str := string(trace)

    f, err := os.OpenFile("var/log/reference-logging", os.O_APPEND|os.O_WRONLY, 0600)
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
    // map with levels and messages 
    var levels map[string]string
    levels = make(map[string]string)

    levels["Fatal"] = "We're going doooowwwwnnnnn!!!!!!"
    levels["Panic"] = "This parachute is a napsack!"
    levels["Error"] = "Negatory...does not compute."
    levels["Warn"] = "Hey buddy - think again!"
    levels["Debug"] = "Dude. Get to work."
    levels["Trace"] = "Happy hunting."

    // randomizer to pick random log
    // fmt.Printf(levels)
    fmt.Printf("%s", levels)

    // print log to file, stdout and stderr
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/stack_traces", stack_traces)
    http.HandleFunc("/levels", levels)
    http.ListenAndServe(":8080", nil)
}