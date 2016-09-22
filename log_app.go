package main

import (
    "io"
    "net/http"
    "html/template"
    // "fmt"
    "io/ioutil"
    "os"
    // "math/rand"
    "log"
    "time"
    "strconv"
    "regexp"
)

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
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, "templates/index.html")
    // loop("./models")
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

    f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0600)
    check(err)

    defer f.Close()

    if _, err = f.WriteString(str + "\n\n"); err != nil {
        panic(err)
    }

    log.SetOutput(io.MultiWriter(f, os.Stdout, os.Stderr))
    log.Println(str)

}

// function to loop through files in directory, return slice of file names
func loop(filepath string) (files []string) {
    dirname := filepath
    d, err := os.Open(dirname)
    check(err)

    defer d.Close()

    fi, err := d.Readdir(-1)
    check(err)

    // slice := make([]string, 3)
    slice := []string{}
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            // fmt.Println(fi.Name())
            file := string(fi.Name())
            slice = append(slice, file)
        }
    }
    return slice
}

// currently only returns go stack trace
func stack_traces(w http.ResponseWriter, r *http.Request) {
    // TODO get array of file names in model and regex for "stack"
    process_file("models/go_stack_trace")

    http.Redirect(w, r, "/", http.StatusFound)
}

// prints log with random level to Stdout Stderr and reference-logging
func levels(w http.ResponseWriter, r *http.Request) {
    // array of levels and level messages
    // var levels_array = [6][2]string{ 
    //             {"Fatal", "We're going doooowwwwnnnnn!!!!!!"}, 
    //             {"Panic", "This parachute is a napsack!"}, 
    //             {"Error", "Negatory...does not compute."}, 
    //             {"Warn", "Hey buddy - think again!"},
    //             {"Debug", "Dude. Get to work."},
    //             {"Trace", "Happy hunting."},
    //         }

    // // picks random level and message
    // message := levels_array[rand.Intn(len(levels_array))]

    // // opens reference-logging file
    // f, err := os.OpenFile("var/log/reference-logging", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    // check(err)

    // defer f.Close()

    // // sets output to print to file, stdout, stderr
    // log.SetOutput(io.MultiWriter(f, os.Stdout, os.Stderr))

    // // print log with level and message
    // log.Println("[" + message[0] + "] " + message[1])


    // get array of file names loop
    files := loop("models")

    levels := []string{}
    // iterate over file names  
    for _, file := range files {
        // process file include "level" keyword
        r, _ := regexp.Compile("level")
        if r.FindStringSubmatch(file) != nil {
            levels = append(levels, file)
        }        
    }

    // process files 
    for _, file := range levels {
        process_file("models/" + file)
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

// send batch logs all at once (currently only stack traces)
func batch(w http.ResponseWriter, r *http.Request) {
    str_num := r.URL.Query().Get("num")
    // convert num string to int
    num, err := strconv.Atoi(str_num)
    check(err)

    // send num number of logs
    for i := 0; i < num; i++ {
        process_file("models/go_stack_trace")

        http.Redirect(w, r, "/", http.StatusFound)
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