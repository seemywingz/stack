package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

var loading = spinner.New(spinner.CharSets[39], 150*time.Millisecond)

// EoE : exit with error code 1 and print if err is notnull
func EoE(msg string, err error) {
	if err != nil {
		fmt.Printf("\n❌  %s\n   %v\n", msg, err)
		os.Exit(1)
	}
}

// GetHomeDir : returns a full path to user's home dorectory
func GetHomeDir() string {
	usr, err := user.Current()
	EoE("Failed to get Current User", err)
	if usr.HomeDir != "" {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// Confirm : return confirmation based on user input
func Confirm(q string) bool {
	a := GetInput(q + " (Y/n) ")
	var res bool
	switch a {
	case "":
		fallthrough
	case "y":
		fallthrough
	case "Y":
		res = true
	case "n":
	case "N":
		res = false
	default:
		return Confirm(q)
	}
	return res
}

// GetInput : return string of user input
func GetInput(q string) string {
	print(q)
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	return strings.TrimRight(ans, "\n")
}

// SelectFromArray : select an element in the provided array
func SelectFromArray(a []string) string {
	fmt.Println("Choices:")
	for i := range a {
		fmt.Println("[", i, "]: "+a[i])
	}
	sel, err := strconv.Atoi(GetInput("Enter Number of Selection: "))
	EoE("Error Getting Integer Input from User", err)
	if sel <= len(a)-1 {
		return a[sel]
	}
	return SelectFromArray(a)
}

// SetFromInput : set value of provided var to the value of user input
func SetFromInput(a *string, q string) {
	*a = strings.TrimRight(GetInput(q), "\n")
}

// SendRequest : send http request to provided url
func SendRequest(req *http.Request) []byte {
	client := http.Client{}
	res, err := client.Do(req)
	EoE("Error Getting HTTP Response", err)
	resData, err := ioutil.ReadAll(res.Body)
	EoE("Error Parsing HTTP Response", err)
	return resData
}

func sh(cmdStr string) string {
	// Perfsorm Shell Command
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmdOut, err := cmd.CombinedOutput()
	EoE("Error runnign command: "+string(cmdOut), err)
	return string(cmdOut)
}
