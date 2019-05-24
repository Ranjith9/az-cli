package azurecompute

import (
	"fmt"
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-03-01/compute"
        "github.com/Azure/go-autorest/autorest"
)

func getDisksClient() compute.DisksClient {
	disksClient := compute.NewDisksClient(subscription)
	disksClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return disksClient
}

type DisksIn struct {
	ResourceGroup string
	DiskName  string      `json:"snapshotname,omitempty"`
	Location  string      `json:"location,omitempty"`
}

func (d DisksIn) DeleteDisk() (ar autorest.Response, err error) {

        disksClient := getDisksClient()
        future, err := disksClient.Delete(
                ctx,
                d.ResourceGroup,
                d.DiskName,
                )

        if err != nil {
                return ar, fmt.Errorf("cannot delete disk: %v", err)
        }

        err = future.WaitForCompletionRef(ctx, disksClient.Client)
        if err != nil {
                return ar, fmt.Errorf("cannot get the disk delete future response: %v", err)
        }

        return  future.Result(disksClient)
}
