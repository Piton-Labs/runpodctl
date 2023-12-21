package train

import (
	"cli/api"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var CreatePod bool
var AllPods bool
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
			configEnv, err = api.ReadConfigFile(inputFilename)
			if err != nil {
				fmt.Print(err)
				return
			}
		}

		imageFilter := "train"
		if AllPods {
			imageFilter = ""
		}

		var f = &api.PodFilter{
			Image:          imageFilter,
			UnavailableGpu: false,
		}

		pods, err := api.GetFilteredPods(f)
		if err != nil {
			fmt.Println(err)
		}

		if len(pods) > 0 {
			launchExistingPod(pods[0], &configEnv)
		} else if CreatePod {
			createPod(&configEnv)
		} else {
			fmt.Println("No existing pods that matched requirements and -c was not specified to create a pod.")
		}

	},
}

func init() {
	RunConfigCmd.Flags().StringVar(&inputFilename, "filename", "", "The name of the file to read from")
	RunConfigCmd.Flags().StringVar(&outputFilename, "output", ".env", "The name of the file to read from")
	RunConfigCmd.Flags().BoolVarP(&CreatePod, "create", "c", false, "Creates a new pod if required")
	RunConfigCmd.Flags().BoolVarP(&AllPods, "all", "a", false, "Does not restrict to images with train")

	AddTrainFlags(RunConfigCmd)
	AddPodFlags(RunConfigCmd)
}

func createPod(configEnv *[]api.PodEnv) {
	fmt.Println("Create and run new pod")
	// input := &api.CreatePodInput{
	// 	ContainerDiskInGb: containerDiskInGb,
	// 	DeployCost:        deployCost,
	// 	DockerArgs:        dockerArgs,
	// 	GpuCount:          gpuCount,
	// 	GpuTypeId:         gpuTypeId,
	// 	ImageName:         imageName,
	// 	MinMemoryInGb:     minMemoryInGb,
	// 	MinVcpuCount:      minVcpuCount,
	// 	Name:              name,
	// 	TemplateId:        templateId,
	// 	VolumeInGb:        volumeInGb,
	// 	VolumeMountPath:   volumeMountPath,
	// }
	// if len(ports) > 0 {
	// 	input.Ports = strings.Join(ports, ",")
	// }
	// input.Env = make([]*api.PodEnv, len(env))

	// for i, v := range env {
	// 	e := strings.Split(v, "=")
	// 	if len(e) != 2 {
	// 		cobra.CheckErr(fmt.Errorf("wrong env value: %s", e))
	// 	}
	// 	input.Env[i] = &api.PodEnv{Key: e[0], Value: e[1]}
	// }

	// input.CloudType = "SECURE"

	// if len(outputFilename) > 0 {
	// 	api.WritePodInput(outputFilename, input)
	// }

	// pod, err := api.CreatePod(input)
	// cobra.CheckErr(err)

	// if pod["desiredStatus"] == "RUNNING" {
	// 	fmt.Printf(`pod "%s" created for $%.3f / hr`, pod["id"], pod["costPerHr"])
	// 	fmt.Println()
	// } else {
	// 	cobra.CheckErr(fmt.Errorf(`pod "%s" start failed; status is %s`, args[0], pod["desiredStatus"]))
	// }
}

func launchExistingPod(pod *api.Pod, configEnv *[]api.PodEnv) {
	fmt.Println("Launch existing pod")
	// Map the existing environment variables
	newEnv := map[string]string{}
	for _, item := range pod.Env {
		parts := strings.Split(item, "=")
		if len(parts) == 2 {
			fmt.Printf("%s: %s\n", parts[0], parts[1])
			newEnv[parts[0]] = parts[1]
		} else {
			fmt.Println("Invalid env variable " + item)
		}
	}

	// Overwrite or add any new variables
	for _, item := range *configEnv {
		newEnv[item.Key] = item.Value
	}

	// Convert the variables to the required format for
	var strEnv []string
	for key, val := range newEnv {
		strEnv = append(strEnv, fmt.Sprintf("{ key: \"%s\", value: \"%s\" }\n", key, val))
	}
	pod.Env = strEnv

	api.ModifyPod(pod)
}

func PrintPodEnv(pod *api.Pod) {
	fmt.Printf("%s\n- GPU: %s\n- Image: %s\n- GpuCount: %d\n- Env\n", pod.Name, pod.Machine.GpuDisplayName, pod.ImageName, pod.GpuCount)
	for _, val := range pod.Env {
		fmt.Printf("  - %s", val)
	}
}
