# Training using custom runpod

The added functionality for training using runpod helps you setup your config, find a machine to run on, and run the training session.

## Installation

Download code from github and execute ```make dev``` if you're on an Apple machine run ```make dev-mac``` this compiles the program and puts it in the <current directory</bin and is named ```runpodctl```

## Usage

To view the help for the new training functionality run ```bin/runpodctl train```


### List Pods

To see the list of pods run:
```
$ bin/runpodctl train pod

Filters:
- Available GPU Only: true

There are 1 pods, 0 are filtered out
ID              NAME                            AVAIL/REQUIRED GPU      IMAGE NAME                              STATUS
m5gy9l66zopb67  jeff_test_V2_public-Smaller     6/1 H100 80GB PCIe      maxpitonlabs/sd-docker-image-jeff-train EXITED
```

Note: By default only pods with available GPU's whose status is "EXITED" will be shown as this is a pod that can currently be run.

adding:

  -a, --allfields            include all fields in output
  -e, --available            Only include exited pods (default true)
      --getAll               ignore filters completely
  -g, --gpu string           filter out gpus that don't match the value
  -h, --help                 help for pod
  -i, --image string         filter out images that don't match the value
  -u, --includeUnavailable   include unavailable gpus
  -t, --train                include only images with 'train' in their name

### Run

You can reuse training pod when there are existing "Exited" pod instances. If you wish to create a pod start by using Runpod's web interface to create the pod and stop it as soon as the container is created. Eventually that will be added to this utility but due to an issue with the Runpod API and specifying docker registry credentials this is not currently possible. You should also create a config file with the environment variables needed to run a training job. See config below.

All of the filters available when listing pods are available when running a job.

```
./bin/runpodctl train run --config ~/training_runs/env
```

### Config

To create a default config file run:

```
./bin/runpodctl train config
2024/01/04 17:00:27 wrote environment to file: .env
```

You can also set the items in the config file with flags if you want, however, it's probably just easier to write out a default config and modify it in your favorite editor.

--caption_extension string               Caption extension (default ".txt")
--clip_skip int                          Clip Skip (default 2)
--filename string                        The name of the file to save to (default ".env")
--gradient_accumulation_steps int        Gradient Accumulation Steps (default 1)
-h, --help                                   help for config
--learning_rate float32                  Learning Rate (default 0.0001)
--logging_dir string                     Logging Directory (default "log")
--lr_scheduler string                    Lr Schedule (default "cosine_with_restarts")
--max_bucket_reso int                    Max Bucket Resolution (default 768)
--max_data_loader_n_workers string       Maximum data Loader Workers (default "12")
--max_train_epochs string                Maximum Train Epochs (default "80")
--min_bucket_reso int                    Min Bucket Resolution (default 320)
--mixed_precision string                 Mixed precision (default "fp16")
--network_alpha string                   Network Alpha (default "128")
--network_dim int                        Network Dim (default 160)
--network_module string                  Network Module (default "networks.lora")
--num_cpu_threadsPerProcess int          Number Cpu Threads PerP rocess (default 20)
--num_machines int                       Number of machines (default 1)
--num_processes int                      Number of processes (default 1)
--output_dir string                      Output Directory (default "/stable/training_outputs")
--output_name string                     Lora Name (default "$LORA_NAME")
--pretrained_model_name_or_path string   Pretrained Model Name or Path (default "/stable/base_models/")
--prior_loss_weight int                  Prior Loss Weight (default 1)
--resolution string                      Resolution (default "512")
--save_every_n_epochs string             Save Every N epochs (default "5")
--save_model_as string                   Save Model As (default "safetensors")
--save_precision string                  Save Precision (default "fp16")
--t_encoder_lr float32                   T Encoder LR (default 4e-05)
--train_batch_Size string                Training Batch Size (default "16")
--train_data_dir string                  Training Data dir (default "/stable/training_data")
--unet_lr float32                        Unet Lr (default 0.0001)
