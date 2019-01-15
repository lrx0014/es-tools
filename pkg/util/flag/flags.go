package flag

import (
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		glog.V(2).Infof("FlAG: --%s=%q", flag.Name, flag.Value)
	})
}
