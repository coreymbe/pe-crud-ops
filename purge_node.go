package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CertState struct {
	DesiredState string `json:"desired_state"`
}
type PurgePayload struct {
	Certname string `json:"certname"`
}
type PurgeNode struct {
	Command string       `json:"command"`
	Version int          `json:"version"`
	Payload PurgePayload `json:"payload"`
}

func main() {
	help := flag.Bool("help", false, "Optional, prints usage info")
	pe_console := flag.String("pe_console", "", "Required, the hostname of the PE Console")
	agent := flag.String("agent", "", "Required, the hostname of the agent node")
	flag.Parse()

	usage := `usage:
	
purge_node -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME> [-help]
	
Options:
	-help        Optional, Prints this message.
	-pe_console  Required, The hostname of the PE Console.
	-agent       Required, The hostname of the agent node to add.
 `

	if *help == true {
		fmt.Println(usage)
		return
	}

	if (*pe_console == "") || (*agent == "") {
		log.Fatalf("The pe_console, agent, token, and ssh_key options are required:\n%s", usage)
	}

	// Requires a custom certificate called "pe-crud-ops-delete.pem".
	// The certificate needs to be added the allowlist for both the Puppet CA and PuppetDB APIs.
	cert, err := tls.LoadX509KeyPair("/etc/puppetlabs/puppet/ssl/certs/pe-crud-ops-delete.pem", "/etc/puppetlabs/puppet/ssl/private_keys/pe-crud-ops-delete.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile("/etc/puppetlabs/puppet/ssl/certs/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		},
	}

	revoke_data := CertState{
		DesiredState: "revoked",
	}
	revoke_payload, err := json.Marshal(revoke_data)
	if err != nil {
		log.Fatal(err)
	}
	put_body := bytes.NewReader(revoke_payload)

	purge_data := PurgeNode{
		Command: "deactivate node",
		Version: 3,
		Payload: PurgePayload{
			Certname: *agent,
		},
	}
	purge_payload, err := json.Marshal(purge_data)
	if err != nil {
		log.Fatal(err)
	}
	post_body := bytes.NewReader(purge_payload)

	client := &http.Client{Transport: tr}

	revoke_req, err := http.NewRequest("PUT", "https://"+*pe_console+":8140/puppet-ca/v1/certificate_status/"+*agent, put_body)
	if err != nil {
		log.Fatal(err)
	}
	revoke_req.Header.Set("Content-Type", "application/json")

	revoke_resp, err := client.Do(revoke_req)
	if err != nil {
		log.Fatal(err)
	}
	defer revoke_resp.Body.Close()

	revoke_status, err := revoke_resp.Status, nil
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(revoke_status)

	delete_req, err := http.NewRequest("DELETE", "https://"+*pe_console+":8140/puppet-ca/v1/certificate_status/"+*agent, nil)
	if err != nil {
		log.Fatal(err)
	}
	delete_req.Header.Set("Content-Type", "application/json")

	delete_resp, err := client.Do(delete_req)
	if err != nil {
		log.Fatal(err)
	}
	defer delete_resp.Body.Close()

	delete_status, err := delete_resp.Status, nil
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(delete_status)

	purge_req, err := http.NewRequest("POST", "https://"+*pe_console+":8081/pdb/cmd/v1", post_body)
	if err != nil {
		log.Fatal(err)
	}
	purge_req.Header.Set("Content-Type", "application/json")

	purge_resp, err := client.Do(purge_req)
	if err != nil {
		log.Fatal(err)
	}
	defer purge_resp.Body.Close()

	purge_status, err := purge_resp.Status, nil
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(purge_status)
}
