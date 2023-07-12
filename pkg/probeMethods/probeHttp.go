package probemethods

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/probe"
	httpProbe "k8s.io/kubernetes/pkg/probe/http"
)

// ProbeHTTPGet performs an HTTPGet probe
func ProbeHTTPGet(podIP string, container v1.Container, probeSpec *corev1.Probe) (probe.Result, string, error) {

	httpProber := httpProbe.New(false) // Don't follow redirects to a different hostname

	timeout := time.Duration(probeSpec.TimeoutSeconds) * time.Second
	req, err := httpProbe.NewRequestForHTTPGetAction(probeSpec.HTTPGet, &container, podIP, "probe")
	if err != nil {
		return probe.Unknown, "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	port := req.URL.Port()
	host := req.URL.Hostname()
	path := req.URL.Path
	scheme := req.URL.Scheme
	headers := probeSpec.HTTPGet.HTTPHeaders

	klog.V(4).Info("HTTP-Probe: scheme = ", scheme,
		", host = ", host,
		", port = ", port,
		", path = ", path,
		", timeout = ", timeout,
		", headers = ", headers)

	return httpProber.Probe(req, time.Duration(probeSpec.TimeoutSeconds)*time.Second)
}
