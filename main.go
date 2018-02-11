package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	storageAccountName = getEnvVarOrExit("STORAGE_ACCOUNT")
	resourceGroup      = getEnvVarOrExit("RESOURCE_GROUP")

	tlsWebhook = getEnvVarOrExit("TLS_WEHBOOK")
)

func main() {
	c, err := getEventGridClient()
	if err != nil {
		log.Fatalf("cannot get eventgrid client: %v", err)
	}

	scope := fmt.Sprintf("/subscriptions/%s/%s/providers/Microsoft.Storage/storageAccounts/%s", subscriptionID, resourceGroup, storageAccountName)
	subscriptionName := fmt.Sprintf("%seventsubscription", storageAccountName)

	subscription := eventgrid.EventSubscription{
		EventSubscriptionProperties: &eventgrid.EventSubscriptionProperties{
			Destination: eventgrid.WebHookEventSubscriptionDestination{
				EndpointType: eventgrid.EndpointTypeWebHook,
				WebHookEventSubscriptionDestinationProperties: &eventgrid.WebHookEventSubscriptionDestinationProperties{
					EndpointURL: to.StringPtr(tlsWebhook),
				},
			},
		},
	}

	ctx := context.Background()

	f, err := c.CreateOrUpdate(ctx, scope, subscriptionName, subscription)
	if err != nil {
		log.Fatalf("cannot create event subscription: %v", err)
	}

	err = f.WaitForCompletion(ctx, c.Client)
	if err != nil {
		log.Fatalf("cannot get the subscription create or update future response: %v", err)
	}
}
