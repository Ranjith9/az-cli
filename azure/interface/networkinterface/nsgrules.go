package azurenetwork

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

func getNsgRuleClient() network.SecurityRulesClient {
	nsgRuleClient := network.NewSecurityRulesClient(subscription)
	nsgRuleClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return nsgRuleClient
}

// SecurityRuleIn struct is to work with the NSGrules
type SecurityRuleIn struct {
	ResourceGroup string
	NsgName       string `json:"nsgname,omitempty"`
	RuleName      string `json:"rulename,omitempty"`
	Port          string `json:"port,omitempty"`
	Priority      int32  `json:"priority,omitempty"`
}

//  CreateNetworkSecurityRule method is to create the nsgrules in a resourcegroup
func (rule SecurityRuleIn) CreateNetworkSecurityRule() (nsgrule network.SecurityRule, err error) {
	nsgRuleClient := getNsgRuleClient()
	future, err := nsgRuleClient.CreateOrUpdate(
		ctx,
		rule.ResourceGroup,
		rule.NsgName,
		rule.RuleName,
		network.SecurityRule{
			Name: to.StringPtr(rule.RuleName),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
				SourcePortRange:          to.StringPtr("1-65535"),
				DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
				DestinationPortRange:     to.StringPtr(rule.Port),
				Access:                   network.SecurityRuleAccessAllow,
				Direction:                network.SecurityRuleDirectionInbound,
				Priority:                 to.Int32Ptr(rule.Priority),
			},
		},
	)

	if err != nil {
		return nsgrule, fmt.Errorf("cannot create nsgRule: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, nsgRuleClient.Client)
	if err != nil {
		return nsgrule, fmt.Errorf("cannot get nsgRule create or update future response: %v", err)
	}

	return future.Result(nsgRuleClient)
}
