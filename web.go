package main

import (
	"fmt"
	"net/http"
	"time"
)

func threadWeb() {
	log.Info("Starting web thread")
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", *webPort),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//http.HandleFunc("/control", httpControlHandler)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("web"))))
	httpServer.ListenAndServe()
}
