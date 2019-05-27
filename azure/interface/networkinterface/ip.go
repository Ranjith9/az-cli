package azurenetwork

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

func getIPClient() network.PublicIPAddressesClient {
	ipClient := network.NewPublicIPAddressesClient(subscription)
	ipClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return ipClient
}

// This struct is to work with the ips
type IpIn struct {
	ResourceGroup string
	IpName        string `json:"ipname,omitempty"`
	Location      string `json:"location,omitempty"`
}

// This method is to create the ips in a resourcegroup
func (pubip IpIn) CreatePublicIP() (ip network.PublicIPAddress, err error) {
	ipClient := getIPClient()
	future, err := ipClient.CreateOrUpdate(
		ctx,
		pubip.ResourceGroup,
		pubip.IpName,
		network.PublicIPAddress{
			Name:     to.StringPtr(pubip.IpName),
			Location: to.StringPtr(pubip.Location),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPv4,
				PublicIPAllocationMethod: network.Static,
				DNSSettings: &network.PublicIPAddressDNSSettings{
					DomainNameLabel: to.StringPtr(pubip.IpName),
				},
			},
		},
	)

	if err != nil {
		return ip, fmt.Errorf("cannot create public ip address: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	if err != nil {
		return ip, fmt.Errorf("cannot get public ip address create or update future response: %v", err)
	}

	return future.Result(ipClient)
}

// This method is to delete the ips in a resourcegroup
func (pubip IpIn) DeletePublicIP() (ar autorest.Response, err error) {
	ipClient := getIPClient()
	future, err := ipClient.Delete(
		ctx,
		pubip.ResourceGroup,
		pubip.IpName,
	)
	if err != nil {
		return ar, fmt.Errorf("cannot delete ip: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get ip delete future response: %v", err)
	}

	return future.Result(ipClient)
}

// This method is to get the ips in a resourcegroup
func (pubip IpIn) GetPublicIP() (ip network.PublicIPAddress, err error) {
	ipClient := getIPClient()
	future, err := ipClient.Get(
		ctx,
		pubip.ResourceGroup,
		pubip.IpName,
		"")

	if err != nil {
		return ip, fmt.Errorf("cannot list ip: %v", err)
	}

	return future, err
}

// This method is to lists the ips in a resourcegroup
func (pubip IpIn) ListPublicIP() (ip []network.PublicIPAddress, err error) {
	ipClient := getIPClient()
	future, err := ipClient.List(
		ctx,
		pubip.ResourceGroup,
	)

	if err != nil {
		return ip, fmt.Errorf("cannot list IPs: %v", err)
	}

	return future.Values(), err
}

// This method is to lists the ips
func (pubip IpIn) ListAllPublicIP() (ip []network.PublicIPAddress, err error) {
	ipClient := getIPClient()
	future, err := ipClient.ListAll(
		ctx,
	)

	if err != nil {
		return ip, fmt.Errorf("cannot list ip: %v", err)
	}

	return future.Values(), err
}
