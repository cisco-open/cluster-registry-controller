# Resource Sync Rule

### Limit namespace cache
In larger cluster with namespaces more than 30-40, it was observed that CR-controller watches/caches all the namespaces into the pod memory,
which caused frequent OOMKilled issue while attaching peer clusters.

To reduce memory caching by cluster-registry-controller we can use `namespaces` field in `ResourceSyncRule` to allow caching on selected
namespaces.

custom-resource
`spec.rules.match.namespaces`
```
  namespaces:
    items:
      type: string
    type: array
```
