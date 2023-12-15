package cmd

import (
	"cli/cmd/train"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var numProcesses int

var trainCmd = &cobra.Command{
	Use:   "train [command]",
	Short: "train config",
	Long:  "get resources for pods",
}

func init() {
	train.AddTrainFlags(trainCmd)

	trainCmd.AddCommand(train.ConfigCmd)
	trainCmd.AddCommand(train.GetPodCmd)
	trainCmd.AddCommand(train.RunConfigCmd)
}

func GetPodEnv() {
	trainCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("Name: %s, Value: %s\n", flag.Name, flag.Value)
	})
}
