# Cluster registry controller

The cluster registry controller helps to form a group of Kubernetes clusters and synchronize
any k8s resources across those clusters arbitrarily.

## Overview

This repo contains the Kubernetes controller implementation for the
[Cluster Registry API](https://github.com/banzaicloud/cluster-registry).

The `Cluster` resource represents a Kubernetes cluster.
The cluster registry controller fills the status of the `Cluster` CR with cluster related metadata.

The controller is mostly useful in multi-cluster scenarios.
In this scenario, the cluster registry controller is deployed to all Kubernetes clusters.
The same `Cluster` CRs should be placed on all participating Kubernetes clusters as well.
Also, the credentials for all clusters should be distributed to all clusters (these are usually stored in k8s secrets).

> The cluster registry controller syncs the `Cluster` CR and related secret resources across clusters
> to help bootstrap the cluster group itself.

You can define your own `ResourceSyncRule` resources to sync k8s resources between these clusters.

In such a multi-cluster setup, here is how the cluster registry controller works:
- The controller only writes to the local cluster where it is deployed to
- The controller only reads from peer clusters

By default, the required resources are kept in sync between all clusters.
It can be further adjusted, from which clusters and to which clusters certain resource should be synced.

## Quickstart

### Form cluster group with two clusters

1. Install cluster registry controller on the first cluster. The following command installs the cluster registry controller on your cluster, creates a Cluster CR with the name `FIRST-CLUSTER-NAME`, and it also creates a secret that holds a Kubeconfig with read access to this cluster.

    ```
    helm install --namespace=cluster-registry --create-namespace cluster-registry-controller deploy/charts/cluster-registry --set localCluster.name=<FIRST-CLUSTER-NAME>
    ```

    > Tip: Use the `--set apiServerEndpointAddress=<PUBLIC-API-SERVER-ADDRESS>` flag, if your Kubernetes cluster API returns private ip for the api server.

2. Install cluster registry controller on the second cluster:

    ```
    helm install --namespace=cluster-registry --create-namespace cluster-registry-controller deploy/charts/cluster-registry --set localCluster.name=<SECOND-CLUSTER-NAME>
    ```

   This command will install the cluster registry controller on your cluster, it will also create a Cluster CR with the name `SECOND-CLUSTER-NAME` and it will create a secret as well that holds a Kubeconfig with read access to this cluster.

3. Copy/paste Cluster and secret resources from first->second and second->first cluster:
    
    Secret is needed so that the cluster registry controller from one cluster can read from the other cluster.

    From first cluster to second cluster:
    ```
    kubectl get cluster <FIRST-CLUSTER-NAME> -o yaml | pbcopy      pbpaste | kubectl apply -f -
    kubectl get secret <FIRST-CLUSTER-NAME> -o yaml | pbcopy       pbpaste | kubectl apply -f -
    ```

    From second cluster to first cluster:
    ```
    kubectl get cluster <SECOND-CLUSTER-NAME> -o yaml | pbcopy    pbpaste | kubectl apply -f -
    kubectl get secret <SECOND-CLUSTER-NAME> -o yaml | pbcopy     pbpaste | kubectl apply -f -
    ```

4. Check the status of the Cluster CRs:

    Both Cluster CRs should show `Synced` state.
    On the first cluster, the <FIRST-CLUSTER-NAME> Cluster CR should be type local in the status, the
    <SECOND-CLUSTER-NAME> should be peer.
    On the second cluster, the <SECOND-CLUSTER-NAME> Cluster CR should be type local in the status, the
    <FIRST-CLUSTER-NAME> should be peer.
    
    > The type in the Cluster status is determined by the clusterID field in the Cluster spec and by the 
    > `kube-system` namespace uid. If they match, the cluster is local, otherwise it is a peer cluster.

The cluster group is successfully formed at this point.

### Attach additional clusters to the group

1. Install the cluster registry controller on the new cluster as shown above.

2. Choose one cluster from the existing cluster group and the the Cluster and secret resource swap with the new cluster:

   From first cluster to third cluster:
    ```
    kubectl get cluster <FIRST-CLUSTER-NAME> -o yaml | pbcopy      pbpaste | kubectl apply -f -
    kubectl get secret <FIRST-CLUSTER-NAME> -o yaml | pbcopy       pbpaste | kubectl apply -f -
    ```

   From third cluster to first cluster:
    ```
    kubectl get cluster <THIRD-CLUSTER-NAME> -o yaml | pbcopy    pbpaste | kubectl apply -f -
    kubectl get secret <THIRD-CLUSTER-NAME> -o yaml | pbcopy     pbpaste | kubectl apply -f -
    ```
   
    Important to note, that **you don't need to do the Cluster and secret resource swap between the new cluster and 
    any other cluster in the cluster group.**
    The cluster registry controller automatically syncs these Cluster and secret resources between the clusters in the 
    cluster group and once a new cluster is added to the cluster group in any of the clusters, then it will be synced
    between the other clusters automatically.

    So, if you have 10 clusters in a cluster group, you'll still only need to do the Cluster and secret resource swap once
    with one cluster from the cluster group, the rest should synchronize automatically.

5. Check the status of the Cluster CRs:

   All Cluster CRs should show `Synced` state.
   If so, then the cluster group is successfully expanded.

### ResourceSyncRule example usage

#### Sync everywhere

1. Create a sample secret on the third cluster, which will be copied around:

    ````yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: test-secret
    data: {}
    ````

2. Create a `ResourceSyncRule` on the first cluster to synchronize the secret to all clusters:

    ```yaml
    apiVersion: clusterregistry.k8s.cisco.com/v1alpha1
    kind: ResourceSyncRule
    metadata:
      name: test-secret-sink
    spec:
      groupVersionKind:
        kind: Secret
        version: v1
      rules:
      - match:
        - objectKey:
          name: test-secret
          namespace: cluster-registry
    ```
   
    This `ResourceSyncRule` resource itself and the `secret` resource as well should appear shortly on all 
    clusters of the cluster group.
    
    At this point, if a secret from any of the clusters (except from the third one when it originates from) is deleted
    or modified, it will be synced back immediately by the cluster registry controller.

#### Sync to a set of clusters

Cluster registry controller can be configured, so that it only sync to some specific clusters in the cluster group (instead of all
of them).

1. Add the following annotation to the `ResourceSyncRule` on the first cluster:

    ```yaml
    annotations:
      cluster-registry.k8s.cisco.com/resource-sync-disabled: "true"
    ```

2. Delete the `ResourceSyncRule` from the second cluster.

    The `ResourceSyncRule` resource will not be recreated because of the annotation, which was just added.

    > If the annotation is not added as described in the previous step, then the `ResourceSyncRule` will be recreated.

3. Delete the `test-secret` from the second cluster.

    The secret will not be recreated because the `ResourceSyncRule` resource does not exist on the second cluster.

#### Sync from a set of clusters

Cluster registry controller can be configured, so that it only sync from some specific clusters in the cluster group (instead of all
of them).

1. Add the following field to the `ResourceSyncRule` spec on the first cluster:

    ```yaml
    clusterFeatureMatch:
    - featureName: test-secret-feature
    ```
   
    This causes that the secret will only be synced from clusters where there are `ClusterFeature` resources defined.

    > At this point, there is no `ClusterFeature` present on any cluster, so if the secret would be deleted now from 
    > the first cluster, it would not be recreated.

2. Apply the following `ClusterFeature` to the third cluster:

    ```yaml
    apiVersion: clusterregistry.k8s.cisco.com/v1alpha1
    kind: ClusterFeature
    metadata:
      name: test-secret-source
    spec:
      featureName: test-secret-feature
    ```

3. Delete the `test-secret` from the first cluster.

    It should be recreated now, because it can sync the secret from the third cluster.

## RBAC considerations

The cluster registry controller only writes to local clusters and only reads from peer clusters.
By default, it has access to read `namespace`, `node` and `secret` resources.
The quickstart example worked, because the controller was allowed to read the secret from the remote cluster.

If other resources should be synced, then the RBAC rules of the operator should be expanded.

- On the cluster, where the resources are read from (usually where `ClusterFeature` resources are present)
  a ClusterRole should be defined with the correct read roles and the following label should be added:
  
  ```
  labels:
    cluster-registry.k8s.cisco.com/reader-aggregated: "true"
  ```
  
- On the cluster, where the resources are written to (Usually where `ResourceSyncRule` resources are present)
  a ClusterRole should be defined with the correct write roles and the following label should be added:

  ```
  labels:
    cluster-registry.k8s.cisco.com/controller-aggregated: "true"
  ```

## Contributing

If you find this project useful, help us:

- Support the development of this project and star this repo! :star:
- If you use Cluster registry controller, add yourself to the list of [adopters](ADOPTERS.md).:metal: <br>
- Help new users with issues they may encounter :muscle:
- Send a pull request with your new features and bug fixes :rocket: Check out the [developer docs](docs/development.md) for that.
