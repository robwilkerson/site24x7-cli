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
	"os"
	"site24x7/api"

	"github.com/spf13/cobra"
)

// validateMuteDuration ensures an allowable duration value
func validateMuteDuration(d int) bool {
	validDurations := []int{5, 15, 30, 45, 60, 120, 180, 360, 720, 1440}

	for _, vd := range validDurations {
		if d == vd {
			return true
		}
	}

	return false
}

// validateMutedResources ensures that the passed options make sense.
//
// This logic represents my interpretation of the API documentation located at
// https://www.site24x7.com/help/api/#mute-monitor-alerts. As I understand it:
//
// * The "all" flag represents an "A" value for the category attribute. If
//   passed, both the monitor list and the group list must be empty.
// * If the monitor list is not empty, the group list must be empty and the
//   "all" flag must be false.
// * If the group list is not empty, the monitory list must be empty and the
//   "all" flag must be false.
func validateMutedResources(m []string, g []string, a bool) error {
	// We never want to surprise anyone by changing their input, so let's just
	// focus on communicating any problems
	if a && (len(m) > 0 || len(g) > 0) {
		return fmt.Errorf("If all monitors are to be muted, do not send --monitors or --groups flags\n")
	} else if len(m) > 0 && len(g) > 0 {
		return fmt.Errorf("Either the --monitors or --groups flag can be sent; not both\n")
	}

	return nil
}

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "alert <command>",
	Short: "Performs alert actions",
	Long:  `Performs alert actions.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// authenticate before all non-config commands
		api.Authenticate()
	},
	// Run: func(cmd *cobra.Command, args []string) {
	//  NOOP - requires subcommand
	// 	fmt.Println("alert command called; subcommand required")
	// },
}

// alertMuteCmd represents the alert mute subcommand
var alertMuteCmd = &cobra.Command{
	Use:   "mute",
	Short: "Suppresses one or more monitor alerts",
	Long: `Suppress monitor alerts for a particular monitor, monitor group, 
or all the resources for a specified time duration.`,
	Aliases: []string{"silence", "quiet", "suppress"},
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetInt("duration")
		m, _ := cmd.Flags().GetStringSlice("monitors")
		g, _ := cmd.Flags().GetStringSlice("groups")
		a, _ := cmd.Flags().GetBool("all")
		// r, _ := cmd.Flags().GetString("reason")
		// e, _ := cmd.Flags().GetBool("extend")
		// n, _ := cmd.Flags().GetBool("notify")

		if ok := validateMuteDuration(d); !ok {
			fmt.Printf("Invalid duration (%d)\n", d)
			os.Exit(1)
		}
		if err := validateMutedResources(m, g, a); err != nil {
			fmt.Println("Invalid resources selected for muting:")
			fmt.Printf(" * %s", err)
			os.Exit(1)
		}

		fmt.Println("Called to mute an alert")
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertMuteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	alertMuteCmd.Flags().IntP("duration", "d", 5, "Duration of the mute time in minutes")
	alertMuteCmd.Flags().StringSliceP("monitors", "m", []string{}, "A list of monitor IDs to be muted")
	alertMuteCmd.Flags().StringSliceP("groups", "g", []string{}, "A list of monitor groups to be muted")
	alertMuteCmd.Flags().StringP("reason", "r", "", "The reason for muting the alerts")
	alertMuteCmd.Flags().BoolP("all", "a", false, "Mutes all monitors")
	alertMuteCmd.Flags().BoolP("extend", "e", false, "Extend an existing mute time period")
	alertMuteCmd.Flags().BoolP("notify", "n", true, "Notify administrators that alarms are muted")
}
