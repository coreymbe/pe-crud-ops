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

type Target struct {
	Hostnames []string `json:"hostnames"`
	User      string   `json:"user"`
	Transport string   `json:"transport"`
}
type Scope struct {
	Nodes []string `json:"nodes"`
}
type Params struct {
	Server string `json:"server"`
}
type Payload struct {
	Task   string   `json:"task"`
	Params Params   `json:"params"`
	Scope  Scope    `json:"scope"`
	Target []Target `json:"target"`
}

func main() {
	task_data := Payload{
		Task: "pe_bootstrap",
		Params: Params{
			Server: "{PE_CONSOLE}",
		},
		Scope: Scope{
			Nodes: []string{"{AGENT_HOSTNAME}"},
		},
		Target: []Target{
			Target{
				Hostnames: []string{"{AGENT_HOSTNAME}"},
				User:      "{SSH_USER}",
				Transport: "ssh",
			},
		},
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

	req, err := http.NewRequest("POST", os.ExpandEnv("https://${HOSTNAME}:8143/orchestrator/v1/command/task"), post_body)
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
