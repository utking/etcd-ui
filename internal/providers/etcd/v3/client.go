package v3

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

// AuthEnabled will use a random username/password and analyze the reposne
// of reading a random key. Seems to be the best way to check if auth is en/dis-abled.
// We still need to know it we use TLS or not though.
func CheckAuthEnabled(
	endpoints []string,
	sslCertFile, sslKeyFile, sslCAFile string,
) (bool, error) {
	var (
		username = fmt.Sprintf("tmp-%x", time.Now().Unix())
		password = fmt.Sprintf("pwd-%x", time.Now().Unix())
		key      = fmt.Sprintf("/fake/test/auth/key-%x", time.Now().Unix())
	)

	tmpClient, clntErr := New(endpoints, sslCertFile, sslKeyFile, sslCAFile, username, password)
	if clntErr != nil {
		switch {
		case strings.Contains(clntErr.Error(), "authentication failed, invalid user ID or password"):
			return true, nil
		case strings.Contains(clntErr.Error(), "authentication is not enabled"):
			return false, nil
		default:
			// We will have to implement it when it's found
			log.Printf("error %q is not implemented. Please report it\n", clntErr.Error())

			return false, clntErr
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.GetOpTimeout())
	defer cancel()

	getKeyResp, err := tmpClient.client.Get(ctx, key)

	if err == nil {
		// The chance to guess name/pass/key is too slim to get no error unless there is no auth
		// The response should contain no keys
		if getKeyResp != nil && len(getKeyResp.Kvs) == 0 {
			return false, nil
		}

		// Either the response is nil of there are keys in the response. Edge case to investigate later
		return false, fmt.Errorf("this should not have happened, retry please")
	}

	// We should not have gone here, so let's check the error again and see. Investigate when happend
	switch {
	case strings.Contains(err.Error(), "authentication failed, invalid user ID or password"):
		return true, nil
	case strings.Contains(err.Error(), "authentication is not enabled"):
		return false, nil
	default:
		// We will have to implement it when it's found
		log.Printf("error %q is not implemented. Please report it\n", err.Error())

		return false, err
	}
}
