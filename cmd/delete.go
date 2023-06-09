package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove an agent node from a PE installation.",
	Long: `Remove an agent node from a PE installation.

	Example usage:
	  pe-crud-ops delete agent
		`,
	Run: func(cmd *cobra.Command, args []string) {

	},

	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println("error: Must use delete with subcommand agent")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteAgentCmd)

	deleteCmd.PersistentFlags().StringP("agent", "a", "", "FQDN of the new agent node.")
	deleteCmd.MarkPersistentFlagRequired("agent")
}
