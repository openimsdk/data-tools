# Data conversion v3 version MYSQL -> MONGODB

## OPENIM-SERVER
* `https://github.com/openimsdk/open-im-server`
* v3 version, versions below v3.5 are upgraded to v3.5 and above.

```shell
cd openim
go build -o openim main.go
./openim -c config.yaml
```

The output of `run success` indicates that the conversion is successful and the program exits.


## CHAT
* `https://github.com/openimsdk/chat`
* Compatible with `https://github.com/openimsdk/open-im-server` v3 version
* v1.6 version, versions below v1.6 can be upgraded to v1.6 and above.
* Only use `https://github.com/openimsdk/open-im-server`. Self-implemented business servers do not need to be converted.

```shell
cd chat
go build -o chat main.go
./chat -c config.yaml
```

The output of `run success` indicates that the conversion is successful and the program exits.
