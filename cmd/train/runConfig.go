package train

import (
	"cli/api"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var CreatePod bool
var AllPods bool
var ForceCreate bool
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
		var configEnv []*api.PodEnv
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

		if !ForceCreate && len(pods) > 0 {
			launchExistingPod(pods[0], configEnv)
		} else if ForceCreate || CreatePod {
			createPod(configEnv)
		} else {
			fmt.Println("No existing pods that matched requirements and -c or --forceCreate were not specified to create a pod.")
		}

	},
}

func init() {
	RunConfigCmd.Flags().StringVar(&inputFilename, "filename", "", "The name of the file to read from")
	RunConfigCmd.Flags().StringVar(&outputFilename, "output", ".env", "The name of the file to read from")
	RunConfigCmd.Flags().BoolVarP(&CreatePod, "create", "c", false, "Creates a new pod if one of the existing pods won't work")
	RunConfigCmd.Flags().BoolVarP(&AllPods, "all", "a", false, "Does not restrict to images with train")
	RunConfigCmd.Flags().BoolVar(&ForceCreate, "forceCreate", false, "Forces Creation of a new pod")

	AddTrainFlags(RunConfigCmd)
	AddPodFlags(RunConfigCmd)
}

func createPod(configEnv []*api.PodEnv) {
	fmt.Println("Create and run new pod")
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
	input.Env = configEnv

	input.CloudType = "SECURE"

	if len(outputFilename) > 0 {
		api.WritePodInput(outputFilename, input)
	}

	missingEnv := input.CheckEnv()
	if len(missingEnv) > 0 {
		fmt.Printf("ERROR: The following environment variables must be loaded from a config or set with a flag:\n %s\n", strings.Join(missingEnv, "\n "))
		return
	}

	isImageAllowed := len(input.ImageName) > 0 && len(input.GpuTypeId) > 0
	isTemplateAllowed := len(input.TemplateId) > 0
	if !isImageAllowed && !isTemplateAllowed {
		fmt.Print("To create a pod an image or templateId must be provided, if an image is specified you must also set the GpuTypeId\n\n")
		return
	}

	fmt.Printf("isImageAllowed %t isTemplateAllowed %t %s", isImageAllowed, isTemplateAllowed, input.TemplateId)

	PrintPodInputEnv(input)

	pod, err := api.CreatePod(input)
	cobra.CheckErr(err)

	if pod["desiredStatus"] == "RUNNING" {
		fmt.Printf(`pod "%s" created for $%.3f / hr`, pod["id"], pod["costPerHr"])
		fmt.Println()
	} else {
		cobra.CheckErr(fmt.Errorf(`pod "%s" start failed; status is %s`, pod["id"], pod["desiredStatus"]))
	}
}

func launchExistingPod(pod *api.Pod, configEnv []*api.PodEnv) {
	fmt.Println("Launch existing pod")
	// Map the existing environment variables
	newEnv := map[string]string{}
	for _, item := range pod.Env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			newEnv[parts[0]] = strings.Join(parts[1:], "=")
		} else {
			fmt.Printf("Invalid env variable %s len %d\n", item, len(parts))
		}
	}

	// Overwrite or add any new variables
	for _, item := range configEnv {
		newEnv[item.Key] = item.Value
	}

	// Convert the variables to the required format for
	var strEnv []string
	for key, val := range newEnv {
		strEnv = append(strEnv, fmt.Sprintf("{ key: \"%s\", value: \"%s\" }\n", key, val))
	}
	pod.Env = strEnv

	// missingEnv := pod.CheckEnv()
	// if len(missingEnv) > 0 {
	// 	fmt.Printf("ERROR: The following environment variables must be loaded from a config or set with a flag:\n %s\n", strings.Join(missingEnv, "\n "))
	// 	fmt.Print(pod.Env)
	// 	return
	// }

	err := api.ModifyPod(pod)
	if err != nil {
		fmt.Println("Error creating pod")
		fmt.Println(err)
		return
	}

	fmt.Printf("Pod %s updated successfully. Starting it now\n", pod.Id)

	restartPod, err := api.StartOnDemandPod(pod.Id)

	if err != nil {

	}

	PrintPodEnv(pod)
	fmt.Println(restartPod)

}

func PrintPodEnv(pod *api.Pod) {
	fmt.Printf("%s\n- GPU: %s\n- Image: %s\n- GpuCount: %d\n- Env\n", pod.Name, pod.Machine.GpuDisplayName, pod.ImageName, pod.GpuCount)
	for _, val := range pod.Env {
		fmt.Printf("  - %s", val)
	}
}

func PrintPodInputEnv(pod *api.CreatePodInput) {
	fmt.Printf("%s\n- GPU: %s\n- Image: %s\n- GpuCount: %d\n- Env\n", pod.Name, pod.GpuTypeId, pod.ImageName, pod.GpuCount)
	for _, val := range pod.Env {
		fmt.Printf("  - %s=%s\n", val.Key, val.Value)
	}
}
