// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/asciifaceman/gopen/pkg/logging"
	"github.com/asciifaceman/gopen/pkg/nmap"

	"github.com/spf13/cobra"
)

var targets []string
var flags []string

var (
	poolSize = 1
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := logging.Init()
		if err != nil {
			fmt.Printf("Failed to get logger: %v", err)
			return
		}
		logger := logging.Logger()

		logger.Info("Starting gopen...")

		s, _ := nmap.NewScanner(&nmap.Config{
			Logger:   logger,
			Size:     poolSize,
			ExecPath: "/usr/local/bin/nmap",
		})

		// Make sure quotes are clean
		//var cleanFlags []string
		//for _, flag := range flags {
		//	cleanFlags = append(cleanFlags, flag)
		//}

		if len(flags) > 0 {
			for _, target := range targets {
				s.NewTaskWithOptions(
					nmap.WithTarget(target),
					nmap.WithFlags(flags),
				)
			}
		} else {
			for _, target := range targets {
				s.NewTaskWithDefaultOptions(target)
			}
		}

		// Close channel to signify no new tasks will be sent
		close(s.Pending)

		s.ExecPool.Wait()

	},
}

func init() {
	nmapCmd.AddCommand(runCmd)
	nmapCmd.PersistentFlags().StringSliceVarP(&targets, "targets", "t", targets, "repeated or comma delimited target hosts")
	nmapCmd.PersistentFlags().StringArrayVarP(&flags, "flags", "f", flags, "repeated or comma delimited nmap flags")

	nmapCmd.PersistentFlags().IntVarP(&poolSize, "pool", "p", poolSize, "Whole int for number of workers")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
