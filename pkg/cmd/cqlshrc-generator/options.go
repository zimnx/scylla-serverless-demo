package cqlshrc_generator

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type CqlshrcGeneratorOptions struct {
	NodeDomain     string
	CASecretName   string
	CertSecretName string
	Namespace      string
	KubeConfigPath string
	Username       string
	Password       string

	RestConfig *restclient.Config
}

func NewCqlshrcGeneratorOptions() *CqlshrcGeneratorOptions {
	return &CqlshrcGeneratorOptions{
		Username:       "cassandra",
		Password:       "cassandra",
		KubeConfigPath: path.Join(os.Getenv("HOME"), ".kube/config"),
	}
}

func (o *CqlshrcGeneratorOptions) Validate(args []string) error {
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

func (o *CqlshrcGeneratorOptions) Complete(args []string) error {
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

func (o *CqlshrcGeneratorOptions) Run(ctx context.Context) error {
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

	if err := ioutil.WriteFile("ca.crt", caCertData, 0666); err != nil {
		return fmt.Errorf("can't write ca cert file: %w", err)
	}

	if err := ioutil.WriteFile("client.crt", clientCertData, 0666); err != nil {
		return fmt.Errorf("can't write client cert file: %w", err)
	}

	if err := ioutil.WriteFile("client.key", clientKeyData, 0666); err != nil {
		return fmt.Errorf("can't write client key file: %w", err)
	}

	type params struct {
		Username string
		Password string
		Hostname string
		Port     int
		CAPath   string
		KeyPath  string
		CertPath string
	}

	t, err := template.New("").Parse(
		`
[authentication]
username = {{ .Username }}
password = {{ .Password }}

[connection]
hostname = {{ .Hostname }}
port = {{ .Port }}

[ssl]
certfile = {{ .CAPath }}
userkey = {{ .KeyPath }} 
usercert = {{ .CertPath }}
validate = true 
`)

	if err != nil {
		return fmt.Errorf("can't parse cqlshrc template: %w", err)
	}

	err = t.Execute(os.Stdout, params{
		Username: o.Username,
		Password: o.Password,
		Hostname: fmt.Sprintf("any.%s", o.NodeDomain),
		Port:     443,
		CAPath:   "ca.crt",
		KeyPath:  "client.key",
		CertPath: "client.crt",
	})

	return nil
}
