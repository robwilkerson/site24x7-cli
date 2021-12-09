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

	"github.com/spf13/cobra"
)

// monitorGroupCmd represents the monitorGroup command
var monitorGroupCmd = &cobra.Command{
	Use:     "monitor_group <command>",
	Short:   "Performs monitor group actions",
	Long:    `Performs monitor group actions.`,
	Aliases: []string{"mg", "mgroup"},
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
			return fmt.Errorf("Expected %d arguments, received %d", expectedArgLen, actualArgLen)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		mg := api.MonitorGroup{Name: args[0]}

		err := mg.Create()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorGroupCmd)
	monitorGroupCmd.AddCommand(monitorGroupCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorGroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorGroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
