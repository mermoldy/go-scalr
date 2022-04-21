package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderConfigurationParameterCreate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	configuration, deleteConfiguration := createProviderConfiguration(
		t, client, "kubernetes", "kubernetes_dev",
	)
	defer deleteConfiguration()

	t.Run("success", func(t *testing.T) {
		options := ProviderConfigurationParameterCreateOptions{
			Key:         String("config_path"),
			Sensitive:   Bool(false),
			Value:       String("~/.kube/config"),
			Description: String("A path to a kube config file."),
		}
		parameter, err := client.ProviderConfigurationParameters.Create(ctx, configuration.ID, options)
		if err != nil {
			t.Fatal(err)
		}

		parameter, err = client.ProviderConfigurationParameters.Read(ctx, parameter.ID)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, parameter.Key)
		assert.Equal(t, *options.Sensitive, parameter.Sensitive)
		assert.Equal(t, *options.Value, parameter.Value)
		assert.Equal(t, *options.Description, parameter.Description)
	})

	t.Run("success sensitive", func(t *testing.T) {
		options := ProviderConfigurationParameterCreateOptions{
			Key:       String("client_certificate"),
			Sensitive: Bool(true),
			Value:     String("-----BEGIN CERTIFICATE-----\nMIIB9TCCAWACAQAwgbgxGTAXB"),
		}
		parameter, err := client.ProviderConfigurationParameters.Create(ctx, configuration.ID, options)
		if err != nil {
			t.Fatal(err)
		}

		parameter, err = client.ProviderConfigurationParameters.Read(ctx, parameter.ID)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, parameter.Key)
		assert.Equal(t, *options.Sensitive, parameter.Sensitive)
		assert.Equal(t, "", parameter.Value)
	})
}

func TestProviderConfigurationParameterCreateMany(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
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
		createdParameters, err := client.ProviderConfigurationParameters.CreateMany(ctx, configuration.ID, optionsList)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(optionsList), len(createdParameters))

		optionsMap := make(map[string]ProviderConfigurationParameterCreateOptions)
		for _, createOption := range optionsList {
			optionsMap[*createOption.Key] = createOption
		}

		for _, createdParameter := range createdParameters {
			parameter, err := client.ProviderConfigurationParameters.Read(ctx, createdParameter.ID)
			require.NoError(t, err)

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
		_, err := client.ProviderConfigurationParameters.CreateMany(ctx, configuration.ID, optionsList)
		require.Error(t, err)
		assert.EqualError(t, err, "Invalid Attribute\n\nCan not create parameter. Key 'config_path' has already been taken.")
	})

}

func TestProviderConfigurationParametersList(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes dev",
		)
		defer removeConfiguration()

		providerTestingDataSet := []struct {
			Key       string
			Value     string
			Sensitive bool
		}{
			{Key: "config_path", Value: "~/.kube/config", Sensitive: false},
			{Key: "config_context", Value: "my-context", Sensitive: false},
			{Key: "client_certificate", Value: "--BEGIN CERTIFICATE--\nMIIB9", Sensitive: true},
		}

		var createdParameterIDs []string
		for _, parameterData := range providerTestingDataSet {
			parameter, err := client.ProviderConfigurationParameters.Create(ctx, configuration.ID, ProviderConfigurationParameterCreateOptions{
				Key:       String(parameterData.Key),
				Sensitive: Bool(parameterData.Sensitive),
				Value:     String(parameterData.Value),
			})
			if err != nil {
				t.Fatal(err)
			}

			createdParameterIDs = append(createdParameterIDs, parameter.ID)
		}

		parametersList, err := client.ProviderConfigurationParameters.List(ctx, configuration.ID, ProviderConfigurationParametersListOptions{})

		require.NoError(t, err)
		assert.Equal(t, 3, len(parametersList.Items))

		var resultIDs []string
		for _, configuration := range parametersList.Items {
			resultIDs = append(resultIDs, configuration.ID)
		}

		for _, parameterID := range createdParameterIDs {
			assert.Contains(t, resultIDs, parameterID)
		}
	})
}

func TestProviderConfigurationParameterUpdate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	t.Run("success all attributes", func(t *testing.T) {
		configuration, removeConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes dev",
		)
		defer removeConfiguration()

		parameter, err := client.ProviderConfigurationParameters.Create(ctx, configuration.ID, ProviderConfigurationParameterCreateOptions{
			Key:       String("config_context"),
			Sensitive: Bool(true),
			Value:     String("my-context"),
		})
		if err != nil {
			t.Fatal(err)
		}

		options := ProviderConfigurationParameterUpdateOptions{
			Key:         String("config_path"),
			Sensitive:   Bool(false),
			Value:       String("~/.kube/config"),
			Description: String("A path to a kube config file."),
		}
		updatedParameter, err := client.ProviderConfigurationParameters.Update(
			ctx, parameter.ID, options,
		)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, updatedParameter.Key)
		assert.Equal(t, *options.Sensitive, updatedParameter.Sensitive)
		assert.Equal(t, *options.Value, updatedParameter.Value)
		assert.Equal(t, *options.Description, updatedParameter.Description)
	})
}

func TestProviderConfigurationParameterDelete(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	configuration, removeConfiguration := createProviderConfiguration(t, client, "kubernetes", "kubernetes dev")
	defer removeConfiguration()

	t.Run("success", func(t *testing.T) {
		parameter, err := client.ProviderConfigurationParameters.Create(ctx, configuration.ID, ProviderConfigurationParameterCreateOptions{
			Key:       String("config_context"),
			Sensitive: Bool(true),
			Value:     String("my-context"),
		})
		require.NoError(t, err)

		err = client.ProviderConfigurationParameters.Delete(ctx, parameter.ID)
		require.NoError(t, err)

		// Try loading the configuration - it should fail.
		_, err = client.ProviderConfigurationParameters.Read(ctx, parameter.ID)
		assert.Equal(
			t,
			ResourceNotFoundError{
				Message: fmt.Sprintf("ProviderConfigurationParameter with ID '%s' not found or user unauthorized", parameter.ID),
			}.Error(),
			err.Error(),
		)
	})
}
