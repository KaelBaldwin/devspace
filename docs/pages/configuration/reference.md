---
title: Full reference
---

## version
```yaml
version: v1beta1                   # string   | Version of the config
```

<details>
<summary>
### List of supported versions
</summary>
- v1beta1   ***latest***
- v1alpha4
- v1alpha3
- v1alpha2 
- v1alpha1
</details>

---
## images
```yaml
images:                             # map[string]struct | Images to be built and pushed
  image1:                           # string   | Name of the image
    image: dscr.io/username/image   # string   | Image repository and name 
    tag: v0.0.1                     # string   | Image tag
    createPullSecret: true          # bool     | Create a pull secret containing your Docker credentials (Default: true)
    insecure: false                 # bool     | Allow push/pull to/from insecure registries (Default: false)
    skipPush: false                 # bool     | Skip pushing image to registry, recommended for minikube (Default: false)
    build: ...                      # struct   | Build options for this image
  image2: ...
```
[Learn more about building images with DevSpace.](/docs/image-building/overview)

### images[*].build
```yaml
build:                              # struct   | Build configuration for an image
  disabled: false                   # bool     | Disable image building (Default: false)
  dockerfile: ./Dockerfile          # string   | Relative path to the Dockerfile used for building (Default: ./Dockerfile)
  context: ./                       # string   | Relative path to the context used for building (Default: ./)
  kaniko: ...                       # struct   | Build image with kaniko and set options for kaniko
  docker: ...                       # struct   | Build image with docker and set options for docker
  options: ...                      # struct   | Set build options that are independent of of the build tool used
```
Notice:
- Setting `docker` or `kaniko` will define the build tool for this image.
- You **cannot** use `docker` and `kaniko` in combination. 
- If neither `docker` nor `kaniko` is specified, `docker` will be used by default.

### images[*].build.docker
```yaml
docker:                             # struct   | Options for building images with Docker
  preferMinikube: true              # bool     | If available, use minikube's in-built docker daemon instaed of local docker daemon (default: true)
```

### images[*].build.kaniko
```yaml
kaniko:                             # struct   | Options for building images with kaniko
  cache: true                       # bool     | Use caching for kaniko build process
  snapshotMode: "time"              # string   | Type of snapshotMode for kaniko build process (compresses layers)
  flags: []                         # string[] | Array of flags for kaniko build command
  namespace: ""                     # string   | Kubernetes namespace to run kaniko build pod in (Default: "" = deployment namespace)
  pullSecret: ""                    # string   | Mount this Kubernetes secret instead of creating one to authenticate to the registry (default: "")
```
> It is recommended to use Docker for building images when using DevSpace Cloud.

### images[*].build.options
```yaml
build:                              # struct   | Options for building images
  target: ""                        # string   | Target used for multi-stage builds
  network: ""                       # string   | Network mode used for building the image
  buildArgs: {}                     # map[string]string | Key-value map specifying build arguments that will be passed to the build tool (e.g. docker)
```


---
## deployments
```yaml
deployments:                        # struct[] | Array of deployments
- name: my-deployment               # string   | Name of the deployment
  namespace: ""                     # string   | Namespace to deploy to (Default: "" = namespace of the active namespace/Space)
  component: ...                    # struct   | Deploy a DevSpace component chart using helm
  helm: ...                         # struct   | Use Helm as deployment tool and set options for Helm
  kubectl: ...                      # struct   | Use "kubectl apply" as deployment tool and set options for kubectl
```
Notice:
- Setting `component`, `helm` or `kubectl` will define the type of deployment and the deployment tool to be used.
- You **cannot** use `component`, `helm` and `kubectl` in combination.

### deployments[*].component
```yaml
component:                          # struct   | Options for deploying a DevSpace component
  containers: ...                   # struct   | Relative path
  replicas: 1                       # int      | Number of replicas (Default: 1)
  autoScaling: ...                  # struct   | AutoScaling configuration
  rollingUpdate: ...                # struct   | RollingUpdate configuration
  volumes: ...                      # struct   | Component volumes
  service: ...                      # struct   | Component service
  serviceName: my-service           # string   | Service name for headless service (for StatefulSets)
  podManagementPolicy: OrderedReady # enum     | "OrderedReady" or "Parallel" (for StatefulSets)
  pullSecrets: ...                  # string[] | Array of PullSecret names
```
[Learn more about configuring component deployments.](/docs/deployment/components/what-are-components)

### deployments[*].component.containers
```yaml
containers:                         # struct   | Options for deploying a DevSpace component
- name: my-container                # string   | Container name (optional)
  image: dscr.io/username/image     # string   | Image name (optionally with registry URL)
  command:                          # string[] | ENTRYPOINT override
  - sleep
  args:                             # string[] | ARGS override
  - 99999
  env:                              # map[string]string | Kubernetes env definition for containers
  - name: MY_ENV_VAR
    value: "my-value"
  volumeMounts: ...                 # struct   | VolumeMount Configuration
  resources: ...                    # struct   | Kubernestes resource limits and requests
  livenessProbe: ...                # struct   | Kubernestes livenessProbe
  redinessProbe: ...                # struct   | Kubernestes redinessProbe
```

