# Installing Devtron

[![Join Discord](https://img.shields.io/badge/Join%20us%20on-Discord-e01563.svg)](https://discord.gg/jsRG5qx2gp)

## Devtron Installation

### Prerequisites

You will need to be ready with following prerequisites before Devtron installation

* A Kubernetes cluster \(preferably K8s 1.16 or above\). Check [Creating a Production grade EKS cluster using EKSCTL](https://devtron.ai/blog/creating-production-grade-kubernetes-eks-cluster-eksctl/)

### Installing Devtron

* [Install with Helm3 - Recommended](install-devtron-helm-3.md)
* [Install with Helm2](install-devtron-helm-2.md)
* [Install with kubectl](install-devtron-using-kubectl.md)

#### Installation status

Run following command

```bash
kubectl -n devtroncd get installers installer-devtron -o jsonpath='{.status.sync.status}'
```

The install commands initiates Devtron-operator which spins up all the Devtron micro-services one by one in about 20 mins. You can use the above command to check the status of the installation if the installation is still in progress, it will print `Downloaded`. When the installation is complete, it prints `Applied`.

#### Access devtron dashboard

**Login credentials**

For login use username:`admin` and for password run command mentioned below.

```bash
kubectl -n devtroncd get secret devtron-secret -o jsonpath='{.data.ACD_PASSWORD}' | base64 -d
```

#### Access grafana

**Grafana URL**

grafana is available at `BASE_URL/grafana`

**Login credentials**

To log into grafana use username: `admin` and for password run command mentioned below.

```bash
kubectl -n devtroncd get secret devtron-grafana-cred-secret -o jsonpath='{.data.admin-password}' | base64 -d
```

#### Configuration

**Configure Secrets**

For `helm` installation this section referes to _**secrets.env**_ section of `values.yaml`. For `kubectl` based installation it refers to `kind: secret` in _**install/devtron-operator-configs.yaml**_.

Following properties should be configured

| Parameter | Description | Default |
| ---: | :--- | :--- |
| **POSTGRESQL\_PASSWORD** | Using this parameter the auto-generated password for postgres can be edited as per requirement(Used by devtron to store the app information) | |
| **WEBHOOK\_TOKEN** | If you want to continue using jenkins for CI then provide this for authentication of requests, should be base64 encoded |  |

**Configure ConfigMaps**

For `helm` installation this section refers to _**configs**_ section of `values.yaml`. For `kubectl` based installation it refers to `kind: ConfigMap` in _**install/devtron-operator-configs.yaml**_.

Following properties should be configured

| Parameter | Description | Default |
| ---: | :--- | :--- |
| **BASE\_URL\_SCHEME** | either of http or https \(required\) | http |
| **BASE\_URL** | url without scheme and trailing slash, this is the domain pointing to the cluster on which devtron platform is being installed. For example if you have directed domain `devtron.example.com` to the cluster and ingress controller is listening on port `32080` then url will be `devtron.example.com:32080` \(required\) | `change-me` |
| **DEX\_CONFIG** | dex config if you want to integrate login with SSO \(optional\) for more information check [Argocd documentation](https://argoproj.github.io/argo-cd/operator-manual/user-management/) |  |
| **EXTERNAL\_SECRET\_AMAZON\_REGION** | AWS region for secret manager to pick \(required\) |  |
| **PROMETHEUS\_URL** | url of prometheus where all cluster data is stored, if this is wrong, you will not be able to see application metrics like cpu, ram, http status code, latency and throughput \(required\) |  |

#### Storage for Logs and Cache

AWS SPECIFIC

While installing devtron and using AWS-S3 bucket for storing the logs and caches, below parameters are to be used in the ConfigMap.
```NOTE : For using s3 bucket it is important to add s3 permission policy to the IAM role attached to the nodes of the cluster.```

| Parameter | Description | Default |
| ---: | :--- | :--- |
| **DEFAULT\_CACHE\_BUCKET** | AWS bucket to store docker cache, this should be created before hand \(required\) |  |
| **DEFAULT\_BUILD\_LOGS\_BUCKET** | AWS bucket to store build logs, this should be created before hand \(required\) |  |
| **DEFAULT\_CACHE\_BUCKET\_REGION** | AWS region of S3 bucket to store cache \(required\) |  |
| **DEFAULT\_CD\_LOGS\_BUCKET\_REGION** | AWS region of S3 bucket to store CD logs \(required\) |  |


AZURE SPECIFIC 

While installing devtron using Azure Blob Storage for storing logs and caches, below parameters will be used in the ConfigMap.


| Parameter | Description | Default |
| ---: | :--- | :--- |
| **AZURE\_ACCOUNT\_NAME** | Account name for AZURE Blob Storage |  |
| **AZURE\_BLOB\_CONTAINER\_CI\_LOG** | AZURE Blob storage container for storing ci-logs after running the CI pipeline |  |
| **AZURE\_BLOB\_CONTAINER\_CI\_CACHE** | AZURE Blob storage container for storing ci cache after running the CI pipeline |  |

To convert string to base64 use

```bash
echo -n "string" | base64 -d
```

**Please Note:** 
1\) Ensure that the **cluster has read and write access** to the S3 buckets/Azure Blob storage container mentioned in DEFAULT\_CACHE\_BUCKET, DEFAULT\_BUILD\_LOGS\_BUCKET or AZURE\_BLOB\_CONTAINER\_CI\_LOG, AZURE\_BLOB\_CONTAINER\_CI\_CACHE

2\) Ensure that cluster has **read access** to AWS secrets backends \(SSM & secrets manager\)

#### Cleanup

Run following commands to delete all the components installed by devtron

```bash
cd devtron-installation-script/
kubectl delete -n devtroncd -f yamls/
kubectl delete -n devtroncd -f charts/devtron/templates/devtron-installer.yaml
kubectl delete -n devtroncd -f charts/devtron/templates/install.yaml
kubectl delete -n devtroncd -f charts/devtron/crds
kubectl delete ns devtroncd
```

#### Troubleshooting steps

**1**. How do I know when installation is complete?

Run following command

```bash
kubectl -n devtroncd get installers installer-devtron -o jsonpath='{.status.sync.status}'
```

Once installation process is complete, above command will print `Applied` It may take around 30 mins for installation to complete

**2**. How do I track progress of installation?

Run following command to check logs of the pod

```bash
pod=$(kubectl -n devtroncd get po -l app=inception -o jsonpath='{.items[0].metadata.name}')&& kubectl -n devtroncd logs -f $pod
```

**3**. devtron installer logs have error and I want to restart installation.

Run following command to clean up components installed by devtron installer

```bash
cd devtron-installation-script/
kubectl delete -n devtroncd -f yamls/
kubectl -n devtroncd patch installer installer-devtron --type json -p '[{"op": "remove", "path": "/status"}]'
```

In case you are still facing issues please feel free to reach out to us on [discord](https://discord.gg/jsRG5qx2gp)

