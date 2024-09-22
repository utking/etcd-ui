package v3

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/utking/etcd-ui/internal/helpers/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func createPemPool(certFile, keyFile, caFile string) (*tls.Config, error) {
	if keyFile == "" || certFile == "" {
		return nil, fmt.Errorf("KeyFile and CertFile must both be present[key: %v, cert: %v]", keyFile, certFile)
	}

	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cer},
		InsecureSkipVerify: utils.TLSSkipVerify(), //nolint:gosec //we need to skip for self-signed
	}

	if caFile != "" {
		caCertPool, err := newCertPool([]string{caFile})
		if err != nil {
			return nil, err
		}

		tlsConfig.RootCAs = caCertPool
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		Certificates:       []tls.Certificate{cer},
		InsecureSkipVerify: utils.TLSSkipVerify(), //nolint:gosec //we need to skip for self-signed
	}

	return tlsConfig, nil
}

// NewCertPool creates x509 certPool with provided CA files.
func newCertPool(caFiles []string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, CAFile := range caFiles {
		pemByte, err := os.ReadFile(CAFile)
		if err != nil {
			return nil, err
		}

		for {
			var block *pem.Block

			block, pemByte = pem.Decode(pemByte)
			if block == nil {
				break
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}

			certPool.AddCert(cert)
		}
	}

	return certPool, nil
}

// New, if no errors, returns a new client.
//   - If sslCertFile, sslKeyFile, and sslCAFile are empty
//     the TLS part won't be configured
//   - Username and password can be also empty if auth is disabled for the cluster
func New(
	endpoints []string,
	sslCertFile, sslKeyFile, sslCAFile string,
	username, password string,
) (*Client, error) {
	var (
		err       error
		tlsErr    error
		tlsConfig *tls.Config
	)

	if sslCertFile != "" && sslKeyFile != "" {
		tlsConfig, tlsErr = createPemPool(sslCertFile, sslKeyFile, sslCAFile)
		if tlsErr != nil {
			return nil, tlsErr
		}
	}

	clnt := &Client{
		opTimeout: utils.GetOpTimeout(),
	}

	clnt.client, err = clientv3.New(
		clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: clnt.opTimeout,
			TLS:         tlsConfig,
			Username:    username,
			Password:    password,
		},
	)

	if err != nil {
		return nil, err
	}

	return clnt, nil
}

func (c *Client) Close() {
	c.client.Close()
}
