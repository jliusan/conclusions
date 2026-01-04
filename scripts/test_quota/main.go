package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota/v2"
)
func main(){
	ExampleGroupQuotasClient_BeginCreateOrUpdate()
}

func main1() {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	factory, err := armquota.NewClientFactory(
		cred,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	client := factory.NewGroupQuotasClient()

	pager := client.NewListPager(
		"Microsoft.MachineLearningServices",
		nil,
	)

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for _, q := range page.Value {
			fmt.Println("ID:", *q.ID)
			fmt.Println("LocalizedValue:", *q.Properties.DisplayName)
			fmt.Println("---")
		}
	}

}

func ExampleGroupQuotasClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupQuotasClient().BeginCreateOrUpdate(ctx, "judytest-managementGroupID", "groupquota1", armquota.GroupQuotasEntity{
		Properties: &armquota.GroupQuotasEntityProperties{
			DisplayName: to.Ptr("GroupQuota1"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armquota.GroupQuotasClientCreateOrUpdateResponse{
	// 	GroupQuotasEntity: &armquota.GroupQuotasEntity{
	// 		Name: to.Ptr("groupquota1"),
	// 		Type: to.Ptr("Microsoft.Quota/groupQuotas"),
	// 		ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1"),
	// 		Properties: &armquota.GroupQuotasEntityProperties{
	// 			DisplayName: to.Ptr("GroupQuota1"),
	// 			GroupType: to.Ptr(armquota.GroupTypeAllocationGroup),
	// 			ProvisioningState: to.Ptr(armquota.RequestStateAccepted),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-09-01/GroupQuotas/DeleteGroupQuotas.json
func ExampleGroupQuotasClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupQuotasClient().BeginDelete(ctx, "E7EC67B3-7657-4966-BFFC-41EFD36BAA09", "groupquota1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: 2025-09-01/GroupQuotas/GetGroupQuotas.json
func ExampleGroupQuotasClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGroupQuotasClient().Get(ctx, "E7EC67B3-7657-4966-BFFC-41EFD36BAA09", "groupquota1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armquota.GroupQuotasClientGetResponse{
	// 	GroupQuotasEntity: &armquota.GroupQuotasEntity{
	// 		Name: to.Ptr("groupquota1"),
	// 		Type: to.Ptr("Microsoft.Quota/groupQuotas"),
	// 		ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1"),
	// 		Properties: &armquota.GroupQuotasEntityProperties{
	// 			DisplayName: to.Ptr("GroupQuota1"),
	// 			GroupType: to.Ptr(armquota.GroupTypeAllocationGroup),
	// 			ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-09-01/GroupQuotas/ListGroupQuotas.json
func ExampleGroupQuotasClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGroupQuotasClient().NewListPager("E7EC67B3-7657-4966-BFFC-41EFD36BAA09", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page = armquota.GroupQuotasClientListResponse{
		// 	GroupQuotaList: armquota.GroupQuotaList{
		// 		NextLink: to.Ptr("https://yourLinkHere.com"),
		// 		Value: []*armquota.GroupQuotasEntity{
		// 			{
		// 				Name: to.Ptr("groupquota1"),
		// 				Type: to.Ptr("Microsoft.Quota/groupQuotas"),
		// 				ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1"),
		// 				Properties: &armquota.GroupQuotasEntityProperties{
		// 					DisplayName: to.Ptr("GroupQuota1"),
		// 					GroupType: to.Ptr(armquota.GroupTypeAllocationGroup),
		// 					ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}
