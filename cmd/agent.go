package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/coreymbe/pe-crud-ops/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Check the agent status.",
	Long: `Check the installation status of the new agent node.

	Example usage:
	  pe-crud-ops agent status
		`,
	Run: func(cmd *cobra.Command, args []string) {

	},

	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println("error: Must use agent with subcommand status")

		return nil
	},
}

var createAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Add a node to a PE installation.",
	Long: `Add a new puppet agent node to a PE installation.

	Example usage:
	pe-crud-ops create agent --agent puppet.agent.com --ssh_creds /path/to/ssh/key.pem
		`,
	Run: func(cmd *cobra.Command, args []string) {
		pe_console := viper.GetString("PE.Console")
		token := viper.GetString("PE.Token")
		ca_cert := viper.GetString("PE.CACert")
		host_cert := viper.GetString("PE.HostCert")
		host_key := viper.GetString("PE.PrivKey")
		ssh_creds, _ := cmd.Flags().GetString("ssh_creds")
		agent, _ := cmd.Flags().GetString("agent")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")

		c, err := client.NewClient(pe_console, token, ca_cert, host_cert, host_key)
		if err != nil {
			log.Fatal(err)
		}

		conn_cred, err := ioutil.ReadFile(ssh_creds)
		if err != nil {
			log.Fatal(err)
		}
		pk := string(conn_cred)

		bstp, err := c.Bootstrap(agent, user, pass, pk)
		if err != nil {
			log.Fatal(err)
		}

		jobID := bstp.Job.Name

		response := fmt.Sprintf("Task ID: %s", jobID)
		fmt.Println(response)
	},
}

var deleteAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Remove an agent node from a PE installation.",
	Long: `Remove the agent from PuppetDB and the certificate
	from the Puppet CA.

	Example usage:
	pe-crud-ops delete agent --agent puppet.agent.com
		`,
	Run: func(cmd *cobra.Command, args []string) {
		pe_console := viper.GetString("PE.Console")
		token := viper.GetString("PE.Token")
		ca_cert := viper.GetString("PE.CACert")
		host_cert := viper.GetString("PE.HostCert")
		host_key := viper.GetString("PE.PrivKey")
		agent, _ := cmd.Flags().GetString("agent")

		c, err := client.NewClient(pe_console, token, ca_cert, host_cert, host_key)
		if err != nil {
			log.Fatal(err)
		}

		purge_err := c.PurgeNode(agent)
		if purge_err != nil {
			log.Fatal(purge_err)
		} else {
			response := fmt.Sprintf("Agent node '%s' removed from %s.", agent, pe_console)
			fmt.Println(response)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(statusCmd)
}
