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
	"site24x7/logger"
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
	Aliases: []string{"configure", "cfg"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		updateRefreshTokenOnly, _ := cmd.Flags().GetBool("refresh-token")

		// If a config file already exists, verify that the user wants to
		// overwrite that file (or a given item in it).

		// TODO: Move work to impl package and test
		var overwrite string

		if f := viper.ConfigFileUsed(); f != "" {
			if !updateRefreshTokenOnly {
				fmt.Print("A config file already exists, do you want to overwrite it? [y/N]: ")
			} else {
				fmt.Print("A config file already exists, do you want to overwrite its refresh_token value? [y/N]: ")
			}
			fmt.Scanln(&overwrite)

			if strings.ToUpper(overwrite) != "Y" {
				fmt.Println("No changes were made; exiting.")
				os.Exit(0)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(cmd.Flags())

		updateRefreshTokenOnly, _ := cmd.Flags().GetBool("refresh-token")

		// TODO: Move work to impl package and test
		var clientID string
		var clientSecret string
		var grantToken string

		// Request the client id and secret
		if !updateRefreshTokenOnly {
			fmt.Print("Site24x7 Client ID [None]: ")
			fmt.Scanln(&clientID)
			fmt.Print("Site24x7 Client Secret [None]: ")
			fmt.Scanln(&clientSecret)

			if clientID == "" || clientSecret == "" {
				logger.Out("At least one empty value provided; nothing to do.")
				logger.Out("Without providing both a client id and secret, this tool is useless.")
				os.Exit(0)
			}
		}

		// Request the grant token from the user
		fmt.Print("Site24x7 Grant Token [None]: ")
		fmt.Scanln(&grantToken)

		if grantToken == "" {
			logger.Out("No grant token provided; nothing to do\n")

			return nil
		}

		// Update the client values if we're not dealing with a refresh token
		// only call
		if !updateRefreshTokenOnly {
			viper.Set("auth.client_id", clientID)
			viper.Set("auth.client_secret", clientSecret)
		}

		// Exchange the grant token for a refresh token
		refreshToken, err := api.Configure(grantToken)
		if err != nil {
			logger.Warn("Unable to exchange the grant token provided for a refresh token.")
			return fmt.Errorf("%s", err)
		}

		viper.Set("auth.refresh_token", refreshToken)
		if err = viper.WriteConfig(); err != nil {
			return fmt.Errorf("unable to complete configuration (%s)", err)
		}

		logger.Out("Configuration complete!")

		return nil
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
	configCmd.Flags().BoolP("refresh-token", "r", false, "Updates the refresh token only")
}