### deployments[\*].component.containers[*].volumeMounts
```yaml
volumeMounts: 
  containerPath: /my/path           # string   | Mount path within the container
  volume:                           # struct   | Volume to mount
    name: my-volume                 # string   | Name of the volume to be mounted
    subPath: /in/my/volume          # string   | Path inside to volume to be mounted to the containerPath
    readOnly: false                 # bool     | Mount volume as read-only (Default: false)
```

### deployments[*].component.autoScaling
```yaml
autoScaling: 	                    # struct   | Auto-Scaling configuration
  horizontal:                       # struct   | Configuration for horizontal auto-scaling
    maxReplicas: 5                  # int      | Max replicas to deploy
    averageCPU: 800m                # string   | Target value for CPU usage
    averageMemory: 1Gi              # string   | Target value for memory (RAM) usage
```

### deployments[*].component.rollingUpdate
```yaml
rollingUpdate: 	                    # struct   | Rolling-Update configuration
  enabled: false                    # bool     | Enable/Disable rolling update (Default: disabled)
  maxSurge: "25%"                   # string   | Max number of pods to be created above the pod replica limit
  maxUnavailable: "50%"             # string   | Max number of pods unavailable during update process
  partition: 1                      # int      | For partitioned updates of StatefulSets
```

### deployments[*].component.volumes
```yaml
volumes: 	                        # struct   | Array of volumes to be created
- name: my-volume                   # string   | Volume name
  size: 10Gi                        # string   | Size of the volume in Gi (Gigabytes)
  configMap: ...                    # struct   | Kubernetes ConfigMapVolumeSource
  secret: ...                       # struct   | Kubernetes SecretVolumeSource
```

### deployments[*].component.service
```yaml
service: 	                        # struct   | Component service configuration
  name: my-service                  # string   | Name of the service
  type: NodePort                    # string   | Type of the service (default: NodePort)
  ports:                            # array    | Array of service ports
  - port: 80                        # int      | Port exposed by the service
    containerPort: 3000             # int      | Port of the container/pod to redirect traffic to
    protocol: tcp                   # string   | Traffic protocol (tcp, udp)
```

### deployments[*].helm
```yaml
helm:                               # struct   | Options for deploying with Helm
  chart: ...                        # struct   | Relative path 
  wait: true                        # bool     | Wait for pods to start after deployment (Default: true)
  rollback: true                    # bool     | Rollback if deployment failed (Default: true)
  force: false                      # bool     | Force deleting and re-creating Kubernetes resources during deployment (Default: false)
  timeout: 40                       # int      | Timeout to wait for pods to start after deployment (Default: 40)
  tillerNamespace: ""               # string   | Kubernetes namespace to run Tiller in (Default: "" = same a deployment namespace)
  devSpaceValues: true              # bool     | If DevSpace CLI should append pullSecrets and set images to values.yaml before deployment (Default: true)
  valuesFiles:                      # string[] | Array of paths to values files
  - ./chart/my-values.yaml          # string   | Path to a file to override values.yaml with
  values: {}                        # struct   | Any object with Helm values to override values.yaml during deployment
```
[Learn more about configuring deployments with Helm.](/docs/deployment/helm-charts/what-are-helm-charts)

### deployments[*].helm.chart
```yaml
chart:                              # struct   | Chart to deploy
  name: my-chart                    # string   | Chart name
  version: v1.0.1                   # string   | Chart version
  repo: "https://my-repo.tld/"      # string   | Helm chart repository
  username: "my-username"           # string   | Username for Helm chart repository
  password: "my-password"           # string   | Password for Helm chart repository
```

### deployments[*].kubectl
```yaml
kubectl:                            # struct   | Options for deploying with "kubectl apply"
  cmdPath: ""                       # string   | Path to the kubectl binary (Default: "" = detect automatically)
  manifests: []                     # string[] | Array containing glob patterns for the Kubernetes manifests to deploy using "kubectl apply" (e.g. kube/* or manifests/service.yaml)
  kustomize: false                  # bool     | Use kustomize when deploying manifests via "kubectl apply" (Default: false)
  flags: []                         # string[] | Array of flags for the "kubectl apply" command
```
[Learn more about configuring deployments with Helm.](/docs/deployment/kubernetes-manifests/what-are-manifests)


---
## dev
```yaml
dev:                                # struct   | Options for "devspace dev"
  overrideImages: []                # struct[] | Array of override settings for image building
  terminal: ...                     # struct   | Options for the terminal proxy
  ports: []                         # struct[] | Array of port-forwarding settings for selected pods
  sync: []                          # struct[] | Array of file sync settings for selected pods
  autoReload: ...                   # struct   | Options for auto-reloading (i.e. re-deploying deployments and re-building images)
  selectors: []                     # struct[] | Array of selectors used to select Kubernetes pods (used within terminal, ports and sync)
```
[Learn more about development with DevSpace.](/docs/development/workflow)

