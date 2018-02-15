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
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func main() {
	login()
}

func login() {
	resp, err := http.PostForm("https://backback.ceit.aut.ac.ir/api/v1/login", url.Values{"email": {"sepehr.sabour@gmail.com"}, "password": {"1234567"}})
	if err != nil {
		log.Fatalf("Login: %s", err)
	}
	log.Infoln(resp)
}

func createProject() {
}

func createThingProfile() {
}

func createThing() {
}

func activateThing() {
}
