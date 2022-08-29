# Serverless demo

The demo consist of following binaries:
* bundle-generator - generates serverless connection bundle file from k8s secrets
* cqlsh-generator - generates cqlshrc config file from k8s secrets
* producer - sample application producing rows
* consumer - sample application consuming rows

Producer is responsible for simulating a stock price changes, and consumer is supposed to watch these changes and print out the history change.


## Build

`make build`

## Demo
Repeat for every cluster, choosing different stock name to differentiate:

1. Generate a connection bundle:
    ```
    ./bundle-generator --namespace scyllacluster-0 --ca-secret-name scyllacluster-0-ca --cert-secret-name scyllacluster-0-client-admin --node-domain nodes.scyllacluster-0.apps.eu-west-3.gke.k8s-scylla.dev-test.scylla-operator.scylladb.com > scyllacluster-0-bundle.yaml`
    ```
2. Generate cqlshrc config:
   ```
   ./cqlshrc-generator --namespace scyllacluster-0 --ca-secret-name scyllacluster-0-ca --cert-secret-name scyllacluster-0-client-admin --node-domain nodes.scyllacluster-0.apps.eu-west-3.gke.k8s-scylla.dev-test.scylla-operator.scylladb.com > cqlshrc`
   ```
3. Start producer
    ```
    ./producer --bundle-path scyllacluster-0-bundle.yaml --stock-name BTC`
    ```
4. Start consumer
    ```
    ./consumer --bundle-path scyllacluster-0-bundle.yaml `
    ```
5. Inspect table using cqlshrc
   ```
   cqlsh --cqlshrc="cqlshrc"
   SELECT * FROM stocks.history limit 10;
   ```
