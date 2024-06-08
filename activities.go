package yeet

import (
	"context"
	"fmt"
	"io"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.temporal.io/sdk/activity"
)

type Git struct {
	path string
}

type GitParam struct {
	Host       string
	Owner      string
	Repository string
	Revision   string
}

type GitResult struct {
	Path string
}

func (a *Git) Clone(ctx context.Context, param GitParam) (*GitResult, error) {
	tempDir, err := os.MkdirTemp("", "yeet-")
	if err != nil {
		return nil, err
	}

	a.path = tempDir
	_, err = git.PlainClone(a.path, false, &git.CloneOptions{
		URL:           fmt.Sprintf("https://%s/%s/%s", param.Host, param.Owner, param.Repository),
		ReferenceName: plumbing.ReferenceName(param.Revision),
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}

	result := &GitResult{
		Path: a.path,
	}
	return result, nil
}

func (a *Git) Clean(ctx context.Context, param GitParam) error {
	os.RemoveAll(a.path)

	return nil
}

type BuildParam struct {
	Path  string
	Image string
	Tag   string
}

type BuildResult struct {
	Image string
	Tag   string
}

type Build struct {
}

func (a *Build) Buildpacks(ctx context.Context, param BuildParam) (*BuildResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Path:", param.Path)
	logger.Info("Image:", param.Image)
	logger.Info("Tag:", param.Tag)

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer docker.Close()

	reader, err := docker.ImagePull(ctx, "docker.io/buildpacksio/pack", image.PullOptions{})
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	// cli.ImagePull is asynchronous.
	// The reader needs to be read completely for the pull operation to complete.
	// If stdout is not required, consider using io.Discard instead of os.Stdout.
	io.Copy(os.Stdout, reader)

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image:      "docker.io/buildpacksio/pack",
		Cmd:        []string{"build", "--builder=heroku/builder:22", fmt.Sprintf("%s:%s", param.Image, param.Tag)},
		Tty:        false,
		WorkingDir: "/workspace",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
			{
				Type:   mount.TypeBind,
				Source: param.Path,
				Target: "/workspace",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := docker.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	_, err = docker.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}

	result := &BuildResult{
		Image: param.Image,
		Tag:   param.Tag,
	}
	return result, nil
}

type Deploy struct {
	Config DeployConfig
}

type DeployConfig map[string]EventConfig

type EventConfig struct {
	Stages []DeployStage
}

type DeployStage struct {
	Name   string
	Groups []string
	Wait   string
}

type DeployParam struct {
	RepoPath string
	SubPath  string
}

type DeployResult struct {
}

func (a *Deploy) GetConfig(ctx context.Context, param DeployParam) (DeployConfig, error) {
	cuectx := cuecontext.New()
	instances := load.Instances([]string{
		fmt.Sprintf("%s/apps/%s/yeet.cue", param.RepoPath, param.SubPath),
	}, nil)
	if len(instances) == 0 {
		fmt.Println("No instances loaded")
		return nil, nil
	}
	instance := instances[0]

	if err := instance.Err; err != nil {
		fmt.Println("Failed to load instance:", err)
		return nil, err
	}

	// Build a value from the Cue instance
	value := cuectx.BuildInstance(instance)

	if err := value.Err(); err != nil {
		fmt.Println("Failed to build value from instance:", err)
		return nil, err
	}

	// Fill the Go struct with the Cue value
	if err := value.Decode(&a.Config); err != nil {
		fmt.Println("Failed to decode Cue value into Go struct:", err)
		return nil, err
	}

	// Print the parsed configuration
	fmt.Printf("Parsed ConfigMap: %+v\n", a.Config)

	return a.Config, nil
}

func (a *Deploy) ProcessStage(ctx context.Context, param DeployParam, stage DeployStage) error {
	// TODO place holder
	// update based on groups
	// commit
	// push
	filePath := fmt.Sprintf("%s/apps/%s/dev/main/bundle.cue", param.RepoPath, param.SubPath)
	cuectx := cuecontext.New()
	instances := load.Instances([]string{
		filePath,
	}, nil)
	instance := instances[0]

	tagPath := cue.MakePath(
		cue.Str("bundle"),
		cue.Str("instances"),
		cue.Str("podinfo"),
		cue.Str("values"),
		cue.Str("controllers"),
		cue.Str("main"),
		cue.Str("containers"),
		cue.Str("app"),
		cue.Str("image"),
		cue.Str("tag"),
	)
	newTag := "new-tag@sha256:newshahash"
	value := cuectx.BuildInstance(instance)
	// TODO https://github.com/cue-lang/cue/issues/2170
	value = value.FillPath(tagPath, newTag)

	updatedCueBytes, _ := format.Node(value.Syntax())
	err := os.WriteFile(filePath, updatedCueBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
