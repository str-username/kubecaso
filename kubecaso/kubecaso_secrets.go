package kubecaso

import (
	"log"
)

type SecretsData struct {
	Pod           string
	Namespace     string
	SecretName    string
	SecretVersion string
}

// Secrets return secrets
func (kc *KubeCaso) Secrets() ([]SecretsData, error) {
	namespaces, err := kc.Cli.Namespaces()
	if err != nil {
		return nil, err
	}

	var secretsList []SecretsData

	for _, ns := range namespaces {
		pods, err := kc.Cli.Pods(ns.Name, kc.cfg.Watch.Label)
		if err != nil {
			return nil, err
		}
		for _, pod := range pods {
			for _, vol := range pod.Spec.Volumes {
				if vol.Secret != nil {
					sVersion, err := kc.Cli.SecretResourceVersion(ns.Name, vol.Secret.SecretName)
					if err != nil {
						log.Printf("Error getting Secret version for %s: %v", vol.Secret.SecretName, err)
						continue
					}
					data := SecretsData{
						Pod:           pod.Name,
						Namespace:     ns.Name,
						SecretName:    vol.Secret.SecretName,
						SecretVersion: sVersion,
					}
					secretsList = append(secretsList, data)
				}
			}
			for _, container := range pod.Spec.Containers {
				for _, env := range container.EnvFrom {
					if env.SecretRef != nil {
						sVersion, err := kc.Cli.SecretResourceVersion(ns.Name, env.SecretRef.Name)
						if err != nil {
							log.Printf("Error getting Secret version for %s: %v", env.SecretRef.Name, err)
							continue
						}
						data := SecretsData{
							Pod:           pod.Name,
							Namespace:     ns.Name,
							SecretName:    env.SecretRef.Name,
							SecretVersion: sVersion,
						}
						secretsList = append(secretsList, data)
					}
				}
			}
		}
	}
	return secretsList, nil
}
