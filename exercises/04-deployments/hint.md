# Hint for Exercise 04: Deployments

## Hint 1

Create the Deployment with `kubectl create deployment`, then scale it:

```bash
kbx kubectl create deployment nginx-deployment --image=nginx:latest -n kbx-04
```

Then check what was created:

```bash
kbx kubectl get deployment -n kbx-04
kbx kubectl get pods -n kbx-04
```

By default, `create deployment` starts with 1 replica. You need to scale it to 3.

## Hint 2

Scale the Deployment to 3 replicas:

```bash
kbx kubectl scale deployment nginx-deployment --replicas=3 -n kbx-04
```

Verify all 3 replicas are ready:

```bash
kbx kubectl get deployment nginx-deployment -n kbx-04
# READY column should show 3/3
```

## Check your solution

```bash
kbx check 04
```

---

**Key concept:** The Deployment controller continuously reconciles the actual state (running pods) with the desired state (replicas: 3). Delete a pod and watch it come back automatically.
