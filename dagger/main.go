// A generated module for JsonSchemaAsciidoc functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"

	"dagger/json-schema-asciidoc/internal/dagger"
)

type JsonSchemaAsciidoc struct{}

// Returns a container that echoes whatever string argument is provided
func (m *JsonSchemaAsciidoc) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

func (m *JsonSchemaAsciidoc) Publish(ctx context.Context, source *dagger.Directory, binaryName string) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.Build(source, binaryName).
		Publish(ctx, fmt.Sprintf("ttl.sh/json-schema-renderer-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
}

// Return the result of running unit tests
func (m *JsonSchemaAsciidoc) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *JsonSchemaAsciidoc) BuildEnv(source *dagger.Directory) *dagger.Container {

	return dag.Container().
		From("golang:1.23").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		WithExec([]string{"go", "mod", "download"})
}

func (m *JsonSchemaAsciidoc) Build(source *dagger.Directory, binaryName string) *dagger.Container {

	build := m.BuildEnv(source).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("gomod")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("CGO_ENABLED", "0").
		WithDirectory("/src", source).
		WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", binaryName, "."}).
		File("/src/" + binaryName)

	return dag.Container().
		From("alpine:latest").
		WithFile("/usr/bin/"+binaryName, build).
		WithEntrypoint([]string{binaryName})
}
