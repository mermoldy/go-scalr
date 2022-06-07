package scalr

import (
	"context"
	"fmt"
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAzureTestingCreds(t *testing.T) (armClientId, armClientSecret, armSubscriptionId, armTenantId string) {
	armClientId = os.Getenv("TEST_ARM_CLIENT_ID")
	armClientSecret = os.Getenv("TEST_ARM_CLIENT_SECRET")
	armSubscriptionId = os.Getenv("TEST_ARM_SUBSCRIPTION_ID")
	armTenantId = os.Getenv("TEST_ARM_TENANT_ID")
	if len(armClientId) == 0 ||
		len(armClientSecret) == 0 ||
		len(armSubscriptionId) == 0 ||
		len(armTenantId) == 0 {
		t.Skip("Please set TEST_ARM_CLIENT_ID, TEST_ARM_CLIENT_SECRET, TEST_ARM_SUBSCRIPTION_ID and TEST_ARM_TENANT_ID env variables to run this test.")
	}
	return
}

func TestProviderConfigurationCreateAzurerm(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	armClientId, armClientSecret, armSubscriptionId, armTenantId := getAzureTestingCreds(t)

	t.Run("success azurerm", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:               &Account{ID: defaultAccountID},
			Name:                  String("azurerm_dev"),
			ProviderName:          String("azurerm"),
			ExportShellVariables:  Bool(false),
			AzurermClientId:       String(armClientId),
			AzurermClientSecret:   String(armClientSecret),
			AzurermSubscriptionId: String(armSubscriptionId),
			AzurermTenantId:       String(armTenantId),
		}
		pcfg, err := client.ProviderConfigurations.Create(ctx, options)
		if err != nil {
			t.Fatal(err)
		}
		defer client.ProviderConfigurations.Delete(ctx, pcfg.ID)

		pcfg, err = client.ProviderConfigurations.Read(ctx, pcfg.ID)
		require.NoError(t, err)

		assert.Equal(t, options.Account.ID, pcfg.Account.ID)
		assert.Equal(t, *options.Name, pcfg.Name)
		assert.Equal(t, *options.ProviderName, pcfg.ProviderName)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.AzurermClientId, pcfg.AzurermClientId)
		assert.Equal(t, "", pcfg.AzurermClientSecret)
		assert.Equal(t, *options.AzurermSubscriptionId, pcfg.AzurermSubscriptionId)
		assert.Equal(t, *options.AzurermTenantId, pcfg.AzurermTenantId)
	})
}

func TestProviderConfigurationCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	t.Run("success aws", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("AWS dev account us-east-1"),
			ProviderName:         String("aws"),
			ExportShellVariables: Bool(false),
			AwsAccessKey:         String("my-access-key"),
			AwsSecretKey:         String("my-secret-key"),
		}
		pcfg, err := client.ProviderConfigurations.Create(ctx, options)
		if err != nil {
			t.Fatal(err)
		}
		defer client.ProviderConfigurations.Delete(ctx, pcfg.ID)

		pcfg, err = client.ProviderConfigurations.Read(ctx, pcfg.ID)
		require.NoError(t, err)

		assert.Equal(t, options.Account.ID, pcfg.Account.ID)
		assert.Equal(t, *options.Name, pcfg.Name)
		assert.Equal(t, *options.ProviderName, pcfg.ProviderName)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.AwsAccessKey, pcfg.AwsAccessKey)
		assert.Equal(t, "", pcfg.AwsSecretKey)
	})
	t.Run("success google", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("AWS dev account us-east-1"),
			ProviderName:         String("google"),
			ExportShellVariables: Bool(false),
			GoogleProject:        String("my-google-project"),
			GoogleCredentials:    String("my-google-credentials"),
		}
		pcfg, err := client.ProviderConfigurations.Create(ctx, options)
		if err != nil {
			t.Fatal(err)
		}
		defer client.ProviderConfigurations.Delete(ctx, pcfg.ID)

		pcfg, err = client.ProviderConfigurations.Read(ctx, pcfg.ID)
		require.NoError(t, err)

		assert.Equal(t, options.Account.ID, pcfg.Account.ID)
		assert.Equal(t, *options.Name, pcfg.Name)
		assert.Equal(t, *options.ProviderName, pcfg.ProviderName)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.GoogleProject, pcfg.GoogleProject)
		assert.Equal(t, "", pcfg.GoogleCredentials)
	})
}

func TestProviderConfigurationRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with parameters", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes_dev",
		)
		defer removeConfiguration()

		optionsList := []ProviderConfigurationParameterCreateOptions{
			{
				Key:         String("config_path"),
				Sensitive:   Bool(false),
				Value:       String("~/.kube/config"),
				Description: String("A path to a kube config file."),
			},
			{
				Key:       String("client_certificate"),
				Sensitive: Bool(true),
				Value:     String("-----BEGIN CERTIFICATE-----\nMIIB9TCCAWACAQAwgbgxGTAXB"),
			},
			{
				Key:   String("host"),
				Value: String("my-host"),
			},
		}
		for _, option := range optionsList {
			_, err := client.ProviderConfigurationParameters.Create(
				ctx, configuration.ID, option,
			)
			require.NoError(t, err)
		}

		configuration, err := client.ProviderConfigurations.Read(ctx, configuration.ID)
		require.NoError(t, err)
		assert.Equal(t, len(optionsList), len(configuration.Parameters))

		includedParameters := make(map[string]ProviderConfigurationParameter)
		for _, p := range configuration.Parameters {
			includedParameters[p.Key] = *p
		}

		for _, option := range optionsList {
			includedParameter := includedParameters[*option.Key]
			assert.Equal(t, *option.Key, includedParameter.Key)
			var description string
			if option.Description != nil {
				description = *option.Description
			}
			assert.Equal(t, description, includedParameter.Description)
			assert.Equal(t, option.Sensitive != nil && *option.Sensitive, includedParameter.Sensitive)
		}
	})
}