### dev.overrideImages
```yaml
overrideImages:                     # struct[] | Array of override settings for image building
- name: default                     # string   | Name of the image to apply this override rule to
  entrypoint: []                    # string[] | Array defining with the entrypoint that should be used instead of the entrypoint defined in the Dockerfile
  dockerfile: default               # string   | Relative path of the Dockerfile that should be used instead of the one originally defined
  context: default                  # string   | Relative path of the context directory that should be used instead of the one originally defined
```
[Learn more about image overriding.](/docs/development/overrides)

### dev.terminal
```yaml
terminal:                           # struct   | Options for the terminal proxy
  disabled: false                   # bool     | Disable terminal proxy / only start port-forwarding and code sync if defined (Default: false)
  selector:                         # TODO
  command: []                       # string[] | Array defining the shell command to start the terminal with (Default: ["sh", "-c", "command -v bash >/dev/null 2>&1 && exec bash || exec sh"])
```
[Learn more about configuring the terminal proxy.](/docs/development/terminal)

### dev.ports
```yaml
ports:                              # struct[] | Array of port forwarding settings for selected pods
- selector:                         # TODO
  forward:                          # struct[] | Array of ports to be forwarded
  - port: 8080                      # int      | Forward this port on your local computer
    remotePort: 3000                # int      | Forward traffic to this port exposed by the pod selected by "selector" (TODO)
    bindAddress: ""                 # string   | Address used for binding / use 0.0.0.0 to bind on all interfaces (Default: "localhost" = 127.0.0.1)
```
[Learn more about port forwarding.](/docs/development/port-forwarding)

### dev.sync
```yaml
sync:                               # struct[] | Array of file sync settings for selected pods
- selector:                         # TODO
  localSubPath: ./                  # string   | Relative path to a local folder that should be synchronized (Default: "./" = entire project)
  containerPath: /app               # string   | Absolute path in the container that should be synchronized with localSubPath
  waitInitialSync: false            # bool     | Wait until initial sync is completed before continuing (Default: false)
  excludePaths: []                  # string[] | Paths to exclude files/folders from sync in .gitignore syntax
  downloadExcludePaths: []          # string[] | Paths to exclude files/folders from download in .gitignore syntax
  uploadExcludePaths: []            # string[] | Paths to exclude files/folders from upload in .gitignore syntax
  bandwidthLimits:                  # struct   | Bandwidth limits for the synchronization algorithm
    download: 0                     # int64    | Max file download speed in kilobytes / second (e.g. 100 means 100 KB/s)
    upload: 0                       # int64    | Max file upload speed in kilobytes / second (e.g. 100 means 100 KB/s)
```
[Learn more about confguring the code synchronization.](/docs/development/synchronization)


### dev.autoReload
```yaml
autoReload:                         # struct   | Options for auto-reloading (i.e. re-deploying deployments and re-building images)
  paths: []                         # string[] | Array containing glob patterns of files that are watched for auto-reloading (i.e. reload when a file matching any of the patterns changes)
  deployments: []                   # string[] | Array containing names of deployments to watch for auto-reloading (i.e. reload when kubectl manifests or files within the Helm chart change)
  images: []                        # string[] | Array containing names of images to watch for auto-reloading (i.e. reload when the Dockerfile changes)
```

### dev.selectors
```yaml
selectors:                          # struct[] | Array of selectors used to select Kubernetes pods (used within terminal, ports and sync)
- name: default                     # string   | Name of this pod selector (used to reference this selector within terminal, ports and sync)
  namespace: ""                     # string   | Namespace to select pods in (Default: "" = namespace of the active Space)
  labelSelector: {}                 # map[string]string | Key-value map of Kubernetes labels used to select pods
  ContainerName: ""                 # string   | Name of the container within the selected pod (Default: "" = first container in the pod)
```

---
## cluster
> **Warning:** Change the cluster configuration only if you *really* know what you are doing. Editing this configuration can lead to issues with when running DevSpace CLI commands.

```yaml
cluster:                            # struct   | Cluster configuration
  kubeContext: ""                   # string   | Name of the Kubernetes context to use (Default: "" = current Kubernetes context used by kubectl)
  namespace: ""                     # string   | Namespace for deploying applications
  apiServer: ""                     # string   | URL of your Kubernetes API server (master)
  caCert: ""                        # string   | CA Certificate of your Kubernetes API server
  user:                             # struct   | Options for user authentication
    clientCert: ""                  # string   | Use certificate-based authentication using this client certificate
    clientKey: ""                   # string   | Use certificate-based authentication using this client key
    token: ""                       # string   | Use token-based authentication using this token
```
Notice:
- You **cannot** use `clientCert` and `clientKey` in combination with `token`.

> If you want to work with self-managed Kubernetes clusters, it is highly recommended to connect an external cluster to DevSpace Cloud or run your own instance of DevSpace Cloud (coming soon) instead of using the following configuration options.
