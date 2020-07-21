---
title: KubeEdge debug (Issue 324)
status: Pending
authors:
    - "@shenkonghui"
    - "@qingchen1203"
approvers:
    - 
creation-date: 2020-07-15
last-updated:  2020-07-16
---

# Motivation

Many users shared their feedback that kubeEdge edge nodes do not have good debugging and diagnostic methods, which may prevent people from trying kubeEdge.
There should be a simple and clear way to help operation and maintenance personnel to ensure the stable operation of kubeEdge, so that users can focus more on using it immediately.

#### Goal

- Alpha

Collect the full amount of information related to kubeedge in the current environment, and provide it to O&M personnel to locate and solve difficult problems.
- Beta

1. Diagnose specific fault scenarios in an all-round way and locate the cause of the fault.
2. Check whether the system specific items meet the requirements of edgecore installation and operation.

# Proposal

KubeEdge should have simple commands to debug and troubleshoot edge components.
Therefore, it is recommended to use the following commands during operation and maintenance of KubeEdge.

## Inscope

1. To support first set of basic commands (listed below) to debug edge (node) components.

For edge, commands shall be:

- `keadm help`
  - `keadm debug diagnose`
  - `keadm debug collect`
  - `keadm debug check`
  - `keadm debug get`

# Scope of commands

## Design of the commands

**NOTE**: All the below steps are executed as root user, to execute as sudo user. Please add `sudo` before all commands.

### keadm debug  --help

```
keadm help command provide debug function to help diagnose the cluster.

Usage:
  keadm debug [command]

Examples:

Available Commands:
  diagnose        Diagnose specific fault scenarios in an all-round way and locate the cause of the fault.
  collect         Obtain all the data of the current node, and provide it to the operation and maintenance personnel to locate the problem.
  check           Check whether the system specific items meet the requirements of edgecore installation and operation.
  get             Get and format data of available resource types in the local database of the edge node.

Flags:
  -h, --help      help for keadm debug

```

### keadm debug diagnose --help

```
keadm diagnose command can be help to diagnose specific fault scenarios in an all-round way and locate the cause of the fault.

Usage:
  keadm debug diagnose [command]

Examples:

# view the running status of node (key components such as sqlite, edgehub, metamanager, edged and many more)
keadm debug diagnose node

Available Commands:
  all           All resource
  node          Troubleshoot the cause of edge node failure with installed software
  pod           Troubleshooting specific container application instances on nodes
  installation  It is same as "keadm debug check all"

```

### keadm debug check --help

```
keadm check command can be check whether the system specific items meet the requirements of edgecore installation and operation.

Usage:
  keadm debug check [command]

Available Commands:
  all      Check all
  arch     Determine the node hardware architecture whether support or not. arch support amd64,arm64v8,arm32v7,i386 and s390x. qemu_arch support x86_64,aarch64,arm,i386 and s390x
  cpu      Determine if the NUMBER of CPU cores meets the requirement
  memory   Check the system memory size and the amount of memory left
  disk     Check whether the disk meets the requirements
  dns      Check whether the node domain name resolution function is normal
  runtime  Check whether the node Container runtime is normal, can use parameter `--runtime` to set container runtime,  default is docker
  network  Check whether the node can communicate with the endpoint on the cloud,can use parameter `--ip` to set test ip, default to ping clusterdns
  pid      Check if the current number of processes in the environment is too many. If the number of available processes is less than 5%, the number of processes is considered insufficient
  
Flags:
  -h, --help   help for keadm debug check

```

### keadm debug collect --help

```
Obtain all the data of the current node, and then provide it to the operation and maintenance personnel to locate the problem.

Usage:
  keadm debug collect [flags]

Examples:

keadm debug collect --output-path

Flags:
  --config       Specify configuration file, defalut is /etc/kubeedge/config/edgecore.yaml
  --output-path  Cache data and store data compression packages in a directory that defaults to the current directory
  --detail       Whether to print internal log output

Flags:
  -h, --help   help for keadm debug collect

```

### keadm debug get --help

```
Prints a table of the most important information about the specified resource from the local database of the edge node.

Usage:
  keadm debug get [resource] [flags]

Examples:

# List all pod
keadm debug get pod -A

# List all pod in namespace test
keadm debug get pod -n test

# List a single configmap  with specified NAME
keadm debug get configmap web -n default

# List the complete information of the configmap with the specified name in the yaml output format
keadm debug get configmap web -n default -o yaml

Available resource:
  all
  pod
  node
  service
  secret
  configmap
  endpoint

Flags:
  -h, --help            Help for keadm debug get
  -A, --all-namespaces  List the requested object(s) across all namespaces
  -n, --namespace=''    List the requested object(s) in specified namespaces
  -o, --output=''       Output format. One of:json|yaml|jsonpath=...
  -l, --selector=''     Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)

```

