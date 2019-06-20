package main

import (
  "net/http"
  "log"
)

func main() {
  // http.HandleFunc("/signin", Signin)
  http.HandleFunc("/signup", Signup)


  log.Fatal(http.ListenAndServe(":8000", nil))
}
