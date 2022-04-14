package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderConfigurationLinkCreate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	t.Run("workspace link", func(t *testing.T) {
		environment, removeEnvironment := createEnvironment(t, client)
		defer removeEnvironment()

		configuration, deleteConfiguration := createProviderConfiguration(
			t, client, "kubernetes", "kubernetes_dev",
		)
		defer deleteConfiguration()

		workspace, deleteWorkspace := createWorkspace(t, client, environment)
		defer deleteWorkspace()

		environmentOptions := ProviderConfigurationEnvironmentLinkCreateOptions{
			Default: Bool(true), ProviderConfiguration: configuration,
		}

		environmentLink, err := client.ProviderConfigurationLinks.EnvironmentCreate(
			ctx, environment.ID, environmentOptions,
		)
		require.NoError(t, err)

		environmentLink, err = client.ProviderConfigurationLinks.Read(ctx, environmentLink.ID)
		require.NoError(t, err)

		assert.Equal(t, *environmentOptions.Default, environmentLink.Default)

		workspaceOptions := ProviderConfigurationWorkspaceLinkCreateOptions{
			Alias: String("dev"), ProviderConfiguration: configuration,
		}

		workspacelink, err := client.ProviderConfigurationLinks.WorkspaceCreate(
			ctx, workspace.ID, workspaceOptions,
		)
		require.NoError(t, err)

		workspacelink, err = client.ProviderConfigurationLinks.Read(ctx, workspacelink.ID)
		require.NoError(t, err)

		assert.Equal(t, *workspaceOptions.Alias, workspacelink.Alias)
	})
}

func TestProviderConfigurationLinkUpdate(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	environment, removeEnvironment := createEnvironment(t, client)
	defer removeEnvironment()

	configuration, deleteConfiguration := createProviderConfiguration(
		t, client, "kubernetes", "kubernetes_dev",
	)
	defer deleteConfiguration()

	t.Run("environment link", func(t *testing.T) {
		environmentOptions := ProviderConfigurationEnvironmentLinkCreateOptions{
			Default: Bool(true), ProviderConfiguration: configuration,
		}

		environmentLink, err := client.ProviderConfigurationLinks.EnvironmentCreate(
			ctx, environment.ID, environmentOptions,
		)
		require.NoError(t, err)

		updateOptions := ProviderConfigurationLinkUpdateOptions{Default: Bool(false)}

		environmentLink, err = client.ProviderConfigurationLinks.Update(ctx, environmentLink.ID, updateOptions)
		require.NoError(t, err)

		assert.Equal(t, *updateOptions.Default, environmentLink.Default)
	})
	t.Run("workspace link", func(t *testing.T) {
		workspace, deleteWorkspace := createWorkspace(t, client, environment)
		defer deleteWorkspace()

		environmentOptions := ProviderConfigurationEnvironmentLinkCreateOptions{
			Default: Bool(true), ProviderConfiguration: configuration,
		}

		_, err := client.ProviderConfigurationLinks.EnvironmentCreate(
			ctx, environment.ID, environmentOptions,
		)
		require.NoError(t, err)

		workspaceOptions := ProviderConfigurationWorkspaceLinkCreateOptions{
			Alias: String("dev"), ProviderConfiguration: configuration,
		}

		workspacelink, err := client.ProviderConfigurationLinks.WorkspaceCreate(
			ctx, workspace.ID, workspaceOptions,
		)
		require.NoError(t, err)

		updateOptions := ProviderConfigurationLinkUpdateOptions{Alias: String("demo"), Default: Bool(false)}

		workspacelink, err = client.ProviderConfigurationLinks.Update(ctx, workspacelink.ID, updateOptions)
		require.NoError(t, err)

		assert.Equal(t, *updateOptions.Alias, workspacelink.Alias)
	})
}

func TestProviderConfigurationLinkDelete(t *testing.T) {
	client := testClient(t)
	client.headers.Set("Prefer", "profile=internal")
	ctx := context.Background()

	environment, removeEnvironment := createEnvironment(t, client)
	defer removeEnvironment()

	configuration, deleteConfiguration := createProviderConfiguration(
		t, client, "kubernetes", "kubernetes_dev",
	)
	defer deleteConfiguration()

	t.Run("success", func(t *testing.T) {
		options := ProviderConfigurationEnvironmentLinkCreateOptions{
			Default: Bool(true), ProviderConfiguration: configuration,
		}

		link, err := client.ProviderConfigurationLinks.EnvironmentCreate(
			ctx, environment.ID, options,
		)
		require.NoError(t, err)

		err = client.ProviderConfigurationLinks.Delete(ctx, link.ID)
		require.NoError(t, err)

		// Try loading the configuration - it should fail.
		_, err = client.ProviderConfigurationLinks.Read(ctx, link.ID)
		assert.Equal(
			t,
			ResourceNotFoundError{
				Message: fmt.Sprintf("ProviderConfigurationLink with ID '%s' not found or user unauthorized", link.ID),
			}.Error(),
			err.Error(),
		)
	})
}
