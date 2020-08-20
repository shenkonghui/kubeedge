package debug

import (
	"fmt"

	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	constant "github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/common"
	types "github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/common"
	"github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/util"
	"github.com/spf13/cobra"
)

var (
	edgeCheckLongDescription = `Obtain all the data of the current node, and then provide it to the operation
and maintenance personnel to locate the problem`

	edgeCheckShortDescription = `Check specific information.`
	edgeCheckExample          = `
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
`
)

type CheckObject types.CheckObject

// NewEdgecheck returns KubeEdge edge check command.
func NewEdgeCheck(out io.Writer, collectOptions *types.CheckOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check",
		Short:   edgeCheckShortDescription,
		Long:    edgeCheckLongDescription,
		Example: edgeCheckExample,
	}
	for _, v := range constant.CheckObjectMap {
		cmd.AddCommand(NewSubEdgeCheck(out, CheckObject(v)))
	}
	return cmd
}

// NewEdgecheck returns KubeEdge edge check subcommand.
func NewSubEdgeCheck(out io.Writer, object CheckObject) *cobra.Command {
	co := NewCheckOptins()
	cmd := &cobra.Command{
		Short: object.Desc,
		Use:   object.Use,
		RunE: func(cmd *cobra.Command, args []string) error {
			return object.ExecuteCheck(object.Use, co)
		},
	}
	switch object.Use {
	case constant.ArgCheckAll:
		cmd.Flags().StringVarP(&co.Domain, "domain", "d", co.Domain, "specify test domain")
		cmd.Flags().StringVarP(&co.IP, "ip", "i", co.IP, "specify test ip")
		cmd.Flags().StringVarP(&co.EdgeHubURL, "edge-hub-url", "e", co.EdgeHubURL, "specify edgehub url,")
		cmd.Flags().StringVarP(&co.Runtime, "runtime", "r", co.Runtime, "specify test runtime")
	case constant.ArgCheckDNS:
		cmd.Flags().StringVarP(&co.Domain, "domain", "d", co.Domain, "specify test domain")
	case constant.ArgCheckNetwork:
		cmd.Flags().StringVarP(&co.IP, "ip", "i", co.IP, "specify test ip")
		cmd.Flags().StringVarP(&co.EdgeHubURL, "edge-hub-url", "e", co.EdgeHubURL, "specify edgehub url,")
	case constant.ArgCheckRuntime:
		cmd.Flags().StringVarP(&co.Runtime, "runtime", "r", co.Runtime, "specify test runtime")
	}

	return cmd
}

// add flags
func NewCheckOptins() *types.CheckOptions {
	co := &types.CheckOptions{}
	co.Runtime = types.DefaultRuntime
	co.Domain = "www.github.com"
	co.Timeout = 1
	return co
}

//Start to check data
func (co *CheckObject) ExecuteCheck(use string, ob *types.CheckOptions) error {
	err := fmt.Errorf("")

	switch use {
	case constant.ArgCheckAll:
		err = CheckAll(ob)
	case constant.ArgCheckArch:
		err = CheckArch()
	case constant.ArgCheckCPU:
		err = CheckCPU()
	case constant.ArgCheckMemory:
		err = CheckMemory()
	case constant.ArgCheckDisk:
		err = CheckDisk()
	case constant.ArgCheckDNS:
		err = CheckDNS(ob.Domain)
	case constant.ArgCheckNetwork:
		err = CheckNetWork(ob.IP, ob.Timeout, ob.EdgeHubURL)
	case constant.ArgCheckRuntime:
		err = CheckRuntime(ob.Runtime)
	case constant.ArgCheckPID:
		err = CheckPid()
	}

	if err != nil {
		util.PrintFail(use, constant.StrCheck)
	} else {
		util.PrintSuccedd(use, constant.StrCheck)
	}

	return err
}

