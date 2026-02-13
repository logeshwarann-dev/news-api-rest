package test

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func forward(t testing.TB, svcNamespace string, svcName string, containerPort int) string {
	t.Helper()

	// Load client config
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingClientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	clientConfig, err := loadingClientConfig.ClientConfig()
	require.NoError(t, err)

	//Use spdy transport and upgrader
	transport, upgrader, err := spdy.RoundTripperFor(clientConfig)
	require.NoError(t, err)

	//Provide buffers for errors and outputs
	var buffErr, buffOut bytes.Buffer
	stopChan := make(chan struct{}, 1)
	readChan := make(chan struct{}, 1)

	//Get a free port on local machine
	freePort := getFreePort(t)
	host := strings.TrimPrefix(clientConfig.Host, "https://")
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", svcNamespace,
		getRandomPod(t, clientConfig, svcNamespace, svcName).GetName())

	//Create port forwarder
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost,
		&url.URL{Scheme: "https", Path: path, Host: host})
	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", freePort, containerPort)},
		stopChan, readChan, &buffOut, &buffErr)
	require.NoError(t, err)

	//Forward the container port to local machine port
	go func() {
		err := forwarder.ForwardPorts()
		require.NoError(t, err)

		<-t.Context().Done()
		close(stopChan)
	}()
	//Wait till port forwarding is done
	<-readChan

	return fmt.Sprintf("http://localhost:%d", freePort)
}

func getFreePort(tb testing.TB) int {
	tb.Helper()

	listerner, err := net.Listen("tcp", "localhost:0")
	require.NoError(tb, err)

	freePort := listerner.Addr().(*net.TCPAddr).Port
	listerner.Close()
	return freePort
}

func getRandomPod(tb testing.TB, config *rest.Config, svcNS string, svcName string) *core.Pod {
	tb.Helper()

	ctx := tb.Context()
	clientset, err := kubernetes.NewForConfig(config)
	require.NoError(tb, err)

	service, err := clientset.CoreV1().Services(svcNS).Get(ctx, svcName, meta.GetOptions{})
	require.NoError(tb, err)

	var labels []string
	for key, value := range service.Spec.Selector {
		labels = append(labels, key+"="+value)
	}

	pods, err := clientset.CoreV1().Pods(svcNS).List(ctx, meta.ListOptions{
		LabelSelector: strings.Join(labels, ","),
		Limit:         1,
	})
	require.NoError(tb, err)
	require.NotNil(tb, pods)

	return &pods.Items[0]

}
