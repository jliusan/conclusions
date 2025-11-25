package main

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain credential: %v", err)
	}

	subscriptionID := "4d042dc6-fe17-4698-a23f-ec6a8d1e98f4"
	client, err := armnetwork.NewWebApplicationFirewallPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	pager := client.NewListAllPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Fatalf("failed to list WAF policies: %v", err) // ERROR OCCURS HERE
		}
		log.Printf("Got %d policies", len(page.Value))
	}
}
