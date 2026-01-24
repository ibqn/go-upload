### Getting started

project was created with

```sh
go mod init go-upload
```

```sh
go mod tidy
```

### docker

- Start database services first (this creates the network)

```sh
docker compose -f compose.yaml up -d
```

-  Start the Go app (joins the existing network)

```sh
docker compose -f compose.go.yaml up -d
```