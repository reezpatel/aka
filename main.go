package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/reezpatel/aka/inmem"
)

type requestBody struct {
	URL string `json:"url"`
	Aka string `json:"aka"`
}

type applicationContext struct {
	db inmem.InMem
}

func (ctx applicationContext) handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := requestBody{}
	err = json.Unmarshal(body, &data)
	method := r.Method

	if method == "POST" || method == "PATCH" || method == "DELETE" {
		auth := r.Header.Get("Authorization")

		if auth != os.Getenv("AUTH_TOKEN") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	}

	switch method {
	case "POST":
		_, found := ctx.db.Get(data.Aka)
		if data.Aka != "" && !found {
			fmt.Println("POST", data.Aka, "-", data.Aka)
			ctx.db.Add(data.Aka, data.URL)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
			break
		}

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Failed"))

	case "PATCH":
		_, found := ctx.db.Get(data.Aka)
		if data.Aka != "" && found {
			fmt.Println("PATCH", data.Aka, "-", data.URL)
			ctx.db.Update(data.Aka, data.URL)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
			break
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Failed"))

	case "DELETE":
		_, found := ctx.db.Get(data.Aka)
		if data.Aka != "" && found {
			fmt.Println("DELETE", data.Aka)
			ctx.db.Remove(data.Aka)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
			break
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Failed"))

	case "GET":
		key := r.URL.Path[1:]
		value, found := ctx.db.Get(key)
		if found {
			fmt.Println("GET", key, "-", value)
			http.Redirect(w, r, value, http.StatusTemporaryRedirect)
			break
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}

}

func main() {
	filename := path.Join(os.Getenv("DATA_PATH"), "db.json")

	db := inmem.New()
	appCtx := applicationContext{db}
	channel := make(chan os.Signal, 1)
	ticker := time.Tick(2 * time.Hour)

	db.Load(filename)
	fmt.Println("DATA LOADED FROM", filename)

	defer db.Persist(filename)

	http.HandleFunc("/", appCtx.handleRequest)

	signal.Notify(channel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		<-channel
		db.Persist(filename)
		fmt.Println("Saving before exiting...")
		os.Exit(1)
	}()

	go func() {
		for range ticker {
			fmt.Println("Persisting...")
			db.Persist(filename)
		}
	}()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
