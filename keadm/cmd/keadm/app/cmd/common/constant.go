/*
Copyright 2019 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

const (
	// KubeEdgeVersion sets the version of KubeEdge to be used
	KubeEdgeVersion = "kubeedge-version"

	// KubernetesVersion sets the version of Kuberneted to be used
	KubernetesVersion = "kubernetes-version"

	// KubeConfig sets the path of kubeconfig
	KubeConfig = "kube-config"

	// Master sets the address of K8s master
	Master = "master"

	// CloudCoreIPPort sets the IP and port of KubeEdge cloud component
	CloudCoreIPPort = "cloudcore-ipport"

	// KubeEdge Node unique idenfitcation string
	EdgeNodeName = "edgenode-name"

	// KubeEdge interface name string
	InterfaceName = "interfacename"

	// KubeEdge remote-runtime-endpoint string
	RemoteRuntimeEndpoint = "remote-runtime-endpoint"

	// CertPath sets the path of the certificates generated by the KubeEdge Cloud component
	CertPath = "certPath"

	// DefaultCertPath is the default certificate path in edge node
	DefaultCertPath = "/etc/kubeedge/certs"

	// DefaultK8SMinimumVersion is the minimum version of K8S
	DefaultK8SMinimumVersion = 11

	// DefaultKubeConfig is the default path of kubeconfig
	DefaultKubeConfig = "/root/.kube/config"

	// DefaultProjectID is default project id
	DefaultProjectID = "e632aba927ea4ac2b575ec1603d56f10"

	// RuntimeType is default runtime type
	RuntimeType = "runtimetype"

	// DefaultKubeEdgeVersion is the default KubeEdge version
	DefaultKubeEdgeVersion = "1.3.1"

	// Token sets the token used when edge applying for the certificate
	Token = "token"

	// HttpServer sets the port where to apply for the edge certificate
	CertPort = "certport"

	AdvertiseAddress = "advertise-address"

	TokenSecretName = "tokensecret"

	TokenDataName = "tokendata"

	EdgecoreConfig = "config"

	// Default edgecore config path
	EdgecoreConfigPath = "/etc/kubeedge/config/edgecore.yaml"

	// cmd to copy file
	CmdCopyFile = "cp -r %s %s/"

	/*system info*/
	CmdDiskInfo    = "df -h > %s/disk"
	CmdArchInfo    = "arch > %s/arch"
	CmdProcessInfo = "ps -axu > %s/process"
	CmdDateInfo    = "date > %s/date"
	CmdUptimeInfo  = "uptime > %s/uptime"
	CmdHistorynfo  = "history -a && cat ~/.bash_history  > %s/history"
	CmdNetworkInfo = "netstat -pan > %s/network"

	PathCpuinfo   = "/proc/cpuinfo"
	PathMemory    = "/proc/meminfo"
	PathHosts     = "/etc/hosts"
	PathDNSResolv = "/etc/resolv.conf"

	/*edgecore info*/
	PathEdgecoreService = "/lib/systemd/system/edgecore.service"
	CmdEdgecoreVersion  = "edgecore  --version > %s/version"

	/*runtime info*/
	CmdDockerVersion    = "docker version > %s/version"
	CmdContainerInfo    = "docker ps -a > %s/containerInfo"
	CmdContainerLogInfo = "journalctl -u docker  > %s/log"
	CmdDockerInfo       = "docker info > %s/info"
	CmdDockerImageInfo  = "docker images > %s/images"
	PathDockerService   = "/lib/systemd/system/docker.service"
)
