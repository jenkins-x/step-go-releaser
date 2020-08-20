package pkg

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-x/jx-api/pkg/util"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/kube"
	opts "github.com/jenkins-x/jx-helpers/pkg/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	githubappLabel             = "jenkins.io/githubapp-owner="
	ownerAnnotation            = "jenkins.io/githubapp-owner"
	nonGithubAppSecretSelector = "jenkins.io/kind=git,jenkins.io/service-kind=github"
	githubPassword             = "password"
)

// Run implements the command
func (o *options) Run() error {
	var err error
	o.KubeClient, o.Namespace, err = kube.LazyCreateKubeClientAndNamespace(o.KubeClient, o.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to create kube Client")
	}

	if o.organisation == "" {
		return opts.MissingOption(organisation)
	}
	if o.revision == "" {
		return opts.MissingOption(revision)
	}
	if o.branch == "" {
		return opts.MissingOption(branch)
	}
	if o.version == "" {
		return opts.MissingOption(version)
	}
	if o.buildDate == "" {
		return opts.MissingOption(buildDate)
	}
	if o.goVersion == "" {
		return opts.MissingOption(goVersion)
	}
	if o.rootPackage == "" {
		return opts.MissingOption(rootPackage)
	}

	return o.goReleaser()

}

func (o *options) goReleaser() error {

	token, err := o.getToken()
	if err != nil {
		return errors.Wrapf(err, "failed to get github token for organisation %s", o.organisation)
	}
	o.Runner.SetEnvVariable("GITHUB_TOKEN", token)
	o.Runner.SetEnvVariable("REV", o.revision)
	o.Runner.SetEnvVariable("BRANCH", o.branch)
	o.Runner.SetEnvVariable("VERSION", o.version)
	o.Runner.SetEnvVariable("BUILDDATE", o.buildDate)
	o.Runner.SetEnvVariable("GOVERSION", o.goVersion)
	o.Runner.SetEnvVariable("ROOTPACKAGE", o.rootPackage)

	o.Runner.SetName("goreleaser")

	args := []string{"release", "--config=.goreleaser.yml", "--rm-dist", "--release-notes=./changelog.md", "--skip-validate", "--timeout", o.timeout}
	o.Runner.SetArgs(args)

	if o.CommandRunner == nil {
		o.CommandRunner = cmdrunner.DefaultCommandRunner
	}
	o.Runner.Out = os.Stdout
	o.Runner.Err = os.Stderr

	_, err = o.CommandRunner(&o.Runner)
	return err
}

func (o *options) getToken() (string, error) {
	client := o.KubeClient
	ns := o.Namespace

	selector := githubappLabel + o.organisation
	listOpts := metav1.ListOptions{
		LabelSelector: selector,
	}
	secretInterface := client.CoreV1().Secrets(ns)
	secrets, err := secretInterface.List(listOpts)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", errors.Wrapf(err, "failed to get secrets for selector: %s", selector)
	}

	for _, s := range secrets.Items {
		token := s.Data[githubPassword]
		if len(token) > 0 {
			return string(token), nil
		}
	}

	// lets try find a non-github app secret
	listOpts = metav1.ListOptions{
		LabelSelector: nonGithubAppSecretSelector,
	}

	secrets, err = secretInterface.List(listOpts)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", errors.Wrapf(err, "failed to get secrets for selector: %s", nonGithubAppSecretSelector)
	}
	for _, s := range secrets.Items {
		token := s.Data[githubPassword]
		if len(token) > 0 {
			return string(token), nil
		}
	}
	return "", errors.Errorf("could not find a secret for selector %s or %s", selector, nonGithubAppSecretSelector)
}

func (o *options) run() (string, error) {
	e := exec.Command(o.Runner.CurrentName(), o.Runner.CurrentArgs()...)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	os.Setenv("PATH", util.PathWithBinary())

	if len(o.Runner.CurrentEnv()) > 0 {
		m := map[string]string{}
		environ := os.Environ()
		for _, kv := range environ {
			paths := strings.SplitN(kv, "=", 2)
			if len(paths) == 2 {
				m[paths[0]] = paths[1]
			}
		}
		for k, v := range o.Runner.CurrentEnv() {
			m[k] = v
		}
		envVars := []string{}
		for k, v := range m {
			envVars = append(envVars, k+"="+v)
		}
		e.Env = envVars
	}

	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	var text string

	err := e.Run()
	if err != nil {
		if err != nil {
			errors.Wrapf(err, "failed to run command")
		}
	}
	return text, err
}
