# cbctl

* Cloud-Barista MCKS Command Line Interface
* MCKS : Multi Cloud Kubernetes Service
* https://github.com/cloud-barista/cb-mcks


## Quick Started

### Installation

* MacOS
```
$ curl -LO "https://github.com/itnpeople/cbctl/releases/download/$(curl -s https://api.github.com/repos/itnpeople/cbctl/releases/latest | grep tag_name | sed -E 's/.*"([^"]+)".*/\1/')/cbctl-darwin-amd64"
$ mv cbctl-darwin-amd64 /usr/local/bin/cbctl
$ chmod +x /usr/local/bin/cbctl
$ cbctl version
```

* Linux
```
$ curl -LO "https://github.com/itnpeople/cbctl/releases/download/$(curl -s https://api.github.com/repos/itnpeople/cbctl/releases/latest | grep tag_name | sed -E 's/.*"([^"]+)".*/\1/')/cbctl-linux-amd64"
$ mv cbctl-linux-amd64 /usr/local/bin/cbctl
$ chmod +x /usr/local/bin/cbctl
$ cbctl version
```

* Windows

```
https://github.com/itnpeople/cbctl/releases/download/v0.0.0/cbctl-windows-amd64.exe
```

### Run the MCKS

* start up

```
$ ./examples/lab/startup.sh
```

* Verify running
```
$ docker ps

CONTAINER ID   IMAGE                                COMMAND                  CREATED      STATUS      PORTS                                        NAMES
3aedcfdb6c8c   cloudbaristaorg/cb-mcks:latest       "/app/cb-mcks"           2 days ago   Up 2 days   0.0.0.0:1470->1470/tcp, 50254/tcp            cb-mcks
3e8f6ad76539   cloudbaristaorg/cb-tumblebug:0.5.0   "/app/src/cb-tumbleb…"   2 days ago   Up 2 days   0.0.0.0:1323->1323/tcp, 50252/tcp            cb-tumblebug
283b91eeb270   cloudbaristaorg/cb-spider:0.5.0      "/root/go/src/github…"   2 days ago   Up 2 days   2048/tcp, 0.0.0.0:1024->1024/tcp, 4096/tcp   cb-spider

$ docker logs cb-mcks -f
```

### Create a Cluster

* Initialize (cb-spider)

```
$ cbctl driver create --csp aws
$ cbctl credential create --csp aws --name crdential-aws --secret-id "$AWS_SECRET_ID" --secret "$AWS_SECRET_KEY"
$ cbctl region create --csp aws --name region-aws-tokyo --region ap-northeast-1 --zone ap-northeast-1a 
$ cbctl connection create --csp aws --name config-aws-tokyo --region region-aws-tokyo --credential credential-aws
```

* Kubernetes cluster provisioning
```
$ cbctl cluster create \
  --name "cb-cluster"\
  --control-plane-connection="config-aws-tokyo"\
  --control-plane-count="1"\
  --control-plane-spec="t2.medium"\
  --worker-connection="config-aws-tokyo"\
  --worker-count="1"\
  --worker-spec="t2.medium"
```

## User Guide

* Commands
```
$ cbctl
$ cbctl cluster
$ cbctl node
$ cbctl driver
$ cbctl credential
$ cbctl region
$ cbctl connection
$ cbctl plugin
```

### CB-MCKS

* Cluster

```
$ cbctl cluster create \
  --name "cb-cluster"\
  --control-plane-connection="config-aws-tokyo"\
  --control-plane-count="1"\
  --control-plane-spec="t2.medium"\
  --worker-connection="config-aws-tokyo"\
  --worker-count="1"\
  --worker-spec="t2.medium"

$ cbctl cluster list 
$ cbctl cluster get --name "cb-cluster"
$ cbctl cluster delete --name "cb-cluster"
```

* Nodes

```
$ cbctl node add \
 --cluster "cb-cluster"\
 --worker-connection="config-aws-tokyo"\
 --worker-count="1"\
 --worker-spec="t2.medium"

$ cbctl node list --cluster "cb-cluster" 
$ cbctl node get --cluster "cb-cluster" --name "w-1-oiq77"
$ cbctl node delete --cluster "cb-cluster" --name "w-1-oiq77"
```

* Kubeconfig

```
$ cbctl cluster update-kubeconfig --name cb-cluster
$ kubectl config  current-context
```

* SSH private-key
```
$ cbctl node get-key --cluster cb-cluster --name w-1-j4j8z > output/w-1-j4j8z.pem
$ chmod 400 output/w-1-j4j8z.pem
$ ssh -i output/w-1-j4j8z.pem cb-user@xxx.xxx.xxx.xxx
```

* Using Yaml File

```
$ cbctl cluster create -f examples/yaml/create-cluster.yaml

$ cbctl node add --cluster cb-cluster -f - <<EOF
worker: 
  - connection: config-aws-tokyo
    count: 1
    spec: t2.medium
EOF
```

* Persistent flags

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


### Initialize Cloud Connection Info.
> cb-spider

* Driver

```
$ cbctl driver create --csp aws
$ cbctl driver list
$ cbctl driver get --csp aws
$ cbctl driver delete --csp aws
```

* Credential
```
$ source ./examples/credentials.sh \
  aws="${HOME}/.aws/credentials" \
  gcp="${HOME}/.ssh/google-credential-cloudbarista.json" \
  azure="${HOME}/.azure/azure-credential-cloudbarista.json" \
  alibaba="${HOME}/.ssh/alibaba_accesskey.csv" \
  tencent="${HOME}/.tccli/default.credential" \
  openstack="${HOME}/.ssh/openstack-openrc.sh"

$ cbctl credential create --csp aws --name crdential-aws --secret-id "$AWS_SECRET_ID" --secret "$AWS_SECRET_KEY"
$ cbctl credential create --csp gcp --name credential-gcp --client-email "$GCP_SA" --project-id "$GCP_PROJECT" --private-key "$GCP_PKEY"
$ cbctl credential create --csp azure --name credential-azure --secret-id "$AZURE_CLIENT_ID" --secret "$AZURE_CLIENT_SECRET" --subscription "$AZURE_SUBSCRIPTION_ID" --tenant "$AZURE_TENANT_ID"
$ cbctl credential create --csp alibaba --name credential-alibaba --secret-id "$ALIBABA_SECRET_ID" --secret "$ALIBABA_SECRET_KEY"
$ cbctl credential create --csp tencent --name credential-tencent --secret-id "$TENCENT_SECRET_ID" --secret "$TENCENT_SECRET_KEY"
$ cbctl credential create --csp ibm --name credential-ibm --api-key "$IBM_API_KEY"

$ cbctl credential list
$ cbctl credential get --name credential-aws
$ cbctl credential delete --name credential-aws
```

* Region
```
$ cbctl region create --csp aws --name region-aws-tokyo --region ap-northeast-1 --zone ap-northeast-1a 
$ cbctl region list
$ cbctl region get --name region-aws-tokyo
$ cbctl region delete --name region-aws-tokyo
```

* Connection Info.
```
$ cbctl connection create --csp aws --name config-aws-tokyo --region region-aws-tokyo --credential credential-aws
$ cbctl connection list
$ cbctl connection get --name config-aws-tokyo
$ cbctl connection delete --name config-aws-tokyo
```

### Plugins

```
$ cbctl plugin list
$ cbctl <plugin name>
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
$ cbctl plugin list
The following compatible plugins are available:

/usr/local/bin/cbctl-foo
```

* execute plugin
```
$ cbctl foo
I am plugin foo
```
