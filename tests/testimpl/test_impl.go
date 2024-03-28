package common

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/lib/azure/configure"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

const terraformDir string = "../../examples/route"
const varFile string = "test.tfvars"

func TestRoutes(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		t.Fatalf("ARM_SUBSCRIPTION_ID must be set for acceptance tests")
	}
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Error getting credentials: %e\n", err)
	}
	routeTablesClientOptions := policy.ClientOptions{}
	clientFactory, err := armnetwork.NewClientFactory(subscriptionID, credential, &routeTablesClientOptions)
	if err != nil {
		t.Fatalf("Error creating client factory: %e\n", err)
	}
	routeTableClient := clientFactory.NewRouteTablesClient()
	terraformOptions := configure.ConfigureTerraform(terraformDir, []string{terraformDir + "/" + varFile})
	t.Run("doesRouteExist", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
		routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "route_table_name")
		routeNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "route_names")
		expectedRouteNames := make(map[string]bool)
		options := armnetwork.RouteTablesClientGetOptions{}

		// Set all expected routes to false
		for _, v := range routeNames {
			expectedRouteNames[v] = false
		}

		routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, &options)
		if err != nil {
			t.Fatalf("Error getting Route Table: %v", err)
		}
		if routeTable.Name == nil {
			t.Fatalf("Route Table does not exist")
		}

		// If expected route matches with actual route, set the map value to true
		for expectedRoute := range expectedRouteNames {
			// routes := routeTable.RouteTablePropertiesFormat.Routes
			routes := routeTable.RouteTable.Properties.Routes

			for _, route := range routes {
				if expectedRoute == *route.Name {
					expectedRouteNames[expectedRoute] = true
					break
				}
			}
		}
		// Check if all routes are true. If not, test fails
		for route, found := range expectedRouteNames {
			assert.True(t, found, "Route is not found: "+route)
		}
		assert.Equal(t, *routeTable.Name, routeTableName)
	})
}
