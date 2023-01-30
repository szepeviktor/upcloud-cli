package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/kubernetes/nodegroup"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "kubernetes create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a Kubernetes cluster",
			`upctl kubernetes create \
				--name my-cluster \
				--network 03e5ca07-f36c-4957-a676-e001e40441eb \
				--node-group count=2,name=my-minimal-node-group,plan=K8S-2xCPU-4GB, \
				--zone de-fra1`,
			`upctl kubernetes create \
				--name my-cluster \
				--network 03e5ca07-f36c-4957-a676-e001e40441eb \
				--node-group count=4,kubelet-arg="log-flush-frequency=5s",label="owner=devteam",label="env=dev",name=my-node-group,plan=K8S-4xCPU-8GB,ssh-key="ssh-ed25519 AAAAo admin@user.com",ssh-key="/path/to/your/public/ssh/key.pub",storage=01000000-0000-4000-8000-000160010100,taint="env=dev:NoSchedule",taint="env=dev2:NoSchedule" \
				--zone de-fra1`,
		),
	}
}

type createParams struct {
	request.CreateKubernetesClusterRequest
	nodeGroups []string
}

func (p *createParams) processParams(exec commands.Executor) error {
	ngs := make([]upcloud.KubernetesNodeGroup, 0)

	for _, v := range p.nodeGroups {
		ng, err := processNodeGroup(v)
		if err != nil {
			return err
		}
		ngs = append(ngs, ng)
	}
	p.NodeGroups = ngs

	networkDetails, err := exec.All().GetNetworkDetails(exec.Context(), &request.GetNetworkDetailsRequest{UUID: p.Network})

	if err != nil || networkDetails == nil || len(networkDetails.IPNetworks) == 0 {
		return fmt.Errorf("invalid network: %w", err)
	}

	p.NetworkCIDR = networkDetails.IPNetworks[0].Address

	return nil
}

func processNodeGroup(in string) (upcloud.KubernetesNodeGroup, error) {
	p := nodegroup.CreateNodeGroupParams{}
	fs := nodegroup.GetCreateNodeGroupFlagSet(&p)
	ng := upcloud.KubernetesNodeGroup{}

	args, err := commands.ParseN(in, 2)
	if err != nil {
		return ng, err
	}

	err = fs.Parse(args)

	if err != nil {
		return ng, err
	}

	return nodegroup.ProcessNodeGroupParams[upcloud.KubernetesNodeGroup](p)
}

type createCommand struct {
	*commands.BaseCommand
	params createParams
	completion.Kubernetes
}

// InitCommand implements Command.InitCommand
func (c *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	c.params = createParams{CreateKubernetesClusterRequest: request.CreateKubernetesClusterRequest{}}

	fs.StringVar(&c.params.Name, "name", "", "Kubernetes cluster name.")
	fs.StringVar(&c.params.Network, "network", "", "Network to use. The value should be UUID of a private network.")
	fs.StringArrayVar(
		&c.params.nodeGroups,
		"node-group",
		[]string{},
		"Node group(s) for running workloads, multiple can be declared.\n"+
			"Usage: `--node-group "+
			"count=8,"+
			"kubelet-arg=\"log-flush-frequency=5s\","+
			"label=\"owner=devteam\","+
			"label=\"env=dev\","+
			"name=my-node-group,"+
			"plan=K8S-2xCPU-4GB,"+
			"ssh-key=\"ssh-ed25519 AAAAo admin@user.com\","+
			"ssh-key=\"/path/to/your/public/ssh/key.pub\","+
			"storage=01000000-0000-4000-8000-000160010100,"+
			"taint=\"env=dev:NoSchedule\","+
			"taint=\"env=dev2:NoSchedule\"`",
	)
	fs.StringVar(&c.params.Zone, "zone", "", "Zone where to create the cluster.")
	c.AddFlags(fs)

	_ = c.Cobra().MarkFlagRequired("name")
	_ = c.Cobra().MarkFlagRequired("network")
	_ = c.Cobra().MarkFlagRequired("node-group")
	_ = c.Cobra().MarkFlagRequired("zone")
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Creating cluster %s", c.params.Name)
	exec.PushProgressStarted(msg)
	if err := c.params.processParams(exec); err != nil {
		return nil, err
	}

	r := c.params.CreateKubernetesClusterRequest

	res, err := svc.CreateKubernetesCluster(exec.Context(), &r)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}