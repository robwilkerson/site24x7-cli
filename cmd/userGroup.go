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
	"site24x7/cmd/impl/usergroup"
	"site24x7/logger"

	"github.com/spf13/cobra"
)

// userGroupCmd represents the `user_group` command
var userGroupCmd = &cobra.Command{
	Use:   "user_group <command>",
	Short: "Performs user group actions",
	Long: `Performs user group actions.
	
https://www.site24x7.com/help/api/#user-groups`,
	Aliases: []string{"ug", "usergroup", "ugroup", "usergru"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// authenticate before all non-config commands
		api.Authenticate()
		// set the log verbosity for any monitor_group command execution
		logger.SetVerbosity(cmd.Flags())
	},
	// Run: func(cmd *cobra.Command, args []string) {
	//  NOOP - requires subcommand
	// 	fmt.Println("user_group command called; subcommand required")
	// },
}

// userGroupCreateCmd represents the `user_group create` command
var userGroupCreateCmd = &cobra.Command{
	Use:   "create <display name>",
	Short: "Creates a new user group",
	Long: `Creates a new user group.

https://www.site24x7.com/help/api/#create-user-group`,
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
		name := args[0]
		json, err := usergroup.Create(name, cmd.Flags())
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

// userGroupGetCmd represents the `user_group get` subcommand
var userGroupGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Retrieves a specific user group",
	Long: `Retrieves a specific user group.

https://www.site24x7.com/help/api/#retrieve-user-group`,
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
		id := args[0]
		j, err := usergroup.Get(id)
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

// userGroupUpdateCmd represents the `user_group update` subcommand
// var userGroupUpdateCmd = &cobra.Command{
// 	Use:   "update <id>",
// 	Short: "Updates an existing user group",
// 	Long: `Updates an existing user group.

// https://www.site24x7.com/help/api/#update-user-group`,
// 	Aliases: []string{"modify"},
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		id := args[0]
// 		json, err := usergroup.Update(id, cmd.Flags())
// 		if err != nil {
// 			// Handle a known error just a bit more cleanly
// 			if err, ok := err.(*api.NotFoundError); ok {
// 				logger.Warn(err.Error())
// 				return nil
// 			}

// 			return err
// 		}

// 		logger.Out(string(json))

// 		return nil
// 	},
// }

// userGroupDeleteCmd represents the `user_group delete` subcommand
var userGroupDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Deletes a specific user group",
	Long: `Deletes a specific user group.

https://www.site24x7.com/help/api/#delete-user-group`,
	Aliases: []string{"del", "rm", "remove"},
	Args: func(cmd *cobra.Command, args []string) error {
		expectedArgLen := 1
		actualArgLen := len(args)
		if actualArgLen != expectedArgLen {
			return fmt.Errorf("expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		err := usergroup.Delete(id)
		if err != nil {
			return err
		}

		logger.Out("User group successfully deleted!")

		return nil
	},
}

// userGroupListCmd represents the `user_group list` subcommand
var userGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieves a list of all user groups",
	Long: `Retrieves a list of all user groups.
	
https://www.site24x7.com/help/api/#list-of-all-user-groups`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		json, err := usergroup.List()
		if err != nil {
			return err
		}

		logger.Out(string(json))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userGroupCmd)
	userGroupCmd.AddCommand(userGroupCreateCmd)
	userGroupCmd.AddCommand(userGroupGetCmd)
	// userGroupCmd.AddCommand(userGroupUpdateCmd)
	userGroupCmd.AddCommand(userGroupDeleteCmd)
	userGroupCmd.AddCommand(userGroupListCmd)

	// Flags for the `user_group create` command
	userGroupCreateCmd.Flags().AddFlagSet(usergroup.GetWriterFlags())
	userGroupCreateCmd.MarkFlagRequired("users")

	// Flags for the `user_group update` command
	// userGroupUpdateCmd.Flags().AddFlagSet(usergroup.GetWriterFlags())
}
