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
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// Global
var (
	region,
	profile,
	cluster string
	sess   *session.Session
	svcECS *ecs.ECS
	svcCF  *cloudformation.CloudFormation
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
	Use:   "events",
	Short: "List Current Events for Provided Service",
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
	RootCmd.PersistentFlags().StringVarP(&profile, "profle", "p", "default", "Set  AWS Profile")
	RootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "Set AWS Region")
	RootCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "default", "Set AWS ECS Cluster Name")

	// Add SubCommand
	RootCmd.AddCommand(statusCmd)

	RootCmd.AddCommand(eventCmd)
	eventCmd.Flags().StringP("service", "s", "default", "Set the name of the service")
}

func getSession() {
	// Specify profile for config and region for requests
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))
	svcECS = ecs.New(sess)
	svcCF = cloudformation.New(sess)
}

func getStatus(cmd *cobra.Command, args []string) {
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(cluster),
		},
	}
	result, err := svcECS.DescribeClusters(input)
	EoE("Error Getting Cluster Status", err)
	fmt.Println(result)
}

func getServeiceEvents(cmd *cobra.Command) {
	input := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(cmd.Flag("service").Value.String()),
		},
	}

	result, err := svcECS.DescribeServices(input)
	EoE("Error Describing Container Instances", err)
	if len(result.Services) == 0 {
		erSrv := cmd.Flag("service").Value.String()
		EoE("Error Finding Events for "+erSrv, errors.New("Please Check Service Name"))
	}
	events := result.Services[0].Events
	fmt.Println(events)
}

func getEvents(cmd *cobra.Command, args []string) {

	if cmd.Flag("service").Value.String() != "default" {
		getServeiceEvents(cmd)
	}

	input := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String("EC2ContainerService-default"),
	}
	result, err := svcCF.DescribeStackEvents(input)
	EoE("Error Describing Stack Events", err)
	fmt.Println(result)
}
