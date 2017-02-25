Golang client for the [Yandex.Disk API](https://tech.yandex.com/disk/api/concepts/about-docpage/).

*Work in progress*

### Usage

```go
import "github.com/chibisov/go-yadisk/yadisk"
```

Construct a new Yandex.Disk client, then use the various services 
on the client to access different parts of the Yandex.Disk API. For example:

```go
client := yadisk.NewClient("ACCESS_TOKEN")
ctx := context.Background()

// get a general information about a user's Disk
disk, response, err = client.Disk.Get(ctx)

// get meta information about file or directory
resources, response, err = client.Resources.Get(ctx, "/")

// get meta information about resources in the trash
resources, response, err = client.Trash.Resources.Get(ctx, "/")

// restore file from the trash
resource, response, err = client.Trash.Resources.Restore(ctx, "/file.jpg")

// get data for download a file
downloadInfo, response, err = client.Resources.Download(ctx, "/file.jpg")

// get data for uploading a file
downloadInfo, response, err = client.Resources.Upload(ctx, "/file.jpg")
```

### Tests

Running only unit tests:

```sh
$ go test ./tests/unit/
```

Running only integration tests:

```sh
$ go test ./tests/integration/
```

For integration tests you need to provide some environment variables for 
credentials:

```sh
$ ACCESS_TOKEN=123 go test ./tests/integration/
```

You can save it to `.test.config` (which is ignored by git) and run like this:

```sh
$ eval $(cat .test.config) go test ./tests/integration/
```

Running all tests:

```sh
$ go test ./tests/*
```