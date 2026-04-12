# Hint for Exercise 02: Labels & Selectors

## 💡 Creating pods with labels

Use the `--labels` flag with `kubectl run`:

```bash
kbx kubectl run <pod-name> --image=<image> --labels="key1=value1,key2=value2" -n <namespace>
```

## 📋 Commands for this exercise

### Create the frontend pod

```bash
kbx kubectl run frontend --image=nginx:latest --labels="app=web,tier=frontend,env=prod" -n kbx-02
```

### Create the backend pod

```bash
kbx kubectl run backend --image=nginx:latest --labels="app=api,tier=backend,env=prod" -n kbx-02
```

### Create the worker pod

```bash
kbx kubectl run worker --image=nginx:latest --labels="app=processor,tier=backend,env=dev" -n kbx-02
```

## 🔍 Verify your work

```bash
# List all pods with labels
kbx kubectl get pods -n kbx-02 --show-labels

# Query by tier
kbx kubectl get pods -l tier=frontend -n kbx-02
kbx kubectl get pods -l tier=backend -n kbx-02

# Query by environment
kbx kubectl get pods -l env=prod -n kbx-02
kbx kubectl get pods -l env=dev -n kbx-02
```

## ✅ Check your solution

```bash
kbx check 02
```

---

**Key concept:** Labels are key-value pairs that help organize and select Kubernetes objects. Selectors let you filter resources based on their labels.