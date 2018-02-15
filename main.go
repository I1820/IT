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
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Disable https certificate validation
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	login()
}

func login() {
	resp, err := http.PostForm("https://backback.ceit.aut.ac.ir/api/v1/login", url.Values{"email": {"sepehr.sabour@gmail.com"}, "password": {"1234567"}})
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

	log.WithFields(log.Fields{
		"Phase": "login",
	}).Infoln(string(data))
}

func createProject() {
}

func createThingProfile() {
}

func createThing() {
}

func activateThing() {
}
