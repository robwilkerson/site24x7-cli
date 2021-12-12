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
	"encoding/json"
	"fmt"
	"os"
	"site24x7/api"
	"strings"

	"github.com/spf13/cobra"
)

func validateUserGetInput(id string, email string) error {
	if id != "" && email != "" {
		return fmt.Errorf("Please include either an ID OR an email address, not both")
	} else if id == "" && email == "" {
		return fmt.Errorf("Either an ID or an email address is required to retrieve a user")
	}

	return nil
}

// userCmd represents the user command
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

// userCreateCmd represents the user create subcommand
var userCreateCmd = &cobra.Command{
	Use:     "create <email address>",
	Short:   "Creates a new user",
	Long:    `Creates a new user.`,
	Aliases: []string{"add", "new"},
	Args: func(cmd *cobra.Command, args []string) error {
		expectedArgLen := 1
		actualArgLen := len(args)
		if actualArgLen != expectedArgLen {
			return fmt.Errorf("Expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		role, _ := cmd.Flags().GetString("role")
		notifyBy, _ := cmd.Flags().GetStringSlice("notify-by")

		u := api.User{EmailAddress: args[0]}
		u.Name, _ = cmd.Flags().GetString("name")
		u.Role, _ = api.UserRoles[role]
		// map notification channels inputs to their Site24x7 ids
		for _, m := range notifyBy {
			id, ok := api.UserNotifyMediums[strings.ToUpper(m)]
			if ok {
				u.NotifyMedium = append(u.NotifyMedium, id)
			}
		}

		if err := u.Create(); err != nil {
			fmt.Println(err)
		}
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves a specified user",
	Long: `Retrieves a specified user.

The Site24x7 API only supports retrieval by their ID, but this CLI will also
support retrieval by email address, albeit less efficient, for improved
usability.`,
	Aliases: []string{"fetch", "retrieve"},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		email, _ := cmd.Flags().GetString("email")

		if err := validateUserGetInput(id, email); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		u := api.User{Id: id, EmailAddress: email}
		if err := u.Get(); err != nil {
			fmt.Println(err)
		}

		out, _ := json.MarshalIndent(u, "", "    ")
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Mandatory values with sensible defaults
	userCreateCmd.Flags().StringP("name", "n", "Unnamed User", "Full name (first last) of the user, e.g. \"Fred Flintstone\"")
	userCreateCmd.Flags().StringP("role", "r", "NoAccess", "Role assigned to the user for Site24x7 access")
	userCreateCmd.Flags().StringSlice("notify-by", []string{"email"}, "Medium by which the user will receive alerts")

	// Optional values; we'll start simple and not support these
	// createCmd.Flags().String("statusiq-role", "", "Role assigned to the user for accessing StatusIQ")
	// createCmd.Flags().String("cloudspend-role", "", "Role assigned to the user for accessing CloudSpend")
	// createCmd.Flags().String("alert-email-format", "HTML", "Email format for alert emails")
	// createCmd.Flags().String("alert-skip-days", "[]", "Days of the week on which the user should not be sent alerts - 0 (Sunday) - 7 (Saturday)")
	// createCmd.Flags().String("alert-skip-days", "[]", "Days of the week on which the user should not be sent alerts")
	// createCmd.Flags().StringP("job-title", "t", "", "Role assigned to the user for accessing CloudSpend")

	userGetCmd.Flags().StringP("id", "i", "", "The Site24x7 user identifier")
	userGetCmd.Flags().StringP("email", "e", "", "A user email address")
}
