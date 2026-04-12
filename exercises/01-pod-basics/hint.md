# Hint for Exercise 01: Pod Basics

## 💡 Step-by-step guide

### Creating a pod

The easiest way to create a pod is using `kubectl run`:

```bash
kbx kubectl run <pod-name> --image=<image-name> -n <namespace>
```

### For this exercise

You need to create a pod with:
- **Name:** `my-first-pod`
- **Image:** `nginx:latest`
- **Namespace:** `kbx-01`

### Example command

```bash
kbx kubectl run my-first-pod --image=nginx:latest -n kbx-01
```

### Verify your work

```bash
kbx kubectl get pods -n kbx-01
kbx kubectl describe pod my-first-pod -n kbx-01
```

### Check your solution

```bash
kbx check 01
```

---

**Still stuck?** The pod name must be exactly `my-first-pod` and it must be in the `kbx-01` namespace.