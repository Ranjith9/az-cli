package environment

import (
	"fmt"
	"az-cli/azure/interface/computeinterface"
	"az-cli/azure/interface/networkinterface"
	"az-cli/azure/interface/resourceinterface"
)

type CreateIN struct {
	Name         string
}

func (c CreateIN) Create() {

	g := azureregroup.GroupsIn{"CLI-group", "CentralIndia"}
	g.CreateResourceGroup()


	v := azurenetwork.VnetIn{"CLI-group", "gloify-net", "192.168.0.0/16", "CentralIndia"}
	v.CreateVirtualNetwork()

	s := azurenetwork.SubnetIn{"CLI-group", "gloify-net", "subnet1", "192.168.10.0/24", ""}
	sub, _ := s.CreateVirtualNetworkSubnet()


	n := azurenetwork.NsgIn{"CLI-group", c.Name + "-nsg", "CentralIndia"}
	nsg, _ := n.CreateNetworkSecurityGroup()
	fmt.Println(*nsg.ID)

	nr := azurenetwork.SecurityRuleIn{"CLI-group", *nsg.Name, "ssh", "22", 123}
	nsgrule, _ := nr.CreateNetworkSecurityRule()
	fmt.Println(*nsgrule.ID)

	i := azurenetwork.IpIn{"CLI-group", c.Name + "-ip", "CentralIndia"}
	ip, _ := i.CreatePublicIP()
	fmt.Println(*ip.ID)

	ic := azurenetwork.NicIn{"CLI-group", c.Name + "-nic", *nsg.ID, *sub.ID, *ip.ID, "CentralIndia"}
	nic, _ := ic.CreateNIC()
	fmt.Println(*nic.ID)

	m := azurecompute.VMIn{"CLI-group", c.Name, *nic.ID, "ubuntu", "ubuntu@12345", "CentralIndia"}
	vm, _ := m.CreateVM()
	fmt.Println(*vm.ID)

}
