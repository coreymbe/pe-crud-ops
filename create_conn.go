package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SensitiveParameters struct {
	Password string `json:"password"`
}
type Parameters struct {
	User string `json:"user"`
}
type Payload struct {
	Certnames           []string            `json:"certnames"`
	Type                string              `json:"type"`
	Parameters          Parameters          `json:"parameters"`
	SensitiveParameters SensitiveParameters `json:"sensitive_parameters"`
	Duplicates          string              `json:"duplicates"`
}

func main() {
	task_data := Payload{
		Certnames: []string{"{AGENT_HOSTNAME}"},
		Type:      "ssh",
		Parameters: Parameters{
			User: "{SSH_USER}",
		},
		SensitiveParameters: SensitiveParameters{
			Password: "{USER_PASS}",
		},
		Duplicates: "replace",
	}

	payloadBytes, err := json.Marshal(task_data)
	if err != nil {
		log.Fatal(err)
	}

	post_body := bytes.NewReader(payloadBytes)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", os.ExpandEnv("https://${PE_CONSOLE}:8143/inventory/v1/command/create-connection"), post_body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authentication", "{PE_TOKEN}")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
