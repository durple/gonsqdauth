package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

type Authorization struct {
	Channels    []string `json:"channels"`
	Permissions []string `json:"permissions"`
	Topic       string   `json:"topic"`
}

type AuthResponse struct {
	Authorizations []Authorization `json:"authorizations"`
	Identity       string          `json:"identity"`
	IdentityURL    string          `json:"identity_url"`
	Ttl            int             `json:"ttl"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func auth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	var authResponse AuthResponse
	var authorization Authorization

	ttl := 3600
	identity := "authUser"
	identityUrl := "no url"

	permissions := []string{"subscribe", "publish"}
	topic := ".*"
	channels := []string{".*"}
	authorization.Channels = channels
	authorization.Permissions = permissions
	authorization.Topic = topic
	authResponse.Authorizations = []Authorization{authorization}
	authResponse.Identity = identity
	authResponse.IdentityURL = identityUrl
	authResponse.Ttl = ttl
	authBytes, _ := json.Marshal(authResponse)
	w.Write(authBytes)
}

func main() {
	var tls = flag.Bool("tls", false, "-tls=\"true/false\"")
	flag.Parse()
	http.HandleFunc("/", handler)
	http.HandleFunc("/auth", auth)
	if *tls {
		log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
		log.Fatal(http.ListenAndServeTLS(":10443", "private/mycert1.cer", "private/mycert1.key", nil))
	} else {
		log.Printf("About to listen on 8080. Go to https://127.0.0.1:8080/")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
