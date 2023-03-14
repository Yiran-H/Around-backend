package main

import (
    "fmt"
    "log"
    "net/http" 
    "around/backend"
    "around/handler"   
)
func main() {
    fmt.Println("started-service")
    backend.InitElasticsearchBackend();
    backend.InitGCSBackend(); //backend是一个package
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}