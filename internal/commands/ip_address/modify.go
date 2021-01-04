package ip_address

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service service.IpAddress
	req     request.ModifyIPAddressRequest
}

func ModifyCommand(service service.IpAddress) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify an ip address"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.MAC, "mac", "", "MAC address of server interface to attach floating IP to.")
	fs.StringVar(&s.req.PTRRecord, "ptr-record", "", "A fully qualified domain name.")
	s.AddFlags(fs)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(address string) interface{} {
				s.req.IPAddress = address
				return &s.req
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyIPAddressRequest).IPAddress },
				MaxActions:    maxIpAddressActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Modifying IP Address",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyIPAddress(req.(*request.ModifyIPAddressRequest))
				},
			},
		}.Send(args)
	}
}