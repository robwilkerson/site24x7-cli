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
	"site24x7/cmd/impl"

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
		var f impl.UserWriterFlags

		// Parse flag values
		f.Name, _ = cmd.Flags().GetString("name")
		f.Role, _ = cmd.Flags().GetInt("role")
		f.NotifyMethod, _ = cmd.Flags().GetIntSlice("notify-by")
		f.MonitorGroups, _ = cmd.Flags().GetStringSlice("monitor-groups")
		f.AlertEmailFormat, _ = cmd.Flags().GetInt("alert-email-format")
		f.AlertSkipDays, _ = cmd.Flags().GetIntSlice("alert-skip-days")
		f.AlertStartTime, _ = cmd.Flags().GetString("alert-start-time")
		f.AlertEndTime, _ = cmd.Flags().GetString("alert-end-time")
		f.AlertMethodsDown, _ = cmd.Flags().GetIntSlice("alert-methods-down")
		f.AlertMethodsTrouble, _ = cmd.Flags().GetIntSlice("alert-methods-trouble")
		f.AlertMethodsUp, _ = cmd.Flags().GetIntSlice("alert-methods-up")
		f.AlertMethodsAppLogs, _ = cmd.Flags().GetIntSlice("alert-methods-applogs")
		f.AlertMethodsAnomaly, _ = cmd.Flags().GetIntSlice("alert-methods-anomaly")
		f.JobTitle, _ = cmd.Flags().GetInt("job-title")
		f.MonitorGroups, _ = cmd.Flags().GetStringSlice("groups")
		f.MobileCountryCode, _ = cmd.Flags().GetString("mobile-country-code")
		f.MobileNumber, _ = cmd.Flags().GetString("mobile-phone-number")
		f.MobileSMSProviderID, _ = cmd.Flags().GetInt("mobile-sms-provider-id")
		f.MobileCallProviderID, _ = cmd.Flags().GetInt("mobile-sms-provider-id")
		f.NonEUAlertConsent, _ = cmd.Flags().GetBool("non-eu-alert-consent")
		f.ResourceType, _ = cmd.Flags().GetInt("resource-type")
		f.StatusIQRole, _ = cmd.Flags().GetInt("statusiq-role")
		f.CloudSpendRole, _ = cmd.Flags().GetInt("cloudspend-role")

		// Do all of the work in a testable custom function
		u := api.User{EmailAddress: args[0]}
		if err := impl.UserCreate(f, &u, u.Create); err != nil {
			return err
		}

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
		var f impl.UserAccessorFlags
		f.ID, _ = cmd.Flags().GetString("id")
		f.EmailAddress, _ = cmd.Flags().GetString("email")

		// Do all of the work in a testable custom function
		u := api.User{}
		json, err := impl.UserGet(f, &u, u.Get)
		if err != nil {
			return err
		}

		fmt.Println(string(json))

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
		var a impl.UserAccessorFlags
		var f impl.UserWriterFlags

		a.ID, _ = cmd.Flags().GetString("id")
		a.EmailAddress, _ = cmd.Flags().GetString("email")

		// Parse flag values
		f.Name, _ = cmd.Flags().GetString("name")
		// f.Role, _ = cmd.Flags().GetInt("role")
		// f.NotifyMethod, _ = cmd.Flags().GetIntSlice("notify-by")
		f.MonitorGroups, _ = cmd.Flags().GetStringSlice("monitor-groups")
		// f.AlertEmailFormat, _ = cmd.Flags().GetInt("alert-email-format")
		// f.AlertSkipDays, _ = cmd.Flags().GetIntSlice("alert-skip-days")
		// f.AlertStartTime, _ = cmd.Flags().GetString("alert-start-time")
		// f.AlertEndTime, _ = cmd.Flags().GetString("alert-end-time")
		// f.AlertMethodsDown, _ = cmd.Flags().GetIntSlice("alert-methods-down")
		// f.AlertMethodsTrouble, _ = cmd.Flags().GetIntSlice("alert-methods-trouble")
		// f.AlertMethodsUp, _ = cmd.Flags().GetIntSlice("alert-methods-up")
		// f.AlertMethodsAppLogs, _ = cmd.Flags().GetIntSlice("alert-methods-applogs")
		// f.AlertMethodsAnomaly, _ = cmd.Flags().GetIntSlice("alert-methods-anomaly")
		// f.JobTitle, _ = cmd.Flags().GetInt("job-title")
		// f.MonitorGroups, _ = cmd.Flags().GetStringSlice("groups")
		// f.MobileCountryCode, _ = cmd.Flags().GetString("mobile-country-code")
		// f.MobileNumber, _ = cmd.Flags().GetString("mobile-phone-number")
		// f.MobileSMSProviderID, _ = cmd.Flags().GetInt("mobile-sms-provider-id")
		// f.MobileCallProviderID, _ = cmd.Flags().GetInt("mobile-sms-provider-id")
		// f.NonEUAlertConsent, _ = cmd.Flags().GetBool("non-eu-alert-consent")
		// f.ResourceType, _ = cmd.Flags().GetInt("resource-type")
		// f.StatusIQRole, _ = cmd.Flags().GetInt("statusiq-role")
		// f.CloudSpendRole, _ = cmd.Flags().GetInt("cloudspend-role")

		// Do all of the work in a testable custom function
		u := api.User{}
		if err := impl.UserUpdate(a, f, &u, u.Update); err != nil {
			return err
		}

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
		var f impl.UserAccessorFlags
		f.ID, _ = cmd.Flags().GetString("id")
		f.EmailAddress, _ = cmd.Flags().GetString("email")

		// Do all of the work in a testable custom function
		u := api.User{}
		err := impl.UserDelete(f, &u, u.Delete)
		if err != nil {
			return err
		}

		fmt.Println("User successfully deleted!")

		return nil
	},
}

