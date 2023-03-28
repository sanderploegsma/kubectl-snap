# kubectl-snap

Perfectly balanced, as all things should be...

## About

`kubectl-snap` is a plugin for `kubectl` that deletes half of the pods in your Kubernetes cluster or namespace. 
Inspired by [honk-ci/kubectl-snap](https://github.com/honk-ci/kubectl-snap), but rewritten in Go based on [kubernetes/sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin).

## Install

### Homebrew

On macOS and Linux you can install using [Homebrew](https://brew.sh):

    brew install sanderploegsma/tap/kubectl-snap

### Manual

Grab the latest release for your platform and place the binaries somewhere in your `PATH`, like `/usr/local/bin`.

## Usage

```
Usage:
  kubectl-snap [flags]

Examples:

  # Snap pods in the kube-system namespace
  kubectl snap -n kube-system
  
  # Snap pods in all namespaces without prompting for confirmation (USE WITH CAUTION)
  kubectl snap -F


Flags:
  -F, --force              If true, do not prompt for confirmation
  -h, --help               help for kubectl-snap
  -n, --namespace string   If present, the namespace scope for this CLI request
  -v, --verbose            Enable verbose output
      --version            version for kubectl-snap
```
