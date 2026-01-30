package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "frush",
	Short: "frush is a tool for quickly operating free-ran-ue and validating 5G core network, free5GC, behavior.",
	Long:  "frush is a tool for quickly operating free-ran-ue and validating 5G core network, free5GC, behavior.",
}

func ExecuteWithArgs(args []string) error {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpTemplate(`{{with .Long}}{{.}}{{end}}

Available Commands:
{{range .Commands}}{{if not .Hidden}}{{printf "    %s  \t%s\n" .Name .Short}}{{end}}{{end}}
Type "help [command]" for more information.
`)

	

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
