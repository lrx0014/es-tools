package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewKubeConfig(server string, certificateAuthorityData, clientCertificateData, clientKeyData []byte) (*rest.Config, error) {
	config := clientcmdapi.Config{
		Preferences: *clientcmdapi.NewPreferences(),
		Clusters: map[string]*clientcmdapi.Cluster{
			"kubernetes": &clientcmdapi.Cluster{
				Server: server,
				CertificateAuthorityData: certificateAuthorityData,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"kubernetes": &clientcmdapi.AuthInfo{
				ClientCertificateData: clientCertificateData,
				ClientKeyData:         clientKeyData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"kubernetes": &clientcmdapi.Context{
				Cluster:  "kubernetes",
				AuthInfo: "kubernetes",
			},
		},
		CurrentContext: "kubernetes",
	}

	return clientcmd.NewNonInteractiveClientConfig(config, config.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
}

func NewKubeClient(server string, certificateAuthorityData, clientCertificateData, clientKeyData []byte) (kubernetes.Interface, error) {
	config, err := NewKubeConfig(server, certificateAuthorityData, clientCertificateData, clientKeyData)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
