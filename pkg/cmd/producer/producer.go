package producer

import (
	"context"

	"github.com/spf13/cobra"
)

func NewProducer(ctx context.Context) *cobra.Command {
	o := NewProducerOptions()

	cmd := &cobra.Command{
		Use: "producer",

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
	cmd.Flags().StringVarP(&o.StockName, "stock-name", "", o.StockName, "name of the stock to produce")

	return cmd
}
