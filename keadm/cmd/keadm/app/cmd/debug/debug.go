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
	"io"

	"github.com/spf13/cobra"
)

var (
	edgeDebugLongDescription = `"keadm debug" command help  provide debug function to help diagnose the cluster`

	edgeDebugShortDescription = `debug function to help diagnose the cluster`
)

// NewEdgeDebug returns KubeEdge edge debug command.
func NewEdgeDebug(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: edgeDebugShortDescription,
		Long:  edgeDebugLongDescription,
	}
	// add subCommand collect
	cmd.AddCommand(NewDiagnose(out, nil))

	return cmd
}
