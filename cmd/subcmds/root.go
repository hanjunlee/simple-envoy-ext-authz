package subcmds

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ext-authz",
		Short: "It help to mocking a external authz server of Envoy.",
		Long: `You can easily mocking a external authz server. 
It helps to set the valid authorization token and headers with the response.`,
	}
)

// Execute execute the root command.
func Execute() error {
	return rootCmd.Execute()
}
