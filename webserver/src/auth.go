/* auth.go */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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
