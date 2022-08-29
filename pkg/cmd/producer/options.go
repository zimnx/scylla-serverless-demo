package producer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gocql/gocql/scyllacloud"
	"k8s.io/apimachinery/pkg/util/errors"
)

type ProducerOptions struct {
	BundlePath string
	StockName  string
}

func NewProducerOptions() *ProducerOptions {
	return &ProducerOptions{}
}

func (o *ProducerOptions) Validate(args []string) error {
	var errs []error

	if len(o.BundlePath) == 0 {
		errs = append(errs, fmt.Errorf("bundle path cannot be empty"))
	}

	if len(o.StockName) == 0 {
		errs = append(errs, fmt.Errorf("stock name cannot be empty"))
	}

	return errors.NewAggregate(errs)
}

func (o *ProducerOptions) Complete(args []string) error {
	return nil
}

func (o *ProducerOptions) Run(ctx context.Context) error {
	cluster, err := scyllacloud.NewCloudCluster(o.BundlePath)
	if err != nil {
		return err
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS stocks WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec()
	if err != nil {
		return fmt.Errorf("can't create keyspace: %w", err)
	}

	err = session.Query("CREATE TABLE IF NOT EXISTS stocks.history (name text, time TIMESTAMP, price int, PRIMARY KEY (name, time)) WITH cdc = {'enabled':true}").Exec()
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	insertQ := session.Query("INSERT INTO stocks.history (name, time, price) VALUES (?, ?, ?)")

	lastPrice := rand.Intn(1000) + 100

	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			sign := rand.Intn(2)
			if sign == 0 {
				sign = -1
			}
			lastPrice += sign * int(2*rand.Float32())
			if err := insertQ.Bind(o.StockName, time.Now(), lastPrice).Exec(); err != nil {
				fmt.Println("cannot insert row: %w", err)
			}
			fmt.Println("Current price of", o.StockName, lastPrice)
		}
	}
}
