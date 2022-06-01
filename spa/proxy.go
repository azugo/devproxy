package spa

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
)

// SpaDevProxy instance.
type SpaDevProxy interface {
	// Start backgroud process and wait when it is ready to accept connections.
	Start(ctx context.Context) error

	// Stop background process by killing it.
	Stop() error

	// DevServerURL returns development server URL.
	DevServerURL() *url.URL
}

// SpaDevProxyOptions options to use for SpaDevProxy.
type SpaDevProxyOptions struct {
	// RunnerType to use. Defaults to NPM.
	RunnerType RunnerType

	// ScriptName is the name of the script to run.
	ScriptName string

	// Args holds command line arguments.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Start runs the command in the
	// calling process's current directory.
	Dir string

	// StartRegexp specifies regular expression used to detect
	// when background process is ready to accept incomming requests.
	StartRegexp *regexp.Regexp

	// Port specifies localhost port that background service is
	// using to accept requests.
	Port int

	// ShowBuildInfo specifies option either to show webpack building
	// progress or not.
	ShowBuildInfo bool
}

type spaDevProxy struct {
	options *SpaDevProxyOptions
	cmd     *exec.Cmd
	url     *url.URL
}

// NewSpaDevProxy creates new SpaDevProxy instance.
func NewSpaDevProxy(options *SpaDevProxyOptions) (SpaDevProxy, error) {
	remote, err := url.Parse(fmt.Sprintf("http://localhost:%d", options.Port))
	if err != nil {
		return nil, err
	}

	return &spaDevProxy{
		options: options,
		// proxy:   httputil.NewSingleHostReverseProxy(remote),
		url: remote,
	}, nil
}

// Start backgroud process and wait when it is ready to accept connections.
func (p *spaDevProxy) Start(ctx context.Context) error {
	if _, err := os.Stat(p.options.Dir); err != nil {
		return err
	}
	path, args := prepareRunner(p.options.RunnerType, p.options.ScriptName, p.options.Args...)
	p.cmd = newCommand(ctx, path, args...)
	p.cmd.Env = append(p.options.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	p.cmd.Dir = p.options.Dir

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	done := forwardOutput(stdout, stderr, p.options.StartRegexp, p.options.ShowBuildInfo)

	err = p.cmd.Start()

	<-done

	return err
}

// Stop background process by killing it.
func (p *spaDevProxy) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	// Using syscall as Process.Kill does not kill child processes
	if err := killCommand(p.cmd); err != nil {
		return err
	}
	return p.cmd.Wait()
}

// DevServerURL returns development server URL.
func (p *spaDevProxy) DevServerURL() *url.URL {
	return p.url
}
