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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
	cf "github.com/crewjam/go-cloudformation"
	cfd "github.com/crewjam/go-cloudformation/deploycfn"
	"github.com/spf13/cobra"
)

// Global
var (
	region,
	profile,
	stack,
	jsonFile string
	n      int
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

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List Current Events for Provided Service",
	Long:  ``,
	Run:   getEvents,
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy new Cloudformation Stack from json",
	Run:   deploy,
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
	RootCmd.PersistentFlags().StringVarP(&stack, "stack", "s", "EC2ContainerService-default", "Set AWS Cloud Formation Stack Name")

	// Add SubCommand
	RootCmd.AddCommand(eventsCmd)
	eventsCmd.Flags().IntVarP(&n, "number", "n", -1, "Number of Events to Output")
	eventsCmd.Flags().String("service", "nil", "Set the name of the ECS container service")

	RootCmd.AddCommand(deployCmd)
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
	if n < 0 {
		n = len(events)
	}

	for i := 0; i < n; i++ {
		println()
		fmt.Println(events[i])
	}

	println()
	fmt.Printf("Events Printed: %v\n", n)
	fmt.Printf("  Events Total: %v\n", len(events))
}

func getEvents(cmd *cobra.Command, args []string) {
	if cmd.Flag("service").Value.String() != "nil" {
		getServeiceEvents(cmd)
	} else { // get Stack Events
		input := &cloudformation.DescribeStackEventsInput{
			StackName: aws.String(stack),
		}
		result, err := svcCF.DescribeStackEvents(input)
		EoE("Error Describing Stack Events", err)
		if n < 0 {
			n = len(result.StackEvents)
		}
		for i := 0; i < n; i++ {
			println()
			fmt.Println("    Resource Status:", *result.StackEvents[i].ResourceStatus)
			fmt.Println("Logical Resource Id:", *result.StackEvents[i].LogicalResourceId)
			fmt.Println("      Resource Type:", *result.StackEvents[i].ResourceType)
			fmt.Println("          Timestamp:", *result.StackEvents[i].Timestamp)
		}
		println()
		fmt.Printf("Events Printed: %v\n", n)
		fmt.Printf("  Events Total: %v\n", len(result.StackEvents))
	}
}

func deploy(cmd *cobra.Command, args []string) {

	jsonFile = "json/test.json"
	jsonObj, err := ioutil.ReadFile(jsonFile)
	EoE("Error Reading Config File:", err)

	t := &cf.Template{}
	json.Unmarshal(jsonObj, t)
	fmt.Printf("AWSTemplateFormatVersion: %s\n", t.AWSTemplateFormatVersion)

	deployInput := cfd.DeployInput{
		Session:        sess,
		StackName:      stack,
		Template:       t,
		Parameters:     nil,
		TemplateBucket: "",
	}
	EoE("Error Deploying New Stack", cfd.Deploy(deployInput))
}
