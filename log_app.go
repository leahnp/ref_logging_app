package main

import (
    // "io"
    "net/http"
    "html/template"
)

func index(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, "templates/index.html")
}

func main() {
    http.HandleFunc("/", index)
    http.ListenAndServe(":8080", nil)
}