package log

import (
	"encoding/json"
	"io"

	v1 "k8s.io/api/core/v1"

	kubeutil "github.com/lrx0014/log-tools/pkg/util/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/glog"
)

type AuthInfo struct {
	ClusterURL               string `json:"clusterUrl"`
	CertificateAuthorityData []byte `json:"certificateAuthorityData"`
	ClientCertificateData    []byte `json:"clientCertificateData"`
	ClientKeyData            []byte `json:"clientKeyData"`
}

func CreateClient() (kubernetes.Interface, error) {
	k8sCert := `
	{
		"clusterUrl": "",
		"certificateAuthorityData": "",
		"clientCertificateData": "",
		"clientKeyData": ""
	}
	`
	var auth AuthInfo
	json.Unmarshal([]byte(k8sCert), auth)
	kubeClient, err := kubeutil.NewKubeClient(auth.ClusterURL, auth.CertificateAuthorityData, auth.ClientCertificateData, auth.ClientKeyData)
	if err != nil {
		glog.Errorf("unable to create a new kubernetes client with the configuration %v\n: %v", auth, err)
		return nil, err
	}
	return kubeClient, nil
}

func openStream(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

func GetLogs(client kubernetes.Interface, namespace string, podID string, container string) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Timestamps: false,
	}
	logStream, err := openStream(client, namespace, podID, logOptions)
	return logStream, err
}
