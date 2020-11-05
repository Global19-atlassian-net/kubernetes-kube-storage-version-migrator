package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	migrationclient "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/clients/clientset"
	"github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/trigger"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	triggerUserAgent = "storage-version-migration-trigger"
)

var (
	kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file specifying the apiserver instance. If unspecified, fallback to in-cluster configuration")
)

func NewTriggerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "kube-storage-migrator-trigger",
		Long: `The Kubernetes storage migrator triggering controller
		detects storage version changes and creates migration requests.
		It also records the status of the storage via the storageState
		API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Run(context.TODO()); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
}

func Run(ctx context.Context) error {
	var err error
	var config *rest.Config
	if *kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
		if err != nil {
			log.Fatalf("Error initializing client config: %v for kubeconfig: %v", err.Error(), *kubeconfigPath)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}
	migration, err := migrationclient.NewForConfig(rest.AddUserAgent(config, triggerUserAgent))
	if err != nil {
		return err
	}
	c := trigger.NewMigrationTrigger(migration)
	c.Run(ctx)
	panic("unreachable")
}
