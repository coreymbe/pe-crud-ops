package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SensitiveParameters struct {
	PrivateKeyContent string `json:"private-key-content"`
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
	help := flag.Bool("help", false, "Optional, prints usage info")
	pe_console := flag.String("pe_console", "", "Required, the hostname of the PE Console")
	agent := flag.String("agent", "", "Required, the hostname of the agent node")
	token := flag.String("token", "", "Required, the RBAC token to create the connection")
	ssh_user := flag.String("ssh_user", "root", "Optional, the RBAC token to create the connection")
	ssh_key := flag.String("ssh_key", "", "Required, path to SSH private key to create the connection")
	flag.Parse()

	usage := `usage:
	
create_conn -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME> -token <PE_TOKEN> -ssh_key <PRIVATE_KEY_PATH> [-ssh_user <SSH_USER> -help]
	
Options:
	-help        Optional, Prints this message.
	-pe_console  Required, The hostname of the PE Console.
	-agent       Required, The hostname of the agent node to add.
	-token       Required, PE RBAC token with appropriate permissions.
	-ssh_user    Optional, The SSH user to access the agent node (Default: "root").
	-ssh_key     Required, The path to the SSH private key for the configured SSH user.
 `

	if *help == true {
		fmt.Println(usage)
		return
	}

	// There has to be a better way to ensure all of the required options are set.
	if *pe_console == "" {
		log.Fatalf("The pe_console option is required:\n%s", usage)
	}
	if *agent == "" {
		log.Fatalf("The agent option is required:\n%s", usage)
	}
	if *token == "" {
		log.Fatalf("The token option is required:\n%s", usage)
	}
	if *ssh_key == "" {
		log.Fatalf("The ssh_key option is required:\n%s", usage)
	}

	private_key, err := ioutil.ReadFile(*ssh_key)
	if err != nil {
		log.Fatalf("Unable to read SSH key: %v", err)
	}
	pk := string(private_key)

	task_data := Payload{
		Certnames: []string{*agent},
		Type:      "ssh",
		Parameters: Parameters{
			User: *ssh_user,
		},
		SensitiveParameters: SensitiveParameters{
			PrivateKeyContent: pk,
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

	req, err := http.NewRequest("POST", "https://"+*pe_console+":8143/inventory/v1/command/create-connection", post_body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authentication", *token)

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
