# confd

[![Integration Tests](https://github.com/haad/confd/actions/workflows/integration-tests.yml/badge.svg)](https://github.com/haad/confd/actions/workflows/integration-tests.yml)
[![CodeQL](https://github.com/haad/confd/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/haad/confd/actions/workflows/codeql-analysis.yml)
[![Super-Linter](https://github.com/haad/confd/actions/workflows/superlinter.yml/badge.svg)](https://github.com/haad/confd/actions/workflows/superlinter.yml)

`confd` is a lightweight configuration management tool focused on:

* keeping local configuration files up-to-date using data stored in [etcd](https://github.com/coreos/etcd),
  [consul](http://consul.io), [dynamodb](http://aws.amazon.com/dynamodb/), [redis](http://redis.io),
  [vault](https://vaultproject.io), [zookeeper](https://zookeeper.apache.org), [aws ssm parameter store](https://aws.amazon.com/ec2/systems-manager/) or env vars and processing [template resources](docs/template-resources.md).
* reloading applications to pick up new config file changes

## Community

* IRC: `#confd` on Freenode
* Mailing list: [Google Groups](https://groups.google.com/forum/#!forum/confd-users)
* Website: [www.confd.io](http://www.confd.io)

## Building

Go 1.10 is required to build confd, which uses the new vendor directory.

```sh
$ mkdir -p $GOPATH/src/github.com/haad
$ git clone https://github.com/haad/confd.git $GOPATH/src/github.com/haad/confd
$ cd $GOPATH/src/github.com/haad/confd
$ make
```

You should now have confd in your `bin/` directory:

```sh
$ ls bin/
confd
```

### Running integration tests

```sh
docker run -it --rm -v $(pwd):/go/src/github.com/haad/confd golang:1.17.6 /go/src/github.com/haad/confd/integration/run.sh
```

## Getting Started

Before we begin be sure to [download and install confd](docs/installation.md).

* [quick start guide](docs/quick-start-guide.md)

## Next steps

Check out the [docs directory](docs) for more docs.
