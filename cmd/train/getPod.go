package train

import (
	"cli/api"
	"cli/format"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var IncludeUnavailableGPU bool
var AvailableStatus bool
var AllFields bool
var TrainOnly bool
var ImageFilter string
var GpuFilter string
var GetAll bool
var StatusFilter string

var GetPodCmd = &cobra.Command{
	Use:   "pod",
	Args:  cobra.MaximumNArgs(0),
	Short: "Get all pods with available gpus",
	Long:  "Gets all  pods with available gpus",
	Run: func(cmd *cobra.Command, args []string) {
		if TrainOnly {
			ImageFilter = "train"
		}

		var f = &api.PodFilter{
			Image:          ImageFilter,
			GpuType:        GpuFilter,
			UnavailableGpu: IncludeUnavailableGPU,
		}

		if f.HasFilter() {
			fmt.Println("\nFilters:")
			if f.IsFilterImage() {
				fmt.Printf("- Image:         %s\n", f.Image)
			}
			if f.IsGpuFilter() {
				fmt.Printf("- GpuType:       %s\n", f.GpuType)
			}
			if f.IsAvailableGpuFilter() {
				fmt.Println("- Available GPU Only: true")
			}

			fmt.Println("")
		}

		pods, err := api.GetFilteredPods(f)
		cobra.CheckErr(err)
		data := make([][]string, len(pods))
		for i, p := range pods {
			row := []string{p.Id, p.Name, fmt.Sprintf("%d/%d %s", p.Machine.GpuAvailable, p.GpuCount, p.Machine.GpuDisplayName), p.ImageName, p.DesiredStatus}
			if AllFields {
				row = append(
					row,
					p.PodType,
					fmt.Sprintf("%d", p.VcpuCount),
					fmt.Sprintf("%d", p.MemoryInGb),
					fmt.Sprintf("%d", p.ContainerDiskInGb),
					fmt.Sprintf("%d", p.VolumeInGb),
					fmt.Sprintf("%.3f", p.CostPerHr),
					fmt.Sprintf("%s", p.Env),
				)
			}
			data[i] = row
		}
		header := []string{"ID", "Name", "Avail/Required GPU", "Image Name", "Status"}
		if AllFields {
			header = append(header, "Pod Type", "vCPU", "Mem", "Container Disk", "Volume Disk", "$/hr", "ENV")
		}

		tb := tablewriter.NewWriter(os.Stdout)
		tb.SetHeader(header)
		tb.AppendBulk(data)
		format.TableDefaults(tb)
		tb.Render()
	},
}

func init() {
	GetPodCmd.Flags().BoolVar(&GetAll, "getAll", false, "ignore filters completely")
	GetPodCmd.Flags().BoolVarP(&IncludeUnavailableGPU, "includeUnavailable", "u", false, "include unavailable gpus")
	GetPodCmd.Flags().BoolVarP(&AllFields, "allfields", "a", false, "include all fields in output")
	GetPodCmd.Flags().BoolVarP(&TrainOnly, "train", "t", false, "include only images with train in their name")
	GetPodCmd.Flags().StringVarP(&ImageFilter, "image", "i", "", "filter out images that don't match the value")
	GetPodCmd.Flags().StringVarP(&GpuFilter, "gpu", "g", "", "filter out gpus that don't match the value")
	GetPodCmd.Flags().BoolVarP(&AvailableStatus, "available", "e", true, "Only include exited pods")
}
