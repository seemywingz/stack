# Stack 
###### We Work -- ECS Stack Master Coding Challenge

## Setup
This CLI Uses the [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/) and requires proper credentials  
Have a look at the [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) to properly setup your creds  
`stack` uses the shared credentials feature, but environment variables can be used to override directives in `~/.aws/credentials`
## Install
To install this
If you have go installed locally  
```
go get github.com/seemywingz/stack
```
Or you can install the pre-built binary for Mac
```
git clone https://github.com/seemywingz/stack.git
./stack/scripts/install.sh
```
the binary will install to `/usr/local/bin/stack`

## Usage  
```
stack [command]

Available Commands:
  deploy      Deploy new Cloudformation Stack from json
  events      List Current Events for Provided Service
  help        Help about any command

Flags:
  -h, --help            help for stack
  -p, --profle string   Set  AWS Profile (default "default")
  -r, --region string   Set AWS Region (default "us-east-1")
  -s, --stack string    Set AWS Cloud Formation Stack Name (default "default")

Use "stack [command] --help" for more information about a command.
```

`stack events`: List Current Events for Provided Service 
```
stack events [flags]

Flags:
  -h, --help             help for events
  -n, --number int       Number of Events to Output (default -1)
      --service string   If provided, will return the events for the provided service, instead of the stack (default "nil")

Global Flags:
  -p, --profle string   Set  AWS Profile (default "default")
  -r, --region string   Set AWS Region (default "us-east-1")
  -s, --stack string    Set AWS Cloud Formation Stack Name (default "default")
```
#### Example: get last 5 events for MyStack
```
stack events -s MyStack -n 5 
```
#### Example: get lat 10 events for MyService
```
stack events --service MyService -n 10
```
  
`stack deploy`: Deploy New Stack from Provided JSON 
```
stack deploy [flags]

Flags:
  -f, --file string   Provide the File Location for Cloud Formation Input (default "stack.json")
  -h, --help          help for deploy

Global Flags:
  -p, --profle string   Set  AWS Profile (default "default")
  -r, --region string   Set AWS Region (default "us-east-1")
  -s, --stack string    Set AWS Cloud Formation Stack Name (default "default")
```
#### Example: Deploy MyStack from ~/path/to/MyStack.json  
```
stack deploy -s MyStack -f ~/path/to/MyStack.json
```

## How much time did you spend?
  * Overall I'd say about 4 hours 
  * To test the cli I made a test IAM user and ECS setup with Cloudformation Test Template  
  * I also spent a lot of time just playing around with output, deciding command syntax and restraining myself from adding features not requested ;)

## What was the most difficult thing for you?
  * The most difficult part was tryng to parse the cloudformation json input, json in go is not exactly fun :'(
  * I decided to use an open source project for this code example [go-cloudformation](https://github.com/crewjam/go-cloudformation)

## What technical debt would you pay if you had one more iteration?
  * I'd like to parse the cloudformation input json myslef
  * Make the cli more robust  
    * more options/defaults
    * better error handling
    * make it smarter