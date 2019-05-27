package azureregroup

import (
	"az-cli/azure/access"
	"context"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-03-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	token, _, subscription = auth.GetServicePrincipalToken()
	ctx                    = context.Background()
)

type GroupsIn struct {
	ResourceGroup string
	Location      string `json:"location,omitempty"`
}

func getGroupsClient() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(subscription)
	groupsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return groupsClient
}

func (g GroupsIn) CreateResourceGroup() (group resources.Group, err error) {
	groupsClient := getGroupsClient()
	return groupsClient.CreateOrUpdate(
		ctx,
		g.ResourceGroup,
		resources.Group{
			Location: to.StringPtr(g.Location),
		})
}

func (g GroupsIn) GetResourceGroup() (resources.Group, error) {
	groupsClient := getGroupsClient()
	return groupsClient.Get(
		ctx,
		g.ResourceGroup,
	)
}

func (g GroupsIn) CheckResourceGroup() (ar autorest.Response, err error) {
	groupsClient := getGroupsClient()
	resp, err := groupsClient.CheckExistence(
		ctx,
		g.ResourceGroup,
	)
	if err != nil {
		return ar, err
	}
	return resp, err
}