func CheckAll(ob *types.CheckOptions) error {
	err := CheckArch()
	if err != nil {
		return err
	}

	err = CheckCPU()
	if err != nil {
		return err
	}

	err = CheckMemory()
	if err != nil {
		return err
	}

	err = CheckDisk()
	if err != nil {
		return err
	}

	err = CheckDNS(ob.Domain)
	if err != nil {
		return err
	}

	err = CheckNetWork(ob.IP, ob.Timeout, ob.EdgeHubURL)
	if err != nil {
		return err
	}

	err = CheckPid()
	if err != nil {
		return err
	}

	err = CheckRuntime(ob.Runtime)
	if err != nil {
		return err
	}
	return nil
}

func CheckArch() error {
	o, err := util.ExecShellFilter(constant.CmdGetArch)
	if !util.IsContain(constant.AllowedValueArch, string(o)) {
		return fmt.Errorf("arch not support: %s", string(o))
	}
	fmt.Printf("arch is : %s\n", string(o))
	return err
}

func CheckCPU() error {
	return util.ComparisonSizeWithCmd(constant.CmdGetCPUNum, constant.AllowedValueCPU, constant.ArgCheckCPU, util.UnitCore)
}

func CheckMemory() error {
	return util.ComparisonSizeWithCmd(constant.CmdGetMenorySize, constant.AllowedValueMemory, constant.ArgCheckMemory, util.UnitMB)
}

func CheckDisk() error {
	return util.ComparisonSizeWithCmd(constant.CmdGetDiskSize, constant.AllowedValueDisk, constant.ArgCheckDisk, util.UnitGB)
}

func CheckDNS(domain string) error {
	r, err := net.LookupHost(domain)
	if err != nil {
		return fmt.Errorf("dns resolution failed, domain: %s err: %s", domain, err)
	}
	if len(r) > 0 {
		fmt.Printf("dns resolution success, domain: %s ip: %s\n", domain, r[0])
	} else {
		fmt.Printf("dns resolution success, domain: %s ip: null\n", domain)
	}
	return err
}

func CheckNetWork(IP string, timeout int, edgeHubURL string) error {
	if IP == "" && edgeHubURL == "" {
		result, err := util.ExecShellFilter(constant.CmdGetDNSIP)
		if err != nil {
			return err
		}
		IP = result
	}
	if edgeHubURL != "" {
		err := CheckHTTP(edgeHubURL)
		if err != nil {
			return err
		}
	}
	if IP != "" {
		result, err := util.ExecShellFilter(fmt.Sprintf(constant.CmdPing, IP, timeout))

		if err != nil {
			return err
		}
		if result != "1" {
			return fmt.Errorf("ping %s timeout", IP)
		}
		fmt.Printf("ping %s success\n", IP)
	}
	return nil
}

func CheckHTTP(url string) error {
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: time.Second * 3}
	response, err := httpClient.Get(url)
	if err != nil {
		if !strings.Contains(err.Error(), "x509") {
			return fmt.Errorf("edgehub url connect fail: %s", err.Error())
		}
	} else {
		fmt.Printf("edgehub url connect success: %s\n", url)
		defer response.Body.Close()
	}
	return nil
}

func CheckRuntime(runtime string) error {
	if runtime == types.DefaultRuntime {
		result, err := util.ExecShellFilter(constant.CmdGetStatusDocker)
		if err != nil {
			return err
		}
		if result != "active" {
			return fmt.Errorf("docker is not running: %s", result)
		}
		fmt.Printf("docker is running\n")
	} else {
		return fmt.Errorf("now only support docker: %s", runtime)
		// TODO
	}
	return nil
}

func CheckPid() error {
	rMax, err := util.ExecShellFilter(constant.CmdGetMaxProcessNum)
	if err != nil {
		return err
	}
	r, err := util.ExecShellFilter(constant.CmdGetProcessNum)
	if err != nil {
		return err
	}
	vMax, err := strconv.ParseFloat(rMax, 32)
	v, err := strconv.ParseFloat(r, 32)
	rate := (1 - v/vMax)
	if rate > constant.AllowedValuePIDRate {
		fmt.Printf("Maximum PIDs: %s; Running processes: %s\n", rMax, r)
		return nil
	}
	return fmt.Errorf("Maximum PIDs: %s; Running processes: %s", rMax, r)
}
