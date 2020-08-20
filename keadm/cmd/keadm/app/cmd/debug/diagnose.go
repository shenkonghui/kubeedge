package debug

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/astaxie/beego/orm"
	kubeedgeTypes "github.com/kubeedge/kubeedge/common/types"
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/dao"
	constant "github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/common"
	types "github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/common"
	"github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/util"
	edgecoreCfg "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

var (
	edgeDiagnoseLongDescription = `keadm debug diagnose command Diagnose relevant information at edge nodes
`
	edgeDiagnoseShortDescription = `Diagnose relevant information at edge nodes`

	edgeDiagnoseExample = `
# Diagnose whether the node is normal
keadm debug diagnose node

# Diagnose whether the pod is normal
keadm debug diagnose pod nginx-xxx -n test

# Diagnose node installation conditions
keadm debug diagnose install 

# Diagnose node installation conditions and specify the detected ip
keadm debug diagnose install -i 192.168.1.2 
`
)

type Diagnose types.DiagnoseObject

// NewDiagnose returns KubeEdge edge debug Diagnose command.
func NewDiagnose(out io.Writer, diagnoseOptions *types.DiagnoseOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "diagnose",
		Short:   edgeDiagnoseShortDescription,
		Long:    edgeDiagnoseLongDescription,
		Example: edgeDiagnoseExample,
	}
	for _, v := range constant.DiagnoseObjectMap {
		cmd.AddCommand(NewSubDiagnose(out, Diagnose(v)))
	}
	return cmd
}

func NewSubDiagnose(out io.Writer, object Diagnose) *cobra.Command {
	do := NewDiagnoseOptins()
	cmd := &cobra.Command{
		Short: object.Desc,
		Use:   object.Use,
		Run: func(cmd *cobra.Command, args []string) {
			object.ExecuteDiagnose(object.Use, do, args)
		},
	}
	switch object.Use {
	case constant.ArgDiagnoseNode:
		cmd.Flags().StringVarP(&do.CheckOptions.Runtime, "runtime", "r", do.CheckOptions.Runtime, "specify the runtime")
	case constant.ArgDiagnosePod:
		cmd.Flags().StringVarP(&do.Namespace, "namespace", "n", do.Namespace, "specify namespace")
	case constant.ArgDiagnoseInstall:
		cmd.Flags().StringVarP(&do.CheckOptions.Domain, "domain", "d", do.CheckOptions.Domain, "specify test domain")
		cmd.Flags().StringVarP(&do.CheckOptions.IP, "ip", "i", do.CheckOptions.IP, "specify test ip")
		cmd.Flags().StringVarP(&do.CheckOptions.EdgeHubURL, "edge-hub-url", "e", do.CheckOptions.EdgeHubURL, "specify edgehub url,")
		cmd.Flags().StringVarP(&do.CheckOptions.Runtime, "runtime", "r", do.CheckOptions.Runtime, "specify the runtime")
	}
	return cmd
}

// add flags
func NewDiagnoseOptins() *types.DiagnoseOptions {
	do := &types.DiagnoseOptions{}
	do.Namespace = "default"
	do.CheckOptions = &types.CheckOptions{
		IP:      "",
		Timeout: 1,
		Runtime: types.DefaultRuntime,
	}
	return do
}

func (da Diagnose) ExecuteDiagnose(use string, ops *types.DiagnoseOptions, args []string) {
	err := fmt.Errorf("")
	switch use {
	case constant.ArgDiagnoseNode:
		err = DiagnoseNode(ops)
	case constant.ArgDiagnosePod:
		if len(args) == 0 {
			fmt.Println("error: You must specify a pod name")
			return
		}
		// diagnose Pod, first diagnose node
		err = DiagnoseNode(ops)
		if err == nil {
			err = DiagnosePod(ops, args[0])
		}
	case constant.ArgDiagnoseInstall:
		err = DiagnoseInstall(ops.CheckOptions)
	}

	if err != nil {
		fmt.Println(err.Error())
		util.PrintFail(use, constant.StrDiagnose)
	} else {
		util.PrintSuccedd(use, constant.StrDiagnose)
	}
}

