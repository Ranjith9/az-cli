package azurecompute

import (
	"az-cli/azure/access"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-03-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	token, _, subscription = auth.GetServicePrincipalToken()
	ctx                    = context.Background()
)

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(subscription)
	vmClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return vmClient
}

type VMIn struct {
	ResourceGroup string
	VmName        string `json:"vmname,omitempty"`
	NicID         string `json:"nicid,omitempty"`
	UserName      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Location      string `json:"location,omitempty"`
}

func (v VMIn) CreateVM() (vm compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.CreateOrUpdate(
		ctx,
		v.ResourceGroup,
		v.VmName,
		compute.VirtualMachine{
			Location: to.StringPtr(v.Location),
			VirtualMachineProperties: &compute.VirtualMachineProperties{
				HardwareProfile: &compute.HardwareProfile{
					VMSize: compute.VirtualMachineSizeTypesBasicA0,
				},
				StorageProfile: &compute.StorageProfile{
					ImageReference: &compute.ImageReference{
						Publisher: to.StringPtr("Canonical"),
						Offer:     to.StringPtr("UbuntuServer"),
						Sku:       to.StringPtr("16.04-LTS"),
						Version:   to.StringPtr("latest"),
					},
				},
				OsProfile: &compute.OSProfile{
					ComputerName:  to.StringPtr(v.VmName),
					AdminUsername: to.StringPtr(v.UserName),
					AdminPassword: to.StringPtr(v.Password),
					LinuxConfiguration: &compute.LinuxConfiguration{
						DisablePasswordAuthentication: to.BoolPtr(false),
					},
				},
				NetworkProfile: &compute.NetworkProfile{
					NetworkInterfaces: &[]compute.NetworkInterfaceReference{
						{
							ID: to.StringPtr(v.NicID),
							NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return vm, fmt.Errorf("cannot create vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return vm, fmt.Errorf("cannot get the vm create or update future response: %v", err)
	}

	return future.Result(vmClient)
}

func (v VMIn) DeleteVM() (ar autorest.Response, err error) {

	vmClient := getVMClient()
	future, err := vmClient.Delete(
		ctx,
		v.ResourceGroup,
		v.VmName,
	)
	if err != nil {
		return ar, fmt.Errorf("cannot delete VM: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the VM delete future response: %v", err)
	}

	return future.Result(vmClient)
}

func (v VMIn) GetVM() (vm compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.Get(
		ctx,
		v.ResourceGroup,
		v.VmName,
		"")

	if err != nil {
		return vm, fmt.Errorf("cannot get virtual VM: %v", err)
	}

	return future, err
}

func (v VMIn) ListVM() (vm []compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.List(
		ctx,
		v.ResourceGroup,
	)

	if err != nil {
		return vm, fmt.Errorf("cannot list the VMs in a resourcegroup: %v", err)
	}

	return future.Values(), err
}
