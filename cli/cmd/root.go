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
	"github.com/bennettaur/diffhook/services/diffhook/trigger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sourcegraph/go-diff/diff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

var cfgFile string
var existingDiffFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diffhook",
	Short: "Analyze diffs and trigger actions based on what's changed",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var diffFile io.Reader
		branch, err := cmd.Flags().GetString("git")
		if err != nil {
			panic(err)
		}
		if len(branch) > 0 {
			err = gitFetch(branch)
			if err != nil {
				panic(err)
			}
			diffFile, err = gitDiff(branch)
			if err != nil {
				panic(err)
			}
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

	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	persistentFlags.StringVar(&existingDiffFile, "diffFile", os.Stdin.Name(), "diff file (default is stdin)")
	persistentFlags.String("git", "", "Run git diff to generate diff")
	persistentFlags.Lookup("git").NoOptDefVal = "origin/main"

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
