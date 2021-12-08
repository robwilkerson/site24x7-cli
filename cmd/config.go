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
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <grant token>",
	Short: "Configures the Site24x7 CLI tool for use",
	Long: `Configures the Site24x7 CLI tool for use
	
Accepts a manually-generated, custom-scoped, short-lived grant token, exchanges
that for a long-lived refreh token and saves that refresh token to the user's
home directory, e.g. ~/<username>/.site24x7cli/config.`,
	Aliases: []string{"configure"},
	Run: func(cmd *cobra.Command, args []string) {
		var grantToken string

		// TODO: get user input for client ID and client secret

		// Request the grant token from the user
		fmt.Print("Site24x7 Grant Token [None]: ")
		fmt.Scanln(&grantToken)
		if grantToken == "" {
			fmt.Printf("No grant token provided; nothing to do\n")
			os.Exit(0)
		}

		// Exchange the grant token for a refresh token
		refreshToken, err := api.Configure(grantToken)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		// Write the refresh token to the config file provided by Cobra
		viper.Set("auth.refresh_token", refreshToken)
		err = viper.WriteConfig()
		if err != nil {
			fmt.Printf("Error complete configuration")
			os.Exit(1)
		}

		fmt.Println("Configuration complete!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
