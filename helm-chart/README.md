# Training Application Helm Chart

## Prerequisites for the value `persistMetaInfo`

- The StorageClass with the name `my-storageclass` has to exist in the cluster

## Prerequisites for the value `ingress.enabled`

- The IngressClass with the name `nginx` has to exist in the cluster
- CertManager has to exist in the cluster
- The ClusterIssuer with the name `letsencrypt-issuer` has to exist in the cluster
