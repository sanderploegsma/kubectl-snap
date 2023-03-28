package cmd

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var commandUse = "%[1]s-snap [flags]"
var commandExample = `
  # Snap pods in the kube-system namespace
  %[1]s snap -n kube-system
  
  # Snap pods in all namespaces without prompting for confirmation (USE WITH CAUTION)
  %[1]s snap -f
`

type SnapOptions struct {
	force     bool
	verbosity int

	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func NewSnapOptions(streams genericclioptions.IOStreams) *SnapOptions {
	return &SnapOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func NewSnapCmd(streams genericclioptions.IOStreams, version string) *cobra.Command {
	o := NewSnapOptions(streams)

	cmd := &cobra.Command{
		Use:          fmt.Sprintf(commandUse, "kubectl"),
		Short:        "Delete half of the pods in a namespace or cluster",
		Example:      fmt.Sprintf(commandExample, "kubectl"),
		Version:      version,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().BoolVarP(&o.force, "force", "f", o.force, "If true, do not prompt for confirmation")
	cmd.Flags().CountVarP(&o.verbosity, "verbose", "v", "")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *SnapOptions) Run() error {
	fmt.Fprintln(o.IOStreams.Out, "When I'm done, half of this cluster will still exist.")
	fmt.Fprintln(o.IOStreams.Out, "Perfectly balanced, as all things should be... I hope they remember you.")

	if !o.confirm() {
		return fmt.Errorf("aborted by user")
	}

	fmt.Fprintln(o.IOStreams.Out, "Hold tight, little one...")

	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(*o.configFlags.Namespace).List(context.TODO(), metav1.ListOptions{})
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

	fmt.Fprintf(o.IOStreams.Out, "\n🤌🌟\n\n")

	for _, pod := range candidates[:len(candidates)/2] {
		if o.verbosity > 0 {
			fmt.Fprintf(o.IOStreams.Out, "-- Deleting pod %s\n", pod.Name)
		}
		err = clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(o.IOStreams.Out, "Balance has been restored.")

	return nil
}

func (o *SnapOptions) confirm() bool {
	if o.force {
		return true
	}

	namespace := "every namespace"
	if *o.configFlags.Namespace != "" {
		namespace = fmt.Sprintf("namespace '%s'", *o.configFlags.Namespace)
	}

	fmt.Fprintf(o.IOStreams.Out, "\nThis will DELETE half the pods in %s\n", namespace)
	fmt.Fprint(o.IOStreams.Out, "Are you sure? (y/N): ")
	scanner := bufio.NewScanner(o.IOStreams.In)
	return scanner.Scan() && strings.HasPrefix(strings.ToLower(scanner.Text()), "y")
}

func (o *SnapOptions) shouldSnapPod(pod v1.Pod) bool {
	// Do not snap pods that are not owned by anything
	if len(pod.OwnerReferences) == 0 {
		return false
	}

	return pod.Status.Phase == v1.PodRunning
}
