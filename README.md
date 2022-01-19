# KubeTempura

Set up the Review App for each Pull Request automatically.

# DESCRIPTION

KubeTempura is the Kubernetes operator for setting up a Review App in a Kubernetes cluster with flexibility.  

# INSTALLATION

## Prerequisites
- A domain name (e.g, kubetempura.example.com)
   - To receive a event from GitHub
- Git
- Make
- Kubectl
- Kustomize

## 1. Setup a new GitHub Webhook
Setup a new Webhook like [this image](https://user-images.githubusercontent.com/174613/149904429-d8d0295c-f6ea-4937-8249-d344920ff842.png).

1. Open the Settings tab on your repository or organization (if you want to setup it organization-wide).
2. Open the Webhooks tab
3. Click the "Add webhook" button
4. Set the "Payload URL" with the value `https://YOURDOMAIN/webhooks`
5. Set the "Content type" with the value `application/json`
6. Generate a secret plain-text token and input it.
7. Choose "Let me select individual events" for the webhook trigger.
8. Enable the event "Pull requests".
9. Finish creating the webhook.


Then, you also need to register that secret as a Kubernetes Secret resource. 

```bash
$ kubectl create ns kubetempura-system
$ kubectl create secret generic -n kubetempura-system github-webhook --from-literal=secret=$YOUR_SECRET
```

## 2. Install CRDs and the controller to your cluster

```bash
$ git clone git@github.com:mercari/kubetempura.git
$ cd kubetempura

$ vi config/crd/default/manager_auth_proxy_patch.yaml
$ make install deploy
```

### 3. Setup the endpoint for the GitHub Webhook
It heavily depends on your cluster. Typically you need to create a Service and a Ingress resources. See below example files:
- [Service](./config/samples/github_service.yaml)
- [Ingress](./config/samples/github_ingress.yaml)

Also you need to register a DNS record. Such as `kubetempura.example.com`

# SYNOPSIS

After the installation, use can create a ReviewApp resource in each namespace. ReviewApp is a template for resources, which you want to create for each PR.

the `resouces` is an array which allows any kind of Kubernetes resources include CRDs.

```yaml
apiVersion: kubetempura.mercari.com/v1
kind: ReviewApp
metadata:
  name: reviewapp-sample
  namespace: default
spec:
  githubRepository: https://github.com/mercari/not-exists
  resources:
    - apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: reviewapp-sample-pr{{PR_NUMBER}}
        namespace: default
      spec:
        selector:
          matchLabels:
            app: reviewapp-sample-pr{{PR_NUMBER}}
        template:
          metadata:
            labels:
              app: reviewapp-sample-pr{{PR_NUMBER}}
          spec:
            containers:
              - name: sample
                image: ghcr.io/stefanprodan/podinfo:pr{{PR_NUMBER}}-{{COMMIT_REF}}
                ...
    - apiVersion: v1
      kind: Service
      metadata:
        name: reviewapp-sample-pr{{PR_NUMBER}}
        namespace: default
      spec:
        ...
        selector:
          app: reviewapp-sample-pr{{PR_NUMBER}}
```

In a template, you can use several variables:
- `{{PR_NUMBER}}`: the number of a PR.
- `{{COMMIT_REF}}`: the commit-ref of a latest (head) commit of a PR. It would be useful to specifying the image tag.
- `{{COMMIT_REF_SHORT}}`: the short version of the commit ref for a compatibility.

## Limitations
- KubeTempura has a limited permission for create/update/delete a resource. If you want to create a resource without one of a kind `Deployment`, `Service`, `ConfigMap`, and `Secrets`, you need to add that resouce in a ClusterRole for KubeTempura.
- KubeTempura works only based on a GitHub Webhook. You need to close and re-open your PR to update a state explicitly when KubeTempura failed to receive a webhook for some reasons.

# CONTRIBUTION

Please read the CLA carefully before submitting your contribution to Mercari. Under any circumstances, by submitting your contribution, you are deemed to accept and agree to be bound by the terms and conditions of the CLA.

https://www.mercari.com/cla/

# LICENSE

Copyright 2022 Mercari, Inc.

Licensed under the Apache License.
