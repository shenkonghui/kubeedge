/*
Copyright 2020 The KubeEdge Authors.
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

package debug

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/spf13/cobra"

	"github.com/kubeedge/beehive/pkg/common/util"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/dao"
	edgecoreCfg "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

const (
	// DefaultErrorExitCode defines exit the code for failed action generally
	DefaultErrorExitCode = 1
	// ResourceTypeAll defines resource type all
	ResourceTypeAll = "all"
)

var (
	debugGetLong = `
Prints a table of the most important information about the specified resource from the local database of the edge node.`
	debugGetExample = `
# List all pod in namespace test
keadm debug get pod -n test
# List a single configmap  with specified NAME
keadm debug get configmap web -n default
# List the complete information of the configmap with the specified name in the yaml output format
keadm debug get configmap web -n default -o yaml
# List the complete information of all available resources of edge nodes using the specified format (default: yaml)
keadm debug get all -o yaml`

	// availableResources Convert flag to currently supports available Resource types in EdgeCore database.
	availableResources = map[string]string{
		"all":        ResourceTypeAll,
		"po":         model.ResourceTypePod,
		"pod":        model.ResourceTypePod,
		"pods":       model.ResourceTypePod,
		"no":         model.ResourceTypeNode,
		"node":       model.ResourceTypeNode,
		"nodes":      model.ResourceTypeNode,
		"svc":        constants.ResourceTypeService,
		"service":    constants.ResourceTypeService,
		"services":   constants.ResourceTypeService,
		"secret":     model.ResourceTypeSecret,
		"secrets":    model.ResourceTypeSecret,
		"cm":         model.ResourceTypeConfigmap,
		"configmap":  model.ResourceTypeConfigmap,
		"configmaps": model.ResourceTypeConfigmap,
		"ep":         constants.ResourceTypeEndpoints,
		"endpoint":   constants.ResourceTypeEndpoints,
		"endpoints":  constants.ResourceTypeEndpoints,
	}
)

// NewCmdDebugGet returns keadm debug get command.
func NewCmdDebugGet(out io.Writer, getOption *GetOptions) *cobra.Command {
	if getOption == nil {
		getOption = NewGetOptions()
	}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Display one or many resources",
		Long:    debugGetLong,
		Example: debugGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := getOption.Validate(args); err != nil {
				CheckErr(err, fatal)
			}
			if err := getOption.Run(args, out); err != nil {
				CheckErr(err, fatal)
			}
		},
	}
	addGetOtherFlags(cmd, getOption)

	return cmd
}

// fatal prints the message if set and then exits.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}

		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// CheckErr formats a given error as a string and calls the passed handleErr
// func with that string and an exit code.
func CheckErr(err error, handleErr func(string, int)) {
	switch err.(type) {
	case nil:
		return
	default:
		handleErr(err.Error(), DefaultErrorExitCode)
	}
}

// addGetOtherFlags
func addGetOtherFlags(cmd *cobra.Command, getOption *GetOptions) {
	cmd.Flags().StringVarP(&getOption.Namespace, "namespace", "n", getOption.Namespace, "List the requested object(s) in specified namespaces")
	cmd.Flags().StringVarP(getOption.PrintFlags.OutputFormat, "output", "o", *getOption.PrintFlags.OutputFormat, "Indicate the output format. Currently supports formats such as yaml|json|wide")
	cmd.Flags().StringVarP(&getOption.LabelSelector, "selector", "l", getOption.LabelSelector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	cmd.Flags().StringVarP(&getOption.DataPath, "input", "i", getOption.DataPath, "Indicate the edge node database path, the default path is \"/var/lib/kubeedge/edgecore.db\"")
	cmd.Flags().BoolVarP(&getOption.AllNamespace, "all-namespaces", "A", getOption.AllNamespace, "List the requested object(s) across all namespaces")
}

// NewGetOptions returns a GetOptions with default EdgeCore database source.
func NewGetOptions() *GetOptions {
	opts := &GetOptions{
		Namespace:  "default",
		DataPath:   edgecoreCfg.DataBaseDataSource,
		PrintFlags: NewGetPrintFlags(),
	}

	return opts
}

// GetOptions contains the input to the get command.
type GetOptions struct {
	AllNamespace  bool
	Namespace     string
	LabelSelector string
	DataPath      string

	PrintFlags *PrintFlags
}

// Run performs the get operation.
func (g *GetOptions) Run(args []string, out io.Writer) error {
	resType := args[0]
	resNames := args[1:]
	_, err := g.queryDataFromDatabase(availableResources[resType], resNames)
	if err != nil {
		return err
	}
	return nil
}

// IsAllowedFormat verification support format
func (g *GetOptions) IsAllowedFormat(f string) bool {
	allowedFormats := g.PrintFlags.AllowedFormats()
	for _, v := range allowedFormats {
		if f == v {
			return true
		}
	}

	return false
}

// Validate checks the set of flags provided by the user.
func (g *GetOptions) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You must specify the type of resource to get. ")
	}
	if !IsAvailableResources(args[0]) {
		return fmt.Errorf("Unrecognized resource type: %v. ", args[0])
	}
	if len(g.DataPath) == 0 {
		fmt.Printf("Not specified the EdgeCore database path, use the default path: %v. ", g.DataPath)
	}
	if !IsFileExist(g.DataPath) {
		return fmt.Errorf("EdgeCore database file %v not exist. ", g.DataPath)
	}

	if err := InitDB(edgecoreCfg.DataBaseDriverName, edgecoreCfg.DataBaseAliasName, g.DataPath); err != nil {
		return fmt.Errorf("Failed to initialize database: %v ", err)
	}
	if len(*g.PrintFlags.OutputFormat) > 0 {
		format := strings.ToLower(*g.PrintFlags.OutputFormat)
		g.PrintFlags.OutputFormat = &format
		if !g.IsAllowedFormat(*g.PrintFlags.OutputFormat) {
			return fmt.Errorf("Invalid output format: %v, currently supports formats such as yaml|json|wide. ", *g.PrintFlags.OutputFormat)
		}
	}
	if args[0] == ResourceTypeAll && len(args) >= 2 {
		return fmt.Errorf("You must specify only one resource. ")
	}

	return nil
}

func (g *GetOptions) queryDataFromDatabase(resType string, resNames []string) ([]dao.Meta, error) {
	var result []dao.Meta

	switch resType {
	case model.ResourceTypePod:
		pods, err := g.getPodsFromDatabase(g.Namespace, resNames)
		if err != nil {
			return nil, err
		}
		result = append(result, pods...)
	case model.ResourceTypeNode:
		node, err := g.getNodeFromDatabase(g.Namespace, resNames)
		if err != nil {
			return nil, err
		}
		result = append(result, node...)
	case model.ResourceTypeConfigmap, model.ResourceTypeSecret, constants.ResourceTypeEndpoints, constants.ResourceTypeService:
		value, err := g.getSingleResourceFromDatabase(g.Namespace, resNames, resType)
		if err != nil {
			return nil, err
		}
		result = append(result, value...)
	case ResourceTypeAll:
		pods, err := g.getPodsFromDatabase(g.Namespace, resNames)
		if err != nil {
			return nil, err
		}
		result = append(result, pods...)
		node, err := g.getNodeFromDatabase(g.Namespace, resNames)
		if err != nil {
			return nil, err
		}
		result = append(result, node...)

		resTypes := []string{model.ResourceTypeConfigmap, model.ResourceTypeSecret, constants.ResourceTypeEndpoints, constants.ResourceTypeService}
		for _, v := range resTypes {
			value, err := g.getSingleResourceFromDatabase(g.Namespace, resNames, v)
			if err != nil {
				return nil, err
			}
			result = append(result, value...)
		}
	default:
		return nil, fmt.Errorf("Query resource type: %v in namespaces: %v failed. ", resType, g.Namespace)

	}

	return result, nil
}

func (g *GetOptions) getPodsFromDatabase(resNS string, resNames []string) ([]dao.Meta, error) {
	var results []dao.Meta
	podJSON := make(map[string]interface{})
	podStatusJSON := make(map[string]interface{})

	podRecords, err := dao.QueryAllMeta("type", model.ResourceTypePod)
	if err != nil {
		return nil, err
	}
	for _, v := range *podRecords {
		namespaceParsed, _, _, _ := util.ParseResourceEdge(v.Key, model.QueryOperation)
		if namespaceParsed != resNS && !g.AllNamespace {
			continue
		}
		if len(resNames) > 0 && !IsExistName(resNames, v.Key) {
			continue
		}

		podKey := strings.Replace(v.Key, constants.ResourceSep+model.ResourceTypePod+constants.ResourceSep,
			constants.ResourceSep+model.ResourceTypePodStatus+constants.ResourceSep, 1)
		podStatusRecords, err := dao.QueryMeta("key", podKey)
		if err != nil {
			return nil, err
		}
		if len(*podStatusRecords) <= 0 {
			results = append(results, v)
			continue
		}
		if err := json.Unmarshal([]byte(v.Value), &podJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte((*podStatusRecords)[0]), &podStatusJSON); err != nil {
			return nil, err
		}
		podJSON["status"] = podStatusJSON["Status"]
		data, err := json.Marshal(podJSON)
		if err != nil {
			return nil, err
		}
		v.Value = string(data)
		results = append(results, v)
	}

	return results, nil
}

func (g *GetOptions) getNodeFromDatabase(resNS string, resNames []string) ([]dao.Meta, error) {
	var results []dao.Meta
	nodeJSON := make(map[string]interface{})
	nodeStatusJSON := make(map[string]interface{})

	nodeRecords, err := dao.QueryAllMeta("type", model.ResourceTypeNode)
	if err != nil {
		return nil, err
	}
	for _, v := range *nodeRecords {
		namespaceParsed, _, _, _ := util.ParseResourceEdge(v.Key, model.QueryOperation)
		if namespaceParsed != resNS && !g.AllNamespace {
			continue
		}
		if len(resNames) > 0 && !IsExistName(resNames, v.Key) {
			continue
		}

		nodeKey := strings.Replace(v.Key, constants.ResourceSep+model.ResourceTypeNode+constants.ResourceSep,
			constants.ResourceSep+model.ResourceTypeNodeStatus+constants.ResourceSep, 1)
		nodeStatusRecords, err := dao.QueryMeta("key", nodeKey)
		if err != nil {
			return nil, err
		}
		if len(*nodeStatusRecords) <= 0 {
			results = append(results, v)
			continue
		}
		if err := json.Unmarshal([]byte(v.Value), &nodeJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte((*nodeStatusRecords)[0]), &nodeStatusJSON); err != nil {
			return nil, err
		}
		nodeJSON["status"] = nodeStatusJSON["Status"]
		data, err := json.Marshal(nodeJSON)
		if err != nil {
			return nil, err
		}
		v.Value = string(data)
		results = append(results, v)
	}

	return results, nil
}

func (g *GetOptions) getSingleResourceFromDatabase(resNS string, resNames []string, resType string) ([]dao.Meta, error) {
	var results []dao.Meta

	resRecords, err := dao.QueryAllMeta("type", resType)
	if err != nil {
		return nil, err
	}
	for _, v := range *resRecords {
		namespaceParsed, _, _, _ := util.ParseResourceEdge(v.Key, model.QueryOperation)
		if namespaceParsed != resNS && !g.AllNamespace {
			continue
		}
		if len(resNames) > 0 && !IsExistName(resNames, v.Key) {
			continue
		}
		results = append(results, v)
	}

	return results, nil
}

// IsAvailableResources verification support resource type
func IsAvailableResources(rsT string) bool {
	_, ok := availableResources[rsT]
	return ok
}

// IsFileExist check file is exist
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || !os.IsNotExist(err)
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

// IsExistName verify the filed in the resNames exists in the name
func IsExistName(resNames []string, name string) bool {
	value := false
	for _, v := range resNames {
		if strings.Index(name, v) != -1 {
			value = true
		}
	}

	return value
}
