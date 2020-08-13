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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/scheme"
	//v1 "k8s.io/kubernetes/pkg/apis/core"
	k8sprinters "k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	"k8s.io/kubernetes/pkg/printers/storage"

	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/dao"
	edgecoreCfg "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

var (
	debugGetLong = `
Prints a table of the most important information about the specified resource from the local database of the edge node.`
	debugGetShort = `
Get and format data of available resource types in the local database of the edge node.`
	debugGetExample = `
# List all pod
keadm debug get pod -A

# List all pod in namespace test
keadm debug get pod -n test

# List a single configmap  with specified NAME
keadm debug get configmap web -n default

# List the complete information of the configmap with the specified name in the yaml output format
keadm debug get configmap web -n default -o yaml

# List the complete information of all available resources of edge nodes using the specified format (default: yaml)
keadm debug get all -o yaml`

	// allowedFormats Currently supports formats such as yaml|json|wide
	allowedFormats = []string{"yaml", "json", "wide"}

	// availableResources Convert flag to currently supports available Resource types in EdgeCore database.
	availableResources = map[string]string{
		"all":        "'pod','service','secret','configmap','endpoints'",
		"po":         "'pod'",
		"pod":        "'pod'",
		"pods":       "'pod'",
		"svc":        "'service'",
		"service":    "'service'",
		"services":   "'service'",
		"secret":     "'secret'",
		"secrets":    "'secret'",
		"cm":         "'configmap'",
		"configmap":  "'configmap'",
		"configmaps": "'configmap'",
		"ep":         "'endpoints'",
		"endpoint":   "'endpoints'",
		"endpoints":  "'endpoints'",
	}
)

// GetOptions contains the input to the get command.
type GetOptions struct {
	AllNamespace  bool
	Namespace     string
	OutputFormat  string
	LabelSelector string
	DataPath      string

	PrintFlags *PrintFlags
}

// NewCmdDebugGet returns keadm debug get command.
func NewCmdDebugGet(out io.Writer, getOption *GetOptions) *cobra.Command {
	if getOption == nil {
		getOption = NewGetOptions()
	}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   debugGetShort,
		Long:    debugGetLong,
		Example: debugGetExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := getOption.Validate(args)
			if err != nil {
				return err
			}
			return Execute(getOption, args, out)
		},
	}
	addGetOtherFlags(cmd, getOption)

	return cmd
}

func addGetOtherFlags(cmd *cobra.Command, getOption *GetOptions) {
	cmd.Flags().StringVarP(&getOption.Namespace, "namespace", "n", getOption.Namespace, "List the requested object(s) in specified namespaces")
	cmd.Flags().StringVarP(&getOption.OutputFormat, "output", "o", getOption.OutputFormat, "Indicate the output format. Currently supports formats such as yaml|json|wide")
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

//Validate checks the set of flags provided by the user.
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
	if !FileExists(g.DataPath) {
		return fmt.Errorf("EdgeCore database file %v not exist. ", g.DataPath)
	}

	err := InitDB(edgecoreCfg.DataBaseDriverName, edgecoreCfg.DataBaseAliasName, g.DataPath)
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v ", err)
	}
	if len(g.OutputFormat) > 0 {
		g.OutputFormat = strings.ToLower(g.OutputFormat)
		if !IsAllowedFormat(g.OutputFormat) {
			return fmt.Errorf("OutputFormat %v not supportted. Currently supports formats such as yaml|json|wide", g.OutputFormat)
		}
	}
	g.PrintFlags.OutputFormat = &g.OutputFormat

	if args[0] == "all" && len(args) >= 2 {
		return fmt.Errorf("you must specify only one resource")
	}

	return nil
}

