# tilde server

### usage

```env SOCKET=/var/run/user/$(id -u)/web.sock ADMIN=$HOME/tlde.socket PORT=$(id -u) tlde```

### variables

```SOCKET``` variable defines location of web socket, can be empty if PORT is not

```ADMIN``` variable defines location of [diamond](https://github.com/aerth/diamond) admin socket

```PORT``` variable defines what port to listen on, can be empty if SOCKET is not

### tlde is a diamond

Use ```diamond admin``` to connect to ADMIN port

```go get -v -u github.com/aerth/diamond/cmd/diamond-admin```

