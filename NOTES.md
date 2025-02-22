# Notes

This file contains snippets and notes for the project.

## 2025-02-19 16:00:44 - Init Project

This is a project to build a simple ci server using golang and docker api.

## 2025-02-19 16:17:18 - Modules

I've decided to use the following modules cause they look best for use case.

- [github.com/goccy/go-yaml](https://github.com/goccy/go-yaml)
- [github.com/docker/docker/client](https://github.com/docker/docker/client)

## 2025-02-19 16:35:45 - Lets Try the Docker SDK Code Snippet from the Website

Creating a Docker client is pretty easy. I will be using the Docker SDK from [here](https://docs.docker.com/reference/api/engine/sdk/examples/). The SDK is available in Python, Go, and JavaScript.

The golang version requires that you manually fetch all the dependencies as `go mod tidy` results in lot of errors. [googleapis/go-genproto#1015](https://github.com/googleapis/go-genproto/issues/1015#issuecomment-2380323686).

```shell
go get "github.com/docker/docker/api/types/container"
go get "github.com/docker/docker/api/types/image"
go get "github.com/docker/docker/client"
go get "github.com/docker/docker/pkg/stdcopy"
```

We need a working docker daemon and need to run this to pull an `alpine:latest` container.

## 2025-02-19 16:42:30 - For Logging we will use my trusty old favorite Uber Zap

I copied over my `logger.go` from other project along with `env.go`

## 2025-02-19 17:25:17 - Creating a Singleton Client

Creating a Singleton client is not that hard. I will be using the following code snippet.

```go
// Docker Client Instance
var (
	cli  *client.Client
	once sync.Once
)

// GetDockerClient returns the Docker client instance
func GetDockerClient() *client.Client {
	// thread-safe implementation (classic go :P)
	once.Do(func() {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			utils.Logger().Printf("Error creating docker client: %v", err)
			cli = nil
		}
	})
	return cli
}
```

Which is a [common singleton pattern](https://refactoring.guru/design-patterns/singleton/go/example) for go. If you're not from the Java world singletons are single instances of classes mostly used to manage resources like retrofit http clients etc, in our case it's a docker [*client.Client](https://pkg.go.dev/github.com/docker/docker/client#NewClientWithOpts).

More resources here:
- [sync](https://pkg.go.dev/sync)
- [Go sync.Once is Simple... Does It Really?](https://victoriametrics.com/blog/go-sync-once/)
- [Understanding Golang's sync.Once: Practical Examples in 2024](https://cristiancurteanu.com/understanding-go-sync-once/)
- [Just Call Your Code Only Once !!](https://medium.easyread.co/just-call-your-code-only-once-256f69ed39a8)
- [Implementing the Singleton Pattern in Go](https://www.codingexplorations.com/blog/implementing-the-singleton-pattern-in-go)
- ⭐ [Golang Patterns - Creational Singleton](https://github.com/tmrts/go-patterns/blob/master/creational/singleton.md)

## 2025-02-19 17:32:45 - Defining a Workflow Format

We need to define a workflow format to our CI Platform. I personally find the [gitlab ci](https://gitlab.com/gitlab-org/gitlab/-/blob/master/lib/gitlab/ci/templates/Go.gitlab-ci.yml), and [circleci samples](https://circleci.com/docs/sample-config/) and [github actions](https://docs.github.com/en/actions/sharing-automations/creating-workflow-templates-for-your-organization) `uses` to be quite useful.

So I created this abomination of a workflow format:
```yaml
# Example Worflow
version: 1.0

# name of the workflow
name: Build and Test Go Binary

# jobs under the workflow
jobs:
  test:
    image: golang:latest
    steps:
      - use: checkout
      - name: format code
        run: go fmt ./...
      - name: run tests
        run: go test ./...
  build:
    # optional name (else it will be the name of the job. like `test` has name `test`)
    name: Build Binary
    requires: [test]
    image: golang:latest
    steps:
      - use: checkout
      - name: build binary
        run: go build -o main .
```

This will work for now. We have a simple workflow defined with image, steps, and uses.

- `name` is used to give a name to the workflow
- `jobs` define the sequence of the workflow
  - `name` each job has a name and is optional. if no name is provided it will be the name of the job key.
  - `requires` each job has a list of jobs that it depends on. If a job fails the other dependent jobs are not executed.
  - `image` each job has a docker image. The image is used to run the job.
  - `steps` define the sequence of the steps inside a job
    - `uses` keyword is used to call actions or internal functions like `checkout` here is simply a quick way to say `git clone .`
    - `name` each step has a name and is optional. if no name is provided it will be the name of the step key or action.
    - `run` each step has a command to run.

This will serve as baseline for my CI Platform.

## 2025-02-19 18:28:06 - Adding Yaml Parsing

We will use [goccy/go-yaml](https://pkg.go.dev/github.com/goccy/go-yaml@v1.15.23) for parsing yaml as it looks best compatible with weird yaml formats.

## 2025-02-19 18:28:06 - Adding Test Cases to Verify Workflow Format

After reading a lot of Golang Unit Testing Posts, [Testify](https://github.com/stretchr/testify) seems like the best choice when it comes to writing testss (also I like the python3 like assertions). The other libraries look like a lot of work and im lazy.

- https://jdheyburn.co.uk/blog/assertions-in-gotests-test-generation/
- https://speedscale.com/blog/golang-testing-frameworks-for-every-type-of-test/

## 2025-02-20 12:21:20 - Renamed WorkflowType to workflow_type cause of go cache bug

Golang requires consistent naming. If im using camel case I need to define all structs in camel case, else i need to use snake case.

Workflow file also was updated. `use` was renamed to `uses` for consistency.

## 2025-02-20 12:21:20 - Added RunTask struct

Added RunTaskType struct for converting workflows to runtasks that docker api can understand. Updated main to include a simple test example. Added test cases for RunTaskType.


## 2025-02-20 12:27:39 - I Learnt How to Running Containers Interactively

I'm writing this next morning as I slept on completing the work last night and pushing changes to remote. So after a bit of research i found we can run `docker exec` from the docker api, and it also provides us with standard io.

On checking out the `Client.ContainerCreate` [godoc](https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerCreate) we had been using before i found other options that can be passed to ContainerCreate.

```go
func (cli *Client) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
```

With a little searching i found code examples for the interesting `Tty` option in `container.Config`.

query: `context:global lang:go .ContainerCreate( Tty` [\[1\]](https://sourcegraph.com/search?q=context:global+lang:go+.ContainerCreate%28+Tty&patternType=keyword&sm=0) [\[2\]](https://github.com/search?q=lang%3Ago+.ContainerCreate%28+Tty&type=code)

So I cooked up a solution:

```go
// create a container with a shell
resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image:        "alpine:latest",
    Cmd:          []string{"/bin/sh"},
    Tty:          true,
    OpenStdin:    true,
    StdinOnce:    true,
    AttachStdin:  true,
    AttachStdout: true,
    AttachStderr: true,
}, nil, nil, nil, "")
if err != nil {
    return err
}

// start the container
if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
    return err
}

// attach io
waiter, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
    Stdin:  true,
    Stdout: true,
    Stderr: true,
    Stream: true,
})
if err != nil {
    return err
}
defer waiter.Close()

go io.Copy(os.Stdout, waiter.Reader)
go io.Copy(os.Stderr, waiter.Reader)

// run shell commands
commands := []string{
    "apk add vim\n",
    "vim test.txt\n",
    "i",
    "hello",
    "\x1b",
    ":wq\n",
    "exit\n",
}

for _, cmd := range commands {
    _, err = waiter.Conn.Write([]byte(cmd))
    if err != nil {
        return err
    }
    time.Sleep(time.Second)
}

// wait for container to exit
statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
select {
case err := <-errCh:
    if err != nil {
        return err
    }
case <-statusCh:
}
```

## 2025-02-22 10:26:09 - Must and DotEnv Go

Must pattern is new to golang and evalutates conditions that are critical to the program. It panics on error.

I've copied it from this post :P
- https://dev.to/dubjay18/gos-must-pattern-streamline-your-error-handling-27ff

```go
func Must[T any](expr T, err error) T {
	if err != nil {
		panic(err)
	}
	return expr
}
```
