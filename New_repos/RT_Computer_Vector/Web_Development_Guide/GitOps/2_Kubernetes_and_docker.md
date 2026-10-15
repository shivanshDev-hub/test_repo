# Kubernetes and docker (for dummies)

## Docker Tutorial 


### Creating first machine 

how does a basic docker command looks like -> `docker run -it --rm ubuntu`

`this will take you inside the very terminalof the machine you just made`

use exit to get out of the ubuntu container now 

in linux installing docker is pretty easy ->
```bash
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER   # so you don't need sudo every time
```

again i know this is **Not a proper way to install the docker** but remember this is docker for dummies 

*to look what docker images does the Docker offer [click here](https://hub.docker.com/)*

### basic commands of docker 

1. docker pull <image:version> --> pulls a image
2. docker images --> show images 
3. docker run <image> --> creates a docker conatiner 
4. docker run -it <image> --> docker images in the interactive mode 
5. docker ps -a --> show all existing container 
6. docker ps --> show all running docker containers
7. docker stop <container id/name> --> to start a container 
8. docker start <container id/name> --> to stop a conatiner (for non interactive ones)
9. docker rmi <image name> --> to remove a image 
10. docker rm <container name> --> to remove a container 
11. docker run -d <image> --> to run a image in detach mode 
12. docker run --name <container name> -d <image name> --> 
13. docker run -d -e variable=value <image> --> variable addition 
14. --name --> to give a name to a container 


```
to remove a docker container : check processes -> stop it -> remove it 
```

### Docker Image Layers 

```
container layer -> layer_2 -> Layer_1 -> Base layer
```
docker images is made up of layer like this 

**port Binding** `docker run -p8080:3306 <image>` --> port binding port 8080 of main machine is binded to the port 3306 of the docker machine here  

### Troubleshoot commands 

1. docker logs <container> --> checks logs of a container 
2. docker exec -it <conatiner> /bin/bash or /bin/sh --> execute a commands

Docker vs Virtual machine => **docker only virtualise the application layer vs os + applcation virtualization in vm**

### Docker for developers 

**concept of the docker network**

- docker network ls --> to check all docker networks 
- docker create <network-name> --> create a new network
- --network flag when making conatiner to specify their network

> wanna create databases read the doc, config the varaibles and add the tables using the respective commands , remember to keep the port same and the tables names same 

> how do you place a file inside the docker conatiner and make it run ? : 
```bash
docker cp myscript.sh <container_id_or_name>:/root/myscript.sh
docker exec -it <container_id> bash
bash /root/myscript.sh
```

`IMPORTANT NOTE: EACH DOCKER IMAGE IN DOCKER HUB HAVE A DETAILED EXPLANANTION OF HOW TO SET IT UP READ THAT ALSO SEARCH AI FOR THAT TOO`

**Entry of the Docker compose**

~ Docker Compose is a toll for defining and running multi-container applications

```yaml
version: "3.9"

services:
  web:
    build: .
    ports:
      - "5000:5000"
    volumes:
      - .:/app
    environment:
      - DEBUG=true
    depends_on:
      - db

  db:
    image: postgres:16
    environment:
      - POSTGRES_PASSWORD=mysecret
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
```

#### Top-Level Keys
 
| Field | Description |
|---|---|
| `version` | Compose file format version (e.g. `"3.9"`) — often omitted in newer Compose |
| `services` | Defines each container/app in your project |
| `volumes` | Declares named volumes for persistent data |
| `networks` | Declares custom networks for services to communicate over |
| `configs` | Defines external configuration files to mount into services |
| `secrets` | Defines sensitive data (passwords, keys) to mount securely |
 
 
#### Service-Level Fields (under `services: <name>:`)
 
**Build / Image**
 
| Field | Description |
|---|---|
| `build` | Build an image from a Dockerfile (path or object with `context`/`dockerfile`) |
| `build.context` | Directory containing the Dockerfile and build files |
| `build.dockerfile` | Custom Dockerfile name/path |
| `build.args` | Build-time variables passed to the Dockerfile |
| `image` | Use a pre-built image instead of (or in addition to) building |
 
**Networking**
 
| Field | Description |
|---|---|
| `ports` | Maps host ports to container ports (`"host:container"`) |
| `expose` | Exposes ports to other services only, not the host |
| `networks` | Attaches the service to specific custom networks |
| `hostname` | Sets the container's hostname |
| `dns` | Custom DNS servers for the container |
 
**Storage**
 
| Field | Description |
|---|---|
| `volumes` | Mounts host paths or named volumes into the container |
| `tmpfs` | Mounts a temporary in-memory filesystem |
 
**Environment / Config**
 
| Field | Description |
|---|---|
| `environment` | Sets environment variables (list or key-value map) |
| `env_file` | Loads environment variables from a `.env`-style file |
| `command` | Overrides the default command (CMD) from the image/Dockerfile |
| `entrypoint` | Overrides the default ENTRYPOINT from the image |
| `working_dir` | Sets the working directory inside the container |
| `user` | Runs the container process as a specific user |
 
**Dependencies / Lifecycle**
 
| Field | Description |
|---|---|
| `depends_on` | Controls startup order between services |
| `restart` | Restart policy: `no`, `always`, `on-failure`, `unless-stopped` |
| `healthcheck` | Defines a command to check if the container is healthy |
| `stop_signal` | Custom signal used to stop the container |
| `stop_grace_period` | How long to wait before force-killing the container |
 
**Resource Limits**
 
| Field | Description |
|---|---|
| `deploy.resources.limits.cpus` | Max CPU usage |
| `deploy.resources.limits.memory` | Max memory usage |
| `mem_limit` | Memory limit (legacy/simpler alternative to `deploy`) |
| `cpus` | CPU limit (legacy/simpler alternative to `deploy`) |
 
**Metadata / Misc**
 
| Field | Description |
|---|---|
| `container_name` | Sets a custom, fixed name for the container |
| `labels` | Attaches metadata labels to the container |
| `profiles` | Marks a service to only start with specific `--profile` flags |
| `platform` | Specifies target platform (e.g. `linux/amd64`) |
 
#### Top-Level `volumes:` Sub-Fields
 
| Field | Description |
|---|---|
| `driver` | Volume driver to use (default: `local`) |
| `driver_opts` | Options passed to the volume driver |
| `external` | Marks the volume as pre-existing (not managed by Compose) |
 
#### Top-Level `networks:` Sub-Fields
 
| Field | Description |
|---|---|
| `driver` | Network driver (`bridge`, `overlay`, `host`, etc.) |
| `external` | Marks the network as pre-existing (not managed by Compose) |
 
#### Common CLI Commands
 
```bash
docker compose up          # build + start all services
docker compose up -d       # start in background (detached)
docker compose down        # stop + remove containers, networks
docker compose build       # rebuild images
docker compose logs -f     # follow logs
docker compose ps          # list running services
docker compose exec <service> <cmd>   # run a command inside a running service
```
`to run these commands you must be in the project folder`

btw never told you the command to delete a network is `docker network rm <network-name>`

time for the entry of dockerfile -> what is the difference between docker file and docker.yaml`<Right>You typically use both together: the Dockerfile builds your app's image, and Compose starts that image alongside a database, cache, etc., all networked automatically.`

### Writing a Dockerfile

**What it is:** a plain text file named exactly `Dockerfile` (no extension) containing step-by-step instructions to build an image.

---

#### Basic Structure

```dockerfile
FROM node:20-slim
WORKDIR /app
COPY package.json .
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
```

---

#### Core Instructions

#### `FROM`
**Sets the base image** to build on top of. Always the first line.
```dockerfile
FROM python:3.12-slim
```

#### `WORKDIR`
**Sets the working directory** inside the container for all following instructions.
```dockerfile
WORKDIR /app
```

#### `COPY`
**Copies files** from your host machine into the image.
```dockerfile
COPY . .
COPY requirements.txt .
```

#### `ADD`
**Like `COPY`, but also** supports extracting tar files and fetching URLs. Use `COPY` unless you specifically need those extras.

#### `RUN`
**Executes a command** at build time — used for installing dependencies.
```dockerfile
RUN pip install -r requirements.txt
RUN apt-get update && apt-get install -y curl
```

#### `ENV`
**Sets environment variables** available inside the container.
```dockerfile
ENV NODE_ENV=production
```

#### `EXPOSE`
**Documents which port** the app listens on. Doesn't actually publish it — that's done with `docker run -p`.
```dockerfile
EXPOSE 8080
```

#### `CMD`
**Default command** the container runs on startup. Only one `CMD` per Dockerfile (last one wins if repeated).
```dockerfile
CMD ["python", "app.py"]
```

#### `ENTRYPOINT`
**Like `CMD`, but harder to override** — often paired with `CMD` for default arguments.
```dockerfile
ENTRYPOINT ["python"]
CMD ["app.py"]
```

#### `ARG`
**Build-time variable**, only available during the build, not at runtime.
```dockerfile
ARG VERSION=1.0
```

#### `USER`
**Sets which user** runs the container process (avoid running as root).
```dockerfile
USER node
```

#### `VOLUME`
**Marks a path** as a mount point for persistent data.
```dockerfile
VOLUME /data
```

---

### Best Practices

**Order matters for caching** — put things that change least often (like `RUN npm install`) before things that change often (like `COPY . .`). This way Docker reuses cached layers and rebuilds are faster.

```dockerfile
COPY package.json .
RUN npm install
COPY . .
```
*(installing dependencies before copying all code — if only your code changes, `npm install` won't rerun)*

**Use slim/alpine base images** to keep image size small.
```dockerfile
FROM node:20-alpine
```

**Combine RUN commands** to reduce image layers.
```dockerfile
RUN apt-get update && apt-get install -y \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*
```

**Don't run as root** in production — add a `USER` instruction.

**Use a `.dockerignore` file** to exclude unnecessary files (like `node_modules`, `.git`) from the build context.
```
node_modules
.git
*.log
```

---

#### Building and Running

```bash
docker build -t myapp:latest .
docker run -p 3000:3000 myapp:latest
```

### Commonly Used Production Base Images

---

#### Distro-based (general purpose)

**`ubuntu`** — full-featured, familiar, good compatibility. Large (~70MB+).
```dockerfile
FROM ubuntu:24.04
```

**`debian`** — similar to Ubuntu, slightly smaller, very stable.
```dockerfile
FROM debian:bookworm-slim
```

---

#### Minimal / Small

**`alpine`** — extremely small (~5MB), uses `musl` instead of `glibc`. Fast pulls, smaller attack surface, but can cause compatibility issues with some native binaries.
```dockerfile
FROM alpine:3.20
```

**`distroless`** (Google) — contains *only* your app + runtime, no shell, no package manager. Very secure, but harder to debug.
```dockerfile
FROM gcr.io/distroless/static-debian12
```

---

#### Language-Specific

**Node.js**
```dockerfile
FROM node:20-slim      # good balance of size + compatibility
FROM node:20-alpine    # smaller, occasional native-module issues
```

**Python**
```dockerfile
FROM python:3.12-slim
FROM python:3.12-alpine
```

**Java**
```dockerfile
FROM eclipse-temurin:21-jre-jammy   # official OpenJDK replacement image
```

**Go** (often used as build stage only, then copied into a tiny final image)
```dockerfile
FROM golang:1.23 AS builder
...
FROM gcr.io/distroless/base-debian12   # final minimal runtime image
```

**.NET**
```dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:8.0
```

**Ruby**
```dockerfile
FROM ruby:3.3-slim
```

---

#### Web Servers

**Nginx** — commonly used as the final stage for serving static frontend builds (React, Vue, etc.)
```dockerfile
FROM nginx:1.27-alpine
```

**Caddy**
```dockerfile
FROM caddy:2-alpine
```

---

#### Databases (usually run as separate services, not built on top of)

```dockerfile
FROM postgres:16
FROM mysql:8.4
FROM redis:7-alpine
FROM mongo:7
```

---

### Quick Guidance

| Priority | Best choice |
|---|---|
| Smallest image size | `alpine` or `distroless` |
| Best compatibility, easiest debugging | `slim` variants (e.g. `python:3.12-slim`) |
| Maximum security (no shell at all) | `distroless` |
| Familiar full environment | `ubuntu` or `debian` |

**Common production pattern — multi-stage build:** use a full image to build your app, then copy just the compiled output into a tiny final image.
```dockerfile
FROM node:20 AS builder
WORKDIR /app
COPY . .
RUN npm install && npm run build

FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
```

This keeps your final shipped image small even if the build process needs a heavier toolchain.

### Publishing Your Own Docker Image

---

#### Step 1: Create a Docker Hub Account
Sign up at https://hub.docker.com if you don't already have an account.

---

#### Step 2: Log In from the CLI
```bash
docker login
```
Enter your Docker Hub username and password (or access token) when prompted.

---

#### Step 3: Build Your Image
```bash
docker build -t myapp:latest .
```

---

#### Step 4: Tag the Image for Your Docker Hub Repo
Docker Hub requires the format `username/repository:tag`.
```bash
docker tag myapp:latest yourusername/myapp:latest
```

---

#### Step 5: Push the Image
```bash
docker push yourusername/myapp:latest
```

---

#### Step 6: Verify
Go to `https://hub.docker.com/r/yourusername/myapp` to confirm it uploaded.

---

#### Step 7 (Optional): Pull It Anywhere
Anyone can now run it with:
```bash
docker pull yourusername/myapp:latest
docker run yourusername/myapp:latest
```

---

#### Versioning Tip
Push both a version tag and `latest` so users can pin a specific version:
```bash
docker tag myapp:latest yourusername/myapp:1.0.0
docker push yourusername/myapp:1.0.0
docker push yourusername/myapp:latest
```

---

#### Using GitHub Container Registry Instead (Alternative to Docker Hub)
```bash
docker login ghcr.io -u yourusername
docker tag myapp:latest ghcr.io/yourusername/myapp:latest
docker push ghcr.io/yourusername/myapp:latest
```
`this is going to be covered as it own seprate guide GHCR !`

### Docker Volumes

**What they are:** a way to persist data outside a container's filesystem, so data survives container restarts, removals, and rebuilds.

---

### Why Volumes Matter

**Containers are ephemeral** — when a container is removed, everything written inside it is lost. Volumes solve this by storing data outside the container's writable layer.

---

### Types of Volumes

#### Named Volumes
**Managed by Docker**, stored in Docker's own storage area. Best for persistent data like databases.

general format for this: `-v host_path:conatiner_path`
```bash
docker volume create mydata
docker run -v mydata:/var/lib/postgresql/data postgres:16
```

#### Bind Mounts
**Maps a specific folder on your host** directly into the container. Best for development (live code editing).
```bash
docker run -v /home/user/myapp:/app node:20
```

#### tmpfs Mounts
**Stored in memory only**, never written to disk. Best for sensitive/temporary data.
```bash
docker run --tmpfs /app/cache node:20
```

---

### Creating and Managing Named Volumes

#### Create a Volume
```bash
docker volume create mydata
```

#### List Volumes
```bash
docker volume ls
```

#### Inspect a Volume
```bash
docker volume inspect mydata
```

#### Remove a Volume
```bash
docker volume rm mydata
```

#### Remove All Unused Volumes
```bash
docker volume prune
```

---

### Using Volumes with `docker run`

#### Named Volume
```bash
docker run -v mydata:/data ubuntu
```

#### Bind Mount
```bash
docker run -v $(pwd):/app ubuntu
```

#### Read-Only Mount
**Prevents the container** from writing to the mounted path.
```bash
docker run -v $(pwd):/app:ro ubuntu
```

#### Using `--mount` (more explicit, recommended for production)
```bash
docker run --mount source=mydata,target=/data ubuntu
```

---

### Using Volumes in Docker Compose

```yaml
services:
  db:
    image: postgres:16
    volumes:
      - db_data:/var/lib/postgresql/data   # named volume
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql   # bind mount

volumes:
  db_data:
```

---

### Common Use Cases

**Persisting database data**
```bash
docker run -v db_data:/var/lib/postgresql/data postgres:16
```

**Live code editing during development**
```bash
docker run -v $(pwd):/app -p 3000:3000 node:20
```

**Sharing data between multiple containers**
```bash
docker run -v shared_data:/data --name writer ubuntu
docker run -v shared_data:/data --name reader ubuntu
```

---

### Quick Guidance

| Use Case | Best Choice |
|---|---|
| Database storage | Named volume |
| Live development / editing code | Bind mount |
| Sensitive temp data (secrets, cache) | tmpfs |
| Sharing data between containers | Named volume |

---

### Cleaning Up

```bash
docker volume ls              # see what exists
docker volume prune           # remove unused volumes
docker volume rm <name>       # remove a specific volume
```

**Warning:** removing a volume permanently deletes its data. Always confirm nothing important depends on it first.

`volumes are mentioned docker.yaml files`

default location of docker volumes 

- ~spyware~ windows: C:\ProgramData\docker\volumes
- Mac\Linux: /var/lib/docker/volumes 

### Docker Networks

**What they are:** a way for containers to communicate with each other, with the host, and with the outside world, in an isolated and controllable manner.

---

### Why Networks Matter

**By default, each container is isolated.** Docker networks let you control exactly which containers can talk to each other, using container names instead of IP addresses, without exposing everything to the outside world.

---

### Types of Networks

#### Bridge (default)
**Standard network for standalone containers** on a single host. Containers on the same bridge network can reach each other by name.
```bash
docker network create mynetwork
```

#### Host
**Removes network isolation** — the container shares the host's network stack directly. Faster, but less secure and no port mapping needed/possible.
```bash
docker run --network host nginx
```

#### None
**No networking at all.** Useful for fully isolated tasks (e.g. batch processing with no network need).
```bash
docker run --network none ubuntu
```

#### Overlay
**Connects containers across multiple hosts** — used in Docker Swarm or multi-node setups. Not needed for single-machine setups.

---

### Creating and Managing Networks

#### Create a Network
```bash
docker network create mynetwork
```

#### List Networks
```bash
docker network ls
```

#### Inspect a Network
```bash
docker network inspect mynetwork
```

#### Remove a Network
```bash
docker network rm mynetwork
```

#### Remove All Unused Networks
```bash
docker network prune
```

---

### Connecting Containers to a Network

#### Run a Container on a Specific Network
```bash
docker run --network mynetwork --name app1 myapp
```

#### Connect an Already-Running Container
```bash
docker network connect mynetwork app1
```

#### Disconnect a Container
```bash
docker network disconnect mynetwork app1
```

---

### Container-to-Container Communication

**Containers on the same custom network** can reach each other using their container name as the hostname — no IP addresses needed.

```bash
docker network create appnet
docker run -d --network appnet --name db postgres:16
docker run -d --network appnet --name web myapp
```

Inside `web`, you can connect to the database just using the hostname `db`:

postgresql://db:5432/mydb


**Note:** this name-based discovery only works on custom networks, not the default bridge network.

---

### Using Networks in Docker Compose

**Compose creates a network automatically** for every project, and all services can reach each other by service name — no manual setup needed.

```yaml
services:
  web:
    build: .
    depends_on:
      - db

  db:
    image: postgres:16
```
Here, `web` can connect to `db` simply by using `db` as the hostname — Compose handles the network creation behind the scenes.

#### Defining a Custom Network Explicitly
```yaml
services:
  web:
    build: .
    networks:
      - appnet

  db:
    image: postgres:16
    networks:
      - appnet

networks:
  appnet:
```

---

### Exposing Ports to the Outside World

**Networks control container-to-container traffic** — to let the outside world (your browser, other machines) reach a container, you still need port publishing.

```bash
docker run -p 8080:80 nginx
```

| Flag | Meaning |
|---|---|
| `-p 8080:80` | Maps host port 8080 → container port 80 |
| `EXPOSE 80` (in Dockerfile) | Documents the port, doesn't publish it |

---

### Quick Guidance

| Use Case | Best Choice |
|---|---|
| Standalone containers on one machine | Bridge (custom, not default) |
| Multiple services talking to each other | Custom bridge network |
| Max performance, no isolation needed | Host |
| Fully isolated task, no networking | None |
| Multi-host cluster (Swarm) | Overlay |

---

### Cleaning Up

```bash
docker network ls              # see what exists
docker network prune           # remove unused networks
docker network rm <name>       # remove a specific network
```

**Note:** you can't remove a network while containers are still connected to it — disconnect or stop them first.

 ## Kubernetes guide

 ### What is Kubernetes

**Kubernetes (K8s)** is an open-source system for automating deployment, scaling, and management of containerized applications. It takes containers (built with Docker or similar tools) and runs them across a cluster of machines, handling things manually managed with plain Docker: restarting crashed containers, scaling up/down based on load, load balancing traffic, rolling out updates without downtime, and self-healing when something fails.

**The core problem it solves:** Docker runs containers on one machine. Kubernetes runs and manages containers across *many* machines, as if they were one big computer.

**Why "K8s"?** Kubernetes has 8 letters between the "K" and the "s" — hence the shorthand.

---

### K8s Architecture

**A Kubernetes cluster is made of two types of machines (nodes):**

```
                ┌─────────────────────-┐
                │   Control Plane      │  ← the "brain"
                │  (Master Node)       │
                └──────────┬───────────┘
                            │
        ┌───────────────────┼───────────────-────┐
        │                   │                    │
┌───────▼───────┐   ┌───────▼───────-┐   ┌───────▼───-────┐
│  Worker Node 1│   │  Worker Node 2 │   │  Worker Node 3 │
│  (runs pods)  │   │  (runs pods)   │   │  (runs pods)   │
└───────────────┘   └────────────────┘   └────────────────┘
```

#### Control Plane
**Makes global decisions** about the cluster (scheduling, scaling, responding to failures). Contains:

- **API Server** — the front door; every command (`kubectl`) talks to this
- **etcd** — key-value store holding the entire cluster's state/configuration
- **Scheduler** — decides which node a new pod should run on
- **Controller Manager** — watches the cluster state and corrects it when it drifts from what's desired

#### Worker Nodes
**Actually run your application containers.** Each worker node contains:

- **kubelet** — agent that talks to the control plane and manages containers on that node
- **kube-proxy** — handles networking rules so pods can communicate
- **Container runtime** — the software that actually runs containers (e.g. containerd, Docker)

---

### Main K8s Components

| Component | What it is |
|---|---|
| **Pod** | Smallest deployable unit — one or more containers that share storage/network |
| **Node** | A physical or virtual machine running pods |
| **Deployment** | Manages a set of identical pods, handles scaling and rolling updates |
| **Service** | Stable network endpoint to access a group of pods |
| **Namespace** | Virtual cluster used to organize/isolate resources |
| **ConfigMap** | Stores non-sensitive configuration data |
| **Secret** | Stores sensitive data (passwords, tokens) |
| **Volume / PersistentVolume** | Storage that survives beyond a pod's lifecycle |
| **Ingress** | Manages external HTTP/HTTPS access into the cluster |
| **StatefulSet** | Like a Deployment, but for apps that need stable identity/storage (databases) |
| **ReplicaSet** | Ensures a specified number of pod replicas are running (managed by Deployments) |

**Relationship:**
```
Deployment  →  creates & manages  →  ReplicaSet  →  creates & manages  →  Pods
Service     →  routes traffic to  →  Pods (via labels)
```
few points

- `kubernetes nework -> 127.0.0.1:233 (everything acts as one machine) docker's network -> 172.18.0.2 (everything act as different machine)`

- `kubernetes volumes concept is different from the concept of docker volume as it's for one set of cotainers but with kubernetes even these machine can shift so how the volumes move`

- `services basically act as a communicator as well the load balacer for the traffic for the pod dividing it between different indentical conatiner of appilication running on different pods` 

- `application like databases cannot be directly replicated (unlike stateless applications  which can be replicated using deployment) and requires something called statefullset for replication` 

---

### Minikube + kubectl — Local Setup

**Minikube** runs a single-node Kubernetes cluster on your own machine — great for learning and local development.

#### Install kubectl (CLI to control any cluster)
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/
kubectl version --client
```

kubectl is the standard/primary tool for talking to the Control Plane. There are alternatives (raw REST calls to the API Server, the Kubernetes Dashboard GUI, Helm for higher-level operations, or the cloud provider's own console for managed clusters like EKS/GKE) — but under the hood, all of them are just hitting the same API Server that kubectl talks to. kubectl is simply the most direct, most commonly used way to do it.

#### Install Minikube
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

#### Start a Local Cluster
```bash
minikube start
```

#### Check Cluster Status
```bash
minikube status
kubectl get nodes
```

#### Access the Dashboard (optional GUI)
```bash
minikube dashboard
```

#### Common kubectl Commands
```bash
kubectl get pods                  # list pods
kubectl get services              # list services
kubectl get deployments           # list deployments
kubectl describe pod <name>       # detailed info + events
kubectl logs <pod-name>           # view logs
kubectl exec -it <pod-name> -- bash   # shell into a pod
kubectl apply -f myfile.yaml      # create/update from a YAML file
kubectl delete -f myfile.yaml     # delete resources from a YAML file
kuberctl exec -it <node-id> -- bin/bash # get's you in the interactive terminal
kubectl delete deployment <pod-name> # to delete a pod (get pods will now show terminating)
```

#### Stop / Delete the Cluster
```bash
minikube stop
minikube delete
```

---

### K8s YAML Configuration File

**Kubernetes resources are defined declaratively** using YAML files — you describe the desired state, and Kubernetes works to make reality match it.

#### Basic Structure (every K8s YAML has these 4 fields)
```yaml
apiVersion: apps/v1     # which K8s API version this resource uses
kind: Deployment        # what kind of resource this is
metadata:               # name, labels, namespace
  name: myapp
spec:                   # the desired state/configuration
  ...
```

#### Example: A Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-deployment
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: myapp:latest
          ports:
            - containerPort: 3000
```

#### Applying It
```bash
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods
```

#### Key Fields Explained

| Field | Meaning |
|---|---|
| `apiVersion` | Version of the Kubernetes API being used |
| `kind` | Type of object (Pod, Deployment, Service, etc.) |
| `metadata.name` | Name of the resource |
| `spec.replicas` | How many pod copies to run |
| `spec.selector` | How the Deployment finds its pods (via labels) |
| `spec.template` | The pod blueprint (its own mini spec) |
| `containers.image` | Which image to run |
| `containers.ports` | Which port the container listens on |

---

### K8s Namespaces — Organize Your Components

**Namespaces are virtual clusters** within a physical cluster — used to separate environments (dev/staging/prod) or teams, without needing separate physical clusters.

`different namespaces cannot share config maps and secrets ,they can only share services and certain resources like volumes cannot/should not be included in any namespace`

#### Default Namespaces
```bash
kubectl get namespaces
```
Typically includes: `default`, `kube-system`, `kube-public`, `kube-node-lease`.

#### Create a Namespace
```bash
kubectl create namespace dev
```

Or via YAML:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: dev
```

#### Deploy Into a Specific Namespace
```bash
kubectl apply -f deployment.yaml -n dev
```

Or specify it in the YAML itself:
```yaml
metadata:
  name: myapp-deployment
  namespace: dev
```

#### View Resources in a Namespace
```bash
kubectl get pods -n dev
kubectl get all -n dev
```

#### Set a Default Namespace for Your Session
```bash
kubectl config set-context --current --namespace=dev
```

#### Delete a Namespace (and everything inside it!)
```bash
kubectl delete namespace dev
```

#### Managing Active Namespaces with `kubens`
 
**Typing `kubectl config set-context --current --namespace=dev` every time you switch namespaces gets old fast.** `kubens` is a small third-party CLI tool that makes switching your active namespace instant.
 
**Install it (part of the `kubectx` project):**
```bash
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```
Or via a package manager (e.g. Homebrew on macOS/Linux):
```bash
brew install kubectx
```
 
**List all namespaces:**
```bash
kubens
```
 
**Switch your active namespace:**
```bash
kubens dev
```
Every subsequent `kubectl` command now runs against `dev` by default, no `-n dev` needed.
 
**Switch back to the previous namespace:**
```bash
kubens -
```
 
**Check which namespace you're currently on:**
```bash
kubens -c
```
 
`kubens` doesn't do anything `kubectl config set-context` can't — it's purely a shortcut so you're not typing the full command every time you jump between `dev`, `staging`, and `prod`.
 
---
 
### K8s ConfigMaps — Storing Non-Sensitive Configuration
 
**ConfigMaps store configuration data** (URLs, feature flags, app settings) separately from your application code, so you can change config without rebuilding your image.
 
#### Option 1: Create a ConfigMap via `kubectl`
 
```bash
kubectl create configmap app-config \
  --from-literal=APP_ENV=production \
  --from-literal=LOG_LEVEL=info
```
 
From a file instead:
```bash
kubectl create configmap app-config --from-file=config.properties
```
 
Check it:
```bash
kubectl get configmaps
kubectl describe configmap app-config
```
 
#### Option 2: Create a ConfigMap via YAML
 
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
```
 
```bash
kubectl apply -f configmap.yaml
```
 
#### Using a ConfigMap as Environment Variables in a Pod
 
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      env:
        - name: APP_ENV
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: APP_ENV
```
 
**Load every key as an env var at once** (instead of listing them one by one):
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      envFrom:
        - configMapRef:
            name: app-config
```
 
#### Using a ConfigMap as a Mounted Volume
 
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      volumeMounts:
        - name: config-volume
          mountPath: /etc/config
  volumes:
    - name: config-volume
      configMap:
        name: app-config
```
Each key becomes a separate file inside `/etc/config` (e.g. `/etc/config/APP_ENV`) — useful when your app expects a config file rather than env vars.
 

---
### K8s Secrets — Storing Sensitive Data
 
**Secrets store sensitive information** (passwords, API keys, tokens, TLS certs) separately from your application code and pod specs, so you never hard-code credentials directly into a Deployment YAML.
 
**Note:** Secrets are base64-*encoded*, not encrypted, by default — anyone with API access to read the Secret can decode it. For real security, enable **encryption at rest** on etcd or use a dedicated secrets manager (Vault, AWS Secrets Manager, Sealed Secrets) for production-sensitive data.
 
#### Option 1: Create a Secret via `kubectl` (no manual encoding needed)
 
```bash
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=supersecret123
```
 
Check it exists (values are hidden by default):
```bash
kubectl get secrets
kubectl describe secret db-secret
```
 
#### Option 2: Create a Secret via YAML
 
Values must be base64-encoded manually first:
```bash
echo -n 'admin' | base64          # YWRtaW4=
echo -n 'supersecret123' | base64 # c3VwZXJzZWNyZXQxMjM=
```
 
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  username: YWRtaW4=
  password: c3VwZXJzZWNyZXQxMjM=
```
 
```bash
kubectl apply -f secret.yaml
```
 
#### Using a Secret as Environment Variables in a Pod
 
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      env:
        - name: DB_USERNAME
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
```
 
#### Using a Secret as a Mounted Volume (alternative to env vars)
 
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      volumeMounts:
        - name: secret-volume
          mountPath: /etc/secrets
          readOnly: true
  volumes:
    - name: secret-volume
      secret:
        secretName: db-secret
```
Each key in the Secret becomes a separate file inside `/etc/secrets` (e.g. `/etc/secrets/username`, `/etc/secrets/password`).
 
#### Secret vs ConfigMap
 
| | ConfigMap | Secret |
|---|---|---|
| Purpose | Non-sensitive config (URLs, feature flags) | Sensitive data (passwords, keys, tokens) |
| Storage | Plain text | Base64-encoded (not encrypted by default) |
| Typical use | App settings | DB credentials, API keys, TLS certs |
 
#### Viewing/Decoding a Secret (for debugging)
 
```bash
kubectl get secret db-secret -o jsonpath='{.data.password}' | base64 --decode
```
---

### K8s Ingress

**Ingress manages external HTTP/HTTPS access** into services running inside the cluster — it acts like a smart router/reverse proxy, letting you use one entry point with domain/path-based routing instead of exposing each service separately.

#### Why Not Just Use a Service?
A `Service` of type `LoadBalancer` exposes one service per public IP — expensive and messy at scale. Ingress lets multiple services share a single entry point.

#### Prerequisite: Install an Ingress Controller
Ingress rules do nothing without a controller running (e.g. NGINX Ingress Controller).
```bash
minikube addons enable ingress
```

#### Example Ingress YAML
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: myapp.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: myapp-service
                port:
                  number: 80
```

#### Applying It
```bash
kubectl apply -f ingress.yaml
kubectl get ingress
```

#### Path-Based Routing Example
```yaml
paths:
  - path: /api
    pathType: Prefix
    backend:
      service:
        name: api-service
        port:
          number: 80
  - path: /
    pathType: Prefix
    backend:
      service:
        name: frontend-service
        port:
          number: 80
```

---

### Helm — Package Managers (e.g. MongoDB)

**Helm is the package manager for Kubernetes** — instead of writing every YAML file by hand, you install pre-built "charts" (bundles of YAML templates) with one command.

#### Install Helm
```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
helm version
```

#### Add a Chart Repository
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

#### Search for a Chart
```bash
helm search repo mongodb
```

#### Install MongoDB Using Helm
```bash
helm install my-mongo bitnami/mongodb
```

#### Check the Release
```bash
helm list
kubectl get pods
```

#### Customize Values Before Installing
```bash
helm show values bitnami/mongodb > values.yaml
# edit values.yaml as needed
helm install my-mongo bitnami/mongodb -f values.yaml
```

#### Upgrade a Release
```bash
helm upgrade my-mongo bitnami/mongodb -f values.yaml
```

#### Uninstall a Release
```bash
helm uninstall my-mongo
```
#### Helm Release Management
 
**Every `helm install`/`upgrade` creates a new numbered revision** of that release, which Helm tracks — letting you inspect history and roll back if an upgrade breaks something.
 
**List all releases in the current namespace:**
```bash
helm list
```
 
**List releases across all namespaces:**
```bash
helm list --all-namespaces
```
 
**Check a release's current status:**
```bash
helm status my-mongo
```
 
**View a release's revision history:**
```bash
helm history my-mongo
```
This shows every past install/upgrade, with a revision number, status, and description for each.
 
**Roll back to a previous revision** (e.g. an upgrade caused a crash loop):
```bash
helm rollback my-mongo 1
```
This reverts the release back to revision `1`, undoing whatever the last upgrade changed.
 
**See exactly what changed in a specific revision:**
```bash
helm get values my-mongo --revision 2
helm get manifest my-mongo --revision 2
```
 
**Why this matters:** unlike a raw `kubectl apply`, Helm remembers every version of a release, so a bad upgrade is a one-command rollback instead of manually reconstructing the previous YAML state.


**Why Helm matters:** deploying something like MongoDB with proper persistence, secrets, and replication by hand would take dozens of lines of YAML — Helm charts package all of that into one tested, reusable install.

#### Using Your Own YAML Files with Helm
 
**You're not limited to installing other people's charts** — you can turn your own Deployment/Service/ConfigMap YAML files into a Helm chart, which lets you reuse them across environments (dev/staging/prod) by swapping just a `values.yaml` file instead of editing raw YAML every time.
 
```bash
helm create mychart
```
This scaffolds a folder structure:
```
mychart/
├── Chart.yaml          # chart metadata (name, version)
├── values.yaml          # default configurable values
└── templates/           # your YAML files go here
```
 
Move your existing YAML files into `mychart/templates/`, then replace hard-coded values (image name, replica count, etc.) with template variables pulling from `values.yaml`:
```yaml
spec:
  replicas: {{ .Values.replicaCount }}
  containers:
    - image: {{ .Values.image }}
```
 
Install your own chart the same way you'd install any other:
```bash
helm install myapp ./mychart
```
 
Update values without touching the templates:
```bash
helm upgrade myapp ./mychart --set replicaCount=5
```
---

### Volumes — Persisting Data in K8s

**Pods are ephemeral** — when a pod dies, its filesystem is gone. Kubernetes volumes let data persist beyond a pod's lifecycle.

#### Volume Types

| Type | Use Case |
|---|---|
| `emptyDir` | Temporary storage, shared between containers in a pod, deleted when pod dies |
| `hostPath` | Mounts a path from the host node (not recommended for production) |
| `PersistentVolume (PV)` | Cluster-wide storage resource, independent of pod lifecycle |
| `PersistentVolumeClaim (PVC)` | A request for storage — pods use PVCs to claim a PV |

#### Persistent Volume (PV) Example
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: myapp-pv
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/myapp
```
`there can be many storage mediums like cloud , local , folder ... etc`

#### Persistent Volume Claim (PVC) Example
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myapp-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```
#### StorageClasses — Dynamic Volume Provisioning

**Instead of manually creating a PersistentVolume ahead of time**, a StorageClass lets Kubernetes automatically create one on demand whenever a PVC asks for it — this is how most real clusters (especially cloud-managed ones) actually work.

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs   # varies by cloud: aws-ebs, gce-pd, azure-disk, etc.
parameters:
  type: gp3
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

**Referencing it from a PVC** — no PV needs to exist beforehand, it gets created automatically:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myapp-pvc
spec:
  storageClassName: fast-ssd
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

**Check available StorageClasses in your cluster:**
```bash
kubectl get storageclass
```

Most managed clusters (EKS, GKE, AKS) ship with a `default` StorageClass already — if a PVC doesn't specify `storageClassName` at all, Kubernetes uses whichever StorageClass is marked default.

| Field | Meaning |
|---|---|
| `provisioner` | Which storage backend/plugin creates the actual disk |
| `reclaimPolicy` | What happens to the disk when the PVC is deleted (`Delete` or `Retain`) |
| `volumeBindingMode` | `Immediate` (bind right away) or `WaitForFirstConsumer` (wait until a pod actually needs it, then pick the right zone/node) |

**Why it matters:** without a StorageClass, someone has to manually pre-create a PV for every single PVC request — StorageClasses automate that, which is essential once you're not just running one Minikube pod locally.

#### Using a PVC in a Pod
```yaml
spec:
  containers:
    - name: myapp
      image: myapp:latest
      volumeMounts:
        - mountPath: /data
          name: myapp-storage
  volumes:
    - name: myapp-storage
      persistentVolumeClaim:
        claimName: myapp-pvc
```

#### Checking Volumes
```bash
kubectl get pv
kubectl get pvc
```

**Access Modes:**
| Mode | Meaning |
|---|---|
| `ReadWriteOnce` | One node can read/write |
| `ReadOnlyMany` | Many nodes can read |
| `ReadWriteMany` | Many nodes can read/write |

---

### K8s StatefulSet — Deploying Stateful Apps

**StatefulSets manage pods that need a stable identity** — unlike Deployments where pods are interchangeable, StatefulSet pods get:

- **Stable, predictable names** (`myapp-0`, `myapp-1`, `myapp-2` — not random hashes)
- **Stable network identity** (each pod gets its own consistent DNS entry)
- **Stable storage** (each pod keeps its own PersistentVolume, even after restarts)
- **Ordered startup/shutdown** (pod 0 starts before pod 1, etc.)

**Use case:** databases, message queues, anything where each replica needs its own persistent identity and data (e.g. MongoDB replica sets, Kafka, Elasticsearch).

#### Example StatefulSet YAML
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongo
spec:
  serviceName: "mongo"
  replicas: 3
  selector:
    matchLabels:
      app: mongo
  template:
    metadata:
      labels:
        app: mongo
    spec:
      containers:
        - name: mongo
          image: mongo:7
          ports:
            - containerPort: 27017
          volumeMounts:
            - name: mongo-data
              mountPath: /data/db
  volumeClaimTemplates:
    - metadata:
        name: mongo-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

**Key difference from Deployment:** `volumeClaimTemplates` automatically creates a *unique* PersistentVolumeClaim for each pod replica (so `mongo-0`, `mongo-1`, `mongo-2` each get their own dedicated storage).

`all storage like mongo-0 , mongo-1 , mongo-2 are not same they share a master slave relationship: search it it basically works by creating and copying a storage that's the master above it`

#### Applying and Checking
```bash
kubectl apply -f statefulset.yaml
kubectl get statefulsets
kubectl get pods -l app=mongo   # notice ordered names: mongo-0, mongo-1, mongo-2
```

#### Deployment vs StatefulSet

| Feature | Deployment | StatefulSet |
|---|---|---|
| Pod names | Random | Predictable (`app-0`, `app-1`) |
| Storage | Shared or none | Unique per pod |
| Startup order | Any order | Sequential |
| Use case | Stateless apps (web servers, APIs) | Stateful apps (databases) |

---

### K8s Service

**A Service gives a stable network endpoint** to a group of pods — since pods are ephemeral and get new IPs when recreated, Services provide a fixed way to reach them regardless of which underlying pods are currently running. ingress talks to the service and service talks to the pod 

#### How It Works
```
Service (stable IP/DNS name)  →  routes to  →  Pods matching a label selector
```
`IF YOU ARE JUST HERE TO COVERT YOU DOCKER APPLICATION TO KUBERNETES APPLICATION SKIP TO THE LAST PART`

#### Service Types

| Type | Use Case |
|---|---|
| `ClusterIP` (default) | Internal-only access, reachable within the cluster |
| `NodePort` | Exposes the service on a static port on every node's IP |
| `LoadBalancer` | Provisions an external load balancer (cloud provider dependent) |
| `ExternalName` | Maps the service to an external DNS name |
| `headless` | this is specifically for the stateless stuff (just set clusterIp = none) |

#### ClusterIP Example (internal communication)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: ClusterIP
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 3000
```

#### NodePort Example (accessible from outside the cluster)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: NodePort
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 3000
      nodePort: 30036
```

#### LoadBalancer Example (cloud environments)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 3000
```

#### Applying and Checking
```bash
kubectl apply -f service.yaml
kubectl get services
kubectl describe service myapp-service
```
`service file do get complicated in the multiport and advanced systems use AI to get the yaml file as well as refer to the documentation`

#### Key Fields

| Field | Meaning |
|---|---|
| `selector` | Which pods this service routes to (matches pod labels) |
| `port` | Port the Service itself listens on |
| `targetPort` | Port the container inside the pod is listening on |
| `nodePort` | External port opened on every node (NodePort type only) |

--

### Full Mental Model

```
Dockerfile → builds image
     ↓
Deployment → runs pods from that image, keeps desired replica count
     ↓
Service → gives pods a stable way to be reached
     ↓
Ingress → routes external HTTP traffic into the right Service
     ↓
PVC/PV → gives pods persistent storage
     ↓
StatefulSet → same as Deployment, but for apps needing stable identity/storage (databases)
     ↓
Namespace → organizes/isolates all of the above (dev, staging, prod)
     ↓
Helm → packages a full set of these YAMLs into one reusable, installable chart
```

### From Docker to Kubernetes

Assuming you already have a working Docker app (image + maybe a `docker-compose.yml`), here's the path to get it running on Kubernetes.

---

### Step 1: Make Sure Your Image Is in a Registry K8s Can Reach

**Kubernetes doesn't build images** — it only pulls them. Push your image somewhere accessible.

```bash
docker tag myapp:latest yourusername/myapp:latest
docker push yourusername/myapp:latest
```

If you're using Minikube locally, you can skip pushing and just do:
```bash
minikube image load myapp:latest
```

---

### Step 2: Translate `docker run` Flags into a Deployment YAML

**Whatever flags you were using in `docker run` map roughly like this:**

| Docker flag | Kubernetes equivalent |
|---|---|
| `docker run -p 3000:3000` | `Service` (port mapping) |
| `docker run -e KEY=val` | `env:` in the pod spec |
| `docker run -v host:container` | `volumes` + `volumeMounts` (or PVC) |
| `docker run --name` | `metadata.name` |
| `docker run --restart=always` | K8s restarts pods automatically by default |

Example Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: yourusername/myapp:latest
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: "production"
```

---

### Step 3: Expose It with a Service

**`docker run -p 3000:3000` becomes a `Service`**, since pods don't get stable ports/IPs on their own.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 3000
  type: NodePort   # or ClusterIP if only internal, LoadBalancer if on cloud
```

---

### Step 4: Handle Volumes, If You Had Any

**If your `docker run` had a `-v` flag** for persistent data (like a database):

```yaml
volumeMounts:
  - name: myapp-storage
    mountPath: /data
volumes:
  - name: myapp-storage
    persistentVolumeClaim:
      claimName: myapp-pvc
```

Needs a separate PVC + PV, or just a PVC if your cluster has dynamic provisioning via a StorageClass.

---

### Step 5: Apply It

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

Check it's running:
```bash
kubectl get pods
kubectl get services
```

---

### Step 6: If You Had `docker-compose.yml` (Multiple Services)

**Each service in your compose file** becomes its own Deployment + Service pair. So a `web` + `db` compose file turns into 4 YAML files (or one combined file with `---` separators): `web-deployment.yaml`, `web-service.yaml`, `db-deployment.yaml`, `db-service.yaml`.

Services reach each other the same way Compose did — by name:
```
compose:  web talks to "db"
k8s:      web talks to "db-service" (the Service name)
```

**Shortcut:** there's a tool called **Kompose** that auto-converts a `docker-compose.yml` into Kubernetes YAML files for you:
```bash
kompose convert -f docker-compose.yml
```

It won't be perfect (especially for volumes/networking nuances), but it gives you a solid starting point instead of writing everything by hand.

---

### One-Line Summary

```
docker run flags        →  Deployment spec (image, env, ports)
docker run -p            →  Service
docker run -v            →  PVC + volumeMounts
docker-compose.yml       →  one Deployment + Service pair PER service (or use Kompose)
```

