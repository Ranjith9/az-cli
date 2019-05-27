package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"io/ioutil"
	"log"
	"os/user"
)

// Credentials struct is to work with Azure credentials
type Credentials struct {
	Profile        string
	ClientID       string
	SubscriptionID string
	TenantID       string
	ClientSecret   string
}

var (
	result Credentials
)

func init() {

	user, _ := user.Current()
	file := user.HomeDir + "/.azure/credentials" // Please use Json format in credentials file.

	plan, _ := ioutil.ReadFile(file)

	var data []Credentials
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println("Cannot unmarshal the json ", err)
	}
	for _, t := range data {

		if t.Profile == "ranjith" {
			result = t
			break
		} else if t.Profile != "ranjith" {
			continue
		}
	}
	if (Credentials{}) == result {
		fmt.Println("This user or the credentials are not valid")
	}
}

// GetServicePrincipalToken method is to initiate the Oauthtoken in Azure
func GetServicePrincipalToken() (adal.OAuthTokenProvider, string, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, result.TenantID)
	code, err := adal.NewServicePrincipalToken(
		*oauthConfig,
		result.ClientID,
		result.ClientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
	}

	return code, result.SubscriptionID, err
}
