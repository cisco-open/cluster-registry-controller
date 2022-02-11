# Developer Guide

## How to run cluster registry controller in your cluster with your changes

Cluster registry controller is built on the [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) project.

If you make changes and would like to try your own version, create your own image:

    make docker-build docker-push deploy IMG=<YOUR-REGISTRY>/cluster-registry-controller:<VERSION>

Watch the operator's logs with:

    kubectl logs -f -n cluster-registry cluster-registry-controller-manager -c manager
