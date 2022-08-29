package bundle_generator

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/zimnx/serverlessExample/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

type BundleGeneratorOptions struct {
	NodeDomain     string
	CASecretName   string
	CertSecretName string
	Namespace      string
	KubeConfigPath string
	Username       string
	Password       string

	RestConfig *restclient.Config
}

func NewBundleGeneratorOptions() *BundleGeneratorOptions {
	return &BundleGeneratorOptions{
		Username:       "cassandra",
		Password:       "cassandra",
		KubeConfigPath: path.Join(os.Getenv("HOME"), ".kube/config"),
	}
}

func (o *BundleGeneratorOptions) Validate(args []string) error {
	var errs []error

	if len(o.NodeDomain) == 0 {
		errs = append(errs, fmt.Errorf("node domain cannot be empty"))
	}

	if len(o.CASecretName) == 0 {
		errs = append(errs, fmt.Errorf("ca secret name cannot be empty"))
	}

	if len(o.CertSecretName) == 0 {
		errs = append(errs, fmt.Errorf("cert secret name cannot be empty"))
	}

	if len(o.Namespace) == 0 {
		errs = append(errs, fmt.Errorf("namespace cannot be empty"))
	}

	return errors.NewAggregate(errs)
}

func (o *BundleGeneratorOptions) Complete(args []string) error {
	var err error

	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	loader.ExplicitPath = o.KubeConfigPath
	o.RestConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loader,
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return fmt.Errorf("can't create client config: %w", err)
	}

	return nil
}

func (o *BundleGeneratorOptions) Run(ctx context.Context) error {
	client, err := kubernetes.NewForConfig(o.RestConfig)
	if err != nil {
		return err
	}

	caSecret, err := client.CoreV1().Secrets(o.Namespace).Get(ctx, o.CASecretName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can't get CA secret: %w", err)
	}

	clientCertSecret, err := client.CoreV1().Secrets(o.Namespace).Get(ctx, o.CertSecretName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can't get client cert secret: %w", err)
	}

	clientKeyData := clientCertSecret.Data["tls.key"]
	clientCertData := clientCertSecret.Data["tls.crt"]
	caCertData := caSecret.Data["tls.crt"]

	bundle := model.ConnectionConfig{
		Datacenters: map[string]*model.Datacenter{
			"default": {
				CertificateAuthorityData: caCertData,
				Server:                   fmt.Sprintf("any.%s", o.NodeDomain),
				NodeDomain:               o.NodeDomain,
			},
		},
		AuthInfos: map[string]*model.AuthInfo{
			"default": {
				ClientCertificateData: clientCertData,
				ClientKeyData:         clientKeyData,
				Username:              o.Username,
				Password:              o.Password,
			},
		},
		Contexts: map[string]*model.Context{
			"default": {
				DatacenterName: "default",
				AuthInfoName:   "default",
			},
		},
		CurrentContext: "default",
	}

	buf, err := yaml.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("can't marshal connection bundle: %w", err)
	}

	if _, err := os.Stdout.Write(buf); err != nil {
		return fmt.Errorf("can't write to stdout: %w", err)
	}

	return nil
}
