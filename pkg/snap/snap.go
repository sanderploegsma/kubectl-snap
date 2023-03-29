package snap

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type SnapOptions struct {
	Namespace        string
	SnapAllPods      bool
	SnapOrphanedPods bool
	SnapStoppedPods  bool
}

// Snap deletes half the pods in a Kubernetes cluster and returns a list of pod names that were deleted.
func Snap(options *SnapOptions) (deleted []string, err error) {
	deleted = make([]string, 0)

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	pods, err := clientset.CoreV1().Pods(options.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}

	candidates := make([]v1.Pod, 0)
	for _, pod := range pods.Items {
		if options.shouldSnapPod(pod) {
			candidates = append(candidates, pod)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, pod := range candidates[:len(candidates)/2] {
		err = clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return
		}
		deleted = append(deleted, pod.Name)
	}

	return
}

func (o *SnapOptions) shouldSnapPod(pod v1.Pod) bool {
	if o.SnapAllPods {
		return true
	}

	// Do not snap pods that are not owned by anything unless instructed to
	if !o.SnapOrphanedPods && len(pod.OwnerReferences) == 0 {
		return false
	}

	snappablePhases := []v1.PodPhase{v1.PodRunning, v1.PodPending}

	if o.SnapStoppedPods {
		snappablePhases = append(snappablePhases, v1.PodFailed, v1.PodSucceeded)
	}

	return slices.Contains(snappablePhases, pod.Status.Phase)
}
