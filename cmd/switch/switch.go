package _switch

import (
	"github.com/spf13/cobra"
)

// Cmd represents the add command
var Cmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch operation",
	Long:  `Switch a resource, to know more type maestro-cli switch --help.`,
}

func init() {
	Cmd.AddCommand(switchActiveVersionCmd)
}
