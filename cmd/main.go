package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"kube-caso/kubecaso"
	"reflect"
	"sync"
	"time"
)

const timeout = 5

// getChangedSecrets return changed secrets
func getChangedSecrets(prevSecrets, newSecrets []kubecaso.SecretsData) []kubecaso.SecretsData {
	var changedSecrets []kubecaso.SecretsData
	secretMap := make(map[string]string)

	// Create old secrets map
	for _, secret := range prevSecrets {
		secretMap[secret.SecretName] = secret.SecretVersion
	}

	// Compare secrets versions
	for _, secret := range newSecrets {
		if prevVersion, ok := secretMap[secret.SecretName]; !ok || prevVersion != secret.SecretVersion {
			changedSecrets = append(changedSecrets, secret)
		}
	}

	return changedSecrets
}

// getChangedConfigMaps return changed configmaps
func getChangedConfigMaps(prevConfigMaps, newConfigMaps []kubecaso.ConfigMapsData) []kubecaso.ConfigMapsData {
	var changedConfigMaps []kubecaso.ConfigMapsData
	configMapMap := make(map[string]string)

	// Create old configmap map
	for _, configMap := range prevConfigMaps {
		configMapMap[configMap.ConfigMapName] = configMap.ConfigMapVersion
	}

	// Compare configmap map versions
	for _, configMap := range newConfigMaps {
		if prevVersion, ok := configMapMap[configMap.ConfigMapName]; !ok || prevVersion != configMap.ConfigMapVersion {
			changedConfigMaps = append(changedConfigMaps, configMap)
		}
	}

	return changedConfigMaps
}

func main() {
	// Init kubecaso client
	kc, err := kubecaso.NewKubeCasoClient("etc/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing KubeCasoClient")
	}

	// Execution interval
	interval := timeout * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Old state variables
	var prevConfigMaps []kubecaso.ConfigMapsData
	var prevSecrets []kubecaso.SecretsData

	// Execute
	for range ticker.C {
		var wg sync.WaitGroup
		errorChan := make(chan error, 2)

		// Current state variables
		var configMaps []kubecaso.ConfigMapsData
		var secrets []kubecaso.SecretsData

		// Run goroutine configmaps
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmaps, err := kc.ConfigMaps()
			if err != nil {
				errorChan <- fmt.Errorf("error fetching ConfigMaps: %v", err)
				return
			}
			configMaps = cmaps
		}()

		// Run goroutine secrets
		wg.Add(1)
		go func() {
			defer wg.Done()
			secs, err := kc.Secrets()
			if err != nil {
				errorChan <- fmt.Errorf("error fetching Secrets: %v", err)
				return
			}
			secrets = secs
		}()

		// Wait done
		wg.Wait()
		close(errorChan)

		// Error handle
		for err := range errorChan {
			log.Fatal().Err(err).Msg("Error occurred")
		}

		if len(prevConfigMaps) == 0 && len(prevSecrets) == 0 {
			prevConfigMaps = configMaps
			prevSecrets = secrets
		}

		// Delete pod if configmap changed
		if !reflect.DeepEqual(configMaps, prevConfigMaps) {
			log.Info().Msg("ConfigMap was changed")
			changedConfigMaps := getChangedConfigMaps(prevConfigMaps, configMaps)

			for _, data := range changedConfigMaps {
				if err := kc.Cli.PodDelete(data.Namespace, data.Pod); err != nil {
					continue
				}
				log.Info().Msgf("Deleted Pod %s in namespace %s due to ConfigMap change", data.Pod, data.Namespace)
			}
			prevConfigMaps = configMaps
		}

		// Delete pod if secret changed
		if !reflect.DeepEqual(secrets, prevSecrets) {
			log.Info().Msg("Secret was changed")
			changedSecrets := getChangedSecrets(prevSecrets, secrets)

			for _, data := range changedSecrets {
				if err := kc.Cli.PodDelete(data.Namespace, data.Pod); err != nil {
					continue
				}
				log.Info().Msgf("Deleted Pod %s in namespace %s due to Secret change", data.Pod, data.Namespace)
			}
			prevSecrets = secrets
		}
		log.Info().Msg("watch")
	}
}
