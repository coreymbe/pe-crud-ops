package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/coreymbe/pe-crud-ops/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the agent status.",
	Long: `Check the installation status of the new agent node.

	Example usage:
	  pe-crud-ops agent status --task_id 123
		`,
	Run: func(cmd *cobra.Command, args []string) {
		pe_console := viper.GetString("PE.Console")
		token := viper.GetString("PE.Token")
		ca_cert := viper.GetString("PE.CACert")
		host_cert := viper.GetString("PE.HostCert")
		host_key := viper.GetString("PE.PrivKey")
		taskID, _ := cmd.Flags().GetString("task_id")

		c, err := client.NewClient(pe_console, token, ca_cert, host_cert, host_key)
		if err != nil {
			log.Fatal(err)
		}

		for r := 0; r < 5; r++ {
			ts, err := c.TaskStatus(taskID)
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
	},
}

func init() {
	statusCmd.Flags().StringP("task_id", "i", "", "The jobID of the bootstrap task to check.")
	statusCmd.MarkFlagRequired("task_id")
}
