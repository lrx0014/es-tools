package es

import (
	"context"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodInfo struct {
	PodName string `json:"podname"`
	Status  string `json:"status"`
}

type PodList struct {
	List []PodInfo `json:"list"`
}

func GetPodList(k8sClient kubernetes.Interface, namespace string, container string) (*PodList, error) {
	connInfo := getConfigFromEnv()
	esClient, err := CreateESClient(connInfo)
	if err != nil {
		log.Fatalf("Unable to create client: %+v", err)
		return nil, err
	}

	collapse := elastic.NewCollapseBuilder("kubernetes.pod_name")

	q := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("kubernetes.namespace_name", namespace)).
		Must(elastic.NewTermQuery("kubernetes.container_name", container))

	res, err := esClient.Search().
		Query(q).
		Sort("@timestamp", false).
		Collapse(collapse).
		Do(context.Background())

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	podlist := new(PodList)

	var pods []PodInfo

	for _, line := range res.Hits.Hits {
		podname := line.Fields["kubernetes.pod_name"].([]interface{})[0].(string)
		podinfo := PodInfo{
			PodName: podname,
			Status:  "Log_Persistent_in_ES",
		}
		pods = append(pods, podinfo)
	}

	podlist.List = pods

	updatePodStatus(k8sClient, namespace, podlist)

	return podlist, nil
}

func updatePodStatus(k8sClient kubernetes.Interface, namespace string, podlist *PodList) error {
	for i := 0; i < len(podlist.List); i++ {
		_, err := k8sClient.CoreV1().Pods(namespace).Get(podlist.List[i].PodName, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("Unable to get pod info from kubernetes: %v \n", err)
			return err
		}
		podlist.List[i].Status = "Running"
	}
	return nil
}

func convertInterfaceArrayToStringArray(input interface{}) (output []string) {
	for _, temp := range input.([]interface{}) {
		output = append(output, temp.(string))
	}
	return output
}
