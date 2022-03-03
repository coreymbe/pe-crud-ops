# PE CRUD Operations

A series of `go` scripts to manually test CRUD operations for a Puppet Agent node against a Puppet Enterprise installation.

## Setup

### Create

  * `create_conn.go`
  * `pe_bootstrap.go`

To create a new Puppet agent node you will need have the following information before running the scripts listed above:

  * FQDN of the PE Console.
  * FQDN of the node that will become a Puppet agent node.
  * SSH access to the node via a private key.
  * PE RBAC token with admin permissions.

***Note***: You will want to remember the Job ID that is returned as this will be used to check the status of the task.

---

### Read

  * `check_status.go`

To check the status of the bootstrap task you will need to have the following information before running the script listed above:

  * The Job ID of the bootstrap task.
  * PE RBAC token with admin permissions.

---

### Delete

  * `purge_node.go`

To remove a Puppet agent node from the Puppet installation you will need have the following information before running the script listed above:

  * FQDN of the PE Console.
  * FQDN of the node that will become a Puppet agent node.
  * On the primary server run: `/opt/puppetlabs/bin/puppetserver ca generate --certname pe-crud-ops-delete`

The certificate needs to be added to the allowlist for both the Puppet CA and PupppetDB APIs. This can be done by adding the following parameters to the respective node groups.

**PE Infrastructure** -> **PE Certificate Authority** -> **Configuration data**:

```
Class: puppet_enterprise::profile::certificate_authority
Parameter: client_allowlist
Value: ["pe-crud-ops-delete"]
```

**PE Infrastructure** -> **PE PuppetDB** -> **Classes**:

```
Class: puppet_enterprise::profile::puppetdb
Parameter: allowlisted_certnames
Value: ["pe-crud-ops-delete"]
```

***Note***: As this request requires certificate authentication it is currently meant to be ran from the PE Primary Server to access the dummy cert you created.

## Usage

**Create Connection**:

  * `go run create_conn.go -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME> -token <PE_TOKEN> -ssh_key <PRIVATE_KEY_PATH> [-ssh_user <SSH_USER>]`

**Bootstrap Agent**:

  * `go run pe_boostrap.go -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME> -token <PE_TOKEN>`

**Check Status**:

  * `go run check_status.go -jobID <JOB_ID> -token <PE_TOKEN>`

**Purge Node**:

  * `go run purge_node.go -pe_console <PE_CONSOLE_HOSTNAME> -agent <AGENT_HOSTNAME>`

---