// userGetCmd represents the `user get` subcommand
var userListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Retrieves a list of all users",
	Long:    `Retrieves a list of all users.`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		json, err := impl.UserList(api.GetUsers)
		if err != nil {
			return err
		}

		fmt.Println(string(json))

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

	// 22 Flags for the `user create` command:
	// 	 ->  10 "simple flags
	//      + 9 alert settings
	//      + 4 mobile settings
	//      - 1 email address passed as an arg (no sensible default value)
	// https://www.site24x7.com/help/api/#create-new-user
	userCreateCmd.Flags().StringP("name", "n", "Unnamed User", "Full name (first last) of the user, e.g. \"Fred Flintstone\"")
	userCreateCmd.Flags().IntP("role", "r", 0, "Role assigned to the user for Site24x7 access")
	userCreateCmd.Flags().IntSliceP("notify-by", "N", []int{1}, "Medium by which the user will receive alerts")
	userCreateCmd.Flags().StringSliceP("monitor-groups", "g", nil, "List of monitor group identifiers to which the user should be assigned for receiving alerts")
	userCreateCmd.Flags().String("alert-email-format", "html", "Email format for alert emails")
	userCreateCmd.Flags().IntSlice("alert-skip-days", []int{}, "Days of the week on which the user should not be sent alerts: 0 (Sunday)-6 (Saturday) (default none")
	userCreateCmd.Flags().String("alert-start-time", "00:00", "The time of day when the user should start receiving alerts")
	userCreateCmd.Flags().String("alert-end-time", "00:00", "The time of day when the user should stop receiving alerts")
	userCreateCmd.Flags().IntSlice("alert-methods-down", []int{1}, "Preferred notification methods for down alerts")
	userCreateCmd.Flags().IntSlice("alert-methods-trouble", []int{1}, "Preferred notification methods for trouble alerts")
	userCreateCmd.Flags().IntSlice("alert-methods-up", []int{1}, "Preferred notification methods when service is restored")
	userCreateCmd.Flags().IntSlice("alert-methods-applogs", []int{1}, "Preferred notification methods for alerts related to application logs")
	userCreateCmd.Flags().IntSlice("alert-methods-anomaly", []int{1}, "Preferred notification methods for alerts when an anomaly is detected")
	userCreateCmd.Flags().Int("job-title", 0, "Job title of the user")
	userCreateCmd.Flags().Bool("non-eu-alert-consent", false, "Mandatory for EU DC; by passing true, you confirm your consent to transfer alert-related data")
	userCreateCmd.Flags().String("mobile-country-code", "", "Country code for mobile phone number; required if voice and/or sms notifications are requested")
	userCreateCmd.Flags().String("mobile-phone-number", "", "Digits only; required if voice and/or sms notifications are requested")
	userCreateCmd.Flags().Int("mobile-sms-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	userCreateCmd.Flags().Int("mobile-call-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	userCreateCmd.Flags().Int("resource-type", 0, "See https://www.site24x7.com/help/api/#resource_type_constants")
	userCreateCmd.Flags().Int("statusiq-role", 0, "Role assigned to the user for accessing StatusIQ")
	userCreateCmd.Flags().Int("cloudspend-role", 0, "Role assigned to the user for accessing CloudSpend")

	// Flags for the `user get` command
	userGetCmd.Flags().StringP("id", "i", "", "A user identifier")
	userGetCmd.Flags().StringP("email", "e", "", "A user email address")

	// Flags for the `user update` command
	// TODO: The Site24x7 looks to be pretty ham-fisted with this one based
	// solely on a look at the docs. Will have to do some testing flag by flag.
	// Flags specifically for targeting the user we want to update
	userUpdateCmd.Flags().StringP("id", "i", "", "A user identifier")
	userUpdateCmd.Flags().StringP("email", "e", "", "A user email address")
	userUpdateCmd.Flags().StringP("name", "n", "", "Full name (first last) of the user, e.g. \"Fred Flintstone\"")
	// Flags to set key data points
	// userUpdateCmd.Flags().IntP("role", "r", -1, "Role assigned to the user for Site24x7 access")
	// userUpdateCmd.Flags().IntSliceP("notify-by", "N", []int{}, "Medium by which the user will receive alerts")
	// userUpdateCmd.Flags().StringSliceP("monitor-groups", "g", []string{}, "List of monitor group identifiers to which the user should be assigned for receiving alerts")
	// userUpdateCmd.Flags().String("alert-email-format", "html", "Email format for alert emails")
	// userUpdateCmd.Flags().IntSlice("alert-skip-days", []int{}, "Days of the week on which the user should not be sent alerts: 0 (Sunday)-6 (Saturday) (default none")
	// userUpdateCmd.Flags().String("alert-start-time", "00:00", "The time of day when the user should start receiving alerts")
	// userUpdateCmd.Flags().String("alert-end-time", "00:00", "The time of day when the user should stop receiving alerts")
	// userUpdateCmd.Flags().IntSlice("alert-methods-down", []int{1}, "Preferred notification methods for down alerts")
	// userUpdateCmd.Flags().IntSlice("alert-methods-trouble", []int{1}, "Preferred notification methods for trouble alerts")
	// userUpdateCmd.Flags().IntSlice("alert-methods-up", []int{1}, "Preferred notification methods when service is restored")
	// userUpdateCmd.Flags().IntSlice("alert-methods-applogs", []int{1}, "Preferred notification methods for alerts related to application logs")
	// userUpdateCmd.Flags().IntSlice("alert-methods-anomaly", []int{1}, "Preferred notification methods for alerts when an anomaly is detected")
	// userUpdateCmd.Flags().Int("job-title", 0, "Job title of the user")
	// userUpdateCmd.Flags().Bool("non-eu-alert-consent", false, "Mandatory for EU DC; by passing true, you confirm your consent to transfer alert-related data")
	// userUpdateCmd.Flags().String("mobile-country-code", "", "Country code for mobile phone number; required if voice and/or sms notifications are requested")
	// userUpdateCmd.Flags().String("mobile-phone-number", "", "Digits only; required if voice and/or sms notifications are requested")
	// userUpdateCmd.Flags().Int("mobile-sms-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	// userUpdateCmd.Flags().Int("mobile-call-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	// userUpdateCmd.Flags().Int("resource-type", 0, "See https://www.site24x7.com/help/api/#resource_type_constants")
	// userUpdateCmd.Flags().Int("statusiq-role", 0, "Role assigned to the user for accessing StatusIQ")
	// userUpdateCmd.Flags().Int("cloudspend-role", 0, "Role assigned to the user for accessing CloudSpend")

	// Flags for the `user delete` command
	userDeleteCmd.Flags().StringP("id", "i", "", "A user identifier")
	userDeleteCmd.Flags().StringP("email", "e", "", "A user email address")
}
