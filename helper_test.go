package scalr

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
)

const badIdentifier = "! / nope"

func testClient(t *testing.T) *Client {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func createConfigurationVersion(t *testing.T, client *Client, w *Workspace) (*ConfigurationVersion, func()) {
	var wCleanup func()

	if w == nil {
		w, wCleanup = createWorkspace(t, client, nil)
	}

	ctx := context.Background()
	cv, err := client.ConfigurationVersions.Create(
		ctx,
		w.ID,
		ConfigurationVersionCreateOptions{AutoQueueRuns: Bool(false)},
	)
	if err != nil {
		t.Fatal(err)
	}

	return cv, func() {
		if wCleanup != nil {
			wCleanup()
		}
	}
}

func createUploadedConfigurationVersion(t *testing.T, client *Client, w *Workspace) (*ConfigurationVersion, func()) {
	cv, cvCleanup := createConfigurationVersion(t, client, w)

	ctx := context.Background()
	err := client.ConfigurationVersions.Upload(ctx, cv.UploadURL, "test-fixtures/config-version")
	if err != nil {
		cvCleanup()
		t.Fatal(err)
	}

	for i := 0; ; i++ {
		cv, err = client.ConfigurationVersions.Read(ctx, cv.ID)
		if err != nil {
			cvCleanup()
			t.Fatal(err)
		}

		if cv.Status == ConfigurationUploaded {
			break
		}

		if i > 10 {
			cvCleanup()
			t.Fatal("Timeout waiting for the configuration version to be uploaded")
		}

		time.Sleep(1 * time.Second)
	}

	return cv, cvCleanup
}

func createRun(t *testing.T, client *Client, w *Workspace) (*Run, func()) {
	var wCleanup func()

	if w == nil {
		w, wCleanup = createWorkspace(t, client, nil)
	}

	cv, cvCleanup := createUploadedConfigurationVersion(t, client, w)

	ctx := context.Background()
	r, err := client.Runs.Create(ctx, RunCreateOptions{
		ConfigurationVersion: cv,
		Workspace:            w,
	})
	if err != nil {
		t.Fatal(err)
	}

	return r, func() {
		if wCleanup != nil {
			wCleanup()
		} else {
			cvCleanup()
		}
	}
}

func createPlannedRun(t *testing.T, client *Client, w *Workspace) (*Run, func()) {
	r, rCleanup := createRun(t, client, w)

	var err error
	ctx := context.Background()
	for i := 0; ; i++ {
		r, err = client.Runs.Read(ctx, r.ID)
		if err != nil {
			t.Fatal(err)
		}

		switch r.Status {
		case RunPlanned, RunCostEstimated, RunPolicyChecked, RunPolicyOverride:
			return r, rCleanup
		}

		if i > 45 {
			rCleanup()
			t.Fatal("Timeout waiting for run to be planned")
		}

		time.Sleep(1 * time.Second)
	}
}

func createCostEstimatedRun(t *testing.T, client *Client, w *Workspace) (*Run, func()) {
	r, rCleanup := createRun(t, client, w)

	var err error
	ctx := context.Background()
	for i := 0; ; i++ {
		r, err = client.Runs.Read(ctx, r.ID)
		if err != nil {
			t.Fatal(err)
		}

		switch r.Status {
		case RunCostEstimated, RunPolicyChecked, RunPolicyOverride:
			return r, rCleanup
		}

		if i > 45 {
			rCleanup()
			t.Fatal("Timeout waiting for run to be cost estimated")
		}

		time.Sleep(2 * time.Second)
	}
}

func createAppliedRun(t *testing.T, client *Client, w *Workspace) (*Run, func()) {
	r, rCleanup := createPlannedRun(t, client, w)
	ctx := context.Background()

	err := client.Runs.Apply(ctx, r.ID, RunApplyOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; ; i++ {
		r, err = client.Runs.Read(ctx, r.ID)
		if err != nil {
			t.Fatal(err)
		}

		if r.Status == RunApplied {
			return r, rCleanup
		}

		if i > 45 {
			rCleanup()
			t.Fatal("Timeout waiting for run to be applied")
		}

		time.Sleep(1 * time.Second)
	}
}

func createVariable(t *testing.T, client *Client, w *Workspace) (*Variable, func()) {
	var wCleanup func()

	if w == nil {
		w, wCleanup = createWorkspace(t, client, nil)
	}

	ctx := context.Background()
	v, err := client.Variables.Create(ctx, VariableCreateOptions{
		Key:       String(randomString(t)),
		Value:     String(randomString(t)),
		Category:  Category(CategoryTerraform),
		Workspace: w,
	})
	if err != nil {
		t.Fatal(err)
	}

	return v, func() {
		if err := client.Variables.Delete(ctx, v.ID); err != nil {
			t.Errorf("Error destroying variable! WARNING: Dangling resources\n"+
				"may exist! The full error is shown below.\n\n"+
				"Variable: %s\nError: %s", v.Key, err)
		}

		if wCleanup != nil {
			wCleanup()
		}
	}
}

func createWorkspace(t *testing.T, client *Client, org *Environment) (*Workspace, func()) {
	var orgCleanup func()

	if org == nil {
		org, orgCleanup = createEnvironment(t, client)
	}

	ctx := context.Background()
	w, err := client.Workspaces.Create(ctx, org.Name, WorkspaceCreateOptions{
		Name: String(randomString(t)),
	})
	if err != nil {
		t.Fatal(err)
	}

	return w, func() {
		if err := client.Workspaces.Delete(ctx, org.Name, w.Name); err != nil {
			t.Errorf("Error destroying workspace! WARNING: Dangling resources\n"+
				"may exist! The full error is shown below.\n\n"+
				"Workspace: %s\nError: %s", w.Name, err)
		}

		if orgCleanup != nil {
			orgCleanup()
		}
	}
}

func createWorkspaceWithVCS(t *testing.T, client *Client, org *Environment) (*Workspace, func()) {
	var orgCleanup func()

	if org == nil {
		org, orgCleanup = createEnvironment(t, client)
	}

	oc, ocCleanup := createOAuthToken(t, client, org)

	githubIdentifier := os.Getenv("GITHUB_POLICY_SET_IDENTIFIER")
	if githubIdentifier == "" {
		t.Fatal("Export a valid GITHUB_POLICY_SET_IDENTIFIER before running this test!")
	}

	options := WorkspaceCreateOptions{
		Name: String(randomString(t)),
		VCSRepo: &VCSRepoOptions{
			Identifier:   String(githubIdentifier),
			OAuthTokenID: String(oc.ID),
		},
	}

	ctx := context.Background()
	w, err := client.Workspaces.Create(ctx, org.Name, options)
	if err != nil {
		t.Fatal(err)
	}

	return w, func() {
		if err := client.Workspaces.Delete(ctx, org.Name, w.Name); err != nil {
			t.Errorf("Error destroying workspace! WARNING: Dangling resources\n"+
				"may exist! The full error is shown below.\n\n"+
				"Workspace: %s\nError: %s", w.Name, err)
		}

		if ocCleanup != nil {
			ocCleanup()
		}

		if orgCleanup != nil {
			orgCleanup()
		}
	}
}

func randomString(t *testing.T) string {
	v, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	return v
}
