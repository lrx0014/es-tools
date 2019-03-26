package es

import (
	"crypto/tls"
	"crypto/x509"
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

	var pem []byte
	pem = []byte(connInfo.Cert + connInfo.Key)
	cert, err := tls.X509KeyPair(pem, pem)
	if err != nil {
		log.Fatalf("Unable to create tls: %+v", err)
		return nil, err
	}
	caCert := []byte(connInfo.CA)

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

	// esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))

	if err != nil {
		log.Fatalf("Unable to create es client: %+v", err)
	}

	return esClient, nil
}

func getConfigFromEnv() ConnInfo {
	// TODO: Change here to support multi zones
	return ConnInfo{
		Cluster: os.Getenv("ES_CLUSTER"),
		Cert:    os.Getenv("ES_CERT_PATH"),
		Key:     os.Getenv("ES_KEY_PATH"),
		CA:      os.Getenv("ES_CA_PATH"),
	}
}
