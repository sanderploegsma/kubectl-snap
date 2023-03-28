# kubectl-snap

Perfectly balanced, as all things should be...

## About

`kubectl-snap` is a plugin for `kubectl` that deletes half of the pods in your Kubernetes cluster or namespace. 
Inspired by https://github.com/honk-ci/kubectl-snap, but rewritten in Go.

## Install

### Homebrew

On macOS and Linux you can install using [Homebrew](https://brew.sh):

    brew install sanderploegsma/tap/kubectl-snap

### Manual

Grab the latest release for your platform and place the binaries somewhere in your `PATH`, like `/usr/local/bin`.

## Usage

```
Usage:
  kubectl snap [flags]

Examples:

  # Snap pods in the kube-system namespace
  kubectl snap -n kube-system

  # Snap pods in all namespaces without prompting for confirmation (USE WITH CAUTION)
  kubectl snap --force


Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/sander/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -f, --force                          If true, do not prompt for confirmation
  -h, --help                           help for kubectl-snap
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
  -v, --verbose count
```