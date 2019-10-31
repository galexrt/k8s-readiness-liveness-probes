# k8s-readiness-liveness-probes

Kubernetes Readiness and Liveness Probe demo application.

## Files

```console
├── cmd
│   └── application
│       └── main.go
├── Dockerfile
├── LICENSE
├── Makefile
├── manifests
│   ├── deployment.yaml
│   └── service.yaml
└── README.md - This file.
```

## Problem

An application that receives a signal (e.g., when Kubernetes / OpenShift) wants to terminate the Pod, instantly exits / terminates because of the signal.

**Why is that a problem?**

It takes a short amount of time till this change to a Pod is propagated across the whole Kubernetes / OpenShift cluster.
This means that in the "worst" case, user requests are going against an application Pod which is currently being terminated.

## Probes

### Readiness Probe

"Is the application ready to accept and serve requests?"

**Examples**:

* Listener started and ready to process requests.
  * Some applications open the listener before everything is initialized.
* Application is able to process requests.

### Liveness Probe

"Is the application able to handle requests?"

**Examples**:

* Not in a Deadlock.
  * Listener still working.
* Application logic still working.

## Graceful Service Degradation

Application should gracefully handle, e.g., a database which is not available right now.
An application should not panick and exit because of that.
