package completion

import (
	"github.com/UpCloudLtd/cli/internal/service"
	"github.com/spf13/cobra"
)

// Provider should be implemented by a command that can provide argument completion
type Provider interface {
	CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective)
}