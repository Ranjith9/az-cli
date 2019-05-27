package azurenetwork

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

func getSubnetsClient() network.SubnetsClient {
	subnetsClient := network.NewSubnetsClient(subscription)
	subnetsClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return subnetsClient
}

// This struct is to work with subnets
type SubnetIn struct {
	ResourceGroup string
	VnetName      string `json:"vnetname,omitempty"`
	SubnetName    string `json:"subnetname,omitempty"`
	SubnetCidr    string `json:"cidr,omitempty"`
	NsgID         string `json:"nsg,omitempty"`
}

// This method is to create the subnets in a Vnet
func (sub SubnetIn) CreateVirtualNetworkSubnet() (subnet network.Subnet, err error) {

	subnetParams := network.Subnet{
		Name: to.StringPtr(sub.SubnetName),
		SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr(sub.SubnetCidr),
		},
	}

	if sub.NsgID != "" {
		r := subnetParams.SubnetPropertiesFormat
		r.NetworkSecurityGroup = &network.SecurityGroup{
			ID: to.StringPtr(sub.NsgID),
		}
	}
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		sub.ResourceGroup,
		sub.VnetName,
		sub.SubnetName,
		subnetParams,
	)
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// This method is to delete the subnets in a Vnet
func (sub SubnetIn) DeleteVirtualNetworkSubnet() (ar autorest.Response, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.Delete(
		ctx,
		sub.ResourceGroup,
		sub.VnetName,
		sub.SubnetName,
	)

	if err != nil {
		return ar, fmt.Errorf("cannot delete subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the subnet delete future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// This method is to get the subnets in a Vnet
func (sub SubnetIn) GetVirtualNetworkSubnet() (subnet network.Subnet, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.Get(
		ctx,
		sub.ResourceGroup,
		sub.VnetName,
		sub.SubnetName,
		"")
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet: %v", err)
	}

	return future, err
}

// This method is to lists the subnets in a Vnet
func (sub SubnetIn) ListVirtualNetworkSubnet() (subnet []network.Subnet, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.List(
		ctx,
		sub.ResourceGroup,
		sub.VnetName,
	)

	if err != nil {
		return subnet, fmt.Errorf("cannot list subnetwork: %v", err)
	}

	return future.Values(), err
}
