package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {

	redisClient := createCache()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", ShrinkPage)
	mux.HandleFunc("GET /shrunk/", RedirectHandler(redisClient))
	mux.HandleFunc("POST /", ShrinkPageHandler(redisClient))
	mux.HandleFunc("GET /myurls", MyUrlsPage(redisClient))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/styles"))))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
