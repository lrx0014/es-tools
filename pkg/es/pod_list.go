package es

import (
	"context"
	"encoding/json"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

type PodInfo struct {
	PodName string `json:"podname"`
	Status  string `json:"status"`
}

type PodList struct {
	List []PodInfo `json:"list"`
}

func GetPodListFromES(namespace string, container string) (*PodList, error) {
	connInfo := getConfigFromEnv()
	esClient, err := CreateESClient(connInfo)
	if err != nil {
		log.Fatalf("Unable to create client: %+v", err)
		return nil, err
	}

	collapse := elastic.NewCollapseBuilder("kubernetes.pod_name.keyword")

	q := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("kubernetes.namespace_name", namespace)).
		Must(elastic.NewTermQuery("kubernetes.container_name", container))

	res, err := esClient.Search().
		Query(q).
		Collapse(collapse).
		Do(context.Background())

	if err != nil {
		src, err := q.Source()
		if err != nil {
			log.Fatalln(err)
		}
		data, err := json.Marshal(src)
		if err != nil {
			log.Fatalf("marshaling to JSON failed: %v", err)
		}
		got := string(data)
		log.Fatalln(got)
		return nil, err
	}

	podlist := new(PodList)

	var pods []PodInfo

	for _, line := range res.Hits.Hits {
		podname := line.Fields["kubernetes.pod_name.keyword"].(string)
		podinfo := PodInfo{
			PodName: podname,
			Status:  "Log_Persistent_in_ES",
		}
		pods = append(pods, podinfo)
	}

	podlist.List = pods

	return podlist, nil
}
