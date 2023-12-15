package train

import (
	"cli/api"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var CreatePod bool
var inputFilename string
var outputFilename string

var communityCloud bool
var secureCloud bool
var containerDiskInGb int
var deployCost float32

var dockerArgs string
var env []string
var gpuCount int
var gpuTypeId string
var imageName string
var minMemoryInGb int
var minVcpuCount int
var name string
var ports []string
var templateId string
var volumeInGb int
var volumeMountPath string

var RunConfigCmd = &cobra.Command{
	Use:   "run",
	Args:  cobra.MaximumNArgs(1),
	Short: "Run a training config, attempts to use existing pod",
	Long: `Run a training config, attempts to use existing pod.
		If a template id is specified it will create a new pod if an existing pod is not found
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var configEnv []api.PodEnv
		var err error
		if len(inputFilename) > 0 {
			configEnv, err = api.ReadConfigFile(filename)
			if err != nil {
				fmt.Print(err)
				return
			}
			fmt.Print(configEnv)
		}

	},
}

func init() {
	RunConfigCmd.Flags().StringVar(&inputFilename, "filename", "", "The name of the file to read from")
	RunConfigCmd.Flags().StringVar(&outputFilename, "output", ".env", "The name of the file to read from")
	RunConfigCmd.Flags().BoolVarP(&CreatePod, "create", "c", false, "Creates a new pod if required")

	AddTrainFlags(RunConfigCmd)
	AddPodFlags(RunConfigCmd)
}

func createPod(configEnv []*api.PodEnv) {
	input := &api.CreatePodInput{
		ContainerDiskInGb: containerDiskInGb,
		DeployCost:        deployCost,
		DockerArgs:        dockerArgs,
		GpuCount:          gpuCount,
		GpuTypeId:         gpuTypeId,
		ImageName:         imageName,
		MinMemoryInGb:     minMemoryInGb,
		MinVcpuCount:      minVcpuCount,
		Name:              name,
		TemplateId:        templateId,
		VolumeInGb:        volumeInGb,
		VolumeMountPath:   volumeMountPath,
	}
	if len(ports) > 0 {
		input.Ports = strings.Join(ports, ",")
	}
	input.Env = make([]*api.PodEnv, len(env))

	for i, v := range env {
		e := strings.Split(v, "=")
		if len(e) != 2 {
			cobra.CheckErr(fmt.Errorf("wrong env value: %s", e))
		}
		input.Env[i] = &api.PodEnv{Key: e[0], Value: e[1]}
	}

	input.CloudType = "SECURE"

	if len(outputFilename) > 0 {
		api.WritePodInput(outputFilename, input)
	}

	// pod, err := api.CreatePod(input)
	// cobra.CheckErr(err)

	// if pod["desiredStatus"] == "RUNNING" {
	// 	fmt.Printf(`pod "%s" created for $%.3f / hr`, pod["id"], pod["costPerHr"])
	// 	fmt.Println()
	// } else {
	// 	cobra.CheckErr(fmt.Errorf(`pod "%s" start failed; status is %s`, args[0], pod["desiredStatus"]))
	// }
}