// Execute performs the get operation.
func Execute(opts *GetOptions, args []string, out io.Writer) error {
	//var printer printers.ResourcePrinter
	resType := args[0]
	resNames := args[1:]
	results, err := QueryMetaFromDatabase(opts.AllNamespace, opts.Namespace, resType, resNames)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		fmt.Fprintf(out, "No resources found in %v namespace.\n", opts.Namespace)
	}
	results, err = FilterSelector(results, opts.LabelSelector)
	if err != nil {
		return err
	}
	podList, serviceList, secretList, configMapList, endPointsList, err := ParseMetaToList(results)
	if err != nil {
		return err
	}

	opts.PrintFlags.OutputFormat = &opts.OutputFormat
	printer, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	printer, err = printers.NewTypeSetter(scheme.Scheme).WrapToPrinter(printer, nil)
	if opts.OutputFormat == "" || opts.OutputFormat == "wide" {
		if podList != nil {
			talbe, err := ConvertDataToTable(podList)
			if err != nil {
				return err
			}
			printer.PrintObj(talbe, out)
		}
		if serviceList != nil {
			talbe, err := ConvertDataToTable(serviceList)
			if err != nil {
				return err
			}
			printer.PrintObj(talbe, out)
		}
		if secretList != nil {
			talbe, err := ConvertDataToTable(secretList)
			if err != nil {
				return err
			}
			printer.PrintObj(talbe, out)
		}
		if configMapList != nil {
			talbe, err := ConvertDataToTable(configMapList)
			if err != nil {
				return err
			}
			printer.PrintObj(talbe, out)
		}
		if endPointsList != nil {
			talbe, err := ConvertDataToTable(endPointsList)
			if err != nil {
				return err
			}
			printer.PrintObj(talbe, out)
		}
		return nil
	}

	var obj runtime.Object
	if len(podList.Items) != 1 {
		list := v1.List{
			TypeMeta: metav1.TypeMeta{
				Kind:       "List",
				APIVersion: "v1",
			},
			ListMeta: metav1.ListMeta{},
		}
		for _, info := range podList.Items {
			o := info.DeepCopyObject()
			list.Items = append(list.Items, runtime.RawExtension{Object: o})
		}

		listData, err := json.Marshal(list)
		if err != nil {
			return err
		}

		converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, listData)
		if err != nil {
			return err
		}
		obj = converted
	} else {
		obj = podList.Items[0].DeepCopyObject()
	}

	isList := meta.IsListType(obj)
	if isList {
		items, err := meta.ExtractList(obj)
		if err != nil {
			return err
		}

		// take the items and create a new list for display
		list := &unstructured.UnstructuredList{
			Object: map[string]interface{}{
				"kind":       "List",
				"apiVersion": "v1",
				"metadata":   map[string]interface{}{},
			},
		}
		if listMeta, err := meta.ListAccessor(obj); err == nil {
			list.Object["metadata"] = map[string]interface{}{
				"selfLink":        listMeta.GetSelfLink(),
				"resourceVersion": listMeta.GetResourceVersion(),
			}
		}

		for _, item := range items {
			list.Items = append(list.Items, *item.(*unstructured.Unstructured))
		}
		printer.PrintObj(list, out)
	} else {
		var ib map[string]interface{}
		b, _ := json.Marshal(obj)
		json.Unmarshal(b, &ib)
		printer.PrintObj(&unstructured.Unstructured{Object: ib}, out)
	}

	return nil
}

func ConvertDataToTable(obj runtime.Object) (runtime.Object, error) {
	to := metav1.TableOptions{}
	tc := storage.TableConvertor{TableGenerator: k8sprinters.NewTableGenerator().With(printersinternal.AddHandlers)}

	return tc.ConvertToTable(context.TODO(), obj, &to)
}

