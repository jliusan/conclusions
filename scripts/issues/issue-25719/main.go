package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v2"
)

func main() {
	// TestExampleContainerAppsClient_BeginCreateOrUpdate()
	ExampleContainerAppsClient_NewListBySubscriptionPager()
}

func ExampleContainerAppsClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armappcontainers.NewClientFactory("4d042dc6-fe17-4698-a23f-ec6a8d1e98f4", cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Logging: policy.LogOptions{
				// include HTTP body for log
				IncludeBody: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewContainerAppsClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			fmt.Printf("Container App Name: %s v.Identity.Type:%v Identity.UserAssignedIdentities:%v Identity is nil:%v\n", *v.Name, *v.Identity.Type, v.Identity.UserAssignedIdentities, v.Identity == nil)
		}
	}
}

func TestExampleContainerAppsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armappcontainers.NewClientFactory("4d042dc6-fe17-4698-a23f-ec6a8d1e98f4", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// First, create the managed environment
	log.Println("Creating managed environment...")
	envPoller, err := clientFactory.NewManagedEnvironmentsClient().BeginCreateOrUpdate(ctx, "judy_test", "demokube", armappcontainers.ManagedEnvironment{
		Location:   to.Ptr("East US2"),
		Properties: &armappcontainers.ManagedEnvironmentProperties{},
	}, nil)
	if err != nil {
		log.Fatalf("failed to create managed environment: %v", err)
	}
	_, err = envPoller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to finish creating managed environment: %v", err)
	}
	log.Println("Managed environment created successfully")

	// Now create the container app
	log.Println("Creating container app...")

	poller, err := clientFactory.NewContainerAppsClient().BeginCreateOrUpdate(ctx, "judy_test", "testcontainerapp0", armappcontainers.ContainerApp{
		Location: to.Ptr("East US2"),
		Properties: &armappcontainers.ContainerAppProperties{
			Configuration: &armappcontainers.Configuration{
				Dapr: &armappcontainers.Dapr{
					AppPort:            to.Ptr[int32](3000),
					AppProtocol:        to.Ptr(armappcontainers.AppProtocolHTTP),
					EnableAPILogging:   to.Ptr(true),
					Enabled:            to.Ptr(true),
					HTTPMaxRequestSize: to.Ptr[int32](10),
					HTTPReadBufferSize: to.Ptr[int32](30),
					LogLevel:           to.Ptr(armappcontainers.LogLevelDebug),
				},
				Ingress: &armappcontainers.Ingress{
					ClientCertificateMode: to.Ptr(armappcontainers.IngressClientCertificateModeAccept),
					CorsPolicy: &armappcontainers.CorsPolicy{
						AllowCredentials: to.Ptr(true),
						AllowedHeaders: []*string{
							to.Ptr("HEADER1"),
							to.Ptr("HEADER2")},
						AllowedMethods: []*string{
							to.Ptr("GET"),
							to.Ptr("POST")},
						AllowedOrigins: []*string{
							to.Ptr("https://a.test.com"),
							to.Ptr("https://b.test.com")},
						ExposeHeaders: []*string{
							to.Ptr("HEADER3"),
							to.Ptr("HEADER4")},
						MaxAge: to.Ptr[int32](1234),
					},
					External: to.Ptr(true),
					IPSecurityRestrictions: []*armappcontainers.IPSecurityRestrictionRule{
						{
							Name:           to.Ptr("Allow work IP A subnet"),
							Description:    to.Ptr("Allowing all IP's within the subnet below to access containerapp"),
							Action:         to.Ptr(armappcontainers.ActionAllow),
							IPAddressRange: to.Ptr("192.168.1.1/32"),
						},
						{
							Name:           to.Ptr("Allow work IP B subnet"),
							Description:    to.Ptr("Allowing all IP's within the subnet below to access containerapp"),
							Action:         to.Ptr(armappcontainers.ActionAllow),
							IPAddressRange: to.Ptr("192.168.1.1/8"),
						}},
					TargetPort: to.Ptr[int32](3000),
				},
				MaxInactiveRevisions: to.Ptr[int32](10),
				Service: &armappcontainers.Service{
					Type: to.Ptr("redis"),
				},
			},
			EnvironmentID: to.Ptr("/subscriptions/4d042dc6-fe17-4698-a23f-ec6a8d1e98f4/resourceGroups/judy_test/providers/Microsoft.App/managedEnvironments/demokube"),
			Template: &armappcontainers.Template{
				Containers: []*armappcontainers.Container{
					{
						Name:  to.Ptr("testcontainerapp0"),
						Image: to.Ptr("repo/testcontainerapp0:v1"),
						Probes: []*armappcontainers.ContainerAppProbe{
							{
								Type: to.Ptr(armappcontainers.TypeLiveness),
								HTTPGet: &armappcontainers.ContainerAppProbeHTTPGet{
									Path: to.Ptr("/health"),
									HTTPHeaders: []*armappcontainers.ContainerAppProbeHTTPGetHTTPHeadersItem{
										{
											Name:  to.Ptr("Custom-Header"),
											Value: to.Ptr("Awesome"),
										}},
									Port: to.Ptr[int32](8080),
								},
								InitialDelaySeconds: to.Ptr[int32](3),
								PeriodSeconds:       to.Ptr[int32](3),
							}},
					}},
				InitContainers: []*armappcontainers.InitContainer{
					{
						Name: to.Ptr("testinitcontainerapp0"),
						Args: []*string{
							to.Ptr("-c"),
							to.Ptr("while true; do echo hello; sleep 10;done")},
						Command: []*string{
							to.Ptr("/bin/sh")},
						Image: to.Ptr("repo/testinitcontainerapp0:v4"),
						Resources: &armappcontainers.ContainerResources{
							CPU:    to.Ptr[float64](0.2),
							Memory: to.Ptr("100Mi"),
						},
					}},
				Scale: &armappcontainers.Scale{
					MaxReplicas: to.Ptr[int32](5),
					MinReplicas: to.Ptr[int32](1),
					Rules: []*armappcontainers.ScaleRule{
						{
							Name: to.Ptr("httpscalingrule"),
							Custom: &armappcontainers.CustomScaleRule{
								Type: to.Ptr("http"),
								Metadata: map[string]*string{
									"concurrentRequests": to.Ptr("50"),
								},
							},
						}},
				},
				ServiceBinds: []*armappcontainers.ServiceBind{
					{
						Name:      to.Ptr("redisService"),
						ServiceID: to.Ptr("/subscriptions/4d042dc6-fe17-4698-a23f-ec6a8d1e98f4/resourceGroups/judy_test/providers/Microsoft.App/containerApps/redisService"),
					}},
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
	// res.ContainerApp = armappcontainers.ContainerApp{
	// 	Name: to.Ptr("testcontainerApp0"),
	// 	Type: to.Ptr("Microsoft.App/containerApps"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/rg/providers/Microsoft.App/containerApps/testcontainerApp0"),
	// 	Location: to.Ptr("East US"),
	// 	Properties: &armappcontainers.ContainerAppProperties{
	// 		Configuration: &armappcontainers.Configuration{
	// 			Dapr: &armappcontainers.Dapr{
	// 				AppPort: to.Ptr[int32](3000),
	// 				AppProtocol: to.Ptr(armappcontainers.AppProtocolHTTP),
	// 				Enabled: to.Ptr(true),
	// 				HTTPReadBufferSize: to.Ptr[int32](30),
	// 			},
	// 			Ingress: &armappcontainers.Ingress{
	// 				CustomDomains: []*armappcontainers.CustomDomain{
	// 					{
	// 						Name: to.Ptr("www.my-name.com"),
	// 						BindingType: to.Ptr(armappcontainers.BindingTypeSniEnabled),
	// 						CertificateID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-name-dot-com"),
	// 					},
	// 					{
	// 						Name: to.Ptr("www.my-other-name.com"),
	// 						BindingType: to.Ptr(armappcontainers.BindingTypeSniEnabled),
	// 						CertificateID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-other-name-dot-com"),
	// 				}},
	// 				External: to.Ptr(true),
	// 				Fqdn: to.Ptr("testcontainerApp0.demokube-t24clv0g.eastus.containerApps.k4apps.io"),
	// 				IPSecurityRestrictions: []*armappcontainers.IPSecurityRestrictionRule{
	// 					{
	// 						Name: to.Ptr("Allow work IP A subnet"),
	// 						Description: to.Ptr("Allowing all IP's within the subnet below to access containerapp"),
	// 						Action: to.Ptr(armappcontainers.ActionAllow),
	// 						IPAddressRange: to.Ptr("192.168.1.1/32"),
	// 					},
	// 					{
	// 						Name: to.Ptr("Allow work IP B subnet"),
	// 						Description: to.Ptr("Allowing all IP's within the subnet below to access containerapp"),
	// 						Action: to.Ptr(armappcontainers.ActionAllow),
	// 						IPAddressRange: to.Ptr("192.168.1.1/8"),
	// 				}},
	// 				StickySessions: &armappcontainers.IngressStickySessions{
	// 					Affinity: to.Ptr(armappcontainers.AffinitySticky),
	// 				},
	// 				TargetPort: to.Ptr[int32](3000),
	// 				Traffic: []*armappcontainers.TrafficWeight{
	// 					{
	// 						RevisionName: to.Ptr("testcontainerApp0-ab1234"),
	// 						Weight: to.Ptr[int32](80),
	// 					},
	// 					{
	// 						Label: to.Ptr("staging"),
	// 						RevisionName: to.Ptr("testcontainerApp0-ab4321"),
	// 						Weight: to.Ptr[int32](20),
	// 				}},
	// 				Transport: to.Ptr(armappcontainers.IngressTransportMethodAuto),
	// 			},
	// 		},
	// 		EventStreamEndpoint: to.Ptr("testEndpoint"),
	// 		LatestReadyRevisionName: to.Ptr("testcontainerApp0-pjxhsye"),
	// 		LatestRevisionFqdn: to.Ptr("testcontainerApp0-pjxhsye.demokube-t24clv0g.eastus.containerApps.k4apps.io"),
	// 		ManagedEnvironmentID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube"),
	// 		ProvisioningState: to.Ptr(armappcontainers.ContainerAppProvisioningStateSucceeded),
	// 		Template: &armappcontainers.Template{
	// 			Containers: []*armappcontainers.Container{
	// 				{
	// 					Name: to.Ptr("testcontainerApp0"),
	// 					Image: to.Ptr("repo/testcontainerApp0:v4"),
	// 					Resources: &armappcontainers.ContainerResources{
	// 						CPU: to.Ptr[float64](0.2),
	// 						Memory: to.Ptr("100Mi"),
	// 					},
	// 					Probes: []*armappcontainers.ContainerAppProbe{
	// 						{
	// 							Type: to.Ptr(armappcontainers.TypeLiveness),
	// 							HTTPGet: &armappcontainers.ContainerAppProbeHTTPGet{
	// 								Path: to.Ptr("/health"),
	// 								HTTPHeaders: []*armappcontainers.ContainerAppProbeHTTPGetHTTPHeadersItem{
	// 									{
	// 										Name: to.Ptr("Custom-Header"),
	// 										Value: to.Ptr("Awesome"),
	// 								}},
	// 								Port: to.Ptr[int32](8080),
	// 							},
	// 							InitialDelaySeconds: to.Ptr[int32](3),
	// 							PeriodSeconds: to.Ptr[int32](3),
	// 					}},
	// 			}},
	// 			InitContainers: []*armappcontainers.InitContainer{
	// 				{
	// 					Name: to.Ptr("testinitcontainerApp0"),
	// 					Args: []*string{
	// 						to.Ptr("-c"),
	// 						to.Ptr("while true; do echo hello; sleep 10;done")},
	// 						Command: []*string{
	// 							to.Ptr("/bin/sh")},
	// 							Image: to.Ptr("repo/testcontainerApp0:v4"),
	// 							Resources: &armappcontainers.ContainerResources{
	// 								CPU: to.Ptr[float64](0.2),
	// 								Memory: to.Ptr("100Mi"),
	// 							},
	// 					}},
	// 					Scale: &armappcontainers.Scale{
	// 						MaxReplicas: to.Ptr[int32](5),
	// 						MinReplicas: to.Ptr[int32](1),
	// 						Rules: []*armappcontainers.ScaleRule{
	// 							{
	// 								Name: to.Ptr("httpscalingrule"),
	// 								HTTP: &armappcontainers.HTTPScaleRule{
	// 									Metadata: map[string]*string{
	// 										"concurrentRequests": to.Ptr("50"),
	// 									},
	// 								},
	// 						}},
	// 					},
	// 				},
	// 				WorkloadProfileName: to.Ptr("My-GP-01"),
	// 			},
	// 		}
}
