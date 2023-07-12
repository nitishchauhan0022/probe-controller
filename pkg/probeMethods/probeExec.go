package probemethods

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/probe"
)

// ProbeExec performs an Exec probe
func ProbeExec(podIP string, probe corev1.ExecAction) (probe.Result, string, error) {
	//TODO: implement
	return "", "", nil
}
