package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariablesCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	wsTest, wsTestCleanup := createWorkspace(t, client, nil)
	account := &Account{ID: defaultAccountID}

	defer wsTestCleanup()

	t.Run("when options has an empty string value", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(""),
			Category:  Category(CategoryEnv),
			Workspace: wsTest,
		}

		v, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)

		assert.NotEmpty(t, v.ID)
		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.Value, v.Value)
		assert.Equal(t, *options.Category, v.Category)
	})

	t.Run("when options is missing value", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Category:  Category(CategoryEnv),
			Workspace: wsTest,
		}

		v, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)

		assert.NotEmpty(t, v.ID)
		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, "", v.Value)
		assert.Equal(t, *options.Category, v.Category)
	})

	t.Run("when options is missing key", func(t *testing.T) {
		options := VariableCreateOptions{
			Value:     String(randomString(t)),
			Category:  Category(CategoryEnv),
			Workspace: wsTest,
		}

		_, err := client.Variables.Create(ctx, options)
		assert.EqualError(t, err, "key is required")
	})

	t.Run("when options has an empty key", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(""),
			Value:     String(randomString(t)),
			Category:  Category(CategoryEnv),
			Workspace: wsTest,
		}

		_, err := client.Variables.Create(ctx, options)
		assert.EqualError(t, err, "key is required")
	})

	t.Run("when options is missing category", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(randomString(t)),
			Workspace: wsTest,
		}

		_, err := client.Variables.Create(ctx, options)
		assert.EqualError(t, err, "category is required")
	})

	t.Run("when options is missing account", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:         String(randomString(t)),
			Value:       String(randomString(t)),
			Category:    Category(CategoryEnv),
			Environment: wsTest.Environment,
			Workspace:   wsTest,
		}

		_, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)
	})

	t.Run("when options is missing environment", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(randomString(t)),
			Category:  Category(CategoryEnv),
			Account:   account,
			Workspace: wsTest,
		}

		_, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)
	})

	t.Run("when options is missing workspace", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:         String(randomString(t)),
			Value:       String(randomString(t)),
			Category:    Category(CategoryEnv),
			Account:     account,
			Environment: wsTest.Environment,
		}

		_, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)
	})

	t.Run("when options is missing account, environment, workspace", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(randomString(t)),
			Category:  Category(CategoryEnv),
			Workspace: wsTest,
		}

		v, err := client.Variables.Create(ctx, options)
		require.NoError(t, err)

		assert.NotEmpty(t, v.ID)
		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.Value, v.Value)
		assert.Equal(t, *options.Category, v.Category)
		assert.Equal(t, options.Workspace.ID, v.Workspace.ID)
	})

}

func TestVariablesRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	vTest, vTestCleanup := createVariable(t, client, nil)
	defer vTestCleanup()

	t.Run("when the variable exists", func(t *testing.T) {
		v, err := client.Variables.Read(ctx, vTest.ID)
		require.NoError(t, err)
		assert.Equal(t, vTest.ID, v.ID)
		assert.Equal(t, vTest.Category, v.Category)
		assert.Equal(t, vTest.HCL, v.HCL)
		assert.Equal(t, vTest.Key, v.Key)
		assert.Equal(t, vTest.Sensitive, v.Sensitive)
		assert.Equal(t, vTest.Value, v.Value)
	})

	t.Run("when the variable does not exist", func(t *testing.T) {
		v, err := client.Variables.Read(ctx, "nonexisting")
		assert.Nil(t, v)
		assert.Equal(t, ErrResourceNotFound, err)
	})

	t.Run("without a valid variable ID", func(t *testing.T) {
		v, err := client.Variables.Read(ctx, badIdentifier)
		assert.Nil(t, v)
		assert.EqualError(t, err, "invalid value for variable ID")
	})
}

func TestVariablesUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	vTest, vTestCleanup := createVariable(t, client, nil)
	defer vTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := VariableUpdateOptions{
			Key:   String("newname"),
			Value: String("newvalue"),
			HCL:   Bool(true),
		}

		v, err := client.Variables.Update(ctx, vTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.HCL, v.HCL)
		assert.Equal(t, *options.Value, v.Value)
	})

	t.Run("when updating a subset of values", func(t *testing.T) {
		options := VariableUpdateOptions{
			Key: String("someothername"),
			HCL: Bool(false),
		}

		v, err := client.Variables.Update(ctx, vTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.HCL, v.HCL)
	})

	t.Run("with sensitive set", func(t *testing.T) {
		options := VariableUpdateOptions{
			Sensitive: Bool(true),
		}

		v, err := client.Variables.Update(ctx, vTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Sensitive, v.Sensitive)
		assert.Empty(t, v.Value) // Because its now sensitive
	})

	t.Run("without any changes", func(t *testing.T) {
		created, vTestCleanup := createVariable(t, client, vTest.Workspace)
		defer vTestCleanup()

		updated, err := client.Variables.Update(ctx, created.ID, VariableUpdateOptions{})
		require.NoError(t, err)

		assert.Equal(t, created, updated)
	})

	t.Run("with invalid variable ID", func(t *testing.T) {
		_, err := client.Variables.Update(ctx, badIdentifier, VariableUpdateOptions{})
		assert.EqualError(t, err, "invalid value for variable ID")
	})
}

func TestVariablesDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	wTest, wTestCleanup := createWorkspace(t, client, nil)
	defer wTestCleanup()

	vTest, _ := createVariable(t, client, wTest)

	t.Run("with valid options", func(t *testing.T) {
		err := client.Variables.Delete(ctx, vTest.ID)
		assert.NoError(t, err)
	})

	t.Run("with non existing variable ID", func(t *testing.T) {
		err := client.Variables.Delete(ctx, "nonexisting")
		assert.Equal(t, err, ErrResourceNotFound)
	})

	t.Run("with invalid variable ID", func(t *testing.T) {
		err := client.Variables.Delete(ctx, badIdentifier)
		assert.EqualError(t, err, "invalid value for variable ID")
	})
}
