package plugin

import (
	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/pluginhelper/http"
	"github.com/cloudnative-pg/cnpg-i/pkg/lifecycle"
	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/leonardoce/cnpg-i-skip-initdb/internal/identity"
	lifecycleImpl "github.com/leonardoce/cnpg-i-skip-initdb/internal/lifecycle"
	operatorImpl "github.com/leonardoce/cnpg-i-skip-initdb/internal/operator"
)

// NewCmd creates the `plugin` command
func NewCmd() *cobra.Command {
	cmd := http.CreateMainCmd(identity.Implementation{}, func(server *grpc.Server) error {
		// Register the declared implementations
		operator.RegisterOperatorServer(server, operatorImpl.Implementation{})
		lifecycle.RegisterOperatorLifecycleServer(server, lifecycleImpl.Implementation{})
		return nil
	})

	cmd.Use = "plugin"

	return cmd
}
