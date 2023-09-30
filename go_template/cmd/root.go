package cmd

import (
	"context"
	"os"

	// "github.com/TheFriendlyCoder/money/cmd/subcmd"
	// ao "github.com/TheFriendlyCoder/money/lib/applicationOptions"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CobraThemes mapping table that converts one of our theme types from
// a numeric identifier to a set of pre-defined Cobra color definitions
var theme = cc.Config{
	Headings: cc.HiCyan + cc.Bold + cc.Underline,
	Commands: cc.HiYellow + cc.Bold,
	Example:  cc.Italic,
	ExecName: cc.Bold,
	Flags:    cc.Bold,
}

// RootCmd definition for the main root / entry point command for the app
func RootCmd() cobra.Command {
	retval := cobra.Command{
		Use:   "{{app_name}}",
		Short: "{{description}}",
		// By default, cmd will always show the app usage message if the command
		// fails and returns an error. This flag disables that behavior.
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Sanity checks to make sure application is set up properly
		},
	}
	// Add subcommands here
	// retval.AddCommand(subcmd.CreateCmd())
	return retval
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command) error {

	// Initialize config file
	v := viper.GetViper()
	err := initViper(v)
	if err != nil {
		return err
	}

	// Setup application context
	ctx := context.Background()

	// Setup color theme
	cobraConfig := theme
	cobraConfig.RootCmd = cmd
	cobraConfig.NoExtraNewlines = true
	cc.Init(&cobraConfig)

	// Run our command
	return errors.WithStack(cmd.ExecuteContext(ctx))
}

// initViper initializes the Viper app configuration framework
func initViper(v *viper.Viper) error {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Critical application failure: user home folder not found")
	}

	// Search config in home directory
	v.AddConfigPath(home)
	v.SetConfigType("yaml")
	v.SetConfigName("{{config_file}}")

	// If a config file is found, read it in.
	err = v.ReadInConfig()

	// If there is no config file, we ignore that error and assume
	// there is no app config
	if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !os.IsNotExist(err) {
		return errors.WithStack(err)
	}
	return nil
}
