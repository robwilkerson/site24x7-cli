/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"site24x7/api"
	"site24x7/cmd/impl/monitorgroup"
	"site24x7/logger"

	"github.com/spf13/cobra"
)

// monitorGroupCmd represents the monitorGroup command
var monitorGroupCmd = &cobra.Command{
	Use:     "monitor_group <command>",
	Short:   "Performs monitor group actions",
	Long:    `Performs monitor group actions.`,
	Aliases: []string{"mg", "mongroup", "mgroup", "mongru"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// authenticate before all non-config commands
		api.Authenticate()
	},
	// Run: func(cmd *cobra.Command, args []string) {
	//  NOOP - requires subcommand
	// 	fmt.Println("monitor_group command called; subcommand required")
	// },
}

var monitorGroupCreateCmd = &cobra.Command{
	Use:     "create <display name>",
	Short:   "Creates a new monitor group",
	Long:    `Creates a new monitor group.`,
	Aliases: []string{"add", "new"},
	Args: func(cmd *cobra.Command, args []string) error {
		expectedArgLen := 1
		actualArgLen := len(args)
		if actualArgLen != expectedArgLen {
			return fmt.Errorf("expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		name := args[0]
		json, err := monitorgroup.Create(name, cmd.Flags())
		if err != nil {
			// Handle a user already exists error nicely
			if err, ok := err.(*api.ConflictError); ok {
				logger.Warn(err.Error())
				return nil
			}

			return err
		}

		logger.Out(string(json))

		return nil
	},
}

// userGetCmd represents the `monitor_group get` subcommand
var monitorGroupGetCmd = &cobra.Command{
	Use:     "get <id>",
	Short:   "Retrieves a specific monitor group",
	Long:    `Retrieves a specific monitor group.`,
	Aliases: []string{"fetch", "retrieve", "read"},
	Args: func(cmd *cobra.Command, args []string) error {
		expectedArgLen := 1
		actualArgLen := len(args)
		if actualArgLen != expectedArgLen {
			return fmt.Errorf("expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		id := args[0]
		j, err := monitorgroup.Get(id)
		if err != nil {
			if err, ok := err.(*api.NotFoundError); ok {
				logger.Warn(err.Error())
				return nil
			}

			return err
		}

		logger.Out(string(j))

		return nil
	},
}

// monitorGroupListCmd represents the `monitor_group list` subcommand
var monitorGroupListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Retrieves a list of all monitor groups",
	Long:    `Retrieves a list of all monitor groups.`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		json, err := monitorgroup.List()
		if err != nil {
			return err
		}

		logger.Out(string(json))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(monitorGroupCmd)
	monitorGroupCmd.AddCommand(monitorGroupCreateCmd)
	monitorGroupCmd.AddCommand(monitorGroupListCmd)

	// Flags for the `user create` command
	// https://www.site24x7.com/help/api/#create-new-user
	monitorGroupCreateCmd.Flags().AddFlagSet(monitorgroup.GetWriterFlags())
}
