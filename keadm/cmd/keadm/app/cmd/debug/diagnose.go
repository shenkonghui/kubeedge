package debug

import (
	"fmt"
	types "github.com/kubeedge/kubeedge/keadm/cmd/keadm/app/cmd/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	edgeDiagnoseLongDescription = `"keadm debug collect " command obtain all the data of the current node
and then provide it to the operation and maintenance personnel to locate the problem
`
	edgeDiagnoseShortDescription = `Obtain all the data of the current node`

	edgeDiagnoseExample = `
# Check all items and specified as the current directory
keadm debug collect --path .
`
)

// NewEdgeCollect returns KubeEdge edge debug collect command.
func NewEdgeCollect(out io.Writer, diagnoseOptions *types.DiagnoseOptions) *cobra.Command {
	if diagnoseOptions == nil {
		diagnoseOptions = newDiagnoseOptions()
	}
	cmd := &cobra.Command{
		Use:     "diagnose",
		Short:   edgeDiagnoseShortDescription,
		Long:    edgeDiagnoseLongDescription,
		Example: edgeDiagnoseExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteDiagnose()
		},
	}

	//addCollectOtherFlags(cmd, collectOptions)
	return cmd
}

// newCollectOptions returns a struct ready for being used for creating cmd collect flags.
func newDiagnoseOptions() *types.DiagnoseOptions {
	opts := &types.DiagnoseOptions{}
	return opts
}

//Start to collect data
func ExecuteDiagnose() error {
	fmt.Println("Start collecting data")
	return nil
}
