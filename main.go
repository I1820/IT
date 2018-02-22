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
		BaseURL string `default:"http://backback.ceit.aut.ac.ir/api/" env:"backback_base_url"`
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

	createUser()
	login()
	createProject()
}

func createUser() {
	resp, err := http.PostForm(Config.BackBack.BaseURL+Config.BackBack.Version+"/register",
		url.Values{
			"legal":    {"0"},
			"name":     {"Parham Alvani"},
			"email":    {"parham.alvani@yahoo.com"},
			"mobile":   {"09390909540"},
			"password": {"1234567"},
		})
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Phase": "resiter",
		}).Fatalf("StatusCode: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("Body: %s", err)
	}

	if err := resp.Body.Close(); err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("Body: %s", err)
	}

	var response struct {
		Code int
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "register",
	}).Infoln(response)
}

func login() {
	resp, err := http.PostForm(Config.BackBack.BaseURL+Config.BackBack.Version+"/login", url.Values{"email": {"parham.alvani@yahoo.com"}, "password": {"1234567"}})
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("StatusCode: %d", resp.StatusCode)
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
	resp, err := http.PostForm(Config.BackBack.BaseURL+Config.BackBack.Version+"/project", url.Values{"name": {"Me"}, "description": {"This is me"}})
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Fatalf("StatusCode: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Fatalf("Body: %s", err)
	}

	if err := resp.Body.Close(); err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Fatalf("Body: %s", err)
	}

	var response struct {
		Code int
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "create project",
	}).Infoln(response)
}

func createThingProfile() {
}

func createThing() {
}

func activateThing() {
}
