package main

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func ShrinkPage(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/shrink.html"))
	err := t.ExecuteTemplate(w, "base.html", nil)
	if err != nil {
		fmt.Println("Error when executing template", err)
	}
}

func ShrinkPageHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		url := r.FormValue("url")

		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		const keyLength = 6

		shortKey := make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = charset[rand.Intn(len(charset))]
		}

		ctx := context.Background()
		shortUrl := "http://localhost:8080/shrunk/" + string(shortKey)

		key := string(shortKey)

		if err := redisClient.Set(ctx, key, url, 60*time.Second).Err(); err != nil {
			panic(err)
		}

		t := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/result.html"))
		err := t.ExecuteTemplate(w, "base.html", map[string]interface{}{"OriginalUrl": url, "ShortUrl": shortUrl})
		if err != nil {
			fmt.Println("Error when executing template", err)
		}
	}
}

func RedirectHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/shrunk/"):]

		ctx := context.TODO()

		url, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			fmt.Println("Error getting url from cache")
			t := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/expired.html"))
			err := t.ExecuteTemplate(w, "base.html", nil)
			if err != nil {
				fmt.Println("Error when executing template", err)
			}
		} else {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}
	}
}
