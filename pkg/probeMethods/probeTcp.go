package probemethods

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/probe"
	tcpProbe "k8s.io/kubernetes/pkg/probe/tcp"
)

// ProbeTCPSocket performs a TCPSocket probe
func ProbeTCPSocket(podIP string, container v1.Container, probeSpec *corev1.Probe) (probe.Result, string, error) {
	tcpProber := tcpProbe.New() // Create a new TCP prober

	port, err := probe.ResolveContainerPort(probeSpec.TCPSocket.Port, &container)
	if err != nil {
		return probe.Unknown, "", err
	}
	host := probeSpec.TCPSocket.Host
	if host == "" {
		host = podIP
	}
	klog.V(4).Info("TCP-Probe", "host", host, "port", port, "timeout", probeSpec.TimeoutSeconds)
	return tcpProber.Probe(host, port, time.Duration(probeSpec.TimeoutSeconds)*time.Second)

}
