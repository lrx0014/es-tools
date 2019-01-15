package app

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/lrx0014/log-tools/cmd/apiserver/app/options"
	apiserver "github.com/lrx0014/log-tools/pkg/server"
	utilflag "github.com/lrx0014/log-tools/pkg/util/flag"
	"github.com/lrx0014/log-tools/pkg/version"
	"github.com/lrx0014/log-tools/pkg/version/verflag"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "openshift-client-test",
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()
			utilflag.PrintFlags(cmd.Flags())

			// validate options
			if errs := opts.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}
			return Run(opts, stopCh)
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

// Run runs the specified APIServer. This should never exit.
func Run(opts *options.ServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	glog.Infof("Version: %+v", version.Get())

	server, err := apiserver.NewAPIServer(opts)
	if err != nil {
		return err
	}

	return server.Run()
}
