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
			Name:                  String("azurerm dev"),
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
		expectedParameters, err := client.ProviderConfigurations.CreateParameters(
			ctx, configuration.ID, &optionsList,
		)
		require.NoError(t, err)

		configuration, err = client.ProviderConfigurations.Read(ctx, configuration.ID)
		require.NoError(t, err)
		assert.Equal(t, len(expectedParameters), len(configuration.Parameters))

		includedParameters := make(map[string]ProviderConfigurationParameter)
		for _, p := range configuration.Parameters {
			includedParameters[p.ID] = *p
		}

		for _, parameter := range expectedParameters {
			includedParameter := includedParameters[parameter.ID]
			assert.Equal(t, parameter.Key, includedParameter.Key)
			assert.Equal(t, parameter.Description, includedParameter.Description)
			assert.Equal(t, parameter.Sensitive, includedParameter.Sensitive)
		}

	})
}

func TestProviderConfigurationCreateParameters(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		configuration, deleteConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes_dev",
		)
		defer deleteConfiguration()

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
			{
				Key:   String("username"),
				Value: String("my-username"),
			},

			{
				Key:   String("password"),
				Value: String("my-password"),
			},

			{
				Key:   String("insecure"),
				Value: String("my-insecure"),
			},
			{
				Key:   String("config_context"),
				Value: String("my-config_context"),
			},
			{
				Key:   String("config_context_auth_info"),
				Value: String("my-config_context_auth_info"),
			},
			{
				Key:   String("config_context_cluster"),
				Value: String("myconfig_context_cluster"),
			},
			{
				Key:   String("token"),
				Value: String("my-token"),
			},

			{
				Key:   String("proxy_url"),
				Value: String("my-proxy_url"),
			},
			{
				Key:   String("exec"),
				Value: String("my-exec"),
			},
		}
		createdParameters, err := client.ProviderConfigurations.CreateParameters(ctx, configuration.ID, &optionsList)
		require.NoError(t, err)

		assert.Equal(t, len(optionsList), len(createdParameters))

		optionsMap := make(map[string]ProviderConfigurationParameterCreateOptions)
		for _, createOption := range optionsList {
			optionsMap[*createOption.Key] = createOption
		}

		configuration, err = client.ProviderConfigurations.Read(ctx, configuration.ID)
		require.NoError(t, err)

		for _, parameter := range configuration.Parameters {
			createOption := optionsMap[parameter.Key]
			if createOption.Sensitive != nil && *createOption.Sensitive {
				assert.Equal(t, "", parameter.Value)
			} else {
				assert.Equal(t, *createOption.Value, parameter.Value)
			}
		}
	})
	t.Run("key duplication error", func(t *testing.T) {
		configuration, deleteConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes_dev",
		)
		defer deleteConfiguration()

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
				Key:   String("config_path"),
				Value: String("my-another-path"),
			},
		}
		_, err := client.ProviderConfigurations.CreateParameters(ctx, configuration.ID, &optionsList)
		require.Error(t, err)
		assert.EqualError(t, err, "Invalid Attribute\n\nCan not create parameter. Key 'config_path' has already been taken.")
	})

}

func TestProviderConfigurationChangeParameters(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		configuration, deleteConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes_dev",
		)
		defer deleteConfiguration()

		initParameterOptions := []ProviderConfigurationParameterCreateOptions{
			{
				Key:       String("client_certificate"),
				Sensitive: Bool(true),
				Value:     String("-----BEGIN CERTIFICATE-----\nMIIB9TCCAWACAQAwgbgxGTAXB"),
			},
			{
				Key:   String("host"),
				Value: String("my-host"),
			},
			{
				Key:   String("username"),
				Value: String("my-username"),
			},
		}

		toCreate := []ProviderConfigurationParameterCreateOptions{
			{
				Key:   String("insecure"),
				Value: String("my-insecure"),
			},
		}
		var toUpdate []ProviderConfigurationParameterUpdateOptions
		var toDelete []string

		initParameters, err := client.ProviderConfigurations.CreateParameters(ctx, configuration.ID, &initParameterOptions)
		require.NoError(t, err)

		var notChanged ProviderConfigurationParameter

		for _, param := range initParameters {
			if param.Key == "username" {
				notChanged = param
			}
			if param.Key == "host" {
				toUpdate = append(toUpdate, ProviderConfigurationParameterUpdateOptions{
					ID:    param.ID,
					Key:   String("host"),
					Value: String("new-host"),
				})
			} else if param.Key == "client_certificate" {
				toDelete = append(toDelete, param.ID)
			}
		}
		assert.Equal(t, 1, len(toUpdate))
		assert.Equal(t, 1, len(toDelete))

		created, updated, deleted, err := client.ProviderConfigurations.ChangeParameters(
			ctx,
			configuration.ID,
			&toCreate,
			&toUpdate,
			&toDelete,
		)
		require.NoError(t, err)

		assert.ElementsMatch(t, deleted, toDelete)

		assert.Equal(t, len(updated), 1)
		assert.Equal(t, "host", updated[0].Key)
		assert.Equal(t, "new-host", updated[0].Value)

		assert.Equal(t, len(created), 1)
		assert.Equal(t, "insecure", created[0].Key)
		assert.Equal(t, "my-insecure", created[0].Value)

		configuration, err = client.ProviderConfigurations.Read(ctx, configuration.ID)
		require.NoError(t, err)

		expected := []ProviderConfigurationParameter{
			notChanged,
			created[0],
			updated[0],
		}

		var result []ProviderConfigurationParameter
		for _, param := range configuration.Parameters {
			result = append(result, *param)
		}

		assert.ElementsMatch(t, expected, result)
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
