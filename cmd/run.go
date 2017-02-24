package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/drud/dcfg/pkg/dcfg"
	"github.com/spf13/cobra"
)

func inArgs(needle string, args []string) bool {
	for _, arg := range args {
		if arg == needle {
			return true
		}
	}
	return false
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [group_name|all]",
	Short: "run sets of commands from a yaml file",
	Long: `DrudConfig uses an Ansible-like yaml syntax to run batches of commands
	with some additional handy funcitonality.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return fmt.Errorf("File does not exist: %s", cfgFile)
		}

		if len(args) == 0 {
			return fmt.Errorf("You must provide a config group or 'all'.")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running", args, "...")

		fileBytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.Fatalln("Could not read config file:", err)
		}

		groups, err := dcfg.GetTaskSetList(fileBytes)
		if err != nil {
			log.Fatalln(err)
		}

		for _, group := range groups {
			if args[0] == "all" || inArgs(group.Name, args) {
				group.Run()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

}
