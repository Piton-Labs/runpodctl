package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func WriteConfigFile(filename string, envVars []PodEnv) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var buffer bytes.Buffer

	for _, envVar := range envVars {
		fmt.Fprintf(&buffer, "%s=%s\n", envVar.Key, envVar.Value)
	}

	buffer.WriteTo(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("wrote environment to file: %s\n", filename)
}

func ReadConfigFile(filename string) ([]PodEnv, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []PodEnv

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			result = append(result, PodEnv{Key: parts[0], Value: parts[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func WritePodInput(outputFilename string, input *CreatePodInput) error {
	// Write code to output the config here
	return nil
}

func ReadPodInput(filename string) (CreatePodInput, error) {
	// Write code to input the config here
	return CreatePodInput{}, nil
}

type PodFilter struct {
	Image           string
	UnavailableGpu  bool
	GpuType         string
	AvailableStatus bool
}

func (f PodFilter) IsFilterImage() bool {
	return len(f.Image) > 0
}

func (f PodFilter) IsGpuFilter() bool {
	return len(f.GpuType) > 0
}

func (f PodFilter) IsAvailableGpuFilter() bool {
	return !f.UnavailableGpu
}

func (f PodFilter) IsAvailableStatusFilter() bool {
	return f.AvailableStatus
}

func (f PodFilter) HasFilter() bool {
	return f.IsFilterImage() || f.IsGpuFilter() || f.IsAvailableGpuFilter() || f.IsAvailableStatusFilter()
}

func (f PodFilter) IncludePod(p *Pod) bool {
	return f.isImageMatch(p) && f.isGpuMatch(p) && f.isAvailableMatch(p) && f.isAvailableStatus(p)
}

func (f PodFilter) PrintPodMatch(p *Pod) {
	if f.IsFilterImage() {
		fmt.Printf("Image Match:\n %s contains %s: %t\n", p.ImageName, f.Image, f.isImageMatch(p))
	}
	if f.IsGpuFilter() {
		fmt.Printf("GPU Match:\n %s contains %s: %t\n", p.Machine.GpuDisplayName, f.GpuType, f.isGpuMatch(p))
	}
	if f.IsAvailableGpuFilter() {
		fmt.Printf("Available Match:\n %t gpu count %d ok? %t\n", f.UnavailableGpu, p.GpuCount, f.isAvailableMatch(p))
	}
}

func GetFilteredPods(f *PodFilter) (pods []*Pod, err error) {
	var filteredPods []*Pod
	pods, err = GetPods()

	if err != nil {
		return pods, err
	}

	for _, p := range pods {
		if f.IncludePod(p) {
			filteredPods = append(filteredPods, p)
		}
	}

	allPodCount := len(pods)
	filteredPodCount := allPodCount - len(filteredPods)
	fmt.Printf("There are %d pods, %d are filtered out\n", allPodCount, filteredPodCount)
	f.Print()
	return filteredPods, nil
}

func (f PodFilter) isImageMatch(p *Pod) bool {
	return !f.IsFilterImage() || strings.Contains(strings.ToUpper(p.ImageName), strings.ToUpper(f.Image))
}

func (f PodFilter) isGpuMatch(p *Pod) bool {
	return !f.IsGpuFilter() || strings.Contains(strings.ToUpper(p.Machine.GpuDisplayName), strings.ToUpper(f.GpuType))
}

func (f PodFilter) isAvailableMatch(p *Pod) bool {
	return f.UnavailableGpu || p.GpuCount > 0
}

func (f PodFilter) isAvailableStatus(p *Pod) bool {
	return p.DesiredStatus == "EXITED"
}

func (f PodFilter) Print() {
	fmt.Printf("Image: %s\nGpu: %s\nUnavailable: %t\n", f.Image, f.GpuType, f.UnavailableGpu)
}

func TrainEnvFormat(varName string, val pflag.Value) PodEnv {
	varName = fmt.Sprintf("TRAINING_%s", strings.ToUpper(varName))
	return PodEnv{Key: varName, Value: val.String()}
}

func ModifyPod(p *Pod) {
	// This is the proper way to do it, but does not work due to an error in formatting
	// however it's exactly the same as the create syntax so I'll address it later - JPC
	// mod := GetModInputFromPod(p)

	// input := Input{
	// 	Query: `
	// 		mutation {
	// 			podEditJob(input: $input) {
	// 				containerDiskInGb
	// 				costPerHr
	// 				desiredStatus
	// 				dockerArgs
	// 				env
	// 				gpuCount
	// 				id
	// 				imageName
	// 				lastStatusChange
	// 				locked
	// 				machineId
	// 				memoryInGb
	// 				name
	// 				networkVolume {
	// 					dataCenterId
	// 					id
	// 					name
	// 					size
	// 				}
	// 				podType
	// 				ports
	// 				templateId
	// 				uptimeSeconds
	// 				vcpuCount
	// 				volumeInGb
	// 				volumeMountPath
	// 			}
	// 		}
	// 	`,
	// 	Variables: map[string]interface{}{"input": mod},
	// }

	input := Input{
		Query: `
		mutation {
			podEditJob(input: {
			  podId: "` + p.Id + `",
			  containerDiskInGb: ` + fmt.Sprintf("%d", p.ContainerDiskInGb) + `,
			  dockerArgs: "` + p.DockerArgs + `",
			  env: [
				` + strings.Join(p.Env, ", ") + `
			  ]
			  imageName: "` + p.ImageName + `",
			  ports: "` + p.Ports + `",
			  volumeInGb: ` + fmt.Sprintf("%d", p.VolumeInGb) + `,
			  volumeMountPath: "` + p.VolumeMountPath + `",
			}) {
			  containerDiskInGb
			  costPerHr
			  desiredStatus
			  dockerArgs
			  env
			  gpuCount
			  id
			  imageName
			  lastStatusChange
			  locked
			  machineId
			  memoryInGb
			  name
			  networkVolume {
				dataCenterId
				id
				name
				size
			  }
			  podType
			  ports
			  templateId
			  uptimeSeconds
			  vcpuCount
			  volumeInGb
			  volumeMountPath
			}
		  }
		`,
	}

	res, err := Query(input)
	if err != nil {
		fmt.Println("Query Error")
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO Error")
		fmt.Println(err)
		return
	}
	if res.StatusCode != 200 {
		fmt.Println("Status Code error")
		fmt.Println(res.StatusCode)
		fmt.Println(string(rawData))
		err = fmt.Errorf("statuscode %d: %s", res.StatusCode, string(rawData))
		// return
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal(rawData, &data); err != nil {
		fmt.Println("json.Unmarshal Error")
		fmt.Println(err)
		return
	}

	gqlErrors, ok := data["errors"].([]interface{})
	if ok && len(gqlErrors) > 0 {
		firstErr, _ := gqlErrors[0].(map[string]interface{})
		err = errors.New(firstErr["message"].(string))
		return
	}

	gqldata, ok := data["data"].(map[string]interface{})
	if !ok || gqldata == nil {
		err = fmt.Errorf("data is nil: %s", string(rawData))
		return
	}

	pod, ok := gqldata["podEditJob"].(map[string]interface{})
	if !ok || pod == nil {
		err = fmt.Errorf("pod is nil: %s", string(rawData))
		return
	}

	fmt.Println("Pod updated successfully. Starting it now")

	StartOnDemandPod(p.Id)

	return
}

type ModifyPodInput struct {
	PodId             string   `json:"podId"`
	ContainerDiskInGb int      `json:"containerDiskInGb"`
	DockerArgs        string   `json:"dockerArgs"`
	ImageName         string   `json:"imageName"`
	Env               []string `json:"env"`
	GpuCount          int      `json:"gpuCount"`
	Name              string   `json:"name"`
	Ports             string   `json:"ports"`
	VolumeInGb        int      `json:"volumeInGb"`
	VolumeMountPath   string   `json:"volumeMountPath"`
}

func GetModInputFromPod(p *Pod) *ModifyPodInput {
	return &ModifyPodInput{
		PodId:             p.Id,
		ContainerDiskInGb: p.ContainerDiskInGb,
		DockerArgs:        p.DockerArgs,
		ImageName:         p.ImageName,
		Env:               p.Env,
		GpuCount:          p.GpuCount,
		Name:              p.Name,
		Ports:             p.Ports,
		VolumeInGb:        p.VolumeInGb,
		VolumeMountPath:   p.VolumeMountPath,
	}
}
