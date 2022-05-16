package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/coreymbe/pe-crud-ops/client"
)

func main() {

	// Gather Client Parameters
	pe_console := os.Getenv("PE_CONSOLE")
	token := os.Getenv("PE_TOKEN")
	// These params want filepaths
	ca_cert := os.Getenv("PE_CACERT")
	host_cert := os.Getenv("PE_HOSTCERT")
	host_key := os.Getenv("PE_PRIVKEY")

	// Build Client
	c, err := client.NewClient(pe_console, token, ca_cert, host_cert, host_key)
	if err != nil {
		log.Fatal(err)
	}

	// Operations
	create := flag.Bool("create", false, "Bootstrap a node to create a new Puppet agent")
	read := flag.Bool("read", false, "Check the status of the bootstrap task")
	delete := flag.Bool("delete", false, "Remove a Puppet agent node")

	// Create Parameters
	agent := flag.String("agent", "", "Required, the hostname of the agent node")
	ssh_creds := flag.String("ssh_creds", "", "Required, path to SSH private key to create the connection")
	ssh_user := flag.String("ssh_user", "root", "Optional, the SSH user for the connection")
	sudo_pass := flag.String("sudo_pass", "", "Optional, Sudo password for root escalation")

	// Read Parameters
	taskID := flag.String("taskID", "", "Required, the jobID of the bootstrap task.")

	help := flag.Bool("help", false, "Optional, prints usage info")
	flag.Parse()

	usage := `usage:

        pe-crud-ops [-help] [-create] [-read] [-delete] -agent <AGENT_HOSTNAME> -ssh_creds <PRIVATE_KEY_PATH> -taskID <jobID> [-ssh_user <SSH_USER> sudo_pass <SUDO_PASSWORD>]

        Options:
          -help        Prints this message.
          -create      Create a new Puppet agent node.
          -read        Check the status of the bootstrap task.
          -delete      Remove a Puppet agent node from a PDB and revoke the certificate.

        Create:
          -agent       Required, The hostname of the agent node to add.
          -ssh_creds   Required, The path to the SSH private key for the configured SSH user.
          -ssh_user    Optional, The SSH user to access the agent node (Default: "root").
          -sudo_pass   Optional, Password for sudo escalation is not using root user.

        Read:
          -taskID      Required, The jobID of the bootstrap task to check.

        Delete:
          -agent       Required, The hostname of the agent node to add.
`

	if *help {
		fmt.Println(usage)
		return
	}

	if *create {
		conn_cred, err := ioutil.ReadFile(*ssh_creds)
		if err != nil {
			log.Fatal(err)
		}
		pk := string(conn_cred)

		conn, err := c.CreateConn(*agent, ssh_user, pk, sudo_pass)
		if err != nil {
			log.Fatal(err)
		}

		bstp, err := c.Bootstrap(*agent, ssh_user)
		if err != nil {
			log.Fatal(err)
		}

		connID := conn.ConnectionID
		jobID := bstp.Job.ID

		response := fmt.Sprintf("ConnectionID: %s, JobID: %s", connID, jobID)
		fmt.Println(response)
	}

	if *read {
		for r := 0; r < 5; r++ {
			ts, err := c.TaskStatus(*taskID)
			if err != nil {
				log.Fatal(err)
			}
			if ts.State == "running" {
				log.Print("Task is currently running...")
				log.Print("...")
				time.Sleep(12 * time.Second)
				continue
			} else if ts.State == "finished" {
				response := fmt.Sprintf("Task_Created: %s, Task_Finished: %s, Agent_Install: %s", ts.Created, ts.Finished, ts.State)
				fmt.Println(response)
				break
			} else if ts.State == "failed" {
				log.Printf("Bootstrap Task: %s", ts.State)
				break
			}
		}
	}

	if *delete {
		err := c.PurgeNode(*agent)
		if err != nil {
			log.Fatal(err)
		}
	}
}
