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
$ cbctl cluster get "cb-cluster"
$ cbctl cluster delete "cb-cluster"
```

* Nodes

```
$ cbctl node add \
 --cluster "cb-cluster"\
 --worker-connection="config-aws-tokyo"\
 --worker-count="1"\
 --worker-spec="t2.medium"

$ cbctl node list --cluster "cb-cluster" 
$ cbctl node get "w-1-oiq77" --cluster "cb-cluster"
$ cbctl node delete "w-1-oiq77" --cluster "cb-cluster"
```

* Kubeconfig

```
$ cbctl cluster update-kubeconfig "cb-cluster"
$ kubectl config  current-context
```

* SSH private-key
```
$ cbctl node get-key  w-1-j4j8z --cluster cb-cluster  > output/w-1-j4j8z.pem
$ chmod 400 output/w-1-j4j8z.pem
$ ssh -i output/w-1-j4j8z.pem cb-user@xxx.xxx.xxx.xxx
```

* Using Yaml File (create a cluster & add nodes)

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
$ cbctl driver create aws
$ cbctl driver list
$ cbctl driver get aws
$ cbctl driver delete aws
```

* Credential
```
$ source ./examples/credentials.sh \
  aws="${HOME}/.aws/credentials" \
  gcp="${HOME}/.ssh/google-credential-cloudbarista.json" \
  azure="${HOME}/.azure/azure-credential-cloudbarista.json" \
  alibaba="${HOME}/.ssh/alibaba_accesskey.csv" \
  tencent="${HOME}/.tccli/default.credential" \
  openstack="${HOME}/.ssh/openstack-openrc.sh" \
  cloudit="${HOME}/.ssh/cloudit-credential.sh"

$ cbctl credential create credential-aws --csp aws --secret-id "$AWS_SECRET_ID" --secret "$AWS_SECRET_KEY"
$ cbctl credential create credential-gcp --csp gcp --client-email "$GCP_SA" --project-id "$GCP_PROJECT" --private-key "$GCP_PKEY"
$ cbctl credential create credential-azure --csp azure --secret-id "$AZURE_CLIENT_ID" --secret "$AZURE_CLIENT_SECRET" --subscription "$AZURE_SUBSCRIPTION_ID" --tenant "$AZURE_TENANT_ID"
$ cbctl credential create credential-alibaba --csp alibaba --secret-id "$ALIBABA_SECRET_ID" --secret "$ALIBABA_SECRET_KEY"
$ cbctl credential create credential-tencent --csp tencent --secret-id "$TENCENT_SECRET_ID" --secret "$TENCENT_SECRET_KEY"
$ cbctl credential create credential-ibm --csp ibm --api-key "$IBM_API_KEY"
$ cbctl credential create credential-openstack --csp openstack --endpoint "$OS_AUTH_URL" --project-id "$OS_PROJECT_ID" --username "$OS_USERNAME" --password "$OS_PASSWORD" --domain "$OS_USER_DOMAIN_NAME"
$ cbctl credential create credential-cloudit --csp cloudit --endpoint "$CLOUDIT_ENDPOINT" --username "$CLOUDIT_USERNAME" --password "$CLOUDIT_PASSWORD" --token "$CLOUDIT_TOKEN" --tenant "$CLOUDIT_TENANT_ID"

$ cbctl credential list
$ cbctl credential get credential-aws
$ cbctl credential delete credential-aws
```

* Region
```
$ cbctl region create region-aws-tokyo --csp aws --region ap-northeast-1 --zone ap-northeast-1a 
$ cbctl region create region-gcp-tokyo --csp gcp --region asia-northeast1 --zone asia-northeast1-a
$ cbctl region create region-azure-tokyo --csp azure --location japaneast --resource-group cb-mcks
$ cbctl region create region-alibaba-tokyo --csp alibaba --region ap-northeast-1 --zone ap-northeast-1a
$ cbctl region create region-tencent-tokyo --csp tencent --region ap-tokyo --zone ap-tokyo-2
$ cbctl region create region-ibm-tokyo --csp ibm --region jp-tok --zone jp-tok-1

$ cbctl region list
$ cbctl region get region-aws-tokyo
$ cbctl region delete region-aws-tokyo
```

* Connection Info.
```
$ cbctl connection create config-aws-tokyo --csp aws --region region-aws-tokyo --credential credential-aws
$ cbctl connection create config-gcp-tokyo --csp gcp --region region-gcp-tokyo --credential credential-gcp
$ cbctl connection create config-azure-tokyo --csp azure --region region-azure-tokyo --credential credential-azure
$ cbctl connection create config-alibaba-tokyo --csp alibaba --region region-alibaba-tokyo --credential credential-alibaba
$ cbctl connection create config-tencent-tokyo --csp tencent --region region-tencent-tokyo --credential credential-tencent
$ cbctl connection create config-ibm-tokyo --csp ibm --region region-ibm-tokyo --credential credential-ibm

$ cbctl connection list
$ cbctl connection get config-aws-tokyo
$ cbctl connection delete config-aws-tokyo
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
