package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Authorization struct {
	Channels    []string `bson:"channels" json:"channels"`
	Permissions []string `bson:"permissions" json:"permissions"`
	Topic       string   `bson:"topic" json:"topic"`
}

type AuthResponse struct {
	Authorizations []Authorization `bson:"authorizations" json:"authorizations"`
	Identity       string          `bson:"identity" json:"identity"`
	IdentityURL    string          `bson:"identity_url" json:"identity_url"`
	Ttl            int             `bson:"ttl" json:"ttl"`
}

var (
	ttl   *int
	users *mgo.Collection
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func auth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	var authResponse AuthResponse
	auth_secret := req.URL.Query().Get("auth_secret")
	remote_ip := req.URL.Query().Get("remote_ip")
	tls := req.URL.Query().Get("tls")
	log.Printf("%s from %s with tls %s attempted to authenticate", auth_secret, remote_ip, tls)
	err := users.Find(bson.M{"_id": auth_secret}).One(&authResponse)
	if err != nil {
		log.Println("Error occured while quering for user")
		log.Println(err)
	}
	authResponse.Ttl = *ttl
	authBytes, _ := json.Marshal(authResponse)
	w.Write(authBytes)
}

func main() {
	tls := flag.Bool("tls", false, "-tls=\"true/false\"")
	ttl = flag.Int("ttl", 3600, "-ttl=3600")
	mongoUrl := flag.String("mongoserver", "localhost:27017", "-mongoserver=localhost:27017")
	flag.Parse()

	session, err := mgo.Dial(*mongoUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	users = session.DB("nsqauth").C("users")

	http.HandleFunc("/", handler)
	http.HandleFunc("/auth", auth)
	if *tls {
		log.Printf("About to listen on 10443")
		log.Fatal(http.ListenAndServeTLS(":10443", "private/mycert1.cer", "private/mycert1.key", nil))
	} else {
		log.Printf("About to listen on 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
