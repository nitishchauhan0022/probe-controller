/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	probev1alpha1 "probe/api/v1alpha1"
	probenotifier "probe/pkg/notifier"
	probemethods "probe/pkg/probeMethods"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/probe"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProbeReconciler reconciles a Probe object
type ProbeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=probe.probe.k8s,resources=probes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=probe.probe.k8s,resources=probes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=probe.probe.k8s,resources=probes/finalizers,verbs=update

func (r *ProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	klog.Info("Reconciling Probe")
	// Fetch the Probe instance
	probeResource := &probev1alpha1.Probe{}
	if err := r.Get(ctx, req.NamespacedName, probeResource); err != nil {
		if errors.IsNotFound(err) {
			klog.Errorf("Probe resource not found %s", err)
			return ctrl.Result{}, nil
		}
		klog.Errorf("Error getting probe resource %s", err)
		return ctrl.Result{}, err
	}
	RequeueAfter := time.Second * time.Duration(probeResource.Spec.Probe.PeriodSeconds)
	probeResult, msg, err := r.ProbeResource(ctx, &probeResource.Spec)
	if err != nil {
		klog.Errorf("Error probing resource %s", err)
		return ctrl.Result{}, err
	}

	if probeResult == probe.Failure {
		probeResource.Status.FailureCount++

		if probeResource.Status.FailureCount >= 5 {
			// Send Slack notification
			if err := probenotifier.NotifySlack(r.Client, probeResource.Namespace, probeResource.Spec.SlackDetails, "Probe failed for 5 times in a row."); err != nil {
				klog.Errorf("Error sending Slack notification %s", err)
				return ctrl.Result{}, err
			}
		}
	} else {
		probeResource.Status.FailureCount = 0
	}

	// Update the status of the Probe object
	probeResource.Status.Status = string(probeResult)
	probeResource.Status.Message = msg
	if err := r.Status().Update(ctx, probeResource); err != nil {
		klog.Errorf("Error updating probe status %s", err)
		return ctrl.Result{}, err
	}

	// Schedule next probe after probe.Spec.Probe.PeriodSeconds
	return ctrl.Result{RequeueAfter: RequeueAfter}, nil
}

func (r *ProbeReconciler) ProbeResource(ctx context.Context, probeResource *probev1alpha1.ProbeSpec) (probe.Result, string, error) {
	var err error
	var result probe.Result
	var msg string

	// Define a PodList object
	podList := &corev1.PodList{}

	// Define ListOptions
	listOps := &client.ListOptions{LabelSelector: labels.SelectorFromSet(probeResource.Selector.MatchLabels)}

	client := r.Client
	client.List(ctx, podList, listOps)
	// List the pods
	if err = client.List(ctx, podList, listOps); err != nil {
		return "", "", err
	}

	// Loop through the pods and probe
	for _, pod := range podList.Items {
		if probeResource.ContainerName == "" {
			// If no container is specified, use the first container in the pod
			probeResource.ContainerName = pod.Spec.Containers[0].Name
		} else {
			// Check if the container specified in the probe exists in the pod
			for _, container := range pod.Spec.Containers {
				if container.Name == probeResource.ContainerName {
					break
				}
				return "", "", fmt.Errorf("container %s not found in pod %s", probeResource.ContainerName, pod.Name)
			}
		}

		container, err := getContainer(pod, probeResource.ContainerName)
		if err != nil {
			return "", "", err
		}

		switch {
		case probeResource.Probe.HTTPGet != nil:
			// Perform HTTPGet Probe
			if result, msg, err = probemethods.ProbeHTTPGet(pod.Status.PodIP, container, &probeResource.Probe); err != nil {
				return "", "", err
			}
		case probeResource.Probe.TCPSocket != nil:
			// Perform TCPSocket Probe
			if result, msg, err = probemethods.ProbeTCPSocket(pod.Status.PodIP, container, &probeResource.Probe); err != nil {
				return "", "", err
			}
		case probeResource.Probe.Exec != nil:
			// Perform Exec Probe
			if result, msg, err = probemethods.ProbeExec(pod.Status.PodIP, *probeResource.Probe.Exec); err != nil {
				return "", "", err
			}
		}
	}
	return result, msg, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ProbeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&probev1alpha1.Probe{}).
		Complete(r)
}

func getContainer(pod corev1.Pod, name string) (v1.Container, error) {
	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return container, nil
		}
	}
	return v1.Container{}, fmt.Errorf("container %s not found in pod %s", name, pod.Name)
}

func (r *ProbeReconciler) GetClient() client.Client {
	return r.Client
}
