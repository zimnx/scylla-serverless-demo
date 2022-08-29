package cqlshrc_generator

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCqlshrcGenerator(ctx context.Context) *cobra.Command {
	o := NewCqlshrcGeneratorOptions()

	cmd := &cobra.Command{
		Use: "cqlshrc-generator",

		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			err = o.Complete(args)
			if err != nil {
				return err
			}

			err = o.Run(ctx)
			if err != nil {
				return err
			}

			return nil
		},

		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().StringVarP(&o.NodeDomain, "node-domain", "", o.NodeDomain, "domain associated with cluster")
	cmd.Flags().StringVarP(&o.CASecretName, "ca-secret-name", "", o.CASecretName, "name of a secret containing trusted certificate authority")
	cmd.Flags().StringVarP(&o.CertSecretName, "cert-secret-name", "", o.CertSecretName, "name of a secret containing client certificate")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "", o.Namespace, "namespace")
	cmd.Flags().StringVarP(&o.KubeConfigPath, "kube-config", "", o.KubeConfigPath, "path to kube config")
	cmd.Flags().StringVarP(&o.Username, "username", "", o.Username, "username used for CQL authentication")
	cmd.Flags().StringVarP(&o.Password, "password", "", o.Password, "password user for CQL authentication")

	return cmd
}
