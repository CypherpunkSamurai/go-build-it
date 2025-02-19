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