package common

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/configure"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/login"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/network"
	"github.com/nexient-llc/lcaf-component-terratest-common/types"
	"github.com/stretchr/testify/assert"
)

const terraformDir string = "../../examples/route"
const varFile string = "test.tfvars"

func TestRoutes(t *testing.T, ctx types.TestContext) {

	envVarMap := login.GetEnvironmentVariables()
	clientID := envVarMap["clientID"]
	clientSecret := envVarMap["clientSecret"]
	tenantID := envVarMap["tenantID"]
	subscriptionID := envVarMap["subscriptionID"]

	spt, err := login.GetServicePrincipalToken(clientID, clientSecret, tenantID)
	if err != nil {
		t.Fatalf("Error getting Service Principal Token: %v", err)
	}

	routeTableClient := network.GetRouteTablesClient(spt, subscriptionID)
	terraformOptions := configure.ConfigureTerraform(terraformDir, []string{terraformDir + "/" + varFile})
	t.Run("doesRouteExist", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
		routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "route_table_name")
		routeNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "route_names")
		expectedRouteNames := make(map[string]bool)

		// Set all expected routes to false
		for _, v := range routeNames {
			expectedRouteNames[v] = false
		}

		routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, "")
		if err != nil {
			t.Fatalf("Error getting Route Table: %v", err)
		}
		if routeTable.Name == nil {
			t.Fatalf("Route Table does not exist")
		}

		// If expected route matches with actual route, set the map value to true
		for expectedRoute := range expectedRouteNames {
			routes := routeTable.RouteTablePropertiesFormat.Routes

			for _, route := range *routes {
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
