package model

type ConnectionConfig struct {
	// Kind is a string value representing the REST resource this object   represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// In CamelCase.
	// +optional
	Kind string `yaml:"kind,omitempty"`
	// APIVersion defines the versioned schema of this representation of an   object.
	// Servers should convert recognized schemas to the latest internal value,   and may reject unrecognized values.
	// +optional
	APIVersion string `yaml:"apiVersion,omitempty"`
	// Datacenters is a map of referenceable names to datacenter configs.
	Datacenters map[string]*Datacenter `json:"datacenters"`
	// AuthInfos is a map of referenceable names to authentication configs.
	AuthInfos map[string]*AuthInfo `json:"authInfos"`
	// Contexts is a map of referenceable names to context configs.
	Contexts map[string]*Context `json:"contexts"`
	// CurrentContext is the name of the context that you would like to use by default.
	CurrentContext string `json:"currentContext"`
	// Parameters is a struct containing common driver configuration parameters.
	// +optional
	Parameters *Parameters `json:"parameters,omitempty"`
}

type AuthInfo struct {
	// ClientCertificateData contains base64 encoded PEM-encoded data from a client cert file for TLS. Overrides ClientCertificatePath.
	// Either ClientCertificateData or ClientCertificatePath must be provided.
	// +optional
	ClientCertificateData []byte `json:"clientCertificateData,omitempty"`
	// ClientCertificatePath is the path to a client cert file for TLS.
	// +optional
	ClientCertificatePath string `json:"clientCertificatePath,omitempty"`
	// ClientKeyData contains base64 encoded PEM-encoded data from a client key file for TLS. Overrides ClientKeyPath.
	// Either ClientKeyData or ClientKeyPath must be provided.
	// +optional
	ClientKeyData []byte `json:"clientKeyData,omitempty"`
	// ClientKeyPath is the path to a client key file for TLS.
	// +optional
	ClientKeyPath string `json:"clientKeyPath,omitempty"`
	// Username is the username for basic authentication to the Scylla cluster.
	// +optional
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication to the Scylla cluster.
	// +optional
	Password string `json:"password,omitempty"`
}

type Datacenter struct {
	// CertificateAuthorityPath is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthorityPath string `json:"certificateAuthorityPath,omitempty"`
	// CertificateAuthorityData contains base64 encoded PEM-encoded certificate authority certificates. Overrides CertificateAuthorityPath.
	// Either CertificateAuthorityData or CertificateAuthorityPath must be provided.
	// +optional
	CertificateAuthorityData []byte `json:"certificateAuthorityData,omitempty"`
	// Server is the initial contact point of the Scylla cluster.
	// Example: https://hostname:port
	Server string `json:"server"`
	// TLSServerName is used to check server certificates. If TLSServerName is empty, the hostname used to contact the server is used.
	// +optional
	TLSServerName string `json:"tlsServerName,omitempty"`
	// NodeDomain the domain suffix that is concatenated with host_id of the node driver wants to connect to.
	// Example: host_id.<nodeDomain>
	NodeDomain string `json:"nodeDomain"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate. This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `json:"insecureSkipTlsVerify,omitempty"`
	// ProxyURL is the URL to the proxy to be used for all requests made by this
	// client. URLs with "http", "https", and "socks5" schemes are supported. If
	// this configuration is not provided or the empty string, the client
	// attempts to construct a proxy configuration from http_proxy and
	// https_proxy environment variables. If these environment variables are not
	// set, the client does not attempt to proxy requests.
	// +optional
	ProxyURL string `json:"proxyUrl,omitempty"`
}

type Context struct {
	// DatacenterName is the name of the datacenter for this context.
	DatacenterName string `json:"datacenterName"`
	// AuthInfoName is the name of the authInfo for this context.
	AuthInfoName string `json:"authInfoName"`
}

type Parameters struct {
	// DefaultConsistency is the default consistency level used for queries.
	// +optional
	DefaultConsistency Consistency `json:"defaultConsistency,omitempty"`
	// DefaultSerialConsistency is the default consistency level for the serial   part of queries.
	// +optional
	DefaultSerialConsistency Consistency `json:"defaultSerialConsistency,omitempty"`
}

type Consistency string
