# tilde server

```env SOCKET=/var/run/user/$(id -u)/web.sock ADMIN=$HOME/tlde.socket PORT=$(id -u) tlde```

SOCKET variable defines location of web socket
ADMIN variable defines location of [diamond](https://github.com/aerth/diamond) admin socket
PORT variable defines what port to listen on
