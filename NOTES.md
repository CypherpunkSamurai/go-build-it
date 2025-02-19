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
