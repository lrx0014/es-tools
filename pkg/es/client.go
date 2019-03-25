package es

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

type ConnInfo struct {
	Cluster string
	CA      string
	Cert    string
	Key     string
}

func CreateESClient(connInfo ConnInfo) (*elastic.Client, error) {
	cert, err := tls.LoadX509KeyPair(connInfo.Cert, connInfo.Key)
	if err != nil {
		log.Fatalf("Unable to create tls: %+v", err)
		return nil, err
	}
	caCert, err := ioutil.ReadFile(connInfo.CA)
	if err != nil {
		log.Fatalf("Unable to read CA: %+v", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}
	esClient, err := elastic.NewClient(
		elastic.SetHttpClient(httpClient),
		elastic.SetURL(connInfo.Cluster),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetSniff(false))

	if err != nil {
		log.Fatalf("Unable to create es client: %+v", err)
	}

	return esClient, nil
}

func getConfigFromEnv() ConnInfo {
	return ConnInfo{
		Cluster: os.Getenv("ES_CLUSTER"),
		Cert:    os.Getenv("ES_CERT_PATH"),
		Key:     os.Getenv("ES_KEY_PATH"),
		CA:      os.Getenv("ES_CA_PATH"),
	}
}
