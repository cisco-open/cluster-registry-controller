# Cluster-registry chart

The [cluster registry controller](https://github.com/cisco-open/cluster-registry-controller) helps to form a group of Kubernetes clusters and synchronize
any K8s resources across those clusters arbitrarily.

## Prerequisites

- Helm3

## Installing the chart

To install the chart:

```bash
❯ helm repo add cluster-registry https://cisco-open.github.io/cluster-registry-controller
❯ helm install --namespace=cluster-registry --create-namespace cluster-registry cluster-registry/cluster-registry --set localCluster.name=primary
```

## Uninstalling the Chart

To uninstall/delete the `cluster-registry` release:

```bash
❯ helm uninstall cluster-registry
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Cluster Registry chart and their default values.

Parameter | Description | Default
--------- |--| -------
`replicas` | Operator deployment replica count | `1`
`localCluster.name` | Specify to automatically provision the cluster object with this name upon first start | `""`
`localCluster.manageSecret` | If true, automatically provisions the secret for the K8s API server access | `"true"`
`istio.revision` | Sets the `istio.io/rev` label on the deployment | `""`
`podAnnotations` | Operator deployment pod annotations (YAML) | `{}`
`imagePullSecrets` | Operator deployment image pull secrets | `[]`
`podSecurityContext` | Operator deployment pod security context (YAML) | runAsUser: `65534`, runAsGroup: `65534`
`securityContext` | Operator deployment security context (YAML) | allowPrivilegeEscalation: `false`
`image.repository` | Operator container image repository | `ghcr.io/cisco-open/cluster-registry-controller`
`image.tag` | Operator container image tag | `v0.2.11`
`image.pullPolicy` | Operator container image pull policy | `IfNotPresent`
`nodeselector` | Operator deployment node selector (YAML) | `{}`
`affinity` | Operator deployment affinity (YAML) | `{}`
`tolerations` | Operator deployment tolerations | `[]`
`resources` | CPU/Memory resource requests/limits (YAML) | Requests: Memory: `100Mi`, CPU: `100m`, Limits: Memory: `200Mi`, CPU: `300m`
`service.type` | Operator service type | `"ClusterIP"`
`service.port` | Operator service port | `8080`
`serviceAccount.annotations` | Operator service account annotations (YAML) | `{}`
`podDisruptionBudget.enabled` | If true, PodDisruptionBudget is deployed for the operator | `false`
`controller.leaderElection.enabled` | If true, leader election is enabled for the operator deployment | `true`
`controller.leaderElection.name` | Name override for the leader election configmap | `cluster-registry-leader-election`
`controller.log.format` | Format for controller logs | `"json"`
`controller.log.verbosity` | Verbosity for controller logs | `0`
`controller.workers` | Number of workers | `2`
`controller.apiServerEndpointAddress` | Address of the cluster's K8s API server, which is publicly or from the specified network | `""`
`controller.network.name` | Name of the network where the cluster is reachable | `"default"`
`controller.coreResourceSource.enabled` | If true, the core resources (Cluster, ResourceSyncRule, secrets) could be synced from this cluster | `true`
