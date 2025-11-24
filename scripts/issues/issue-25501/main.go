package main

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
)

var (
	subscriptionID    string
	location          = "eastus2"
	resourceGroupName = "SSS3PT_test_rg-sample-resource-group-eastus2"
	serverName        = "sample-pg-server"
	databaseName      = "sample-postgresql-database"
)

var (
	resourcesClientFactory *armresources.ClientFactory
)

var (
	resourceGroupClient *armresources.ResourceGroupsClient
	serversClient       *armpostgresqlflexibleservers.ServersClient
	databasesClient     *armpostgresqlflexibleservers.DatabasesClient
)

func main() {
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set.")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	resourcesClientFactory, err = armresources.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	resourceGroupClient = resourcesClientFactory.NewResourceGroupsClient()

	serversClient, err = armpostgresqlflexibleservers.NewServersClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	databasesClient, err = armpostgresqlflexibleservers.NewDatabasesClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	resourceGroup, err := createResourceGroup(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("resources group:", *resourceGroup.ID)

	server, err := createServer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postgresql server:", *server.ID)

	database, err := createDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postgresql database:", *database.ID)

	keepResource := os.Getenv("KEEP_RESOURCE")
	if len(keepResource) == 0 {
		err = cleanup(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("cleaned up successfully.")
	}
}

func createServer(ctx context.Context) (*armpostgresqlflexibleservers.Server, error) {

	pollerResp, err := serversClient.BeginCreate(
		ctx,
		resourceGroupName,
		serverName,
		armpostgresqlflexibleservers.Server{
			Location: to.Ptr(location),
			Properties: &armpostgresqlflexibleservers.ServerProperties{
				AdministratorLogin:         to.Ptr("dummylogin"),
				AdministratorLoginPassword: to.Ptr("QWE123!@#"),
				Version:                    to.Ptr(armpostgresqlflexibleservers.ServerVersion("16")),
				Storage: &armpostgresqlflexibleservers.Storage{
					StorageSizeGB: to.Ptr[int32](32),
				},
			},
			SKU: &armpostgresqlflexibleservers.SKU{
				Name: to.Ptr("Standard_B1ms"),
				Tier: to.Ptr(armpostgresqlflexibleservers.SKUTierBurstable),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func createDatabase(ctx context.Context) (*armpostgresqlflexibleservers.Database, error) {

	pollerResp, err := databasesClient.BeginCreate(
		ctx,
		resourceGroupName,
		serverName,
		databaseName,
		armpostgresqlflexibleservers.Database{
			Properties: &armpostgresqlflexibleservers.DatabaseProperties{
				Charset:   to.Ptr("UTF8"),
				Collation: to.Ptr("English_United States.1252"),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Database, nil
}

func createResourceGroup(ctx context.Context) (*armresources.ResourceGroup, error) {

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(location),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func cleanup(ctx context.Context) error {

	pollerResp, err := resourceGroupClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
