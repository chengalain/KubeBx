# Exercise 01: Pod Basics

**Type:** Build  
**Difficulty:** Beginner

##  Goal

Create your first Kubernetes pod and understand the basic concepts.

##  Tasks

1. Create a pod named `my-first-pod` in the `kbx-01` namespace
2. Use the `nginx:latest` image
3. Verify the pod is running

##  Success Criteria

- Pod `my-first-pod` exists in the `kbx-01` namespace
- Pod status is `Running`
- Container is using the nginx image

##  Tips

- Use `kbx kubectl run` to create a pod quickly
- Use `kbx kubectl get pods -n kbx-01` to check status
- Use `kbx kubectl describe pod <name> -n kbx-01` for details

##  Stuck?

Run `kbx hint 01` for a hint!