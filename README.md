# Ozone

ozone is the cloud native platform for hosting workshops at google developer groups.

### Demo

https://github.com/user-attachments/assets/e32529ae-b213-4e20-944a-4c1caeed2ae6

### Features
- Continous Integration with ArgoCD
- Fully self hosted Kubernetes on Microsoft Azure
- Built in Kubernetes [client-go](https://github.com/kubernetes/client-go) support for provisioning dynamic resources

### Architecture

Ozone utilizes openvscode-server as its primary way of serving instances to users. Upon a request to start a new instance, the dynamic kubernetes controller built into ozone will provision a new deployment and nodeport service. The current infrastructure depends on [infra](https://github.com/GlennTatum/infra) which utilizes haproxy as a load-balancer to forward any nodeport to the worker plane nodes on the cluster.

### Repository Layout

```
├── frontend
│   └── ui
│       ├── angular.json
│       ├── dev.sh
│       ├── dist
│       ├── Dockerfile.bootstrap
│       ├── Dockerfile.dev
│       ├── node_modules
│       ├── package.json
│       ├── package-lock.json
│       ├── public
│       ├── README.md
│       ├── src
│       ├── start.sh
│       ├── tsconfig.app.json
│       ├── tsconfig.json
│       └── tsconfig.spec.json
└── server
    ├── api
    │   ├── api.go
    │   ├── kubernetes.go
    │   ├── manifests
    │   ├── middleware.go
    │   └── service.go
    ├── argocd.yml
    ├── data
    │   └── init.cql
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── helm
    │   ├── charts
    │   ├── Chart.yaml
    │   ├── envs
    │   └── templates
    ├── main.go
    ├── models
    │   ├── account.go
    │   └── models.go
    ├── setup-secrets.sh
    └── util
        ├── kubernetes.go
        ├── kubernetes_test.go
        ├── kustomization-ingress-patch.yml
        ├── kustomization.yaml
        ├── kustomize.go
        ├── template.go
        └── template_test.go
```