func DiagnoseNode(ops *types.DiagnoseOptions) error {
	osType := util.GetOSInterface()
	isEdgeRuning, err := osType.IsKubeEdgeProcessRunning(util.KubeEdgeBinaryName)
	if err != nil {
		return fmt.Errorf("get edgecore status fail")
	}

	if !isEdgeRuning {
		return fmt.Errorf("edgecore is not running")
	}
	fmt.Println("edgecore is running")

	// need filter current process
	isDockerRuning, err := util.IsProcessRunningWithFilter(ops.CheckOptions.Runtime, "keadm")
	if err != nil {
		return fmt.Errorf("get runtime status fail")
	}

	if !isDockerRuning {
		return fmt.Errorf("runtime is not running")
	}
	fmt.Println("runtime is running")

	return nil
}

func DiagnosePod(ops *types.DiagnoseOptions, podName string) error {
	err := InitDB(edgecoreCfg.DataBaseDriverName, edgecoreCfg.DataBaseAliasName, edgecoreCfg.DataBaseDataSource)
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v ", err)
	}
	fmt.Printf("Database %s is exist \n", edgecoreCfg.DataBaseDataSource)
	podStatus, err := QueryPodFromDatabase(ops.Namespace, podName)
	if err != nil {
		return err
	}

	fmt.Printf("%v phase is %v \n", podName, podStatus.Phase)

	conditions := podStatus.Conditions
	containerConditions := podStatus.ContainerStatuses

	// check conditions
	for _, v := range conditions {
		if v.Status != "True" {
			return fmt.Errorf("%v is not true", v.Type)
		}
	}
	// check containerConditions
	for _, v := range containerConditions {
		if *v.Started {
			return fmt.Errorf("%v is not true", v.Name)
		}
	}

	fmt.Printf("Pod %s is Ready", podName)
	return nil
}

// InitDB Init DB info
func InitDB(driverName, dbName, dataSource string) error {
	if err := orm.RegisterDriver(driverName, orm.DRSqlite); err != nil {
		return fmt.Errorf("Failed to register driver: %v ", err)
	}
	if err := orm.RegisterDataBase(
		dbName,
		driverName,
		dataSource); err != nil {
		return fmt.Errorf("Failed to register db: %v ", err)
	}
	orm.RegisterModel(new(dao.Meta))

	// create orm
	dbm.DBAccess = orm.NewOrm()
	if err := dbm.DBAccess.Using(dbName); err != nil {
		return fmt.Errorf("Using db access error %v ", err)
	}
	return nil
}

func QueryPodFromDatabase(resNamePaces string, podName string) (*v1.PodStatus, error) {
	conditionsPod := fmt.Sprintf("%v/pod/%v",
		resNamePaces,
		podName)
	result, err := dao.QueryMeta("key", conditionsPod)
	if err != nil {
		return nil, fmt.Errorf("read database fail: %s", err.Error())
	}
	if len(*result) == 0 {
		return nil, fmt.Errorf("not find %v in datebase", conditionsPod)
	}
	fmt.Printf("Pod %s is exist \n", podName)

	conditionsStatus := fmt.Sprintf("%v/podstatus/%v",
		resNamePaces,
		podName)
	result, err = dao.QueryMeta("key", conditionsStatus)
	if err != nil {
		return nil, fmt.Errorf("read database fail: %s", err.Error())
	}
	if len(*result) == 0 {
		return nil, fmt.Errorf("not find %v in datebase", conditionsStatus)
	}
	fmt.Printf("PodStatus %s is exist \n", podName)

	r := *result
	podStatus := &kubeedgeTypes.PodStatusRequest{}
	json.Unmarshal([]byte(r[0]), podStatus)
	return &podStatus.Status, nil
}

func DiagnoseInstall(ob *types.CheckOptions) error {
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

	if ob.Domain != "" {
		err = CheckDNS(ob.Domain)
		if err != nil {
			return err
		}
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
