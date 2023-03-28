package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var commandUse = "%[1]s-snap [flags]"
var commandExample = `
  # Snap pods in the kube-system namespace
  %[1]s snap -n kube-system
  
  # Snap pods in all namespaces without prompting for confirmation (USE WITH CAUTION)
  %[1]s snap -F
`

type SnapOptions struct {
	Namespace string
	Force     bool
	Verbosity int

	In  io.Reader
	Out io.Writer
}

func NewSnapCmd(in io.Reader, out io.Writer, version string) *cobra.Command {
	o := &SnapOptions{
		In:  in,
		Out: out,
	}

	cmd := &cobra.Command{
		Use:          fmt.Sprintf(commandUse, "kubectl"),
		Short:        "Delete half of the pods in a Kubernetes cluster or namespace",
		Example:      fmt.Sprintf(commandExample, "kubectl"),
		Version:      version,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "If present, the namespace scope for this CLI request")
	cmd.Flags().BoolVarP(&o.Force, "force", "F", o.Force, "If true, do not prompt for confirmation")
	cmd.Flags().CountVarP(&o.Verbosity, "verbose", "v", "Enable verbose output")

	return cmd
}

func (o *SnapOptions) Run() error {
	fmt.Fprintln(o.Out, "When I'm done, half of this cluster will still exist.")
	fmt.Fprintln(o.Out, "Perfectly balanced, as all things should be... I hope they remember you.")

	if !o.confirm() {
		return fmt.Errorf("aborted by user")
	}

	fmt.Fprintln(o.Out, "Hold tight, little one...")

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(o.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	candidates := make([]v1.Pod, 0)
	for _, pod := range pods.Items {
		if o.shouldSnapPod(pod) {
			candidates = append(candidates, pod)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	fmt.Fprintf(o.Out, "\n🤌🌟\n\n")

	for _, pod := range candidates[:len(candidates)/2] {
		if o.Verbosity > 0 {
			fmt.Fprintf(o.Out, "-- Deleting pod %s\n", pod.Name)
		}
		err = clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(o.Out, "Balance has been restored.")

	return nil
}

func (o *SnapOptions) confirm() bool {
	if o.Force {
		return true
	}

	namespace := "every namespace"
	if o.Namespace != "" {
		namespace = fmt.Sprintf("namespace '%s'", o.Namespace)
	}

	fmt.Fprintf(o.Out, "\nThis will DELETE half the pods in %s\n", namespace)
	fmt.Fprint(o.Out, "Are you sure? (y/N): ")
	scanner := bufio.NewScanner(o.In)
	return scanner.Scan() && strings.HasPrefix(strings.ToLower(scanner.Text()), "y")
}

func (o *SnapOptions) shouldSnapPod(pod v1.Pod) bool {
	// Do not snap pods that are not owned by anything
	if len(pod.OwnerReferences) == 0 {
		return false
	}

	return pod.Status.Phase == v1.PodRunning
}
