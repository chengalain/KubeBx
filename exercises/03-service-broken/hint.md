# Hint for Exercise 03: Service broken

## Hint 1

Look at the labels on the pod and compare them to the selector defined in the Service.

```bash
kbx kubectl get pods -n kbx-03 --show-labels
kbx kubectl describe svc nginx-service -n kbx-03
```

A Service routes traffic only to pods whose labels match its `selector`. If the selector and the pod labels don't match, the Service has no endpoints and traffic goes nowhere.

## Hint 2

The Service selector contains a typo. Check the `Selector` field in `describe svc` output and compare it to the pod's `app` label.

Fix it with:

```bash
kbx kubectl edit svc nginx-service -n kbx-03
```

Update the selector from `app: nginxx` to `app: nginx`, save, then verify:

```bash
kbx kubectl get endpoints nginx-service -n kbx-03
```

You should see an endpoint address instead of `<none>`.

## Check your solution

```bash
kbx check 03
```
