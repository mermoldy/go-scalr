package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestProviderConfigurationCreateScalr(t *testing.T) {
	client := testClient(t)
	scalrHostname := client.baseURL.Host
	scalrToken := client.token
	ctx := context.Background()

	t.Run("success scalr", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:               &Account{ID: defaultAccountID},
			Name:                  String("scalr_dev"),
			ProviderName:          String("scalr"),
			ExportShellVariables:  Bool(false),
			ScalrHostname: 	       String(scalrHostname),
			ScalrToken: 	       String(scalrToken),

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
		assert.Equal(t, *options.ScalrHostname, pcfg.ScalrHostname)
		assert.Equal(t, "", pcfg.ScalrToken)
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
	t.Run("success azurerm", func(t *testing.T) {
		options := ProviderConfigurationCreateOptions{
			Account:               &Account{ID: defaultAccountID},
			Name:                  String("azurermdev"),
			ProviderName:          String("azurerm"),
			ExportShellVariables:  Bool(false),
			AzurermClientId:       String("my-client-id"),
			AzurermClientSecret:   String("my-client-secret"),
			AzurermSubscriptionId: String("my-subscription-id"),
			AzurermTenantId:       String("my-azurerm-tenant-id"),
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
			{Name: "scalr_dev", ProviderName: "scalr"},
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
	t.Run("success azurerm", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "azurerm", "azurerm_dev",
		)
		defer removeConfiguration()

		options := ProviderConfigurationUpdateOptions{
			Name:                  String("azure_dev2"),
			ExportShellVariables:  Bool(true),
			AzurermClientId:       String("my-client-id"),
			AzurermClientSecret:   String("my-client-secret-"),
			AzurermSubscriptionId: String("my-subscription-id"),
			AzurermTenantId:       String("my-tenant-id"),
		}
		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, options,
		)
		require.NoError(t, err)
		assert.Equal(t, *options.Name, updatedConfiguration.Name)
		assert.Equal(t, *options.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *options.AzurermClientId, updatedConfiguration.AzurermClientId)
		assert.Equal(t, "", updatedConfiguration.AzurermClientSecret)
		assert.Equal(t, *options.AzurermSubscriptionId, updatedConfiguration.AzurermSubscriptionId)
		assert.Equal(t, *options.AzurermTenantId, updatedConfiguration.AzurermTenantId)
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

func TestProviderConfigurationUpdateScalr(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	scalrHostname := client.baseURL.Host
	scalrToken := client.token

	t.Run("success scalr", func(t *testing.T) {
		createOptions := ProviderConfigurationCreateOptions{
			Account:              &Account{ID: defaultAccountID},
			Name:                 String("scalr_dev"),
			ProviderName:         String("scalr"),
			ExportShellVariables: Bool(false),
			ScalrHostname: 	  	  String(scalrHostname), 
			ScalrToken:           String(scalrToken),
			
		}
		configuration, err := client.ProviderConfigurations.Create(ctx, createOptions)
		if err != nil {
			t.Fatal(err)
		}
		defer client.ProviderConfigurations.Delete(ctx, configuration.ID)

		updateOptions := ProviderConfigurationUpdateOptions{
			Name:                 String("scalr_prod"),
			ExportShellVariables: Bool(true),
			ScalrHostname: 	  	  String(scalrHostname+"/"),
			ScalrToken: 		  String(scalrToken),
		}
		updatedConfiguration, err := client.ProviderConfigurations.Update(
			ctx, configuration.ID, updateOptions,
		)
		require.NoError(t, err)
		assert.Equal(t, *updateOptions.Name, updatedConfiguration.Name)
		assert.Equal(t, *updateOptions.ExportShellVariables, updatedConfiguration.ExportShellVariables)
		assert.Equal(t, *updateOptions.ScalrHostname, updatedConfiguration.ScalrHostname)
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
