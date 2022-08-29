package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql/scyllacloud"
	scyllacdc "github.com/scylladb/scylla-cdc-go"
	"k8s.io/apimachinery/pkg/util/errors"
)

type ConsumerOptions struct {
	BundlePath string
}

func NewConsumerOptions() *ConsumerOptions {
	return &ConsumerOptions{}
}

func (o *ConsumerOptions) Validate(args []string) error {
	var errs []error

	if len(o.BundlePath) == 0 {
		errs = append(errs, fmt.Errorf("bundle path cannot be empty"))
	}

	return errors.NewAggregate(errs)
}

func (o *ConsumerOptions) Complete(args []string) error {
	return nil
}

func (o *ConsumerOptions) Run(ctx context.Context) error {
	cluster, err := scyllacloud.NewCloudCluster(o.BundlePath)
	if err != nil {
		return err
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	readerConfig := &scyllacdc.ReaderConfig{
		Session:               session,
		TableNames:            []string{"stocks.history"},
		ChangeConsumerFactory: scyllacdc.MakeChangeConsumerFactoryFromFunc(printer),
	}

	reader, err := scyllacdc.NewReader(ctx, readerConfig)
	if err := reader.Run(ctx); err != nil {
		return err
	}

	return nil
}

func printer(ctx context.Context, tableName string, c scyllacdc.Change) error {
	for _, changeRow := range c.Delta {
		nameRaw, ok := changeRow.GetValue("name")
		if !ok {
			return fmt.Errorf("can't get name column")
		}
		priceRaw, ok := changeRow.GetValue("price")
		if !ok {
			return fmt.Errorf("can't get price column")
		}
		timeRaw, ok := changeRow.GetValue("time")
		if !ok {
			return fmt.Errorf("can't get time column")
		}

		name := nameRaw.(*string)
		price := priceRaw.(*int)
		timestamp := timeRaw.(*time.Time)

		fmt.Printf("%s: name: %s, price %d\n", *timestamp, *name, *price)
	}

	return nil
}
