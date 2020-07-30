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
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/dao"
	edgecoreCfg "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	debugGetLongDescription = `
	Prints a table of the most important information about the specified resource from the local database of the edge node
`
	debugGetShortDescription = `Get and format data of available resource types in the local database of the edge node
`
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
	keadm debug get all -o yaml
`
	// allowedFormats Currently supports formats such as yaml|json|wide
	allowedFormats = []string{"yaml", "json", "wide"}

	// availableResources Currently supports available Resource types in EdgeCore database.
	availableResources = []string{
		"all",
		"pod",
		"node",
		"service",
		"secret",
		"configmap",
		"endpoint",
	}
)

// GetOptions contains the input to the get command.
type GetOptions struct {
	AllNamespace  bool
	Namespace     string
	OutputFormat  string
	LabelSelector string
	DataPath      string
}

// NewCmdDebugGet returns keadm debug get command.
func NewCmdDebugGet(out io.Writer, getOption *GetOptions) *cobra.Command {
	if getOption == nil {
		getOption = NewGetOptions()
	}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   debugGetShortDescription,
		Long:    debugGetLongDescription,
		Example: debugGetExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := getOption.Validate(args)
			if err != nil {
				return err
			}
			return Execute(getOption, args)
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
	opts := &GetOptions{}
	opts.DataPath = edgecoreCfg.DataBaseDataSource

	return opts
}

//Validate checks the set of flags provided by the user.
func (g *GetOptions) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You must specify the type of resource to get. ")
	}
	if g.AllNamespace == false && len(g.Namespace) == 0 {
		g.Namespace = "default"
	}
	if len(g.DataPath) == 0 {
		fmt.Printf("Failed to get the EdgeCore database path, will use default path: %v. ", g.DataPath)
	}
	if !FileExists(g.DataPath) {
		return fmt.Errorf("EdgeCore database file %v not exist. ", g.DataPath)
	}

	//TODO: whether to get the parameters from the edgecore configuration file
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

	return nil
}

// Execute performs the get operation.
func Execute(opts *GetOptions, args []string) error {
	if opts.AllNamespace {
		opts.Namespace = "all"
	}

	results, err := QueryMetaFromLocal(opts.Namespace, args[0])
	if err != nil {
		return err
	}
	if len(*results) == 0 {
		return fmt.Errorf("No resources found in namespace: %v ", opts.Namespace)
	}

	fmt.Print(results)
	return nil
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

// QueryMetaFromLocal
func QueryMetaFromLocal(np string, resourceType string) (*[]dao.Meta, error) {
	var err error
	var results *[]dao.Meta

	if np == "all" {
		switch resourceType {
		case "pod":
			results, err = dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.type in ('pod','podlist')",
					dao.MetaTableName,
					dao.MetaTableName))
			if err != nil {
				return nil, err
			}
			return results, nil
		case "service":
			results, err = dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.type in ('service','servicelist')",
					dao.MetaTableName,
					dao.MetaTableName))
			if err != nil {
				return nil, err
			}
			return results, nil
		case "endpoints":
			results, err = dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.type in ('endpoints','endpointslist')",
					dao.MetaTableName,
					dao.MetaTableName))
			if err != nil {
				return nil, err
			}
			return results, nil
		case "all":
			results, err = dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v ",
					dao.MetaTableName))
			if err != nil {
				return nil, err
			}
			return results, nil
		default:
			results, err = dao.QueryMetaByRaw(
				fmt.Sprintf("select * from %v where %v.type = '%v'",
					dao.MetaTableName,
					dao.MetaTableName,
					resourceType))
			if err != nil {
				return nil, err
			}
			return results, nil
		}
	}

	switch resourceType {
	case "pod":
		results, err = dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.type in ('pod','podlist') and %v.key like '%v/%%'",
				dao.MetaTableName,
				dao.MetaTableName,
				dao.MetaTableName,
				np))
		if err != nil {
			return nil, err
		}
		return results, nil
	case "service":
		results, err = dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.type in ('service','servicelist') and %v.key like '%v/%%'",
				dao.MetaTableName,
				dao.MetaTableName,
				dao.MetaTableName,
				np))
		if err != nil {
			return nil, err
		}
		return results, nil
	case "endpoints":
		results, err = dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.type in ('endpoints','endpointslist') and %v.key like '%v/%%'",
				dao.MetaTableName,
				dao.MetaTableName,
				dao.MetaTableName,
				np))
		if err != nil {
			return nil, err
		}
		return results, nil
	case "all":
		results, err = dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.key like '%v/%%'",
				dao.MetaTableName,
				dao.MetaTableName,
				np))
		if err != nil {
			return nil, err
		}
		return results, nil
	default:
		results, err = dao.QueryMetaByRaw(
			fmt.Sprintf("select * from %v where %v.type = '%v' and %v.key like '%v/%%'",
				dao.MetaTableName,
				dao.MetaTableName,
				resourceType,
				dao.MetaTableName,
				np))
		if err != nil {
			return nil, err
		}
		return results, nil
	}

	return results, nil
}
