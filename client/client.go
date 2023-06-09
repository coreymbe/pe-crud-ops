package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client Struct
type Client struct {
	Client     *http.Client
	PE_Console string
	Token      string
	CA_Cert    string
	Host_Cert  string
	Host_Key   string
}

// NewClient
func NewClient(pe_console, pe_token, ca_cert, host_cert, host_key string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := Client{
		Client:     &http.Client{Transport: tr},
		PE_Console: pe_console,
		Token:      pe_token,
		CA_Cert:    ca_cert,
		Host_Cert:  host_cert,
		Host_Key:   host_key,
	}
	return &c, nil
}

// Create Request
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authentication", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Need to get this logic properly implemented.
	// CreateConn returns 201 and Bootstrap returns 202.
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
	}

	return body, err
}

// Delete Request
func (c *Client) purgeRequest(req *http.Request) error {
	//host_cert, err := os.ReadFile(c.Host_Cert)
	//if err != nil {
	//	return err
	//}
	//host_key, err := os.ReadFile(c.Host_Key)
	//if err != nil {
	//	return err
	//}
	cert, err := tls.LoadX509KeyPair(c.Host_Cert, c.Host_Key)
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile(c.CA_Cert)
	if err != nil {
		return err
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

	c.Client = &http.Client{Transport: tr}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if isValid(resp.StatusCode) {
		return nil
	} else {
		return fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
	}
}
