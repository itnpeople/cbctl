# cbctl

* Cloud-Barista MCKS Command Line Interface
* [MCKS : Multi Cloud Kubernetes Service](https://github.com/cloud-barista/cb-mcks)


## Quick Started

* Installation

```
### MacOS
$ curl -LO "https://github.com/itnpeople/cbctl/releases/download/$(curl -s https://api.github.com/repos/itnpeople/cbctl/releases/latest | grep tag_name | sed -E 's/.*"([^"]+)".*/\1/')/cbctl-darwin-amd64"
$ mv cbctl-darwin-amd64 /usr/local/bin/cbctl
$ chmod +x /usr/local/bin/cbctl
$ cbctl version

###Linux
$ curl -LO "https://github.com/itnpeople/cbctl/releases/download/$(curl -s https://api.github.com/repos/itnpeople/cbctl/releases/latest | grep tag_name | sed -E 's/.*"([^"]+)".*/\1/')/cbctl-linux-amd64"
$ mv cbctl-linux-amd64 /usr/local/bin/cbctl
$ chmod +x /usr/local/bin/cbctl
$ cbctl version

### Windows
https://github.com/itnpeople/cbctl/releases/download/v0.0.0/cbctl-windows-amd64.exe
```

* start up the MCKS

```
$ ./examples/lab/startup.sh
```

* verify running
```
$ docker ps

CONTAINER ID   IMAGE                                COMMAND                  CREATED      STATUS      PORTS                                        NAMES
3aedcfdb6c8c   cloudbaristaorg/cb-mcks:latest       "/app/cb-mcks"           2 days ago   Up 2 days   0.0.0.0:1470->1470/tcp, 50254/tcp            cb-mcks
3e8f6ad76539   cloudbaristaorg/cb-tumblebug:0.5.0   "/app/src/cb-tumbleb…"   2 days ago   Up 2 days   0.0.0.0:1323->1323/tcp, 50252/tcp            cb-tumblebug
283b91eeb270   cloudbaristaorg/cb-spider:0.5.0      "/root/go/src/github…"   2 days ago   Up 2 days   2048/tcp, 0.0.0.0:1024->1024/tcp, 4096/tcp   cb-spider

$ docker logs cb-mcks -f
```

* Initialize cloud connection (cb-spider)

```
$ cbctl create driver --csp aws
$ cbctl create credential --csp aws --name crdential-aws --secret-id "$AWS_SECRET_ID" --secret "$AWS_SECRET_KEY"
$ cbctl create region --csp aws --name region-aws-tokyo --region ap-northeast-1 --zone ap-northeast-1a 
$ cbctl create connection --csp aws --name config-aws-tokyo --region region-aws-tokyo --credential credential-aws
```

* Kubernetes cluster provisioning
```
$ cbctl create cluster \
  --name "cb-cluster"\
  --control-plane-connection="config-aws-tokyo"\
  --control-plane-count="1"\
  --control-plane-spec="t2.medium"\
  --worker-connection="config-aws-tokyo"\
  --worker-count="1"\
  --worker-spec="t2.medium"
```

## User Guide

```
$ cbctl create [cluster/node/driver/credential/region/connection]
$ cbctl get [cluster/node/driver/credential/region/connection/mcis]
$ cbctl delete [cluster/node/driver/credential/region/connection/mcis]
$ cbctl update-kubeconfig
$ cbctl get-key
$ cbctl clean [mcis]
```

### Create

* Create a Cluster

```
$ cbctl create cluster \
  --namespace "acornsoft"\
  --name "cb-cluster"\
  --control-plane-connection="config-aws-tokyo"\
  --control-plane-count="1"\
  --control-plane-spec="t2.medium"\
  --worker-connection="config-gcp-tokyo"\
  --worker-count="1"\
  --worker-spec="e2-highcpu-4"
```

* Create Nodes
```
$ cbctl create node \
 --namespace "acornsoft"\
 --cluster "cb-cluster"\
 --worker-connection="config-aws-tokyo"\
 --worker-count="1"\
 --worker-spec="t2.medium"
```

* Create a cloud driver.
```
$ cbctl create driver --csp [CSP]

# examples
$ cbctl create driver --csp aws
$ cbctl create driver --csp gcp
$ cbctl create driver --csp azure
$ cbctl create driver --csp alibaba
$ cbctl create driver --csp tencent
$ cbctl create driver --csp ibm
$ cbctl create driver --csp openstack
$ cbctl create driver --csp cloudit
```

* Create a Region
```
$ cbctl create region [region name] --csp [CSP] --region [CSP region name] --zone [CSP zone name]
$ cbctl create region [region name] --csp azure --location [location] --resource-group [resource group name]     # if azure

# examples
$ cbctl create region region-aws-tokyo --csp aws --region ap-northeast-1 --zone ap-northeast-1a 
$ cbctl create region region-gcp-tokyo --csp gcp --region asia-northeast1 --zone asia-northeast1-a
$ cbctl create region region-azure-tokyo --csp azure --location japaneast --resource-group cb-mcks
$ cbctl create region region-alibaba-tokyo --csp alibaba --region ap-northeast-1 --zone ap-northeast-1a
$ cbctl create region region-tencent-tokyo --csp tencent --region ap-tokyo --zone ap-tokyo-2
$ cbctl create region region-ibm-tokyo --csp ibm --region jp-tok --zone jp-tok-1
```

* Create a credential
```
$ source ./examples/credentials.sh \
  aws="${HOME}/.aws/credentials" \
  gcp="${HOME}/.ssh/google-credential-cloudbarista.json" \
  azure="${HOME}/.azure/azure-credential-cloudbarista.json" \
  alibaba="${HOME}/.ssh/alibaba_accesskey.csv" \
  tencent="${HOME}/.tccli/default.credential" \
  openstack="${HOME}/.ssh/openstack-openrc.sh" \
  cloudit="${HOME}/.ssh/cloudit-credential.sh"

$ cbctl create credential credential-aws --csp aws --secret-id "$AWS_SECRET_ID" --secret "$AWS_SECRET_KEY"
$ cbctl create credential credential-gcp --csp gcp --client-email "$GCP_SA" --project-id "$GCP_PROJECT" --private-key "$GCP_PKEY"
$ cbctl create credential credential-azure --csp azure --secret-id "$AZURE_CLIENT_ID" --secret "$AZURE_CLIENT_SECRET" --subscription "$AZURE_SUBSCRIPTION_ID" --tenant "$AZURE_TENANT_ID"
$ cbctl create credential credential-alibaba --csp alibaba --secret-id "$ALIBABA_SECRET_ID" --secret "$ALIBABA_SECRET_KEY"
$ cbctl create credential credential-tencent --csp tencent --secret-id "$TENCENT_SECRET_ID" --secret "$TENCENT_SECRET_KEY"
$ cbctl create credential credential-ibm --csp ibm --api-key "$IBM_API_KEY"
$ cbctl create credential credential-openstack --csp openstack --endpoint "$OS_AUTH_URL" --project-id "$OS_PROJECT_ID" --username "$OS_USERNAME" --password "$OS_PASSWORD" --domain "$OS_USER_DOMAIN_NAME"
$ cbctl create credential credential-cloudit --csp cloudit --endpoint "$CLOUDIT_ENDPOINT" --username "$CLOUDIT_USERNAME" --password "$CLOUDIT_PASSWORD" --token "$CLOUDIT_TOKEN" --tenant "$CLOUDIT_TENANT_ID"
```


* Create a Connection Info.
```
$ cbctl create connection [connection name] --csp [CSP] --region [region name] --credential [credential name]

# examples
$ cbctl create connection config-aws-tokyo --csp aws --region region-aws-tokyo --credential credential-aws
$ cbctl create connection config-gcp-tokyo --csp gcp --region region-gcp-tokyo --credential credential-gcp
$ cbctl create connection config-azure-tokyo --csp azure --region region-azure-tokyo --credential credential-azure
$ cbctl create connection config-alibaba-tokyo --csp alibaba --region region-alibaba-tokyo --credential credential-alibaba
$ cbctl create connection config-tencent-tokyo --csp tencent --region region-tencent-tokyo --credential credential-tencent
$ cbctl create connection config-ibm-tokyo --csp ibm --region region-ibm-tokyo --credential credential-ibm
```

* Create a cloud-barista namespace

```
$ cbctl create namespace [namespace name]

# example
$ cbctl create namespace acornsoft
```

### Get

* Get clusters
```
$ cbctl get cluster
$ cbctl get cluster [custer name]
$ cbctl get cluster --name [custer name]

# examples
$ cbctl get cluster
$ cbctl get cluster "cb-cluster"
```

* Get nodes
```
$ cbctl get node --cluster [cluster name]
$ cbctl get node [node name] --cluster [cluster name]
$ cbctl get node --name [node name] --cluster [cluster name]


# examples
$ cbctl get node --cluster "cb-cluster"
$ cbctl get node "w-1-j4j8z" --cluster "cb-cluster"
```

* Get drivers

```
$ cbctl get driver
$ cbctl get driver [driver name]
$ cbctl get driver --name [driver name]
$ cbctl get driver --csp [CSP]

# examples
$ cbctl get driver
$ cbctl get driver "aws-driver-v1.0"
$ cbctl get driver --name "aws-driver-v1.0"
$ cbctl get driver --csp "aws"
```

* Get credentials
```
$ cbctl get credential
$ cbctl get credential [credential name]
$ cbctl get credential --name [credential name]

# examples
$ cbctl get credential
$ cbctl get credential "credential-aws"
```

* Get regions
```
$ cbctl get region
$ cbctl get region [region name]
$ cbctl get region --name [region name]

# examples
$ cbctl get region
$ cbctl get region "region-aws-tokyo"
```

* Get connections
```
$ cbctl get connection
$ cbctl get connection [connection anme]
$ cbctl get connection --name [connection anme]

# examples
$ cbctl get connection
$ cbctl get connection "config-aws-tokyo"
```

* Get VM specifications (verify the Connection Info.)

```
$ cbctl get spec --connection [connection name]


# examples
$ cbctl get spec --connection config-aws-tokyo
```

* Get MCISes
```
$ cbctl get mcis
$ cbctl get mcis [mcis name]
$ cbctl get mcis --name [mcis name]

$ cbctl get mcis
$ cbctl get mcis "cb-cluster"
```


### Delete

* Delete the cluster
```
$ cbctl delete cluster [cluster name]
$ cbctl delete cluster --name [cluster name]

# exampels
$ cbctl delete cluster "cb-cluster"
```

* Delete the node

```
$ cbctl delete node [node name] --cluster [cluster name]
$ cbctl delete node --name [node name] --cluster [cluster name]

# exampels
$ cbctl delete node "w-1-j4j8z" --cluster "cb-cluster"
```

* Delete the MCIS
```
$ cbctl delete mcis [mcis name]
$ cbctl delete mcis --name [mcis name]


# examples
$ cbctl delete mcis "cb-cluster"
```


* Delete the MCIRs
```
$ cbctl delete driver [driver name]
$ cbctl delete driver --name [driver name]
$ cbctl delete driver --csp [CSP]

$ cbctl delete credential [credential name]
$ cbctl delete credential --name [credential name]

$ cbctl delete region [region name]
$ cbctl delete region --name [region name]

$ cbctl delete connection [connection info. name]
$ cbctl delete connection --name [connection info. name]

# examples
$ cbctl delete driver "aws-driver-v1.0"
$ cbctl delete driver --csp "aws"
$ cbctl delete credential "credential-aws"
$ cbctl delete region "region-aws-tokyo"
$ cbctl delete connection "config-aws-tokyo"
```

### Update-Kubeconfig

```
$ cbctl update-kubeconfig [cluster-name]
$ cbctl update-kubeconfig --name [cluster-name]

# examples
$ cbctl update-kubeconfig "cb-cluster"
$ kubectl config current-context
```

### Get-Key
```
$ cbctl get-key  [node name] --cluster [cluster-name]
$ cbctl get-key  --name [node name] --cluster [cluster-name]

# examples
$ cbctl get-key  "w-1-j4j8z" --cluster "cb-cluster"  > output/w-1-j4j8z.pem
$ chmod 400 output/w-1-j4j8z.pem
$ ssh -i output/w-1-j4j8z.pem cb-user@xxx.xxx.xxx.xxx
```

### Using Yaml File (filename)
```
$ cbctl create [cluster/node/driver/region/credential/connection/namespace] -f [URL]
$ cbctl create [cluster/node/driver/region/credential/connection/namespace] -f [FILENAME]
$ cbctl create [cluster/node/driver/region/credential/connection/namespace] -f [STDIN]
```

* examples - URL, FILENAME
```
$ cbctl create cluster -f examples/yaml/create-cluster.yaml
$ cbctl create cluster -f https://github.com/itnpeople/cbctl/blob/master/examples/yaml/create-cluster.yaml
```


* examples - STDIN
```
# cluster
$ cbctl create cluster -f - <<EOF
name: cb-cluster
label: lab.
description: create a cluster test
controlPlane:
  - connection: config-aws-tokyo
    count: 1
    spec: t2.medium
worker:
  - connection: config-gcp-tokyo
    count: 1
    spec: e2-highcpu-4
config:
  kubernetes:
    networkCni: calico
    podCidr: 10.244.0.0/16
    serviceCidr: 10.96.0.0/12
    serviceDnsDomain: cluster.local
EOF

# nodes
$ cbctl create node --cluster cb-cluster -f - <<EOF
worker: 
  - connection: config-aws-tokyo
    count: 1
    spec: t2.medium
EOF

# driver
$ cbctl create driver -f - <<EOF
DriverName : "aws-driver-v1.0"
ProviderName : "AWS"
DriverLibFileName : "aws-driver-v1.0.so"
EOF

# region
$ cbctl create region -f - <<EOF
RegionName : "region-aws-tokyo"
ProviderName : "AWS"
KeyValueInfoList :
- Key : "Region"
  Value : "ap-northeast-1"
- Key : "Zone"
  Value : "ap-northeast-1a"
EOF

# credential
$ cbctl create credential -f - <<EOF
CredentialName : "credential-aws"
ProviderName : "AWS"
KeyValueInfoList :
- Key : "ClientId"
  Value : "aaaaaaa"
- Key : "ClientSecret"
  Value : "bbbbbbbbbbbbbbbbbbbbbbbbb"
EOF

# connection
$ cbctl create connection -f - <<EOF
ConfigName : "config-aws-tokyo"
ProviderName : "AWS" 
DriverName : "aws-driver-v1.0" 
CredentialName : "credential-aws"
RegionName : "region-aws-tokyo"
EOF

# namespace
$ cbctl create namespace -f - <<EOF
name : "acornsoft"
description : "acornsoft namespace"
EOF
```

### Plugins

```
$ cbctl plugin
$ cbctl <plugin name>
```

### Clean-up

```
$ cbctl clean mcir
```

### Persistent flags

```
--config [config file path (default:.config)]

--output [json/yaml(default)]
--o [json/yaml(default)]
```

* Optional persistent flags (config)

```
--namespace [cloud-barista namespace (default:acornsoft)]
-n [cloud-barista namespace (default:acornsoft)]
```

#### Using plugin examples

* create a executable plugin (on PATH)
> plugin name = cbctl-foo (prefix : cbctl)

```
$ cat > /usr/local/bin/cbctl-foo <<EOF
#!/bin/bash
echo "I am plugin foo"
EOF

$ chmod +x /usr/local/bin/cbctl-foo
```

* create a executable plugin (on plugin directory)
> plugin name = foo

```
$ mkdir ${HOME}/.cbctl/plugins
$ cat > ${HOME}/.cbctl/plugins/foo <<EOF
#!/bin/bash
echo "I am plugin foo"
EOF

$ chmod +x ${HOME}/.cbctl/plugins/foo
```

* plugin list

```
$ cbctl plugin
The following compatible plugins are available:

/usr/local/bin/cbctl-foo
```

* execute plugin
```
$ cbctl foo
I am plugin foo
```

### Config

```
$ cbctl config
```

* Context
```
$ cbctl config add-context ctx1 \
 --namespace default \
 --url-mcks http://127.0.0.1:1470/mcks \
 --url-spider http://127.0.0.1:1024/spider \
 --url-tumblebug http://127.0.0.1:1323/tumblebug

$ cbctl config list-context
$ cbctl config get-context ctx1
$ cbctl config delete-context ctx1
```

* Current context
```
$ cbctl config current-context
$ cbctl config current-context ctx1
$ cbctl config set-namespace namespace1
```

```
$ cbctl config view
```
