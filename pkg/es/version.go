package es

import (
	"log"
)

func PrintVersion() error {
	connInfo := getConfigFromEnv()
	esClient, err := CreateESClient(connInfo)
	if err != nil {
		log.Fatalf("Unable to create es client: %+v", err)
		return err
	}

	version, err := esClient.ElasticsearchVersion(connInfo.Cluster)
	if err != nil {
		log.Fatalf("Unable to get es version: %+v", err)
		return err
	}

	log.Printf("Connected to ES, the version is: %s \n\n", version)

	return nil
}
