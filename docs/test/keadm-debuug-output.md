
## keadm debug check
```
[root@localhost bin]# ./keadm debug check 
Obtain all the data of the current node, and then provide it to the operation
and maintenance personnel to locate the problem

Usage:
  keadm debug check [command]

Examples:

        # Check all items .
        keadm debug check all

        # Check whether the node arch is supported .
        keadm debug check arch

        # Check whether the node CPU meets  requirements.
        keadm debug check cpu

        # Check whether the node memory meets  requirements.
        keadm debug check mem

        # check whether the node disk meets  requirements.
        keadm debug check disk

        # Check whether the node DNS can resolve a specific domain name.
        keadm debug check dns -d www.github.com

        # Check whether the node network meets requirements.
        keadm debug check network

        # Check whether the number of free processes on the node meets requirements.
        keadm debug check pid

        # Check whether runtime(Docker) is installed on the node.
        keadm debug check runtime


Available Commands:
  all         Check all item
  arch        Check whether the architecture can work
  cpu         Check node CPU requirements
  disk        Check node disk requirements
  dns         Check whether DNS can work
  mem         Check node memory requirements
  network     Check whether the network is normal
  pid         Check node PID requirements
  runtime     Check whether runtime can work

Flags:
  -h, --help   help for check

Use "keadm debug check [command] --help" for more information about a command.
========================================================
[root@localhost bin]# ./keadm debug check  arch
arch is : x86_64

|------------------|
|check arch succeed|
|------------------|
========================================================
[root@localhost bin]# ./keadm debug check  cpu
CPU total: 1 core, Allowed > 1 core
CPU usage rate: 0.00, Allowed rate < 0.9

|-----------------|
|check cpu succeed|
|-----------------|
========================================================
[root@localhost bin]# ./keadm debug check  disk
Disk total: 50268.47 MB, Allowed > 1024 MB
Disk Free total: 39572.15 MB, Allowed > 512MB
Disk usage rate: 0.16, Allowed rate < 0.9

|------------------|
|check disk succeed|
|------------------|
========================================================
[root@localhost bin]# ./keadm debug check  dns
dns resolution success, domain: www.github.com ip: 13.250.177.223

|-----------------|
|check dns succeed|
|-----------------|
========================================================
[root@localhost bin]# ./keadm debug check  dns -h
Check whether DNS can work

Usage:
  keadm debug check dns [flags]

Flags:
  -D, --dns-ip string   specify test dns ip
  -d, --domain string   specify test domain (default "www.github.com")
  -h, --help            help for dns
========================================================
[root@localhost bin]# ./keadm debug check  dns -D 114.114.114.114 -d www.baidu.com
dns resolution success, domain: www.baidu.com ip: 61.135.169.121

|-----------------|
|check dns succeed|
|-----------------|
========================================================
[root@localhost bin]# ./keadm debug check mem
Memory total: 1833.33 MB, Allowed > 256 MB
Memory Free total: 998.42 MB, Allowed > 128 MB
Memory usage rate: 0.13, Allowed rate < 0.9

|-----------------|
|check mem succeed|
|-----------------|
========================================================
[root@localhost bin]# ./keadm debug check network
ping 8.8.8.8 success

|---------------------|
|check network succeed|
|---------------------|
========================================================
[root@localhost bin]# ./keadm debug check network -h
Check whether the network is normal

Usage:
  keadm debug check network [flags]

Flags:
  -e, --edge-hub-url string   specify edgehub url,
  -h, --help                  help for network
  -i, --ip string             specify test ip
========================================================
[root@localhost bin]# ./keadm debug check network -i www.baidu.com
ping www.baidu.com success

|---------------------|
|check network succeed|
|---------------------|
========================================================
[root@localhost bin]# ./keadm debug check pid
Maximum PIDs: 32768; Running processes: 106

|-----------------|
|check pid succeed|
|-----------------|
========================================================

[root@localhost bin]# ./keadm debug check runtime
docker is running

|---------------------|
|check runtime succeed|
|---------------------|
========================================================
[root@localhost bin]# ./keadm debug check all
arch is : x86_64
CPU total: 1 core, Allowed > 1 core
CPU usage rate: 0.01, Allowed rate < 0.9
Memory total: 1833.33 MB, Allowed > 256 MB
Memory Free total: 989.46 MB, Allowed > 128 MB
Memory usage rate: 0.13, Allowed rate < 0.9
Disk total: 50268.47 MB, Allowed > 1024 MB
Disk Free total: 39572.16 MB, Allowed > 512MB
Disk usage rate: 0.16, Allowed rate < 0.9
dns resolution success, domain: www.github.com ip: 13.229.188.59
ping 8.8.8.8 success
Maximum PIDs: 32768; Running processes: 105
docker is running

|-----------------|
|check all succeed|
|-----------------|
```

