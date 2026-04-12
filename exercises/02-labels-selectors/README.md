# Exercise 02: Labels & Selectors

**Type:** Build  
**Difficulty:** Beginner

## 🎯 Goal

Learn how to use labels to organize pods and selectors to query them.

## 📝 Tasks

Create three pods in the `kbx-02` namespace with the following specifications:

### Pod 1: frontend
- **Name:** `frontend`
- **Image:** `nginx:latest`
- **Labels:**
  - `app=web`
  - `tier=frontend`
  - `env=prod`

### Pod 2: backend
- **Name:** `backend`
- **Image:** `nginx:latest`
- **Labels:**
  - `app=api`
  - `tier=backend`
  - `env=prod`

### Pod 3: worker
- **Name:** `worker`
- **Image:** `nginx:latest`
- **Labels:**
  - `app=processor`
  - `tier=backend`
  - `env=dev`

## ✅ Success Criteria

- All three pods exist in the `kbx-02` namespace
- All pods are in `Running` state
- Each pod has the correct labels assigned
- You can query pods using label selectors

## 💡 Tips

- Use `kbx kubectl run` with `--labels` flag to create pods with labels
- Use `kbx kubectl get pods -l <selector>` to query by labels
- Use `kbx kubectl describe pod <name>` to verify labels

### Example label selector queries

```bash
# Get all frontend tier pods
kbx kubectl get pods -l tier=frontend -n kbx-02

# Get all prod environment pods
kbx kubectl get pods -l env=prod -n kbx-02

# Get all backend tier pods
kbx kubectl get pods -l tier=backend -n kbx-02
```

## 🆘 Stuck?

Run `kbx hint 02` for a hint!