func TestProviderConfigurationList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("filtering", func(t *testing.T) {
		type providerTestingData struct {
			Name         string
			ProviderName string
		}
		providerTestingDataSet := []providerTestingData{
			{Name: "aws_prod_us_east_1", ProviderName: "aws"},
			{Name: "aws_prod_us_east_2", ProviderName: "aws"},
			{Name: "aws_dev_us_east1", ProviderName: "aws"},
			{Name: "gc_prod_us_west1_b", ProviderName: "google"},
		}

		for _, providerData := range providerTestingDataSet {
			_, removeConfiguration := createProviderConfiguration(
				t, client, providerData.ProviderName, providerData.Name,
			)

			defer removeConfiguration()
		}

		requestOptions := ProviderConfigurationsListOptions{
			Filter: &ProviderConfigurationFilter{
				ProviderName: "aws",
				Name:         "like:_prod_",
			},
		}
		configurationsList, err := client.ProviderConfigurations.List(ctx, requestOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(configurationsList.Items))

		var resultNames []string
		for _, configuration := range configurationsList.Items {
			resultNames = append(resultNames, configuration.Name)
		}
		assert.Contains(t, resultNames, "aws_prod_us_east_1")
		assert.Contains(t, resultNames, "aws_prod_us_east_2")
	})
}

func TestProviderConfigurationUpdateAzurerm(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	armClientId, armClientSecret, armSubscriptionId, armTenantId := getAzureTestingCreds(t)

	
	t.Run("success", func(t *testing.T) {
		createOptions := ProviderConfigurationCreateOptions{
			Account:               &Account{ID: defaultAccountID},
			Name:                  String("azurerm_dev"),
			ProviderName:          String("azurerm"),
			ExportShellVariables:  Bool(false),
			AzurermClientId:       String(armClientId),
			AzurermClientSecret:   String(armClientSecret),
			AzurermSubscriptionId: String(armSubscriptionId),
			AzurermTenantId:       String(armTenantId),
		}
		configuration, err := client.ProviderConfigurations.Create(ctx, createOptions)
		if err != nil {
			t.Fatal(err)
		}
		defer client.ProviderConfigurations.Delete(ctx, configuration.ID)

		updateOptions := ProviderConfigurationUpdateOptions{
			Name:                 String("azurerm_dev_updated"),
			ExportShellVariables: Bool(true),
			AzurermClientId:       String(armClientId),
			AzurermClientSecret:   String(armClientSecret),
			AzurermSubscriptionId: String(armSubscriptionId),
			AzurermTenantId:       String(armTenantId),
		}

		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, updateOptions,
		)
		require.NoError(t, err)
		assert.Equal(t, *updateOptions.Name, updatedConfiguration.Name)
		assert.Equal(t, *updateOptions.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *updateOptions.AzurermClientId, updatedConfiguration.AzurermClientId)
		assert.Equal(t, "", updatedConfiguration.AzurermClientSecret)
		assert.Equal(t, *updateOptions.AzurermSubscriptionId, updatedConfiguration.AzurermSubscriptionId)
		assert.Equal(t, *updateOptions.AzurermTenantId, updatedConfiguration.AzurermTenantId)
	})
	
}
func TestProviderConfigurationUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("success aws", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "aws", "aws_dev_us_east_1",
		)
		defer removeConfiguration()

		options := ProviderConfigurationUpdateOptions{
			Name:                 String("aws_dev_us_east_2"),
			ExportShellVariables: Bool(true),
			AwsAccessKey:         String("my-aws-access-key1"),
			AwsSecretKey:         String("my-aws-secret-key"),
		}
		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, options,
		)
		require.NoError(t, err)
		assert.Equal(t, *options.Name, updatedConfiguration.Name)
		assert.Equal(t, *options.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *options.AwsAccessKey, updatedConfiguration.AwsAccessKey)
	})
	t.Run("success google", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "google", "google_dev",
		)
		defer removeConfiguration()

		options := ProviderConfigurationUpdateOptions{
			Name:                 String("azurerm_dev2"),
			ExportShellVariables: Bool(true),
			GoogleProject:        String("my-project"),
			GoogleCredentials:    String("my-credentials"),
		}
		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, options,
		)
		require.NoError(t, err)
		assert.Equal(t, *options.Name, updatedConfiguration.Name)
		assert.Equal(t, *options.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *options.GoogleProject, updatedConfiguration.GoogleProject)
		assert.Equal(t, "", updatedConfiguration.GoogleCredentials)
	})
}

func TestProviderConfigurationDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	configuration, _ := createProviderConfiguration(t, client, "aws", "aws_dev_us_east_1")

	t.Run("success", func(t *testing.T) {
		err := client.ProviderConfigurations.Delete(ctx, configuration.ID)
		require.NoError(t, err)

		// Try loading the configuration - it should fail.
		_, err = client.ProviderConfigurations.Read(ctx, configuration.ID)
		assert.Equal(
			t,
			ResourceNotFoundError{
				Message: fmt.Sprintf("ProviderConfiguration with ID '%s' not found or user unauthorized", configuration.ID),
			}.Error(),
			err.Error(),
		)
	})
}