// ParseMetaToList
func ParseMetaToList(results []dao.Meta) (*v1.PodList, *v1.ServiceList, *v1.SecretList, *v1.ConfigMapList, *v1.EndpointsList, error) {
	podList := &v1.PodList{}
	serviceList := &v1.ServiceList{}
	secretList := &v1.SecretList{}
	configMapList := &v1.ConfigMapList{}
	endPointsList := &v1.EndpointsList{}
	value := make(map[string]interface{})

	for _, v := range results {
		switch v.Type {
		case "pod":
			pod := v1.Pod{}

			if err := json.Unmarshal([]byte(v.Value), &value); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			metadata, err := json.Marshal(value["metadata"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			spec, err := json.Marshal(value["spec"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			status, err := json.Marshal(value["status"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(metadata, &pod.ObjectMeta); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(spec, &pod.Spec); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(status, &pod.Status); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			pod.APIVersion = "v1"
			pod.Kind = v.Type
			pod.Name = v.Key
			podList.Items = append(podList.Items, pod)

		case "service":
			svc := v1.Service{}
			if err := json.Unmarshal([]byte(v.Value), &value); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			metadata, err := json.Marshal(value["metadata"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			spec, err := json.Marshal(value["spec"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			status, err := json.Marshal(value["status"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(metadata, &svc.ObjectMeta); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(spec, &svc.Spec); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(status, &svc.Status); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			svc.APIVersion = "v1"
			svc.Kind = v.Type
			svc.Name = v.Key
			serviceList.Items = append(serviceList.Items, svc)
		case "secret":
			secret := v1.Secret{}
			if err := json.Unmarshal([]byte(v.Value), &value); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			metadata, err := json.Marshal(value["metadata"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			data, err := json.Marshal(value["data"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			typeTmp, err := json.Marshal(value["type"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(metadata, &secret.ObjectMeta); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(data, &secret.Data); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(typeTmp, &secret.Type); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			secret.APIVersion = "v1"
			secret.Kind = v.Type
			secret.Name = v.Key
			secretList.Items = append(secretList.Items, secret)
		case "configmap":
			cmp := v1.ConfigMap{}
			if err := json.Unmarshal([]byte(v.Value), &value); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			metadata, err := json.Marshal(value["metadata"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			data, err := json.Marshal(value["data"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(metadata, &cmp.ObjectMeta); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(data, &cmp.Data); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			cmp.APIVersion = "v1"
			cmp.Kind = v.Type
			cmp.Name = v.Key
			configMapList.Items = append(configMapList.Items, cmp)
		case "endpoints":
			ep := v1.Endpoints{}
			if err := json.Unmarshal([]byte(v.Value), &value); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			metadata, err := json.Marshal(value["metadata"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			subsets, err := json.Marshal(value["subsets"])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(metadata, &ep.ObjectMeta); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			if err := json.Unmarshal(subsets, &ep.Subsets); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			ep.APIVersion = "v1"
			ep.Kind = v.Type
			ep.Name = v.Key
			endPointsList.Items = append(endPointsList.Items, ep)
		default:
			return nil, nil, nil, nil, nil, fmt.Errorf("Parsing failed, unrecognized type: %v. ", v.Type)
		}
	}

	return podList, serviceList, secretList, configMapList, endPointsList, nil
}

// FileIsExist check file is exist
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil || !os.IsNotExist(err)
}

// IsAllowedFormat verification support format
func IsAllowedFormat(oFormat string) bool {
	for _, aFormat := range allowedFormats {
		if oFormat == aFormat {
			return true
		}
	}

	return false
}

// IsAvailableResources verification support resource type
func IsAvailableResources(rsT string) bool {
	_, ok := availableResources[rsT]
	return ok
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

// QueryMetaFromDatabase Filter data from the database based on conditions
func QueryMetaFromDatabase(isAllNamespace bool, resNamePaces string, resType string, resNames []string) ([]dao.Meta, error) {
	var results []dao.Meta

	if isAllNamespace {
		if resType == "all" || len(resNames) == 0 {
			results, err := dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.type in (%v)",
					dao.MetaTableName,
					dao.MetaTableName,
					availableResources[resType]))
			if err != nil {
				return nil, err
			}

			return results, nil
		}
		for _, resName := range resNames {
			result, err := dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.key like '%%/%v/%v'",
					dao.MetaTableName,
					dao.MetaTableName,
					strings.ReplaceAll(availableResources[resType], "'", ""),
					resName))
			if err != nil {
				return nil, err
			}
			results = append(results, result...)
		}

		return results, nil
	}
	if resType == "all" || len(resNames) == 0 {
		results, err := dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.key like '%v/%%' and  %v.type in (%v)",
				dao.MetaTableName,
				dao.MetaTableName,
				resNamePaces,
				dao.MetaTableName,
				availableResources[resType]))
		if err != nil {
			return nil, err
		}

		return results, nil
	}
	for _, resName := range resNames {
		result, err := dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.key = '%v/%v/%v'",
				dao.MetaTableName,
				dao.MetaTableName,
				resNamePaces,
				strings.ReplaceAll(availableResources[resType], "'", ""),
				resName))
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}

	return results, nil
}

// FilterSelector Filter data that meets the selector
func FilterSelector(data []dao.Meta, selector string) ([]dao.Meta, error) {
	var results []dao.Meta
	var jsonValue = make(map[string]interface{})

	sLabels, err := SplitSelectorParameters(selector)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		err := json.Unmarshal([]byte(v.Value), &jsonValue)
		if err != nil {
			return nil, err
		}
		vLabel := jsonValue["metadata"].(map[string]interface{})["labels"]
		if vLabel == nil {
			results = append(results, v)
			continue
		}
		flag := true
		for _, sl := range sLabels {
			if !sl.Exist {
				flag = flag && vLabel.(map[string]interface{})[sl.Key] != sl.Value
				continue
			}
			flag = flag && (vLabel.(map[string]interface{})[sl.Key] == sl.Value)

		}

		if flag {
			results = append(results, v)
		}

	}

	return results, nil
}

// Selector
type Selector struct {
	Key   string
	Value string
	Exist bool
}

// SplitSelectorParameters Split selector args (flag: -l)
func SplitSelectorParameters(args string) ([]Selector, error) {
	var results = make([]Selector, 0)
	var sel Selector
	labels := strings.Split(args, ",")
	for _, label := range labels {
		if strings.Contains(label, "==") {
			labs := strings.Split(label, "==")
			if len(labs) != 2 {
				return nil, fmt.Errorf("arguments in selector form may not have more than one \"==\". ")
			}
			sel.Key = labs[0]
			sel.Value = labs[1]
			sel.Exist = true
			results = append(results, sel)
			continue
		}
		if strings.Contains(label, "!=") {
			labs := strings.Split(label, "!=")
			if len(labs) != 2 {
				return nil, fmt.Errorf("arguments in selector form may not have more than one \"!=\". ")
			}
			sel.Key = labs[0]
			sel.Value = labs[1]
			sel.Exist = false
			results = append(results, sel)
			continue
		}
		if strings.Contains(label, "=") {
			labs := strings.Split(label, "=")
			if len(labs) != 2 {
				return nil, fmt.Errorf("arguments in selector may not have more than one \"=\". ")
			}
			sel.Key = labs[0]
			sel.Value = labs[1]
			sel.Exist = true
			results = append(results, sel)
		}

	}
	return results, nil

}
