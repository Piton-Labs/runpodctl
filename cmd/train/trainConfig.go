package train

import (
	"cli/api"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var filename string
var isS3 bool

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.MaximumNArgs(1),
	Short: "Create a training config",
	Long:  "Create a training config and save it locally",
	Run: func(cmd *cobra.Command, args []string) {
		var env []api.PodEnv
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			env = append(env, api.TrainEnvFormat(flag.Name, flag.Value))
		})

		api.WriteConfigFile(filename, env)
	},
}

func init() {
	ConfigCmd.Flags().StringVar(&filename, "filename", ".env", "The name of the file to save to")
	AddTrainFlags(ConfigCmd)
}