## Explaining the commands

### Worker Node (at the Edge) commands

`keadm debug diagnose`

- What is it?
  
- This command will be help to diagnose specific fault scenarios in an all-roundly way and locate the cause of the fault.
  
- What shall be its scope ?
    1. Use command `all` can diagnose all resource

    2. Use command `node` can troubleshoot the cause of edge node failure with installed software
       1. check system resources is enough
       2. check container runtime is running 
       3. check all edgecore components are running
       4. check cloudercore can be connected
       
    3. Use command `pod` can troubleshooting specific container application instances on nodes
       1. check pod probe
       2. check pod storage
       3. check pod network
       4. other more

    4. Use command `installation` is same as "keadm debug check all"

       

`keadm debug check`

- What is it?
  
  - This command will be check whether the system specific items meet the requirements of edgecore installation and operation.
- What shall be its scope ?

  1. Check items include hardware resources or operating system resources (cpu, memory, disk, network, pid limit,etc.)
  2. Use command `arch` can check node hardware architecture:
     - arch: amd64,arm64v8,arm32v7,i386 and s390x
     - qemu_arch: x86_64,aarch64,arm,i386 and s390x
  3. Use command `cpu` can determine if the NUMBER of CPU cores meets the requirement, minimum 1Vcores.
  4. Use command `memory` check the system memory size, and the amount of memory left, requirements minimum 256MB.
  5. Use command `disk` check whether the disk meets the requirements, requirements minimum 1 GB.
  6. Use command `dns` Check whether the node domain name resolution  is normal.
  7. Use command `runtime `  Check whether the node container runtime  is installed, can use parameter `--runtime` to set container runtime,  default is docker
  8. Use command `network `  check whether the node can communicate with the endpoint on the cloud,    can use parameter `--ip` to set test ip, default to ping cloudcore.
  9. Use command `pid ` check if the current number of processes in the environment is too many. If the number of available processes is less than 5%, the number of processes is considered insufficient.

`keadm debug collect`

- What is it?

  - This command will obtain all the data of the current node as  `edge-$date.tar.gz`, and provide it to the operation and maintenance personnel to locate the problem.

- What shall be its scope ?

  1. system data

    - Hardware architecture

      Collect arch command output and determine the type of  installation

    - CPU information

      Parse the /proc/cpuinfo file and output the cpu information file

    - Memory information

      Collect free -h command output

    - Hard disk information

      Collect df -h command output, and mount command output

    - Internet Information

      Collect netstat -anp command output and copy /etc/resolv.conf and /etc/hosts files

    - Process information

      Collect ps -aux command output

    - Time information

      Collect date and uptime command output

    - History command input

      Collect all the commands entered by the current user

  2. Edgecore data

  - database data

    Copy the /var/lib/kubeedge/edgecore.db file

  - log files

    Copy all files under /var/*log*/*kubeedge*/

  - service file

    Copy the edgecore.service files under /lib/systemd/system/

  - software version

  - certificate

    Copy all files under /etc/kubeedge/certs/

  - Edge-Core configuration file in  software

    Copy all files under /etc/kubeedge/config/

  3. Container runtime data

  - runtime version information

  - runtime container information

  - runtime log information

  - runtime configuration and log information

  - runtime image information

`keadm debug get`

- What is it?
  
  - This command will get and format the specified resource`s information from the local database of the edge node.
  
- What shall be its scope ?

  1. Format resource information from the local database, and available resource types:
    - `all`
    - `pod`
    - `node`
    - `service`
    - `secret`
    - `configmap`
    - `endpoint`
  2. Use flag `-n, --namespace=''` to indicate the scope of resource acquisition, if the flag `-A, --all-namespaces` is used, information of the specified resource will be obtained from all ranges
  3. Use flag `-o, --output=''` to indicate output format of the information . support formats `yaml`,`json` and `wide`.
  4. Use flag `-l` to indicate which specified field is used to filter the data in the range
  5. Use flag `-db, --db-path''` to specify the database path,default is `/var/lib/kubeedge/edgecore.db`
