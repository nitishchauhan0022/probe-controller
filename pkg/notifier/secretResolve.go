package notifier

import (
	"context"
	"fmt"
	"probe/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func secretResolver(client client.Client, namespace string, resolvesecret v1alpha1.ResolveFromSecret) (string, error) {
	// Get a Kubernetes client

	// Get a secret
	secret, err := getSecret(client, resolvesecret.SecretName, namespace)
	if err != nil {
		return "", fmt.Errorf("error getting secret %s", err)
	}

	// Print the secret data
	for key, value := range secret.Data {
		if key == resolvesecret.SecretKey {
			return string(value), nil
		}
	}

	return "", fmt.Errorf("key: %s not found in secret %s", resolvesecret.SecretKey, resolvesecret.SecretName)
}

// getSecret retrieves a secret from a specific namespace
func getSecret(client client.Client, secretName, namespace string) (*corev1.Secret, error) {
	ctx := context.Background()

	// Define an instance of the needed object
	secret := &corev1.Secret{}

	// Create a NamespacedName based on the secret name and namespace
	nn := types.NamespacedName{
		Name:      secretName,
		Namespace: namespace,
	}

	// Use the Kubernetes client to get the secret
	if err := client.Get(ctx, nn, secret); err != nil {
		return nil, fmt.Errorf("error getting secret %s", err)
	}

	return secret, nil
}
