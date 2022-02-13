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
	"site24x7/cmd/impl/user"
	"site24x7/logger"

	"github.com/spf13/cobra"
)

// userCmd represents the `user` command
var userCmd = &cobra.Command{
	Use:   "user <command>",
	Short: "Performs user actions",
	Long:  `Performs user actions.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// authenticate before all non-config commands
		api.Authenticate()
	},
	// Run: func(cmd *cobra.Command, args []string) {
	//  NOOP - requires subcommand
	// 	fmt.Println("user command called; subcommand required")
	// },
}

// userCreateCmd represents the `user create` subcommand
var userCreateCmd = &cobra.Command{
	Use:   "create <email address>",
	Short: "Creates a new user",
	Long: `Creates a new user.

Valid roles: https://www.site24x7.com/help/api/#user_constants
Valid Status IQ roles: https://www.site24x7.com/help/api/#user_constants
Valid Cloudspend roles: https://www.site24x7.com/help/api/#user_constants
Valid job titles: https://www.site24x7.com/help/api/#job_title
User notification methods: https://www.site24x7.com/help/api/#alerting_constants
Valid email formats: https://www.site24x7.com/help/api/#alerting_constants
Valid resource types: https://www.site24x7.com/help/api/#resource_type_constants`,
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

		email := args[0]
		json, err := user.Create(email, cmd.Flags())
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

// userGetCmd represents the `user get` subcommand
var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves a specified user",
	Long: `Retrieves a specified user.

The Site24x7 API only supports retrieval by their ID, but this CLI will also
support retrieval by email address, albeit less efficient, for improved
usability.`,
	Aliases: []string{"fetch", "retrieve", "read"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		j, err := user.Get(cmd.Flags())
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

// userUpdateCmd represents the `user update` subcommand
var userUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing user",
	Long: `Updates an existing user.

Valid roles: https://www.site24x7.com/help/api/#user_constants
Valid Status IQ roles: https://www.site24x7.com/help/api/#user_constants
Valid Cloudspend roles: https://www.site24x7.com/help/api/#user_constants
Valid job titles: https://www.site24x7.com/help/api/#job_title
User notification methods: https://www.site24x7.com/help/api/#alerting_constants
Valid email formats: https://www.site24x7.com/help/api/#alerting_constants
Valid resource types: https://www.site24x7.com/help/api/#resource_type_constants`,
	Aliases: []string{"modify"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		json, err := user.Update(cmd.Flags())
		if err != nil {
			// Handle a known error just a bit more cleanly
			if err, ok := err.(*api.NotFoundError); ok {
				logger.Warn(err.Error())
				return nil
			}

			return err
		}

		logger.Out(string(json))

		return nil
	},
}

// userDeleteCmd represents the `user delete` subcommand
var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a specified user",
	Long: `Deletes a specified user.

The Site24x7 API only supports removal by user ID, but this CLI will also
support retrieval by email address, albeit less efficient, for improved
usability.`,
	Aliases: []string{"del", "rm", "remove"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		err := user.Delete(cmd.Flags())
		if err != nil {
			return err
		}

		logger.Out("User successfully deleted!")

		return nil
	},
}

// userGetCmd represents the `user list` subcommand
var userListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Retrieves a list of all users",
	Long:    `Retrieves a list of all users.`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		json, err := user.List()
		if err != nil {
			return err
		}

		logger.Out(string(json))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userUpdateCmd)
	userCmd.AddCommand(userDeleteCmd)
	userCmd.AddCommand(userListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Flags for the `user create` command
	// https://www.site24x7.com/help/api/#create-new-user
	userCreateCmd.Flags().AddFlagSet(user.GetWriterFlags())

	// Flags for the `user get` command
	// https://www.site24x7.com/help/api/#retrieve-user
	userGetCmd.Flags().AddFlagSet(user.GetAccessorFlags())

	// Flags for the `user update` command; updating a user requires us to
	// identify the user that will be updated and identify the data points that
	// will be updated
	// https://www.site24x7.com/help/api/#update-user
	userUpdateCmd.Flags().AddFlagSet(user.GetAccessorFlags())
	userUpdateCmd.Flags().AddFlagSet(user.GetWriterFlags())

	// Flags for the `user delete` command
	// https://www.site24x7.com/help/api/#delete-user
	userDeleteCmd.Flags().AddFlagSet(user.GetAccessorFlags())
}
