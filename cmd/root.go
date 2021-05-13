package cmd

import (
	"fire-up/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// Error message for wrong usage of this sub-command
	rootErr = `No artifact-alias provided!
Usage: fire-up --alias <artifact-alias>`
)

var cfgFile string

var artifactAlias string

var artifactNewName string

var addArtifactFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fire-up",
	Short: "A simple spell to bootstrap any software project",
	Long:
	`A simple spell to bootstrap any software project.
fire-up can look up for the artifacts from a github repository source or a local directory.
No need to have a boilerplate generator for every different type of project!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Bare application ran")
		if artifactAlias == "" {
			log.Fatal(rootErr)
		}
		if addArtifactFlag{
			// TODO: Inject component instead
		} else {
			utils.InitializeProjectFromArtifact(artifactAlias, artifactNewName)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fire-up.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&artifactAlias, "alias", "a", "", "artifact's alias")
	rootCmd.PersistentFlags().StringVarP(&artifactNewName, "rename", "r", "", "artifact's name")
	rootCmd.PersistentFlags().BoolVarP(&addArtifactFlag,"add-component","c",false,"Add a component instead of an artifact")
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

		// Search config in home directory with name ".fire-up" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".fire-up")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
