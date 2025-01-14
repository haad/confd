# DNS SRV Records

SRV records can be used to declare the backend nodes; just use the `-srv-domain` flag.

## Examples

### etcd

```sh
dig SRV _etcd._tcp.confd.io
```

```text
...
;; ANSWER SECTION:
_etcd._tcp.confd.io.    300 IN  SRV 1 100 4001 etcd.confd.io.
```

-

```sh
confd -backend etcd -srv-domain confd.io
```

### consul

```sh
dig SRV _consul._tcp.confd.io
```

```text
...
;; ANSWER SECTION:
_consul._tcp.confd.io.  300 IN  SRV 1 100 8500 consul.confd.io.
```

-

```sh
confd -backend consul -srv-domain confd.io
```

## The backend scheme

By default the `scheme` is set to http; change it with the `-scheme` flag.

```sh
confd -scheme https -srv-domain confd.io
```

Both the SRV domain and scheme can be configured in the confd configuration file. See the [Configuration Guide](configuration-guide.md) for more details.
