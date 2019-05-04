//+build mage

package main

import (
	"log"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build is to group build related tasks
type Build mg.Namespace

// Cloud is to group cloud related tasks
type Cloud mg.Namespace

// Run is to group tasks related to running the application
type Run mg.Namespace

// Test is to group tasks related to testing
type Test mg.Namespace

const (
	cloudRepo = "us.gcr.io/shared-svcs-489885bd"
)

// Bin builds the binary file. Usually used for development.
func (Build) Bin() error {
	mg.Deps(Test.Unit)

	log.Println("Building binary...")
	if err := sh.Run(mg.GoCmd(), "build"); err != nil {
		return err
	}
	return nil
}

// Dev does a quick run for a dev environment. No docker involved. Server will run on localhost.
func (Run) Dev() error {
	mg.Deps(Test.Unit)

	log.Println("Running service via binary...")
	if err := sh.Run(mg.GoCmd(), "run", "main.go", "-e", "dev"); err != nil {
		return err
	}
	return nil
}

// DevDocker runs the locally built docker tagged with current version where version is pulled from git log
func (Run) DevDocker() error {
	mg.Deps(Build.Docker)

	version, err := getVersion()
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Println("Running docker container...")
	if err := sh.Run("docker", "run", "-p", "8080:8080", "--mount", "type=bind,source="+currentDir+"/config,target=/config", "cookbook:"+version); err != nil {
		return err
	}
	return nil
}

// Docker builds the docker container
func (Build) Docker() error {
	mg.Deps(Test.Unit)

	version, err := getVersion()
	if err != nil {
		return err
	}

	log.Println("Bulding docker image...")
	err = sh.Run("docker", "build", "-t", "cookbook:"+version, "--build-arg", "VERSION="+version, ".")
	if err != nil {
		return nil
	}
	return nil
}

// DockerTag builds, names and tags per cloud repo requirements
func (Cloud) DockerTag() error {
	mg.Deps(Build.Docker)

	version, err := getVersion()
	if err != nil {
		return err
	}

	log.Println("Renaming docker image...")
	err = sh.Run("docker", "tag", "cookbook:"+version, cloudRepo+"/cookbook:"+version)
	if err != nil {
		return err
	}
	err = sh.Run("docker", "tag", "cookbook:"+version, cloudRepo+"/cookbook:latest")
	if err != nil {
		return err
	}
	return nil
}

// DockerPush builds, tags and pushes per cloud repo requirements
func (Cloud) DockerPush() error {
	mg.Deps(Cloud.DockerTag)

	version, err := getVersion()
	if err != nil {
		return err
	}

	log.Println("Pushing image to cloud repo...")
	err = sh.Run("docker", "push", cloudRepo+"/cookbook:"+version)
	if err != nil {
		return nil
	}
	err = sh.Run("docker", "push", cloudRepo+"/cookbook:latest")
	if err != nil {
		return nil
	}
	return nil
}

// Dep downloads dependencies
func Dep() error {
	log.Println("Downloading dependencies...")
	if err := sh.Run(mg.GoCmd(), "mod", "download"); err != nil {
		return err
	}
	return nil
}

// All runs all tests for this project
func (Test) All() error {
	mg.Deps(Test.Unit)

	return nil
}

// Unit runs all unit tests for this project
func (Test) Unit() error {
	mg.Deps(Dep)

	log.Println("Running unit tests...")
	if err := sh.Run(mg.GoCmd(), "test", "-v", "./..."); err != nil {
		return err
	}
	return nil
}

// getVersion gets version of current commit.
// This is used when coming up with a docker container tag or used for deployments
func getVersion() (string, error) {
	log.Println("Getting version...")
	version, err := sh.Output("git", "rev-parse", "--short=7", "HEAD")
	if err != nil {
		return "", err
	}
	return version, nil
}
