# PE CRUD Operations

CLI tool for testing Puppet Agent CRUD operations against a Puppet Enterprise installation.

## Setup

```
git clone https://github.com/coreymbe/pe-crud-ops
```

```
cd pe-crud-ops
```

```
go install
```

---

**Note**: To remove a Puppet agent node from the Puppet Enterprise installation you will need create a dummy certificate for authentication. Follow the steps below to create and validate the certificate.

On the primary server run:

```
/opt/puppetlabs/bin/puppetserver ca generate --certname pe-crud-ops-delete
```

You will want to copy the following files to a local directory:

  * `/etc/puppetlabs/puppet/ssl/certs/ca.pem`
  * `/etc/puppetlabs/puppet/ssl/certs/pe-crud-ops-delete.pem`
  * `/etc/puppetlabs/puppet/ssl/private_keys/pe-crud-ops-delete.pem`


The certificate needs to be added to the allowlist for both the Puppet CA and PuppetDB APIs. This can be done by adding the following parameters to the respective node groups.

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

Configure the following settings in the `.cobra.yaml` file before running any `pe-crud-ops` commands.

  * `$PE_CONSOLE`: FQDN of the PE Console.
  * `$PE_TOKEN`: PE RBAC token with admin permissions.
  *	`$PE_CACERT`: Path to a copy of the CA certificate.
  * `$PE_HOSTCERT`: Path to a copy of the host certificate.
  * `$PE_PRIVKEY`: Path to a copy of the host private key.

## Usage

### Create

```
pe-crud-ops create agent -a puppet-agent.example.com -s /path/to/ssh/key.pem
```

***Options***:

  * `--agent (-a)`: The FQDN of the node to become a Puppet agent.
  * `--ssh_creds (-s)`: Path to an SSH private key to access the node.
  * `--user (-u)`: __Optional__ - SSH user with sudo privilege.
  * `--pass (-p)`: __Optional__ - Sudo password for the configured SSH user.

> **Note**: You will want to remember the number of the JobID that is returned as this can be used to check the status of the task.

### Read

```
pe-crud-ops agent status --task_id 1
```

***Options***:

  * `--task_id (-i)`: The Task ID of the bootstrap task.

### Delete

```
pe-crud-ops delete agent --agent puppet-agent.example.com`
```

***Options***:

  * `--agent (-a)`: The certificate name of the Puppet agent node to remove.
