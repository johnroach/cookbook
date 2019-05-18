//+build mage

package main

import (
	"os"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

// Build is to group build related tasks
type Build mg.Namespace

// Cloud is to group cloud related tasks
type Cloud mg.Namespace

// Run is to group tasks related to running the application
type Run mg.Namespace

// Test is to group tasks related to testing
type Test mg.Namespace

// Deploy is to group tasks related to deployments
type Deploy mg.Namespace

// Clean is to group cleanup tasks
type Clean mg.Namespace

// K8sLocal is the data needed for a local/dev k8s deployment
type K8sLocal struct {
	ImageName string
	Version   string
}

const (
	cloudRepo = "us.gcr.io/shared-svcs-489885bd"
	kubeCmd   = "kubectl"
)

// Bin builds the binary file. Usually used for development.
func (Build) Bin() error {
	mg.Deps(Test.Unit)

	logrus.Info("Building binary...")
	if err := sh.RunV(mg.GoCmd(), "build", "-o", "build/cookbook"); err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	logrus.Infof("Binary build complete. You can find the binary at %v/build/cookbook", currentDir)
	return nil
}

// Dev does a quick run for a dev environment. No docker involved. Server will run on localhost.
func (Run) Dev() error {
	mg.Deps(Test.Unit)

	logrus.Info("Running service via binary...")
	if err := sh.RunV(mg.GoCmd(), "run", "main.go", "-e", "dev"); err != nil {
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

	logrus.Info("Running docker container...")
	if err := sh.RunV("docker", "run", "-p", "8080:8080", "--mount", "type=bind,source="+currentDir+"/config,target=/config", "cookbook:"+version); err != nil {
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

	logrus.Info("Building docker images...")
	err = sh.RunV("docker", "build", "-t", "cookbook:"+version, "--build-arg", "VERSION="+version, ".")
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

	logrus.Info("Renaming docker image...")
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

	logrus.Info("Pushing image to cloud repo...")
	err = sh.Run("docker", "push", cloudRepo+"/cookbook:"+version)
	if err != nil {
		return nil
	}
	err = sh.RunV("docker", "push", cloudRepo+"/cookbook:latest")
	if err != nil {
		return nil
	}
	return nil
}

// Dep downloads dependencies
func Dep() error {
	logrus.Info("Downloading dependencies...")
	if err := sh.RunV(mg.GoCmd(), "mod", "download"); err != nil {
		return err
	}
	return nil
}

// All runs all tests for this project
func (Test) All() error {
	logrus.Info("Running all tests...")
	mg.Deps(Test.Unit)

	return nil
}

// Unit runs all unit tests for this project
func (Test) Unit() error {
	mg.Deps(Dep)

	logrus.Info("Running unit tests...")
	if err := sh.RunV(mg.GoCmd(), "test", "-v", "./..."); err != nil {
		return err
	}
	return nil
}

// LocalK8s deploys to local k8s setup. NOT DONE!!
func (Deploy) LocalK8s() error {
	//mg.Deps(Build.Docker)

	version, err := getVersion()
	if err != nil {
		return err
	}

	k8sData := K8sLocal{"cookbook", version}

	// Creating deployment template
	err = deploymentTemplate("deployment/local_deployment.yml", k8sData)
	if err != nil {
		return err
	}

	//make sure we are using local kubernetes (i.e. docker-desktop)
	if err := sh.RunV(kubeCmd, "config", "use-context", "docker-desktop"); err != nil {
		return err
	}

	sh.RunV(kubeCmd, "create", "configmap", "cookbook-config", "--from-file=config/dev.yaml")
	// kubectl create configmap cookbook-config --from-file=config/dev.yaml || kubectl create configmap cookbook-config --from-file config/dev.yaml -o yaml --dry-run | kubectl replace -f -
	// kubectl apply -f deployment/deployment.yml.tmp1

	logrus.Info("Deploying to local kubernetes...")

	return nil
}

// deploymentTemplate uses given data to generate a k8s deployment tmp file
func deploymentTemplate(templatePath string, k8sData K8sLocal) error {

	k8sTemp := template.Must(template.ParseFiles(templatePath))
	deployment, err := os.Create(templatePath + ".tmp")

	if err != nil {
		return err
	}
	err = k8sTemp.Execute(deployment, k8sData)
	if err != nil {
		return err
	}

	return nil
}

// getVersion gets version of current commit.
// This is used when coming up with a docker container tag or used for deployments
func getVersion() (string, error) {
	logrus.Info("Getting version...")
	version, err := sh.Output("git", "rev-parse", "--short=7", "HEAD")
	if err != nil {
		return "", err
	}
	return version, nil
}
