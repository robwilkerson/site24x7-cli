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

// validateUserGetInput validates user data passed to the get command
func validateUserGetInput(id string, email string) error {
	if id != "" && email != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if id == "" && email == "" {
		return fmt.Errorf("either an ID or an email address is required to retrieve a user")
	}

	return nil
}

// lookupIds checks a list of key values against a map and returns a slice
// containing each existing key's value
func lookupIds(list []string, lookup map[string]int) []int {
	var result []int

	for _, i := range list {
		if id, ok := lookup[strings.ToUpper(i)]; ok {
			result = append(result, id)
		}
	}

	return result
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
// TODO: Provide better links for CLI constants
// TODO: Add flags for remaining data points that can be sent
// TODO: Validate all flag values? Perhaps error if an invalid value was
// explicitly passed?
var userCreateCmd = &cobra.Command{
	Use:   "create <email address>",
	Short: "Creates a new user",
	Long: `Creates a new user.

Valid roles: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go
Valid Status IQ roles: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go
Valid Cloudspend roles: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go
Valid job titles: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go
User notification methods: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go
Valid email formats: https://github.com/robwilkerson/site24x7-cli/blob/main/api/user.go`,
	Aliases: []string{"add", "new"},
	Args: func(cmd *cobra.Command, args []string) error {
		expectedArgLen := 1
		actualArgLen := len(args)
		if actualArgLen != expectedArgLen {
			return fmt.Errorf("expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var preferredEmailFormat int

		// Parse flag values
		role, _ := cmd.Flags().GetString("role")
		statusIQRole, _ := cmd.Flags().GetString("statusiq-role")
		cloudspendRole, _ := cmd.Flags().GetString("cloudspend-role")
		notifyBy, _ := cmd.Flags().GetStringSlice("notify-by")
		alertEmailFormat, _ := cmd.Flags().GetString("alert-email-format")
		alertSkipDays, _ := cmd.Flags().GetIntSlice("alert-skip-days")
		alertStartTime, _ := cmd.Flags().GetString("alert-start-time")
		alertEndTime, _ := cmd.Flags().GetString("alert-end-time")
		alertMethodsDown, _ := cmd.Flags().GetStringSlice("alert-methods-down")
		alertMethodsTrouble, _ := cmd.Flags().GetStringSlice("alert-methods-trouble")
		alertMethodsUp, _ := cmd.Flags().GetStringSlice("alert-methods-up")
		alertMethodsAppLogs, _ := cmd.Flags().GetStringSlice("alert-methods-applogs")
		alertMethodsAnomaly, _ := cmd.Flags().GetStringSlice("alert-methods-anomaly")
		jobTitle, _ := cmd.Flags().GetString("job-title")

		// Hydrate a user object with passed values or reasonable defaults
		// Required inputs
		u := api.User{EmailAddress: args[0]}
		u.Name, _ = cmd.Flags().GetString("name")
		if v, ok := api.UserRoles[strings.ToUpper(role)]; ok {
			u.Role = v
		} else {
			// If the user explicitly passes an invalid role, error
			fmt.Printf("ERROR: Invalid role (%s)\n", role)
			cmd.Help()
			os.Exit(1)
		}

		// Optional inputs
		// These all have sensible defaults; if the user explicitly passes an
		// invalid value, we'll just let the API tell us we're wrong so we can
		// keep the code itself reasonably tidy. There's a lot of values here
		// so checking them all is a bit of a chore. May revisit this later.
		if u.NotifyMedium = lookupIds(notifyBy, api.UserNotifyMediums); u.NotifyMedium == nil {
			// At least 1 value was passed, but none were valid
			fmt.Println("ERROR: At least one notification method was passed, but none were valid")
			cmd.Help()
			os.Exit(1)
		}
		if statusIQRole != "" {
			// If a value was explicitly passed, error if it doesn't exist
			if v, ok := api.UserStatusIQRoles[strings.ToUpper(statusIQRole)]; ok {
				u.StatusIQRole = v
			}
		}
		if cloudspendRole != "" {
			// If a value was explicitly passed, error if it doesn't exist
			if v, ok := api.UserCloudspendRoles[strings.ToUpper(cloudspendRole)]; ok {
				u.CloudspendRole = v
			}
		}
		if v, ok := api.UserEmailFormats[alertEmailFormat]; ok {
			preferredEmailFormat = v
		}
		u.AlertSettings = map[string]interface{}{
			"email_format":       preferredEmailFormat,
			"dont_alert_on_days": alertSkipDays,
			"alerting_period": map[string]string{
				"start_time": alertStartTime,
				"end_time":   alertEndTime,
			},
			"down":    lookupIds(alertMethodsDown, api.UserNotifyMediums),
			"trouble": lookupIds(alertMethodsTrouble, api.UserNotifyMediums),
			"up":      lookupIds(alertMethodsUp, api.UserNotifyMediums),
			"applogs": lookupIds(alertMethodsAppLogs, api.UserNotifyMediums),
			"anomaly": lookupIds(alertMethodsAnomaly, api.UserNotifyMediums),
		}
		if v, ok := api.UserJobTitles[strings.ToLower(jobTitle)]; ok {
			u.JobTitle = v
		}

		if err := u.Create(); err != nil {
			fmt.Println(err)
		}
	},
}

// userGetCmd represents the user get subcommand
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

	// Flags for the create command
	userCreateCmd.Flags().StringP("name", "n", "Unnamed User", "Full name (first last) of the user, e.g. \"Fred Flintstone\"")
	userCreateCmd.Flags().StringP("role", "r", "noaccess", "Role assigned to the user for Site24x7 access (--help for valid options)")
	userCreateCmd.Flags().StringSliceP("notify-by", "N", []string{"email"}, "Medium by which the user will receive alerts")
	userCreateCmd.Flags().String("statusiq-role", "", "Role assigned to the user for accessing StatusIQ")
	userCreateCmd.Flags().String("cloudspend-role", "", "Role assigned to the user for accessing CloudSpend")
	userCreateCmd.Flags().StringSliceP("groups", "g", []string{}, "List of group identifiers to which the user should be assigned for receiving alerts")
	userCreateCmd.Flags().String("alert-email-format", "html", "Email format for alert emails")
	userCreateCmd.Flags().IntSlice("alert-skip-days", []int{}, "Days of the week on which the user should not be sent alerts: 0 (Sunday)-6 (Saturday) (default none")
	userCreateCmd.Flags().String("alert-start-time", "00:00", "The time of day when the user should start receiving alerts")
	userCreateCmd.Flags().String("alert-end-time", "00:00", "The time of day when the user should stop receiving alerts")
	userCreateCmd.Flags().StringSlice("alert-methods-down", []string{"email"}, "Preferred notification methods for down alerts")
	userCreateCmd.Flags().StringSlice("alert-methods-trouble", []string{"email"}, "Preferred notification methods for trouble alerts")
	userCreateCmd.Flags().StringSlice("alert-methods-up", []string{"email"}, "Preferred notification methods when service is restored")
	userCreateCmd.Flags().StringSlice("alert-methods-applogs", []string{"email"}, "Preferred notification methods for alerts related to application logs")
	userCreateCmd.Flags().StringSlice("alert-methods-anomaly", []string{"email"}, "Preferred notification methods for alerts when an anomaly is detected")
	userCreateCmd.Flags().String("job-title", "it", "Job title of the user")
	// TODO: Add flags for additional data points

	// Flags for the get command
	userGetCmd.Flags().StringP("id", "i", "", "A user identifier")
	userGetCmd.Flags().StringP("email", "e", "", "A user email address")
}
