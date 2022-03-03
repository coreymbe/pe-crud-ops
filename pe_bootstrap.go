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
	help := flag.Bool("help", false, "Optional, prints usage info")
	pe_console := flag.String("pe_console", "", "Required, the hostname of the PE Console")
	agent := flag.String("agent", "", "Required, the hostname of the agent node")
	token := flag.String("token", "", "Required, the RBAC token to create the connection")
	ssh_user := flag.String("ssh_user", "root", "Optional, the RBAC token to create the connection")
	flag.Parse()

	usage := `usage:
	
pe_bootstrap -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME> -token <PE_TOKEN> [-ssh_user <SSH_USER> -help]
	
Options:
	-help        Optional, Prints this message.
	-pe_console  Required, The hostname of the PE Console.
	-agent       Required, The hostname of the agent node to add.
	-token       Required, PE RBAC token with appropriate permissions.
	-ssh_user    Optional, The SSH user to access the agent node (Default: "root").
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

	task_data := Payload{
		Task: "pe_bootstrap",
		Params: Params{
			Server: *pe_console,
		},
		Scope: Scope{
			Nodes: []string{*agent},
		},
		Target: []Target{
			{
				Hostnames: []string{*agent},
				User:      *ssh_user,
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

	req, err := http.NewRequest("POST", "https://"+*pe_console+":8143/orchestrator/v1/command/task", post_body)
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
