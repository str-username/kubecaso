package kubecaso

import (
    "context"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClientRequest interface
type ClientRequest interface {
    Pods(ns string) ([]v1.Pod, error)
    PodDelete(ns, name string) error
    Namespaces() ([]v1.Namespace, error)
    Configmaps(ns string) ([]v1.ConfigMap, error)
    ConfigmapResourceVersion(namespace, name string) (string, error)
    Secrets(ns string) ([]v1.Secret, error)
    SecretResourceVersion(namespace, name string, version string) (string, error)
}

// Pods return kubernetes pods in namespace
func (k *KubernetesClient) Pods(ns, label string) ([]v1.Pod, error) {
    pods, err := k.This.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
        LabelSelector: label,
    })
    if err != nil {
        return nil, err
    }
    return pods.Items, nil
}

// PodDelete return kubernetes pods in namespace
func (k *KubernetesClient) PodDelete(ns, name string) error {
    err := k.This.CoreV1().Pods(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
    return err
}

// Namespaces return kubernetes namespaces
func (k *KubernetesClient) Namespaces() ([]v1.Namespace, error) {
    namespaces, err := k.This.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return namespaces.Items, nil
}

// Configmaps return kubernetes configmaps by namespace
func (k *KubernetesClient) Configmaps(ns string) ([]v1.ConfigMap, error) {
    configmaps, err := k.This.CoreV1().ConfigMaps(ns).List(context.Background(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return configmaps.Items, nil
}

// ConfigmapResourceVersion return kubernetes configmap by name
func (k *KubernetesClient) ConfigmapResourceVersion(ns, name string) (string, error) {
    configmap, err := k.This.CoreV1().ConfigMaps(ns).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
        return "", err
    }
    return configmap.ResourceVersion, nil
}

// Secrets return kubernetes configmaps by namespace
func (k *KubernetesClient) Secrets(ns string) ([]v1.Secret, error) {
    configmaps, err := k.This.CoreV1().Secrets(ns).List(context.Background(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return configmaps.Items, nil
}

// SecretResourceVersion return kubernetes configmap by name
func (k *KubernetesClient) SecretResourceVersion(ns, name string) (string, error) {
    configmap, err := k.This.CoreV1().Secrets(ns).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
        return "", err
    }
    return configmap.ResourceVersion, nil
}
