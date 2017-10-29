// Copyright © 2017 Kevin Jayne
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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// Global Flags
var (
	tail,
	verbose,
	dryRun bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stack",
	Short: "Amazon ECS Interface",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rocks")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		getStatus()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(GetConfig)

	// Parse Global Flagss
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Vebose output")
	RootCmd.PersistentFlags().BoolVarP(&tail, "tail", "t", false, "Tail the logs")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Perform Dry Run, don't actually execute any actions")

	// Add SubCommand
	RootCmd.AddCommand(statusCmd)
}

func getStatus() {
	fmt.Println("Status")

	// Create a Session with a custom region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := ecs.New(sess)
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String("default"),
		},
	}
	result, err := svc.DescribeClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
