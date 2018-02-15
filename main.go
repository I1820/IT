/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 15-02-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

// Config represents main configuration
var Config = struct {
	BackBack struct {
		BaseURL string `default:"https://backback.ceit.aut.ac.ir/api/" env:"backback_base_url"`
		Version string `default:"v1" env:"backback_version"`
	}
}{}

// JWT Token
var jwtToken token

func main() {
	// Disable https certificate validation
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	login()
}

func login() {
	resp, err := http.PostForm(Config.BackBack.BaseURL+Config.BackBack.Version+"/login", url.Values{"email": {"sepehr.sabour@gmail.com"}, "password": {"1234567"}})
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("StatusCode: %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Body: %s", err)
	}

	if err := resp.Body.Close(); err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Body: %s", err)
	}

	var response struct {
		Code   int
		Result struct {
			User  user
			Token token
		}
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "login",
	}).Infoln(response)

	jwtToken = response.Result.Token
}

func createProject() {
}

func createThingProfile() {
}

func createThing() {
}

func activateThing() {
}
