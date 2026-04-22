# Exercise 04: Deployments

**Type:** Build  
**Difficulty:** Beginner

## Context

A bare pod is fragile: if it crashes, nothing restarts it. A Deployment manages pods for you — it ensures the desired number of replicas are always running, handles restarts, and makes scaling trivial.

## Goal

Create a Deployment `nginx-deployment` in the `kbx-04` namespace with:
- **3 replicas**
- Image: `nginx:latest`
- Label: `app=nginx` on the pod template

## Success criteria

- Deployment `nginx-deployment` exists in `kbx-04`
- Spec declares 3 replicas
- All 3 replicas are Ready
- Containers use the `nginx:latest` image

## Useful commands

```bash
# Create a deployment
kbx kubectl create deployment nginx-deployment --image=nginx:latest -n kbx-04

# Check the deployment
kbx kubectl get deployment -n kbx-04

# Check the pods created by the deployment
kbx kubectl get pods -n kbx-04

# Describe the deployment for details
kbx kubectl describe deployment nginx-deployment -n kbx-04
```

## Stuck?

Run `kbx hint 04` for a hint.
