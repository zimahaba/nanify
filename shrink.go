package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

func ShrinkPage(w http.ResponseWriter, r *http.Request) {
	handleCookie(w, r)
	render(w, pages["shrink"], nil)
}

func ShrinkPageHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := handleCookie(w, r)

		r.ParseForm()

		longUrl := r.FormValue("url")
		_, err := url.ParseRequestURI(longUrl)
		if err != nil {
			render(w, pages["shrink"], map[string]string{"ErrorMsg": "Invalid URL."})
			return
		}

		alias := r.FormValue("alias")

		key := alias
		if alias == "" {
			key = generateKey()
		}

		ctx := context.Background()
		shortUrl := "http://localhost:8080/shrunk/" + key

		if err := redisClient.Set(ctx, key, longUrl, 24*time.Hour).Err(); err != nil {
			panic(err)
		}

		hasSet := redisClient.Exists(ctx, cookie.Value).Val() > 0
		redisClient.SAdd(ctx, cookie.Value, key)
		if !hasSet {
			redisClient.Expire(ctx, cookie.Value, 168*time.Hour)
		}

		render(w, pages["result"], map[string]interface{}{"OriginalUrl": longUrl, "ShortUrl": shortUrl})
	}
}

func RedirectHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := handleCookie(w, r)
		fmt.Println("cvalue: " + cookie.Value)
		key := r.URL.Path[len("/shrunk/"):]

		ctx := context.TODO()

		url, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			fmt.Println("Error getting url from cache")
			t := template.Must(template.ParseFiles(pages["base"], pages["expired"]))
			err := t.ExecuteTemplate(w, "base.html", nil)
			if err != nil {
				fmt.Println("Error when executing template", err)
			}
		} else {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}
	}
}

func MyUrlsPage(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := handleCookie(w, r)
		urls := []interface{}{}

		values, err := redisClient.SMembers(context.TODO(), cookie.Value).Result()
		if err == nil {
			for _, v := range values {
				url, _ := redisClient.Get(context.TODO(), v).Result()
				shortUrl := "http://localhost:8080/shrunk/" + v
				urls = append(urls, ShrinkUrl{LongUrl: url, ShortUrl: shortUrl})
			}
		}

		render(w, pages["myurls"], map[string]interface{}{"Urls": urls})
	}
}

type ShrinkUrl struct {
	LongUrl  string
	ShortUrl string
}