## keadm debug diagnose 
```
[root@localhost bin]# ./keadm debug diagnose node
edgecore is running
edge config is exists: /etc/kubeedge/config/edgecore.yaml
docker is running
dataSource is exists: /var/lib/kubeedge/edgecore.db
cloudcore websocket connection success
|---------------------|
|diagnose node succeed|
|---------------------|


========================================================
[root@localhost bin]# ./keadm debug diagnose install  -h
Diagnose install

Usage:
  keadm debug diagnose install [flags]

Flags:
  -D, --dns-ip string         specify test dns server ip
  -d, --domain string         specify test domain
  -e, --edge-hub-url string   specify edgehub url,
  -h, --help                  help for install
  -i, --ip string             specify test ip
  -r, --runtime string        specify the runtime (default "docker")
========================================================
[root@localhost bin]# ./keadm debug diagnose install -d baidu.com
arch is : x86_64
CPU total: 1 core, Allowed > 1 core
CPU usage rate: 0.05, Allowed rate < 0.9
Memory total: 1833.33 MB, Allowed > 256 MB
Memory Free total: 454.20 MB, Allowed > 128 MB
Memory usage rate: 0.19, Allowed rate < 0.9
Disk total: 50268.47 MB, Allowed > 1024 MB
Disk Free total: 39571.23 MB, Allowed > 512MB
Disk usage rate: 0.16, Allowed rate < 0.9
dns resolution success, domain: baidu.com ip: 220.181.38.148
ping 172.17.0.1 success
Maximum PIDs: 32768; Running processes: 124
docker is running

|------------------------|
|diagnose install succeed|
|------------------------|
========================================================
[root@localhost ~]# keadm debug diagnose pod nginx-ds-85jch -n test
edgecore is running
edge config is exists: /etc/kubeedge/config/edgecore.yaml
docker is running
dataSource is exists: /var/lib/kubeedge/edgecore.db
cloudcore websocket connection successDatabase /var/lib/kubeedge/edgecore.db is exist 
Pod nginx-ds-85jch is exist 
PodStatus nginx-ds-85jch is exist 
pod nginx-ds-85jch phase is Running 
containerConditions nginx-ds is ready
Pod nginx-ds-85jch is Ready
+-------------------+
|diagnose pod succeed.|
+-------------------+
=======================no find pod in datebase=============================
[root@localhost ~]# keadm debug diagnose pod xxx
edgecore is running
edge config is exists: /etc/kubeedge/config/edgecore.yaml
docker is running
dataSource is exists: /var/lib/kubeedge/edgecore.db
cloudcore websocket connection successDatabase /var/lib/kubeedge/edgecore.db is exist 
not find default/pod/xxx in datebase

+-------------------+
|diagnose pod failed|
+-------------------+
=======================pod error command=============================
[root@localhost ~]# keadm debug diagnose pod nginx-deployment-dbbffc676-6r6f5
edgecore is running
edge config is exists: /etc/kubeedge/config/edgecore.yaml
docker is running
dataSource is exists: /var/lib/kubeedge/edgecore.db
cloudcore websocket connection successDatabase /var/lib/kubeedge/edgecore.db is exist 
Pod nginx-deployment-dbbffc676-6r6f5 is exist 
PodStatus nginx-deployment-dbbffc676-6r6f5 is exist 
pod nginx-deployment-dbbffc676-6r6f5 phase is Running 
conditions is not true, type: Ready ,message: containers with unready status: [nginx] ,reason: ContainersNotReady 
containerConditions nginx Terminated, message: oci runtime error: container_linux.go:235: starting container process caused "exec: \"/abc\": stat /abc: no such file or directory"
, reason: ContainerCannotRun, RestartCount: 101 
Pod nginx-deployment-dbbffc676-6r6f5 is not Ready

+-------------------+
|diagnose pod failed|
+-------------------+
```
# keadm debug collect 
```
[root@localhost ~]# keadm debug collect -h
Obtain all the data of the current node, and then provide it to the operation
and maintenance personnel to locate the problem

Usage:
  keadm debug collect [flags]

Examples:

# Check all items and specified as the current directory
keadm debug collect --output-path .


Flags:
  -c, --config string        Specify configuration file, defalut is /etc/kubeedge/config/edgecore.yaml (default "/etc/kubeedge/config/edgecore.yaml")
  -d, --detail               Whether to print internal log output
  -h, --help                 help for collect
  -l, --log-path string      Specify log file (default "/var/log/kubeedge/")
  -o, --output-path string   Cache data and store data compression packages in a directory that default to the current directory (default ".")
========================================================
[root@localhost ~]# keadm debug collect -d
Start collecting data
I1103 19:23:44.717381   19135 collect.go:90] create tmp file: /tmp/edge_2020_1103_192344
I1103 19:23:44.717772   19135 collect.go:162] create tmp file: /tmp/edge_2020_1103_192344/system
I1103 19:23:44.717829   19135 collect.go:252] Execute Shell: arch > /tmp/edge_2020_1103_192344/system/arch
I1103 19:23:44.721803   19135 collect.go:241] Copy File: cp -r /proc/cpuinfo /tmp/edge_2020_1103_192344/system/
I1103 19:23:44.724550   19135 collect.go:241] Copy File: cp -r /proc/meminfo /tmp/edge_2020_1103_192344/system/
I1103 19:23:44.727088   19135 collect.go:252] Execute Shell: df -h > /tmp/edge_2020_1103_192344/system/disk
I1103 19:23:44.733329   19135 collect.go:241] Copy File: cp -r /etc/hosts /tmp/edge_2020_1103_192344/system/
I1103 19:23:44.735910   19135 collect.go:241] Copy File: cp -r /etc/resolv.conf /tmp/edge_2020_1103_192344/system/
I1103 19:23:44.738109   19135 collect.go:252] Execute Shell: ps -axu > /tmp/edge_2020_1103_192344/system/process
I1103 19:23:44.757019   19135 collect.go:252] Execute Shell: date > /tmp/edge_2020_1103_192344/system/date
I1103 19:23:44.759452   19135 collect.go:252] Execute Shell: uptime > /tmp/edge_2020_1103_192344/system/uptime
I1103 19:23:44.763852   19135 collect.go:252] Execute Shell: history -a && cat ~/.bash_history  > /tmp/edge_2020_1103_192344/system/history
I1103 19:23:44.767664   19135 collect.go:252] Execute Shell: netstat -pan > /tmp/edge_2020_1103_192344/system/network
I1103 19:23:44.794834   19135 collect.go:93] collect systemd data finish
I1103 19:23:44.796511   19135 util.go:397] Looking for default routes with IPv4 addresses
I1103 19:23:44.796520   19135 util.go:402] Default route transits interface "eth0"
I1103 19:23:44.796711   19135 util.go:212] Interface eth0 is up
I1103 19:23:44.796763   19135 util.go:259] Interface "eth0" has 3 addresses :[10.211.55.10/8 fdb2:2c26:f4e4:0:21c:42ff:fe3c:fcf1/64 fe80::21c:42ff:fe3c:fcf1/64].
I1103 19:23:44.796790   19135 util.go:228] Checking addr  10.211.55.10/8.
I1103 19:23:44.796796   19135 util.go:235] IP found 10.211.55.10
I1103 19:23:44.796804   19135 util.go:265] Found valid IPv4 address 10.211.55.10 for interface "eth0".
I1103 19:23:44.796808   19135 util.go:408] Found active IP 10.211.55.10 
I1103 19:23:44.797831   19135 collect.go:191] create tmp file: /tmp/edge_2020_1103_192344/edgecore
I1103 19:23:44.798466   19135 collect.go:241] Copy File: cp -r /var/lib/kubeedge/edgecore.db /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.801488   19135 collect.go:241] Copy File: cp -r /var/log/kubeedge/ /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.803945   19135 collect.go:241] Copy File: cp -r /lib/systemd/system/edgecore.service /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.805966   19135 collect.go:241] Copy File: cp -r /etc/kubeedge/config/ /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.807969   19135 collect.go:241] Copy File: cp -r /etc/kubeedge/certs/server.crt /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.809956   19135 collect.go:241] Copy File: cp -r /etc/kubeedge/certs/server.key /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.812697   19135 collect.go:241] Copy File: cp -r /etc/kubeedge/ca/rootCA.crt /tmp/edge_2020_1103_192344/edgecore/
I1103 19:23:44.815439   19135 collect.go:252] Execute Shell: edgecore  --version > /tmp/edge_2020_1103_192344/edgecore/version
I1103 19:23:44.873923   19135 collect.go:101] collect edgecore data finish
I1103 19:23:44.873938   19135 collect.go:228] create tmp file: /tmp/edge_2020_1103_192344/runtime
I1103 19:23:44.874008   19135 collect.go:241] Copy File: cp -r /lib/systemd/system/docker.service /tmp/edge_2020_1103_192344/runtime/
I1103 19:23:44.876349   19135 collect.go:252] Execute Shell: docker version > /tmp/edge_2020_1103_192344/runtime/version
I1103 19:23:44.905132   19135 collect.go:252] Execute Shell: docker info > /tmp/edge_2020_1103_192344/runtime/info
I1103 19:23:44.933076   19135 collect.go:252] Execute Shell: docker images > /tmp/edge_2020_1103_192344/runtime/images
I1103 19:23:44.955286   19135 collect.go:252] Execute Shell: docker ps -a > /tmp/edge_2020_1103_192344/runtime/containerInfo
I1103 19:23:44.973366   19135 collect.go:252] Execute Shell: journalctl -u docker  > /tmp/edge_2020_1103_192344/runtime/log
I1103 19:23:45.688942   19135 collect.go:106] collect runtime data finish
I1103 19:23:45.784169   19135 collect.go:116] compress data finish
collecting data finish, path: /root/edge_2020_1103_192344.tar.gz
========================================================
[root@localhost ~]# tar -xvzf edge_2020_1103_192023.tar.gz  -C ./edge_data
.
edgecore
edgecore/config
edgecore/config/edgecore.yaml
edgecore/edgecore.db
edgecore/edgecore.service
edgecore/kubeedge
edgecore/kubeedge/edgecore.log
edgecore/rootCA.crt
edgecore/server.crt
edgecore/server.key
edgecore/version
runtime
runtime/containerInfo
runtime/docker.service
runtime/images
runtime/info
runtime/log
runtime/version
system
system/arch
system/cpuinfo
system/date
system/disk
system/history
system/hosts
system/meminfo
system/network
system/process
system/resolv.conf
system/uptime
========================================================
```