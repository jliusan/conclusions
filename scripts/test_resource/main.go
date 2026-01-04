package main

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func main() {
	//az ts create --name TemplateSpec-Name --version v1 --resource-group judytest --location eastus2 --template-file D:\GoProject\conclusions\scripts\test_resource\template.json
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armresources.NewClientFactory("4d042dc6-fe17-4698-a23f-ec6a8d1e98f4", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	jsonString := `{
  "id": { "value": 42 },
  "name": { "value": "Alice" },
  "is_active": { "value": true }
}`
	poller, err := clientFactory.NewDeploymentsClient().BeginCreateOrUpdate(ctx, "judytest", "my-deployment", armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			Parameters: jsonString,
			TemplateLink: &armresources.TemplateLink{
				ID: to.Ptr("/subscriptions/4d042dc6-fe17-4698-a23f-ec6a8d1e98f4/resourceGroups/judytest/providers/Microsoft.Resources/TemplateSpecs/Microsoft.Template-20260104111510/versions/v1"),
			},
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
	// res.DeploymentExtended = armresources.DeploymentExtended{
	// 	Name: to.Ptr("my-deployment"),
	// 	Type: to.Ptr("Microsoft.Resources/deployments"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/my-resource-group/providers/Microsoft.Resources/deployments/my-deployment"),
	// 	Properties: &armresources.DeploymentPropertiesExtended{
	// 		CorrelationID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		Dependencies: []*armresources.Dependency{
	// 		},
	// 		Duration: to.Ptr("PT22.8356799S"),
	// 		Mode: to.Ptr(armresources.DeploymentModeIncremental),
	// 		OutputResources: []*armresources.ResourceReference{
	// 			{
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/my-resource-group/providers/Microsoft.Storage/storageAccounts/my-storage-account"),
	// 		}},
	// 		Parameters: map[string]any{
	// 		},
	// 		Providers: []*armresources.Provider{
	// 			{
	// 				Namespace: to.Ptr("Microsoft.Storage"),
	// 				ResourceTypes: []*armresources.ProviderResourceType{
	// 					{
	// 						Locations: []*string{
	// 							to.Ptr("eastus")},
	// 							ResourceType: to.Ptr("storageAccounts"),
	// 					}},
	// 			}},
	// 			ProvisioningState: to.Ptr(armresources.ProvisioningStateSucceeded),
	// 			TemplateHash: to.Ptr("0000000000000000000"),
	// 			TemplateLink: &armresources.TemplateLink{
	// 				ContentVersion: to.Ptr("1.0.0.0"),
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/my-resource-group/providers/Microsoft.Resources/TemplateSpecs/TemplateSpec-Name/versions/v1"),
	// 			},
	// 			Timestamp: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-05T01:20:01.723Z"); return t}()),
	// 		},
	// 	}
}
