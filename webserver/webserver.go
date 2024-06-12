package webserver

import (
	"log"
	"net/http"
)

func Init(addr string )  {
    log.Println("Starting our webmonitor http server.")

    // Registering our handler functions, and creating paths.
    log.Println("Started on", addr)
    
    // Spinning up the server.
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal(err)
    }
	 
}