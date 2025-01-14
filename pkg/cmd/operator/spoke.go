package operator

import (
	"github.com/spf13/cobra"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"open-cluster-management.io/registration-operator/pkg/operators/klusterlet"
	"open-cluster-management.io/registration-operator/pkg/version"
)

// NewKlusterletOperatorCmd generatee a command to start klusterlet operator
func NewKlusterletOperatorCmd() *cobra.Command {

	options := klusterlet.Options{}
	cmd := controllercmd.
		NewControllerCommandConfig("klusterlet", version.Get(), options.RunKlusterletOperator).
		NewCommand()
	cmd.Use = "klusterlet"
	cmd.Short = "Start the klusterlet operator"

	cmd.Flags().BoolVar(&options.SkipPlaceholderHubSecret, "skip-placeholder-hub-secret", false,
		"If set, will skip ensuring a placeholder hub secret which is originally intended for pulling "+
			"work image before approved")

	return cmd
}
