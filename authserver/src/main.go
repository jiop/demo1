package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const PORT = "8002"

type user struct {
	Username string `json:"username"`
	Password string `json:"username"`
	Token    string `json:"username"`
}

var users = make(map[string]string)

var seedUsers = []user{
	user{
		Username: "user1",
		Password: "pass1",
	},
	user{
		Username: "user2",
		Password: "pass2",
	},
	user{
		Username: "user3",
		Password: "pass3",
	},
}

func main() {
	log.Printf("AuthServer started. Listening on port %s", PORT)
	http.HandleFunc("/login", httpLoginRequest)
	http.HandleFunc("/logout", httpLogoutRequest)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func httpLoginRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("Login request")
	username := r.FormValue("username")
	password := r.FormValue("password")
	if token := validateUser(username, password); token == "" {
		log.Print("Unauthorized user trying to connect")
		json.NewEncoder(w).Encode(map[string]string{
			"testval":  "OK",
			"username": username,
			"token":    "",
		})
		log.Printf("username: %s, testval: %s, token: %s", "OK", username, "")
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"testval":  "OK",
			"username": username,
			"token":    token,
		})
		log.Printf("username: %s, testval: %s, token: %s", "OK", username, token)
		log.Print("Authorized user connected")
	}
}

func validateUser(username string, password string) string {
	for _, u := range seedUsers {
		if username == u.Username {
			if password == u.Password {
				return generateSessionToken()
			} else {
				return ""
			}
		}
	}
	return ""
}

func generateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}

func httpLogoutRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("Logout request")
	fmt.Fprint(w, "Logout")
}
