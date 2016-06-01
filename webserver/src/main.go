/* main.go */

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const WEBSERVER_PORT = "8001"
const AUTHSERVER_URL = "http://demo1_authserver1_1"
const AUTHSERVER_PORT = "8002"

var auth = Auth{}

func main() {
	log.Printf("WebServer started. Listening on port %s", WEBSERVER_PORT)
	http.HandleFunc("/login", httpLoginRequest)
	http.HandleFunc("/logout", httpLogoutRequest)
	http.HandleFunc("/unprotected-content", httpUnprotectedContentRequest)
	http.HandleFunc("/protected-content", httpProtectedContentRequest)
	log.Fatal(http.ListenAndServe(":"+WEBSERVER_PORT, nil))
}

func httpLoginRequest(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != "" && password != "" {
		token := auth.Login(username, password)
		if token != "" {
			cookie := http.Cookie{
				Name:    "username",
				Value:   username,
				Expires: time.Now().Add(time.Hour),
			}
			http.SetCookie(w, &cookie)

			cookie = http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(time.Hour),
			}
			http.SetCookie(w, &cookie)

			json.NewEncoder(w).Encode(map[string]string{
				"username": username,
				"token":    token,
				"logged":   "true",
			})
			log.Printf("User %s Logged in successfully.", username)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"err": "Login failure",
	})
	log.Printf("Failed Login. Username used %s", username)
}

func httpLogoutRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("Starting logout request")

	username, err1 := r.Cookie("username")
	token, err2 := r.Cookie("token")

	if err1 == nil && err2 == nil && auth.Logout(username.Value, token.Value) {
		cookie := http.Cookie{
			Name:    "username",
			Value:   "",
			Expires: time.Now(),
		}
		http.SetCookie(w, &cookie)

		cookie = http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now(),
		}
		http.SetCookie(w, &cookie)

		json.NewEncoder(w).Encode(map[string]string{
			"username": username.Value,
			"token":    token.Value,
			"logged":   "false",
		})
		log.Printf("User %s logged out successfully", username)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"err": "Logout failure",
	})
	log.Printf("Failed logout. Username %s", username)
}

func httpUnprotectedContentRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"content": "unprotected content",
	})
	log.Print("Unprotected Content Request")
}

func httpProtectedContentRequest(w http.ResponseWriter, r *http.Request) {
	username, err1 := r.Cookie("username")
	token, err2 := r.Cookie("token")

	log.Printf("COOKIE CONTENT : %s %s", username, token)

	if err1 == nil && err2 == nil && auth.VerifyToken(username.Value, token.Value) {
		json.NewEncoder(w).Encode(map[string]string{
			"content": "protected content",
		})
		log.Printf("User %s accessed protected successfully", username)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"err": "protected content access failure",
	})
	log.Print("Failed access to protected Content")
}
