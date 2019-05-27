package azurenetwork

import (
	"az-cli/azure/access"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-02-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	token, _, subscription = auth.GetServicePrincipalToken()
	ctx                    = context.Background()
)

type VnetIn struct {
	ResourceGroup string
	VnetName      string `json:"vnetname,omitempty"`
	Cidr          string `json:"cidr,omitempty"`
	Location      string `json:"location,omitempty"`
}

func getVnetClient() network.VirtualNetworksClient {
	vnetClient := network.NewVirtualNetworksClient(subscription)
	vnetClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return vnetClient
}

func (net VnetIn) CreateVirtualNetwork() (vnet network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.CreateOrUpdate(
		ctx,
		net.ResourceGroup,
		net.VnetName,
		network.VirtualNetwork{
			Location: to.StringPtr(net.Location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{net.Cidr},
				},
			},
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vnetClient.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(vnetClient)
}

func (net VnetIn) GetVirtualNetwork() (vnet network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.Get(
		ctx,
		net.ResourceGroup,
		net.VnetName,
		"")

	if err != nil {
		return vnet, fmt.Errorf("cannot get virtual network: %v", err)
	}

	return future, err
}

func (net VnetIn) DeleteVirtualNetwork() (ar autorest.Response, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.Delete(
		ctx,
		net.ResourceGroup,
		net.VnetName,
	)

	if err != nil {
		return ar, fmt.Errorf("cannot delete virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vnetClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(vnetClient)
}

func (net VnetIn) ListVirtualNetwork() (vnet []network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.List(
		ctx,
		net.ResourceGroup)

	if err != nil {
		return vnet, fmt.Errorf("cannot list virtual network: %v", err)
	}

	return future.Values(), err
}

func ListAllVirtualNetwork() (vnet []network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.ListAll(
		ctx)

	if err != nil {
		return vnet, fmt.Errorf("cannot list virtual networks: %v", err)
	}

	return future.Values(), err
}
