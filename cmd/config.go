// Copyright © 2017 Kevin Jayne <kevin.jayne@adp.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

const configDir = ".stack"
const configFileName = "config"

var homeDir, configFile string

type jsonConfig struct {
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Email string `json:"email"`
}

var config jsonConfig

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interface for to congiguring stack",
	Long:  ``,
	Run:   Configure,
}

func init() {
	RootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolP("list", "l", false, "List Current Config")
}

// SaveConfig : writes the current config to disk
func SaveConfig() {
	data, jsoEerr := json.Marshal(config)
	EoE("Error Parsing Json:", jsoEerr)
	if !dryRun {
		EoE("Error Saving Config File: "+configFile, ioutil.WriteFile(configFile, data, 0644))
	}
	if verbose {
		println("\n✨  Configuration File Saved Successfully")
	}
}

// ListConfig : prints the current config
func ListConfig() {
	println("")
	println("\n📖  Config")
	println("First Name:📓 ", config.Fname)
	println(" Last Name:📓 ", config.Lname)
	println("     Email:📧 ", config.Email)
}

// GetConfig : Check to see if there is a config file, if not create one
func GetConfig() {

	usr, err := user.Current()
	EoE("Error Getting User", err)
	homeDir := usr.HomeDir
	configFile = filepath.Join(homeDir, configDir, configFileName)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		println("❗  Stack CONFIG NOT FOUND")
		if Confirm("⚙  Want to Create one now?") {
			err := os.MkdirAll(filepath.Join(homeDir, configDir), os.ModePerm)
			EoE("Error Creating Config Directory:", err)
			buildConfig()
		} else {
			println("🏳  Skipping Configuration File Creation")
			os.Exit(10)
		}
	} else { // config exists
		jsonFile, err := ioutil.ReadFile(configFile)
		EoE("Error Reading Config File:", err)
		json.Unmarshal(jsonFile, &config)
	}
}

func buildConfig() {

	println("📝  Writing", configFile)
	SetFromInput(&config.Fname, "\nFirst Name: ")
	SetFromInput(&config.Lname, " Last Name: ")
	SetFromInput(&config.Email, "     Email: ")

	if Confirm("Save Configuratuon File?") {
		SaveConfig()
		println("\n✨  Configuration File Saved")
		os.Exit(0)
	} else {
		println("\n🚫  Configuration File Not Saved")
	}
}

// Configure : Gather User Informaton and save it to config file
func Configure(cmd *cobra.Command, args []string) {
	switch {
	case cmd.Flag("list").Value.String() == "true":
		ListConfig()
		return
	default:
		buildConfig()
	}
}
