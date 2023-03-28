package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sanderploegsma/kubectl-snap/pkg/snap"
	"github.com/spf13/cobra"
)

var (
	// CLI flags
	namespace string
	force     bool
	verbose   bool

	example = `
  # Snap pods in the kube-system namespace
  kubernetes snap -n kube-system
  
  # Snap pods in all namespaces without prompting for confirmation (USE WITH CAUTION)
  kubernetes snap -F
`

	RootCmd = &cobra.Command{
		Use:          "kubectl-snap [flags]",
		Short:        "Delete half of the pods in a Kubernetes cluster or namespace",
		Example:      example,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return execute()
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "If present, the namespace scope for this CLI request")
	RootCmd.Flags().BoolVarP(&force, "force", "F", force, "If true, do not prompt for confirmation")
	RootCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Enable verbose output")
}

func execute() error {
	fmt.Println("When I'm done, half of this cluster will still exist.")
	fmt.Println("Perfectly balanced, as all things should be... I hope they remember you.")

	if !confirm() {
		return fmt.Errorf("aborted by user")
	}

	fmt.Println("Hold tight, little one...")
	fmt.Println()
	fmt.Println("🤌🌟")
	fmt.Println()

	options := &snap.SnapOptions{Namespace: namespace}
	deleted, err := snap.Snap(options)

	if verbose {
		for _, pod := range deleted {
			fmt.Printf("-- Deleted %s\n", pod)
		}
	}

	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Balance has been restored.")

	return nil
}

func confirm() bool {
	if force {
		return true
	}

	selectedNamespace := "every namespace"
	if namespace != "" {
		selectedNamespace = fmt.Sprintf("namespace '%s'", namespace)
	}

	fmt.Printf("\nThis will DELETE half the pods in %s\nAre you sure? (y/N): ", selectedNamespace)

	scanner := bufio.NewScanner(os.Stdin)
	return scanner.Scan() && strings.HasPrefix(strings.ToLower(scanner.Text()), "y")
}
