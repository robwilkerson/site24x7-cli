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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configures the Site24x7 CLI tool for use",
	Long: `Configures the Site24x7 CLI tool for use.

Requests and stores authentication details that are required to access the
Site24x7 API for a given account. This data is stored in a config file located
at $HOME/<username>/.site24x7.yaml.`,
	Aliases: []string{"configure"},
	PreRun: func(cmd *cobra.Command, args []string) {
		// If a config file already exists, verify that the user wants to
		// overwrite that file.

		var overwrite string

		if f := viper.ConfigFileUsed(); f != "" {
			fmt.Print("A config file already exists, do you want to overwrite it? [y/N]: ")
			fmt.Scanln(&overwrite)

			if strings.ToUpper(overwrite) != "Y" {
				fmt.Println("Existing file will not be overwritten; exiting.")
				os.Exit(0)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var clientID string
		var clientSecret string
		var grantToken string

		// Request the client id and secret
		fmt.Print("Site24x7 Client ID [None]: ")
		fmt.Scanln(&clientID)
		fmt.Print("Site24x7 Client Secret [None]: ")
		fmt.Scanln(&clientSecret)

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
		viper.Set("auth.client_id", clientID)
		viper.Set("auth.client_secret", clientSecret)
		viper.Set("auth.refresh_token", refreshToken)
		if err = viper.WriteConfig(); err != nil {
			fmt.Printf("Error completing configuration")
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
