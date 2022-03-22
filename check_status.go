package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	help := flag.Bool("help", false, "Optional, prints usage info")
	jobID := flag.String("jobID", "", "Required, the JobID of the bootstrap task")
	token := flag.String("token", "", "Required, the RBAC token to create the connection")
	flag.Parse()

	usage := `usage:
	
check_status -jobID <JOB_ID> -token <PE_TOKEN> [-help]
	
Options:
	-help		Optional, Prints this message.
	-jobID	Required, The Job ID of the boostrap task to check.
	-token	Required, PE RBAC token with appropriate permissions.
 `

	if *help == true {
		fmt.Println(usage)
		return
	}

	if (*jobID == "") || (*token == "") {
		log.Fatalf("The jobID and token options are required:\n%s", usage)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", *jobID, nil)
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
