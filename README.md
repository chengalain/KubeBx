# KubeBx

**KubeBx** is an open-source CLI tool to learn Kubernetes through hands-on, interactive exercises on your local machine. The name stands for **KubeBuilderX** — build your Kubernetes experience step by step, learn by doing, debug by breaking.

## Features

- **Zero friction setup** — only Docker required, everything else is managed automatically
- **Learn by doing** — practical exercises, not just theory
- **Mix of build & debug** — create resources AND fix broken environments
- **Instant validation** — check your solutions via Kubernetes API
- **Smart hints** — get help when stuck
- **Fully offline** — works without internet after initial setup

## Exercise Types

**Build exercises** — Learn by creating resources step by step  
**Debug exercises** — Real-world scenarios where you fix broken deployments

## Prerequisites

- **Docker** installed and running
- That's it! KubeBx handles the rest (Kind, kubectl)

## Quick Start

### 1. Install KubeBx

```bash
# Download the latest release
# TODO: Add release URL when published

# Or build from source
git clone https://github.com/cheng-alain/kubebx.git
cd kubebx
go build -o kbx
```

### 2. Initialize your cluster

```bash
./kbx init
```

This will:
- Check that Docker is running
- Download Kind and kubectl (stored in `~/.kubebx/bin/`)
- Create a local Kubernetes cluster

### 3. Start learning!

```bash
# List available exercises
./kbx list

# Start the first exercise
./kbx start 01

# Read the instructions
cat exercises/01-pod-basics/README.md

# Work on the exercise using kbx kubectl
./kbx kubectl run my-first-pod --image=nginx:latest -n kbx-01

# Check your solution
./kbx check 01

# Need help?
./kbx hint 01

# Move to next exercise
./kbx next
```

## Available Commands

```bash
kbx init              # Initialize the local Kubernetes cluster
kbx list              # List all available exercises
kbx start <id>        # Start an exercise
kbx check <id>        # Validate your solution
kbx hint <id>         # Get a hint
kbx next              # Move to the next exercise
kbx progress          # Show your learning progress
kbx kubectl [args]    # Run kubectl commands
kbx clean <id>        # Clean up exercise resources
kbx clean --all       # Clean all exercises
kbx reset             # Delete the cluster completely
```

## Roadmap

**v1.0 Exercises:**
- 01 - Pod basics (build)
- 02 - Labels & Selectors (build)
- 03 - Service broken (debug)
- 04 - Deployments (build)
- 05 - Pod crashloop (debug)
- 06 - ConfigMaps & Secrets (build)
- 07 - Pending pod (debug)
- 08 - Ingress (build)
- 09 - RBAC misconfigured (debug)
- 10 - Full broken app (debug - boss level)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Uses [Kind](https://kind.sigs.k8s.io/) for local clusters
- Powered by [client-go](https://github.com/kubernetes/client-go)
