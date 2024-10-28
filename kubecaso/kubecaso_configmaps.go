package kubecaso

import (
	"log"
)

type ConfigMapsData struct {
	Pod              string
	Namespace        string
	ConfigMapName    string
	ConfigMapVersion string
}

// ConfigMaps return ConfigMaps
func (kc *KubeCaso) ConfigMaps() ([]ConfigMapsData, error) {
	namespaces, err := kc.Cli.Namespaces()
	if err != nil {
		return nil, err
	}

	var cMaps []ConfigMapsData

	for _, ns := range namespaces {
		pods, err := kc.Cli.Pods(ns.Name, kc.cfg.Watch.Label)
		if err != nil {
			return nil, err
		}
		for _, pod := range pods {
			for _, vol := range pod.Spec.Volumes {
				if vol.ConfigMap != nil {
					cMapsVersion, err := kc.Cli.ConfigmapResourceVersion(ns.Name, vol.ConfigMap.Name)
					if err != nil {
						log.Printf("Error getting ConfigMap version for %s: %v", vol.ConfigMap.Name, err)
						continue
					}
					data := ConfigMapsData{
						Pod:              pod.Name,
						Namespace:        ns.Name,
						ConfigMapName:    vol.ConfigMap.Name,
						ConfigMapVersion: cMapsVersion,
					}
					cMaps = append(cMaps, data)
				}
			}
		}
	}
	return cMaps, nil
}
