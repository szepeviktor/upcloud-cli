package resolver

import (
	"fmt"
	internal "github.com/UpCloudLtd/cli/internal/service"
)

// CachingServer implements resolver for servers, caching the results
type CachingServer struct{}

var _ ResolutionProvider = CachingServer{}

// Get implements ResolutionProvider.Get
func (s CachingServer) Get(svc internal.AllServices) (Resolver, error) {
	servers, err := svc.GetServers()
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, server := range servers.Servers {
			if server.Title == arg || server.Hostname == arg || server.UUID == arg {
				if rv != "" {
					return "", fmt.Errorf("'%v' is ambiguous, found multiple servers matching", arg)
				}
				rv = server.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", fmt.Errorf("no server found matching '%v'", arg)
	}, nil
}
