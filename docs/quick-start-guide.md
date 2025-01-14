# Quick Start Guide

Before we begin be sure to [download and install confd](installation.md).

## Select a backend

confd supports the following backends:

* etcd
* consul
* vault
* environment variables
* file
* redis
* zookeeper
* dynamodb
* rancher
* ssm (AWS Simple Systems Manager Parameter Store)
* secretsmanager ( AWS secrets Manager)

### Add keys

This guide assumes you have a working [etcd](https://github.com/coreos/etcd#getting-started), or [consul](http://www.consul.io/intro/getting-started/install.html) server up and running and the ability to add new keys.

#### etcd

```sh
etcdctl set /myapp/database/url db.example.com
etcdctl set /myapp/database/user rob
```

#### consul

```sh
curl -X PUT -d 'db.example.com' http://localhost:8500/v1/kv/myapp/database/url
curl -X PUT -d 'rob' http://localhost:8500/v1/kv/myapp/database/user
```

#### vault
```sh
vault mount -path myapp generic
vault write myapp/database url=db.example.com user=rob
```

#### environment variables

```sh
export MYAPP_DATABASE_URL=db.example.com
export MYAPP_DATABASE_USER=rob
```

#### file

myapp.yaml
```yaml
myapp:
  database:
    url: db.example.com
    user: rob
```

#### redis

```sh
redis-cli set /myapp/database/url db.example.com
redis-cli set /myapp/database/user rob
```

#### zookeeper

```text
[zk: localhost:2181(CONNECTED) 1] create /myapp ""
[zk: localhost:2181(CONNECTED) 2] create /myapp/database ""
[zk: localhost:2181(CONNECTED) 3] create /myapp/database/url "db.example.com"
[zk: localhost:2181(CONNECTED) 4] create /myapp/database/user "rob"
```

#### dynamodb

First create a table with the following schema:

```sh
aws dynamodb create-table \
    --region <YOUR_REGION> --table-name <YOUR_TABLE> \
    --attribute-definitions AttributeName=key,AttributeType=S \
    --key-schema AttributeName=key,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
```

Now create the items. The attribute value `value` must be of type [string](http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html):

```sh
aws dynamodb put-item --table-name <YOUR_TABLE> --region <YOUR_REGION> \
    --item '{ "key": { "S": "/myapp/database/url" }, "value": {"S": "db.example.com"}}'
aws dynamodb put-item --table-name <YOUR_TABLE> --region <YOUR_REGION> \
    --item '{ "key": { "S": "/myapp/database/user" }, "value": {"S": "rob"}}'
```

#### Rancher

This backend consumes the [Rancher](https://www.rancher.com) metadata service. For available keys, see the [Rancher Metadata Service docs](http://docs.rancher.com/rancher/rancher-services/metadata-service/).

#### ssm

```sh
aws ssm put-parameter --name "/myapp/database/url" --type "String" --value "db.example.com"
aws ssm put-parameter --name "/myapp/database/user" --type "SecureString" --value "rob"
```

#### secretsmanager

```
aws secretsmanager  create-secret --name "/myapp/database/url" --secret-string "db.example.com"
aws secretsmanager  create-secret --name "/myapp/database/user" --secret-string "rob"
```

### Create the confdir

The confdir is where template resource configs and source templates are stored.

```sh
sudo mkdir -p /etc/confd/{conf.d,templates}
```

### Create a template resource config

Template resources are defined in [TOML](https://github.com/mojombo/toml) config files under the `confdir`.

/etc/confd/conf.d/myconfig.toml
```TOML
[template]
src = "myconfig.conf.tmpl"
dest = "/tmp/myconfig.conf"
keys = [
    "/myapp/database/url",
    "/myapp/database/user",
]
make_directories = true # optional, whether to create dest directory if not exist
```

### Create the source template

Source templates are [Golang text templates](http://golang.org/pkg/text/template/#pkg-overview).

/etc/confd/templates/myconfig.conf.tmpl
```go
[myconfig]
database_url = {{getv "/myapp/database/url"}}
database_user = {{getv "/myapp/database/user"}}
```

### Process the template

confd supports two modes of operation daemon and onetime. In daemon mode confd polls a backend for changes and updates destination configuration files if necessary.

#### etcd

```sh
confd -onetime -backend etcd -node http://127.0.0.1:2379
```

#### consul

```sh
confd -onetime -backend consul -node 127.0.0.1:8500
```

#### vault
```sh
ROOT_TOKEN=$(vault read -field id auth/token/lookup-self)

confd -onetime -backend vault -node http://127.0.0.1:8200 \
      -auth-type token -auth-token $ROOT_TOKEN
```

#### dynamodb

```sh
confd -onetime -backend dynamodb -table <YOUR_TABLE>
```

#### env

```sh
confd -onetime -backend env
```

#### file

```sh
confd -onetime -backend file -file myapp.yaml
```

#### redis

```sh
confd -onetime -backend redis -node 192.168.255.210:6379
```
or if you want to connect to a specific redis database (4 in this example):

```sh
confd -onetime -backend redis -node 192.168.255.210:6379/4
```

#### rancher

```sh
confd -onetime -backend rancher -prefix /2015-07-25
```

*Note*: The metadata API prefix can be defined on the cli, or as part of your keys in the template toml file.

Output:
```text
2014-07-08T20:38:36-07:00 confd[16252]: INFO Target config /tmp/myconfig.conf out of sync
2014-07-08T20:38:36-07:00 confd[16252]: INFO Target config /tmp/myconfig.conf has been updated
```

The `dest` configuration file should now be in sync.

```sh
cat /tmp/myconfig.conf
```

Output:
```text
# This a comment
[myconfig]
database_url = db.example.com
database_user = rob
```

#### ssm

```sh
confd -onetime -backend ssm
```

#### secretsmanager

```
confd -onetime -backend secretsmanager
```

## Advanced Example

In this example we will use confd to manage two nginx config files using a single template.

### Add keys

#### etcd

```sh
etcdctl set /myapp/subdomain myapp
etcdctl set /myapp/upstream/app2 "10.0.1.100:80"
etcdctl set /myapp/upstream/app1 "10.0.1.101:80"
etcdctl set /yourapp/subdomain yourapp
etcdctl set /yourapp/upstream/app2 "10.0.1.102:80"
etcdctl set /yourapp/upstream/app1 "10.0.1.103:80"
```

#### consul

```sh
curl -X PUT -d 'myapp' http://localhost:8500/v1/kv/myapp/subdomain
curl -X PUT -d '10.0.1.100:80' http://localhost:8500/v1/kv/myapp/upstream/app1
curl -X PUT -d '10.0.1.101:80' http://localhost:8500/v1/kv/myapp/upstream/app2
curl -X PUT -d 'yourapp' http://localhost:8500/v1/kv/yourapp/subdomain
curl -X PUT -d '10.0.1.102:80' http://localhost:8500/v1/kv/yourapp/upstream/app1
curl -X PUT -d '10.0.1.103:80' http://localhost:8500/v1/kv/yourapp/upstream/app2
```

### Create template resources

/etc/confd/conf.d/myapp-nginx.toml

```TOML
[template]
prefix = "/myapp"
src = "nginx.tmpl"
dest = "/tmp/myapp.conf"
owner = "nginx"
mode = "0644"
keys = [
  "/subdomain",
  "/upstream",
]
check_cmd = "/usr/sbin/nginx -t -c {{.src}}"
reload_cmd = "/usr/sbin/service nginx reload"
make_directories = true # optional, whether to create dest directory if not exist
```

/etc/confd/conf.d/yourapp-nginx.toml

```TOML
[template]
prefix = "/yourapp"
src = "nginx.tmpl"
dest = "/tmp/yourapp.conf"
owner = "nginx"
mode = "0644"
keys = [
  "/subdomain",
  "/upstream",
]
check_cmd = "/usr/sbin/nginx -t -c {{.src}}"
reload_cmd = "/usr/sbin/service nginx reload"
make_directories = true # optional, whether to create dest directory if not exist
```

### Create the source template

/etc/confd/templates/nginx.tmpl
```go
upstream {{getv "/subdomain"}} {
{{range getvs "/upstream/*"}}
    server {{.}};
{{end}}
}

server {
    server_name  {{getv "/subdomain"}}.example.com;
    location / {
        proxy_pass        http://{{getv "/subdomain"}};
        proxy_redirect    off;
        proxy_set_header  Host             $host;
        proxy_set_header  X-Real-IP        $remote_addr;
        proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
   }
}
```
