package serverfirewall

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFirewallRuleCommand(t *testing.T) {

	var Server1 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	for _, test := range []struct {
		name        string
		args        []string
		expectedReq request.CreateFirewallRuleRequest
		error       string
	}{
		{
			name:  "Empty info",
			args:  []string{},
			error: "direction is required",
		},
		{
			name: "Action is required",
			args: []string{
				Server1.UUID,
				"--direction", "in",
			},
			error: "action is required",
		},
		{
			name: "Family is required",
			args: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
			},
			error: "family (IPv4/IPv6) is required",
		},
		{
			name: "FirewallRule, accept incoming IPv6",
			args: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
				"--family", "IPv6",
			},
			expectedReq: request.CreateFirewallRuleRequest{
				FirewallRule: upcloud.FirewallRule{
					Direction: "in",
					Action:    "accept",
					Family:    "IPv6",
				},
				ServerUUID: Server1.UUID,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.MockService{}
			mService.On("GetServers", mock.Anything).Return(servers, nil)

			cc := commands.BuildCommand(CreateCommand(&mService, &mService), nil, config.New())
			err1 := cc.SetFlags(test.args)
			if err1 != nil {
				panic(err1)
			}

			_, err := cc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			}
		})
	}
}
