# Exercise 03: Service broken

**Type:** Debug  
**Difficulty:** Beginner

## Situation

A pod `nginx-pod` is running in the `kbx-03` namespace. A Service `nginx-service` also exists and should route traffic to that pod — but it doesn't.

Your job is to find why the Service is not routing traffic and fix it.

## Goal

Make the Service `nginx-service` correctly route traffic to the `nginx-pod` pod.

## Useful commands

```bash
# Check the pod
kbx kubectl get pods -n kbx-03 --show-labels

# Check the Service
kbx kubectl get svc -n kbx-03

# Inspect the Service in detail
kbx kubectl describe svc nginx-service -n kbx-03

# Check if the Service has any endpoints
kbx kubectl get endpoints nginx-service -n kbx-03
```

## Success criteria

- Pod `nginx-pod` is `Running`
- Service `nginx-service` exists
- The Service selector matches the pod labels
- The Service has at least one active endpoint

## Stuck?

Run `kbx hint 03` for a hint.
