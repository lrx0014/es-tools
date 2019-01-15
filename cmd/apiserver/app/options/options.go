package options

import (
	"github.com/spf13/pflag"
)

// ServerRunOptions runs an apiserver
type ServerRunOptions struct {
	// address is the IP address for the APIServer to serve on (set to 0.0.0.0
	// for all interfaces)
	Address string
	// port is the port for the APIServer to serve on.
	Port uint
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		Address: "0.0.0.0",
		Port:    8080,
	}
	return &s
}

// AddFlags adds flags to fs and binds them to options.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.Address, "address", "", s.Address, "The IP address for the APIServer to serve on (set to 0.0.0.0 for all interfaces).")
	fs.UintVarP(&s.Port, "port", "", s.Port, "The port for the APIServer to serve on.")
}
