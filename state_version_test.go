package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetStateVersionCreateOptions(workspace *Workspace, run *Run) StateVersionCreateOptions {
	return StateVersionCreateOptions{
		Force:   true,
		Serial:  0,
		Lineage: "3ab79560-29ab-6eb4-ca13-1ee3dafa1fc7",
		State:   "eyJ2ZXJzaW9uIjogNCwgInRlcnJhZm9ybV92ZXJzaW9uIjogIjAuMTQuMiIsICJzZXJpYWwiOiAwLCAibGluZWFnZSI6ICIzYWI3OTU2MC0yOWFiLTZlYjQtY2ExMy0xZWUzZGFmYTFmYzciLCAib3V0cHV0cyI6IHsicmVzIjogeyJ2YWx1ZSI6ICJ9cUd3dDJQXVpMLWg/STZtXVgiLCAidHlwZSI6ICJzdHJpbmcifX0sICJyZXNvdXJjZXMiOiBbeyJtb2RlIjogIm1hbmFnZWQiLCAidHlwZSI6ICJyYW5kb21fc3RyaW5nIiwgIm5hbWUiOiAicmFuZG9tIiwgInByb3ZpZGVyIjogInByb3ZpZGVyW1wicmVnaXN0cnkudGVycmFmb3JtLmlvL2hhc2hpY29ycC9yYW5kb21cIl0iLCAiaW5zdGFuY2VzIjogW3sic2NoZW1hX3ZlcnNpb24iOiAxLCAic2Vuc2l0aXZlX2F0dHJpYnV0ZXMiOiBbXSwgInByaXZhdGUiOiAiZXlKelkyaGxiV0ZmZG1WeWMybHZiaUk2SWpFaWZRPT0ifV19XX0=",
		MD5:     "980638c48ef16d8ed6eaaeeaddacd910",
		Size:    815,
		Resources: []*Resource{{
			Type:    "null",
			Address: "test",
		}},
		Workspace: workspace,
		Run:       run,
	}
}

func TestStateVersionCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	envTest, envTestCleanup := createEnvironment(t, client)
	defer envTestCleanup()

	wsTest, wsTestCleanup := createWorkspace(t, client, envTest)
	defer wsTestCleanup()

	cvTest, cvTestCleunup := createConfigurationVersion(t, client, wsTest)
	defer cvTestCleunup()

	runTest, runTestCleanup := createRun(t, client, wsTest, cvTest)
	defer runTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := GetStateVersionCreateOptions(wsTest, runTest)
		client.headers.Set("Prefer", "profile=internal")
		sv, err := client.StateVersions.Create(ctx, options)
		client.headers.Set("Prefer", "profile=preview")
		require.NoError(t, err)

		// // Get a refreshed view from the API.
		refreshed, err := client.StateVersions.ReadByID(ctx, sv.ID)
		require.NoError(t, err)

		for _, item := range []*StateVersion{
			sv,
			refreshed,
		} {
			assert.NotEmpty(t, item.ID)
			assert.Equal(t, options.Force, item.Force)
			assert.Equal(t, options.Serial, item.Serial)
			assert.Equal(t, options.Run.ID, item.Run.ID)
			assert.Equal(t, options.Workspace.ID, item.Workspace.ID)
		}
	})

	t.Run("with invalid state", func(t *testing.T) {
		options := GetStateVersionCreateOptions(wsTest, runTest)
		options.State = "12"
		_, err := client.StateVersions.Create(ctx, options)
		require.Error(t, err)
	})
}

func TestReadCurrentFromWorkspace(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	envTest, envTestCleanup := createEnvironment(t, client)
	defer envTestCleanup()

	wsTest, wsTestCleanup := createWorkspace(t, client, envTest)
	defer wsTestCleanup()

	cvTest, cvTestCleunup := createConfigurationVersion(t, client, wsTest)
	defer cvTestCleunup()

	runTest, runTestCleanup := createRun(t, client, wsTest, cvTest)
	defer runTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := GetStateVersionCreateOptions(wsTest, runTest)
		client.headers.Set("Prefer", "profile=internal")
		sv, err := client.StateVersions.Create(ctx, options)
		client.headers.Set("Prefer", "profile=preview")
		require.NoError(t, err)

		// // Get a refreshed view from the API.
		refreshed, err := client.StateVersions.ReadCurrentFromWorkspace(ctx, wsTest.ID)
		require.NoError(t, err)

		for _, item := range []*StateVersion{
			sv,
			refreshed,
		} {
			assert.NotEmpty(t, item.ID)
			assert.Equal(t, options.Force, item.Force)
			assert.Equal(t, options.Serial, item.Serial)
			assert.Equal(t, options.Run.ID, item.Run.ID)
			assert.Equal(t, options.Workspace.ID, item.Workspace.ID)
		}
	})

	t.Run("with invalid workspace", func(t *testing.T) {
		_, err := client.StateVersions.ReadCurrentFromWorkspace(ctx, "invalid")
		require.Error(t, err)
	})
}
