---
title: KubeEdge debug (Issue 324)
status: implementable
authors:
    - "@shenkonghui"
    - "@qingchen1203"
approvers:
creation-date: 2020-07-15
last-updated:  2020-07-15
---

# Motivation

Many users shared their feedback that kubeEdge edge nodes do not have good debugging and diagnostic methods, which may prevent people from trying kubeEdge. There should be a simplified method to help operation and maintenance personnel to ensure the stable operation of kubeEdge, so that users can focus more on using it immediately

# Proposal

KubeEdge should have simple commands to debug and troubleshoot edge components.
Therefore, it is recommended to use the following commands during operation and maintenance of KubeEdge.

## Inscope

1. To support first set of basic commands (listed below) to debug edge (node) components in different VM's or hosts.



For edge, commands shall be:

- `keadm debug diagnose`
- `keadm debug collect`
- `keadm debug check`

# Scope of commands

## Design of the commands

**NOTE**: All the below steps are executed as root user, to execute as sudo user. Please add `sudo` infront of all the commands.

### keadm diagnose --help

```
keadm diagnose command can be help to diagnose specific fault scenarios in an all-round way and locate the cause of the fault.

Usage:
  keadm debug diagnose [command]

Examples:

# view the running status of node (key components such as sqlite, edgehub, metamanager, edged and many more)
keadm debug analysis node

Available Commands:
  all           All resource
  node          Troubleshoot the cause of edge node failure with installed software
  pod           Troubleshooting specific container application instances on nodes
  installation  It is same as "keadm check all"

```

### keadm check --help

```
keadm check command can be check whether the system specific items meet the requirements of edgecore installation and operation.

Usage:
  keadm debug check [command]

Available Commands:
  all      Check all
  arch     Determine the node hardware architecture whether support or not
  cpu      Determine if the NUMBER of CPU cores meets the requirement
  memory   Check the system memory size and the amount of memory left
  disk     Check whether the disk meets the requirements
  dns      Check whether the node domain name resolution function is normal
  runtime  Check whether the node Container runtime is normal
  network  Check whether the node can communicate with the endpoint on the cloud
  pid      Check if the current number of processes in the environment is too many. If the number of available processes is less than 5%, the number of processes is considered insufficient
  
Flags:
  -h, --help   help for keadm check

Use "keadm debug check [command] --help" for more information about a command
```



### keadm collect --help

```
Obtain all data of the current node, and then locate and use operation personnel.

Usage:
  keadm debug collect [flags]

Examples:

keadm debug collect --path . 

Flags:
  --path    Cache data and store data compression packages in a directory that defaults to the current directory
  --detail  Whether to print internal log output

```

## Explaining the commands

### Worker Node (at the Edge) commands

`keadm debug diagnose`

- What is it?
  
- This command will be help to diagnose specific fault scenarios in an all-round way and locate the cause of the fault
  
- What shall be its scope ?
    1. Use command `all` can diagnose all resource
    2. Use command `node` can roubleshoot the cause of edge node failure with installed software
       1. check system resources is enought
       2. check container runtime is runting 
       3. check that all edgecore components are running
       4. check that cloudercore can be connected
    3. Use command `pod` can troubleshooting specific container application instances on nodes
       1. check pod Is the configuration correct 
       2. check pod  image can be right pull
       3. check pod schecule 
       4. check pod probe
       5. and many more
    4. Use command `installation` is same as "keadm check all"

`keadm debug check`

- What is it?
  
  - This command will be check whether the system specific items meet the requirements of edgecore installation and operation.
  
- What shall be its scope ?

  1. Check items include hardware resources or operating system resources (cpu, memory, disk, network, pid limit,etc.)
2. Use command `arch` can check node hardware architecture:
  
   - x86_64 architecture
       Ubuntu 16.04 LTS (Xenial Xerus), Ubuntu 18.04 LTS (Bionic Beaver), CentOS 7.x and RHEL 7.x, Galaxy Kylin 4.0.2, ZTE new fulcrum v5.5, winning the bid Kylin v7.0
  
   - armv7i (arm32) architecture
       Raspbian GNU/Linux 9 (stretch)
  
   - aarch64 (arm64) architecture
       Ubuntu 18.04.2 LTS (Bionic Beaver)
  3. Use command `cpu` can cetermine if the NUMBER of CPU cores meets the requirement, minimum 1Vcores.
4. Use command `memory` check the system memory size and the amount of memory left, requirements minimum 256MB.
  5. Use command `disk` check whether the disk meets the requirements, requirements minimum 1 GB.
6. Use command `dns` Check whether the node domain name resolution function is normal.
  7. Use command `runtime `  Check whether the node container runtime function is installed, can use parameter `--runtime` to set container runtime,  default is docker
8. Use command `network `  check whether the node can communicate with the endpoint on the cloud,    can use parameter `--ip` to set test ip, default to ping clusterdns.
  11. Use command `pid ` check if the current number of processes in the environment is too many. If the number of available processes is less than 5%, the number of processes is considered insufficient.

`keadm debug collect`

- What is it?

  - This command will be obtain all related data of the current node, and then locate and use  operation personnel.

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

    Copy the edgecore.service, edgelogger.service, edgemonitor.service, edgedaemon.service files under /lib/systemd/system/

  - software version

  - certificate

    Copy all files under /etc/kubeedge/certs/

  - Edge-Core configuration file in  software (including Edge-daemon)

    Copy all files under /etc/kubeedge/config/

  3. Container runtime data

  - runtime version information

  - runtime container information

  - runtime log information

  - runtime container information

  - runtime configuration and log information

  - runtime image information


