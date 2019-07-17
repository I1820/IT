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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

// Config represents main configuration
var Config = struct {
	BackBack struct {
		BaseURL string `default:"http://127.0.0.1:7070/api/" env:"backback_base_url"`
		Version string `default:"v1" env:"backback_version"`
	}
}{}

// JWT Token
var jwtToken token

var projectID = "5b96bdf969ccb0000a1bb24a"
var thingID = "0000000000000088"

var concurrentRequests = 500
var pipelineRequests = 1

// BackBack is an I1820 sjd-backend client
type BackBack struct {
	client *resty.Client
}

func main() {
	bk := BackBack{
		client: resty.New(),
	}

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// createUser()
	bk.login()
	// createProject()

	failed := 0
	success := 0
	var responseTime float64
	var responseTimeMax float64
	responseTimeStream := make(chan float64, concurrentRequests*pipelineRequests)
	var wg sync.WaitGroup
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for i := 0; i < pipelineRequests; i++ {
				before := time.Now()
				if err := bk.fetchData(); err != nil {
					failed++
					continue
				}
				success++
				interval := time.Now().Sub(before)
				fmt.Printf("%d took %s on loop %d\n", index, interval, i)
				responseTimeStream <- interval.Seconds()
			}
		}(i)
	}
	wg.Wait()
	close(responseTimeStream)
	fmt.Println("Fetch data finished")

	for t := range responseTimeStream {
		responseTime += t
		if t > responseTimeMax {
			responseTimeMax = t
		}
	}

	fmt.Printf("Total: %d, Failed: %d, Success: %d\nRatio: %g%%\n", success+failed, failed, success, float64(failed*100)/float64(success+failed))
	fmt.Printf("Response Time Avg. %gs\n", responseTime/float64(success))
	fmt.Printf("Response Time Max. %gs\n", responseTimeMax)
}

func (bk BackBack) createUser() {
	resp, err := bk.client.R().
		SetFormData(map[string]string{
			"legal":    "0",
			"name":     "Parham Alvani",
			"email":    "parham.alvani@gmail.com",
			"mobile":   "09390909540",
			"password": "123123",
		}).
		Post(Config.BackBack.BaseURL + Config.BackBack.Version + "/register")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "resiter",
		}).Fatalf("StatusCode: %d", resp.StatusCode())
	}

	var response struct {
		Code int
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "register",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "register",
	}).Infoln(response)
}

func (bk BackBack) login() {
	resp, err := bk.client.R().
		SetFormData(map[string]string{
			"email":    "parham.alvani@gmail.com",
			"password": "123123",
		}).
		Post(Config.BackBack.BaseURL + Config.BackBack.Version + "/login")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("StatusCode: %d", resp.StatusCode())
	}

	var response struct {
		Code   int `json:"code"`
		Result struct {
			User        user  `json:"user"`
			AccessToken token `json:"access_token"`
		} `json:"result"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "login",
	}).Infoln(response)

	jwtToken = response.Result.AccessToken
}

func (bk BackBack) fetchData() error {
	resp, err := bk.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetFormData(map[string]string{
			"project_id": projectID,
			"since":      "0",
			"thing_ids":  fmt.Sprintf(`{"ids": [%s]}`, thingID),
		}).
		Post(Config.BackBack.BaseURL + Config.BackBack.Version + "/things/data")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "fetch data",
		}).Errorf("Request: %s", err)
		return err
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "fetch data",
		}).Errorf("StatusCode: %d", resp.StatusCode())
		return fmt.Errorf("StatusCode %d", resp.StatusCode())
	}

	var response interface{}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "fetch data",
		}).Errorf("JSON Unmarshal: %s", err)
		return err
	}

	log.WithFields(log.Fields{
		"Phase": "fetch data",
	}).Infoln(response)
	return nil
}

func (bk BackBack) createProject() {
	resp, err := bk.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetFormData(map[string]string{
			"name":        "Me",
			"description": "This is me",
		}).
		Post(Config.BackBack.BaseURL + Config.BackBack.Version + "/project")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Errorf("Request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Errorf("StatusCode: %d", resp.StatusCode())
	}

	var response struct {
		Code int
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "create project",
		}).Errorf("JSON Unmarshal: %s", err)
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
