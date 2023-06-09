package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Purge Node Payload Struct
type CertState struct {
	DesiredState string `json:"desired_state"`
}
type Payload struct {
	Certname string `json:"certname"`
}
type PurgeNode struct {
	Command string  `json:"command"`
	Version int     `json:"version"`
	Payload Payload `json:"payload"`
}

// Bootstrap Payload Struct
type Targets struct {
	Transport         string   `json:"transport"`
	RunAs             string   `json:"run-as"`
	User              string   `json:"user"`
	PrivateKeyContent string   `json:"private-key-content"`
	SudoPassword      string   `json:"sudo-password"`
	Hostnames         []string `json:"hostnames"`
}
type Scope struct {
	Nodes []string `json:"nodes"`
}
type Params struct {
	Server string `json:"server"`
}
type TaskPayload struct {
	Task    string    `json:"task"`
	Params  Params    `json:"params"`
	Scope   Scope     `json:"scope"`
	Targets []Targets `json:"targets"`
}

// Connection Payload Struct
type SensitiveParameters struct {
	PrivateKeyContent string `json:"private-key-content"`
	SudoPassword      string `json:"sudo-password"`
}
type Parameters struct {
	User  string `json:"user"`
	RunAs string `json:"run-as"`
}
type ConnPayload struct {
	Certnames           []string            `json:"certnames"`
	Type                string              `json:"type"`
	Parameters          Parameters          `json:"parameters"`
	SensitiveParameters SensitiveParameters `json:"sensitive_parameters"`
	Duplicates          string              `json:"duplicates"`
}

// Bootstrap Struct
type Job struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Bootstrap struct {
	Job Job `json:"job"`
}

// CreateConn Struct
type Connection struct {
	ConnectionID string `json:"connection_id"`
}

// TaskStatus Struct
type Status struct {
	ID       string `json:"id"`
	Created  string `json:"created_timestamp"`
	State    string `json:"state"`
	Finished string `json:"finished_timestamp"`
}

// TaskStatus - Return information on the status of the bootstrap task
func (c *Client) TaskStatus(jobID string) (*Status, error) {
	req, err := http.NewRequest("GET", "https://"+c.PE_Console+":8143/orchestrator/v1/jobs/"+jobID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var task Status
	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// CreateConn - Create new connection
func (c *Client) CreateConn(agent string, transport_user string, transport_cred string, sudo_pass string) (*Connection, error) {
	conn_data := ConnPayload{
		Certnames: []string{agent},
		Type:      "ssh",
		Parameters: Parameters{
			User:  transport_user,
			RunAs: "root",
		},
		SensitiveParameters: SensitiveParameters{
			PrivateKeyContent: transport_cred,
			SudoPassword:      sudo_pass,
		},
		Duplicates: "replace",
	}

	connBytes, err := json.Marshal(conn_data)
	if err != nil {
		return nil, err
	}
	req_body := bytes.NewReader(connBytes)

	req, err := http.NewRequest("POST", "https://"+c.PE_Console+":8143/inventory/v1/command/create-connection", req_body)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var conn Connection
	err = json.Unmarshal(body, &conn)
	if err != nil {
		return nil, err
	}

	return &conn, err
}

// Bootstrap - Run the pe_bootstrap task
func (c *Client) Bootstrap(agent string, ssh_user string, sudo_pass string, transport_cred string) (*Bootstrap, error) {
	task_data := TaskPayload{
		Task: "pe_bootstrap",
		Params: Params{
			Server: c.PE_Console,
		},
		Scope: Scope{
			Nodes: []string{agent},
		},
		Targets: []Targets{
			{
				Transport:         "ssh",
				RunAs:             "root",
				User:              ssh_user,
				SudoPassword:      sudo_pass,
				PrivateKeyContent: transport_cred,
				Hostnames:         []string{agent},
			},
		},
	}

	taskBytes, err := json.Marshal(task_data)
	if err != nil {
		return nil, err
	}

	req_body := bytes.NewReader(taskBytes)

	req, err := http.NewRequest("POST", "https://"+c.PE_Console+":8143/orchestrator/v1/command/task", req_body)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var j Bootstrap
	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// PurgeNode - Remove the agent from the Puppet CA and PuppetDB
// This could probably be broken up into 3 seperate function methods.
func (c *Client) PurgeNode(agent string) error {
	revoke_data := CertState{
		DesiredState: "revoked",
	}
	revoke_payload, err := json.Marshal(revoke_data)
	if err != nil {
		return err
	}
	put_body := bytes.NewReader(revoke_payload)

	revoke_req, err := http.NewRequest("PUT", "https://"+c.PE_Console+":8140/puppet-ca/v1/certificate_status/"+agent, put_body)
	if err != nil {
		return err
	}
	revoke_err := c.purgeRequest(revoke_req)
	if revoke_err != nil {
		return revoke_err
	}

	delete_req, err := http.NewRequest("DELETE", "https://"+c.PE_Console+":8140/puppet-ca/v1/certificate_status/"+agent, nil)
	if err != nil {
		return err
	}
	delete_err := c.purgeRequest(delete_req)
	if delete_err != nil {
		return delete_err
	}

	purge_data := PurgeNode{
		Command: "deactivate node",
		Version: 3,
		Payload: Payload{
			Certname: agent,
		},
	}
	purge_payload, err := json.Marshal(purge_data)
	if err != nil {
		return err
	}
	post_body := bytes.NewReader(purge_payload)

	purge_req, err := http.NewRequest("POST", "https://"+c.PE_Console+":8081/pdb/cmd/v1", post_body)
	if err != nil {
		return err
	}
	purge_req.Header.Set("Content-Type", "application/json")
	purge_err := c.purgeRequest(purge_req)
	if purge_err != nil {
		return purge_err
	}

	return nil
}
