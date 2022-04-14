package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderConfigurationCreate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()
	t.Run("success aws", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("AWS dev account us-east-1"),
			ProviderType:         String("aws"),
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
		assert.Equal(t, *options.ProviderType, pcfg.ProviderType)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.AwsAccessKey, pcfg.AwsAccessKey)
		assert.Equal(t, "", pcfg.AwsSecretKey)
	})
	t.Run("success azurerm", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("azure dev"),
			ProviderType:         String("azure"),
			ExportShellVariables: Bool(false),
			AzureClientId:        String("my-client-id"),
			AzureClientSecret:    String("my-client-secret"),
			AzureSubscriptionId:  String("my-subscription-id"),
			AzureTenantId:        String("my-azure-tenant-id"),
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
		assert.Equal(t, *options.ProviderType, pcfg.ProviderType)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.AzureClientId, pcfg.AzureClientId)
		assert.Equal(t, "", pcfg.AzureClientSecret)
		assert.Equal(t, *options.AzureSubscriptionId, pcfg.AzureSubscriptionId)
		assert.Equal(t, *options.AzureTenantId, pcfg.AzureTenantId)
	})
	t.Run("success google", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("AWS dev account us-east-1"),
			ProviderType:         String("google"),
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
		assert.Equal(t, *options.ProviderType, pcfg.ProviderType)
		assert.Equal(t, *options.ExportShellVariables, pcfg.ExportShellVariables)
		assert.Equal(t, *options.GoogleProject, pcfg.GoogleProject)
		assert.Equal(t, "", pcfg.GoogleCredentials)
	})
}

func TestProviderConfigurationList(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	t.Run("filtering", func(t *testing.T) {
		type providerTestingData struct {
			Name         string
			ProviderType string
		}
		providerTestingDataSet := []providerTestingData{
			{Name: "aws_prod_us_east_1", ProviderType: "aws"},
			{Name: "aws_prod_us_east_2", ProviderType: "aws"},
			{Name: "aws_dev_us_east1", ProviderType: "aws"},
			{Name: "gc_prod_us_west1_b", ProviderType: "google"},
		}

		for _, providerData := range providerTestingDataSet {
			_, removeConfiguration := createProviderConfiguration(
				t, client, providerData.ProviderType, providerData.Name,
			)

			defer removeConfiguration()
		}

		requestOptions := ProviderConfigurationsListOptions{
			Filter: &ProviderConfigurationFilter{
				ProviderType: "aws",
				Name:         "like:_prod_",
			},
		}
		configurationsList, err := client.ProviderConfigurations.List(ctx, requestOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(configurationsList.Items))

		resultNames := make([]string, 2)
		for _, configuration := range configurationsList.Items {
			resultNames = append(resultNames, configuration.Name)
		}
		assert.Contains(t, resultNames, "aws_prod_us_east_1")
		assert.Contains(t, resultNames, "aws_prod_us_east_2")
	})
}

func TestProviderConfigurationUpdate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
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
	t.Run("success azurerm", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "azure", "azurerm_dev",
		)
		defer removeConfiguration()

		options := ProviderConfigurationUpdateOptions{
			Name:                 String("azure_dev2"),
			ExportShellVariables: Bool(true),
			AzureClientId:        String("my-client-id"),
			AzureClientSecret:    String("my-client-secret-"),
			AzureSubscriptionId:  String("my-subscription-id"),
			AzureTenantId:        String("my-tenant-id"),
		}
		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, options,
		)
		require.NoError(t, err)
		assert.Equal(t, *options.Name, updatedConfiguration.Name)
		assert.Equal(t, *options.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *options.AzureClientId, updatedConfiguration.AzureClientId)
		assert.Equal(t, "", updatedConfiguration.AzureClientSecret)
		assert.Equal(t, *options.AzureSubscriptionId, updatedConfiguration.AzureSubscriptionId)
		assert.Equal(t, *options.AzureTenantId, updatedConfiguration.AzureTenantId)
	})
	t.Run("success google", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "google", "google_dev",
		)
		defer removeConfiguration()

		options := ProviderConfigurationUpdateOptions{
			Name:                 String("azure_dev2"),
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
	client.headers.Set("Prefer", "profile=internal")
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
