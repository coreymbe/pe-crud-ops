package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd is a subcommand to StoreCmd that ads a Benchmark to the store.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a node to a PE installation.",
	Long: `Add a new puppet agent node to a PE installation.

	Example usage:
	  pe-crud-ops create agent
		`,
	Run: func(cmd *cobra.Command, args []string) {

	},

	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println("error: Must use create with subcommand agent")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createAgentCmd)

	createCmd.PersistentFlags().StringP("agent", "a", "", "FQDN of the new agent node.")
	createCmd.MarkPersistentFlagRequired("agent")
	createCmd.PersistentFlags().StringP("ssh_creds", "s", "", "Path to SSH private key to create the connection.")
	createCmd.MarkPersistentFlagRequired("ssh_creds")
	createCmd.PersistentFlags().StringP("user", "u", "root", "Optional, the SSH user for the connection.")
	createCmd.PersistentFlags().StringP("pass", "p", "", "Optional, sudo password for root escalation.")
}
