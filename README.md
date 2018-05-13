# Site Monitor
## Project structure

cmd/ - Contains the entry point - main.go

monitor/ - Contains the monitor package

infra/ - Contains all the terraform files

monitor.config - Config file containing information on which url to fetch and what to check

## Monitor.Config

Monitor.config contains the information about the website and the checks to be done. It is a json file. The schema file [monitor.config.schema.json](doc_files/monitor.config.schema.json) contains the json schema.

## Log

The location of the log file can be configured through the use of environment variables. Details of all environment variables are present in the [config](##Config) section. Each of the line in the log file is in json format. The schema file for the same is [results.schema.json](doc_files/results.schema.json).


## Config


| Key                      | Default Value    | Comments                                                      |
|--------------------------|------------------|---------------------------------------------------------------|
| MONITOR_CONFIG_FILE_PATH | ./monitor.config | The config file which has all the urls and checks to be done. |
| MONITOR_LOG_DIR_PATH     | .                | Directory path for creating log files.                        |
| MONITOR_USER_AGENT       | monitorv1.0.0    | User-Agent set in the request header.                         |
| MONITOR_TIME_OUT         | 30               | http request time out (in seconds)                            |


## How to run

### docker
```
docker build -t monitor .

docker run monitor

```



### go
```
dep ensure --vendor-only
go run cmd/main.go
```


## Setting up AWS

Go to infra and run the below commands

```
terraform init
terraform get && terraform apply -var-file=environments/dev.tfvars
```

Please change the value set in the dev.tfvars as per requirement.