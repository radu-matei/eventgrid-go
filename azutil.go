package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	defaultLocation = "westeurope"

	defaultActiveDirectoryEndpoint = azure.PublicCloud.ActiveDirectoryEndpoint
	defaultResourceManagerEndpoint = azure.PublicCloud.ResourceManagerEndpoint

	subscriptionID = getEnvVarOrExit("AZ_SUBSCRIPTION_ID")
	tenantID       = getEnvVarOrExit("AZ_TENANT_ID")
	clientID       = getEnvVarOrExit("AZ_CLIENT_ID")
	clientSecret   = getEnvVarOrExit("AZ_CLIENT_SECRET")
)

func getEventGridClient() (eventgrid.EventSubscriptionsClient, error) {
	var subscriptionsClient eventgrid.EventSubscriptionsClient

	oAuthConfig, err := adal.NewOAuthConfig(defaultActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return subscriptionsClient, fmt.Errorf("cannot get oauth config: %v", err)
	}
	token, err := adal.NewServicePrincipalToken(*oAuthConfig, clientID, clientSecret, defaultResourceManagerEndpoint)
	if err != nil {
		return subscriptionsClient, fmt.Errorf("cannot get service principal token: %v", err)
	}

	subscriptionsClient = eventgrid.NewEventSubscriptionsClient(subscriptionID)
	subscriptionsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return subscriptionsClient, nil
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
