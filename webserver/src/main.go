package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const WEBSERVER_PORT = "8001"
const AUTHSERVER_URL = "http://demo1_authserver1_1"
const AUTHSERVER_PORT = "8002"

type Auth struct{}

func (auth *Auth) Login(username string, password string) string {
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)
	resp, err := http.PostForm(AUTHSERVER_URL+":"+AUTHSERVER_PORT+"/login", v)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	r := make(map[string]interface{})
	json.Unmarshal(data, &r)
	token := r["token"].(string)

	if token != "" {
		return token
	}
	return ""
}

func (auth *Auth) Logout(username string, token string) bool {
	v := url.Values{}
	v.Set("username", username)
	v.Set("token", token)
	resp, err := http.PostForm(AUTHSERVER_URL+":"+AUTHSERVER_PORT+"/logout", v)
	if err != nil {
		return false
	}
	json.NewDecoder(resp.Body).Decode(v)
	log.Print(v)
	return true
}

func (auth *Auth) VerifyToken(username string, token string) bool {
	v := url.Values{}
	v.Set("username", username)
	v.Set("token", token)
	resp, err := http.PostForm(AUTHSERVER_URL+":"+AUTHSERVER_PORT+"/verifytoken", v)
	if err != nil {
		return false
	}
	json.NewDecoder(resp.Body).Decode(v)
	log.Print(v)
	return true
}

func main() {
	log.Printf("WebServer started. Listening on port %s", WEBSERVER_PORT)
	http.HandleFunc("/login", httpLoginRequest)
	http.HandleFunc("/logout", httpLogoutRequest)
	http.HandleFunc("/unprotected-content", httpUnprotectedContentRequest)
	log.Fatal(http.ListenAndServe(":"+WEBSERVER_PORT, nil))
}

func httpLoginRequest(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	auth := Auth{}
	if username != "" && password != "" {
		token := auth.Login(username, password)
		if token != "" {
			cookie := http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(time.Hour)}
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
	auth := Auth{}
	if err1 == nil && err2 == nil && auth.Logout(username.Value, token.Value) {
		cookie := http.Cookie{Name: "username", Value: "", Expires: time.Now()}
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
	auth := Auth{}
	if err1 == nil && err2 == nil && auth.VerifyToken(username.Value, token.Value) {
		json.NewEncoder(w).Encode(map[string]string{
			"content": "protected content",
		})
		log.Printf("User %s accessed protected successfully", username)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"err": "protected content access failure",
	})
	log.Print("Failed access to protected Content")
}
