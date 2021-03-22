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
	"bytes"
	"fmt"
	"github.com/bennettaur/changelink/services/changelink/trigger"
	"github.com/sourcegraph/go-diff/diff"
	"io"
	"log"
	"os"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var existingDiffFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "changelink",
	Short: "Analyze diffs and trigger actions based on what's changed",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var diffFile io.Reader
		runGit, err := cmd.Flags().GetBool("run-git")
		if err != nil {
			panic(err)
		}
		if runGit {
			var stdout bytes.Buffer
			gitCmd := exec.Command("git", "diff")
			gitCmd.Stdout = &stdout
			err := gitCmd.Run()
			if err != nil {
				panic(err)
			}
			diffFile = &stdout
		} else if existingDiffFile == os.Stdin.Name() {
			f := os.Stdin
			diffFile = os.Stdin
			defer func() {
				err := f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
		} else {
			f, err := os.Open(existingDiffFile)
			if err != nil {
				panic(err)
			}
			diffFile = f
			defer func() {
				err := f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		r := diff.NewMultiFileDiffReader(diffFile)
		var actionErrors []error
		for _, tw := range trigger.TriggerWatchers(r) {
			log.Printf("Triggering watcher: %v", tw.Watcher.Name)
			for _, action := range *tw.Watcher.Actions {
				err := action.Perform(tw.Watcher.Name, tw.Watcher.FilePath, tw.Reason, tw.TriggeredLines)
				if err != nil {
					actionErrors = append(actionErrors, err)
				}
			}
		}

		if len(actionErrors) > 0 {
			fmt.Printf("Received the following errors:\n %v", actionErrors)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&existingDiffFile, "diffFile", os.Stdin.Name(), "diff file (default is stdin)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("run-git", "g", false, "Run git diff to generate diff")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
