package consumer

import (
	"context"

	"github.com/spf13/cobra"
)

func NewConsumer(ctx context.Context) *cobra.Command {
	o := NewConsumerOptions()

	cmd := &cobra.Command{
		Use: "consumer",

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

	cmd.Flags().StringVarP(&o.BundlePath, "bundle-path", "", o.BundlePath, "path to connection bundle")

	return cmd
}
