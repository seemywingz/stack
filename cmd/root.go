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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// Global
var (
	region,
	profile,
	cluster string
	sess *session.Session
	svc  *ecs.ECS
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stack",
	Short: "Amazon ECS Interface",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Rocks")
	// },
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get Current Cluster Status",
	Long:  ``,
	Run:   getStatus,
}

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "List Current Cluster Events",
	Long:  ``,
	Run:   getEvents,
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
	cobra.OnInitialize(getSession)

	// Parse Global Flagss
	RootCmd.PersistentFlags().StringVarP(&profile, "profle", "p", "default", "Specify which AWS Profile to use")
	RootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "Set AWS Region")
	RootCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "default", "Set AWS ECS Cluster Name")

	// Add SubCommand
	RootCmd.AddCommand(statusCmd)

	RootCmd.AddCommand(eventCmd)
	eventCmd.Flags().StringP("service", "s", "sample-webapp", "Set the name of the service")
}

func getSession() {
	// Specify profile for config and region for requests
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))
	svc = ecs.New(sess)
}

func getStatus(cmd *cobra.Command, args []string) {
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(cluster),
		},
	}
	result, err := svc.DescribeClusters(input)
	EoE("Error Getting Cluster Status", err)
	fmt.Println(result)
}

func getEvents(cmd *cobra.Command, args []string) {
	input := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(cmd.Flag("service").Value.String()),
		},
	}

	result, err := svc.DescribeServices(input)
	EoE("Error Describing Container Instances", err)
	events := result.Services[0].Events
	fmt.Println(events)
}
