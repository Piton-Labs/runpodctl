package train

import "github.com/spf13/cobra"

func AddTrainFlags(cmd *cobra.Command) {
	cmd.Flags().Int("num_processes", 1, "Number of processes")
	cmd.Flags().Int("num_machines", 1, "Number of machines")
	cmd.Flags().Int("num_cpu_threadsPerProcess", 20, "Number Cpu Threads PerP rocess")
	cmd.Flags().Int("prior_loss_weight", 1, "Prior Loss Weight")
	cmd.Flags().Int("gradient_accumulation_steps", 1, "Gradient Accumulation Steps")
	cmd.Flags().Int("network_dim", 160, "Network Dim")
	cmd.Flags().Int("clip_skip", 2, "Clip Skip")
	cmd.Flags().Int("min_bucket_reso", 320, "Min Bucket Resolution")
	cmd.Flags().Int("max_bucket_reso", 768, "Max Bucket Resolution")
	cmd.Flags().Float32("t_encoder_lr", 0.00004, "T Encoder LR")
	cmd.Flags().Float32("learning_rate", 0.0001, "Learning Rate")
	cmd.Flags().Float32("unet_lr", 0.0001, "Unet Lr")
	cmd.Flags().String("mixed_precision", "fp16", "Mixed precision")
	cmd.Flags().String("pretrained_model_name_or_path", "/stable/base_models/", "Pretrained Model Name or Path")
	cmd.Flags().String("train_data_dir", "/stable/training_data", "Training Data dir")
	cmd.Flags().String("resolution", "512", "Resolution")
	cmd.Flags().String("output_dir", "/stable/training_outputs", "Output Directory")
	cmd.Flags().String("logging_dir", "log", "Logging Directory")
	cmd.Flags().String("network_alpha", "128", "Network Alpha")
	cmd.Flags().String("save_model_as", "safetensors", "Save Model As")
	cmd.Flags().String("network_module", "networks.lora", "Network Module")
	cmd.Flags().String("train_batch_Size", "16", "Training Batch Size")
	cmd.Flags().String("output_name", "$LORA_NAME", "Lora Name")
	cmd.Flags().String("lr_scheduler", "cosine_with_restarts", "Lr Schedule")
	cmd.Flags().String("max_train_epochs", "80", "Maximum Train Epochs")
	cmd.Flags().String("save_precision", "fp16", "Save Precision")
	cmd.Flags().String("caption_extension", ".txt", "Caption extension")
	cmd.Flags().String("max_data_loader_n_workers", "12", "Maximum data Loader Workers")
	cmd.Flags().String("save_every_n_epochs", "5", "Save Every N epochs")
}

func AddPodFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&communityCloud, "communityCloud", false, "create in community cloud")
	cmd.Flags().BoolVar(&secureCloud, "secureCloud", false, "create in secure cloud")
	cmd.Flags().IntVar(&containerDiskInGb, "containerDiskSize", 20, "container disk size in GB")
	cmd.Flags().Float32Var(&deployCost, "cost", 0, "$/hr price ceiling, if not defined, pod will be created with lowest price available")
	cmd.Flags().StringSliceVar(&env, "env", nil, "container arguments")
	cmd.Flags().IntVar(&gpuCount, "gpuCount", 1, "number of GPUs for the pod")
	cmd.Flags().StringVar(&gpuTypeId, "gpuType", "", "gpu type id, e.g. 'NVIDIA GeForce RTX 3090'")
	cmd.Flags().StringVar(&imageName, "imageName", "", "container image name")
	cmd.Flags().IntVar(&minMemoryInGb, "mem", 20, "minimum system memory needed")
	cmd.Flags().IntVar(&minVcpuCount, "vcpu", 1, "minimum vCPUs needed")
	cmd.Flags().StringVar(&name, "name", "", "any pod name for easy reference")
	cmd.Flags().StringSliceVar(&ports, "ports", nil, "ports to expose; max only 1 http and 1 tcp allowed; e.g. '8888/http'")
	cmd.Flags().StringVar(&templateId, "templateId", "", "templateId to use with the pod")
	cmd.Flags().IntVar(&volumeInGb, "volumeSize", 1, "persistent volume disk size in GB")
	cmd.Flags().StringVar(&volumeMountPath, "volumePath", "/runpod", "container volume path")

	// cmd.MarkFlagRequired("imageName") //nolint
}
