Scalr Go Client
==============================
This is a fork of the Hashicorp Terraform Enterprise Go client.

## Installation

Installation can be done with a normal `go get`:

```
go get -u github.com/scalr/go-scalr
```

## Documentation

For complete usage of the API client, see the full [package docs](https://godoc.org/github.com/scalr/go-scalr).

## Usage

```go
import scalr "github.com/scalr/go-scalr"
```

Construct a new Scalr client, then use the various endpoints on the client to
access different parts of the Scalr API. For example, to list
all environments:

```go
config := &scalr.Config{
	Token: "insert-your-token-here",
}

client, err := scalr.NewClient(config)
if err != nil {
	log.Fatal(err)
}

orgs, err := client.Environments.List(context.Background(), EnvironmentListOptions{})
if err != nil {
	log.Fatal(err)
}
```

## Examples

The [examples](https://github.com/scalr/go-scalr/tree/master/examples) directory
contains a couple of examples. One of which is listed here as well:

```go
package main

import (
	"log"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	config := &scalr.Config{
		Token: "insert-your-token-here",
	}

	client, err := scalr.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// Create a new environment
	options := scalr.EnvironmentCreateOptions{
		Name:  scalr.String("example"),
		Email: scalr.String("info@example.com"),
	}

	org, err := client.Environments.Create(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	// Delete an environment
	err = client.Environments.Delete(ctx, org.Name)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Running tests

### 1. (Optional) Create a policy sets repo

If you are planning to run the full suite of tests or work on policy sets, you'll need to set up a policy set repository in GitHub.

Your policy set repository will need the following: 
1. A policy set stored in a subdirectory `policy-sets/foo`
1. A branch other than master named `policies`
   
### 2. Set up environment variables

##### Required:
Tests are run against an actual backend so they require a valid backend address
and token.
1. `TFE_ADDRESS` - URL of a Scalr instance to be used for testing, including scheme. Example: `https://scalr.local`
1. `TFE_TOKEN` - A user API token for the Scalr instance being used for testing.

##### Optional:
1. `GITHUB_TOKEN` - [GitHub personal access token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line). Required for running OAuth client tests.
1. `GITHUB_POLICY_SET_IDENTIFIER` - GitHub policy set repository identifier in the format `username/repository`. Required for running policy set tests.

You can set your environment variables up however you prefer. The following are instructions for setting up environment variables using [envchain](https://github.com/sorah/envchain).
   1. Make sure you have envchain installed. [Instructions for this can be found in the envchain README](https://github.com/sorah/envchain#installation).
   1. Pick a namespace for storing your environment variables. I suggest `go-scalr` or something similar.
   1. For each environment variable you need to set, run the following command:
      ```sh
      envchain --set YOUR_NAMESPACE_HERE ENVIRONMENT_VARIABLE_HERE
      ```
      **OR**
    
      Set all of the environment variables at once with the following command:
      ```sh
      envchain --set YOUR_NAMESPACE_HERE TFE_ADDRESS TFE_TOKEN GITHUB_TOKEN GITHUB_POLICY_SET_IDENTIFIER
      ```

### 3. Make sure run queue settings are correct

In order for the tests relating to queuing and capacity to pass, FRQ (fair run queuing) should be
enabled with a limit of 2 concurrent runs per environment on the Terraform Cloud or Terraform Enterprise instance you are using for testing.

### 4. Run the tests

#### Running all the tests
As running the all of the tests takes about ~20 minutes, make sure to add a timeout to your
command (as the default timeout is 10m).

##### With envchain:
```sh
$ envchain YOUR_NAMESPACE_HERE go test ./... -timeout=30m
```

##### Without envchain:
```sh
$ go test ./... -timeout=30m
```
#### Running specific tests

The commands below use notification configurations as an example.

##### With envchain:
```sh
$ envchain YOUR_NAMESPACE_HERE go test -run TestNotificationConfiguration -v ./...
```

##### Without envchain:
```sh
$ go test -run TestNotificationConfiguration -v ./...
```   
