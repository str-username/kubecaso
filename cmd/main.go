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

		//  Run goroutine secrets
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
			log.Fatal().Err(err).Msg("tba")
		}

		if len(prevConfigMaps) == 0 && len(prevSecrets) == 0 {
			log.Info().Msg("initialize state")
			prevConfigMaps = configMaps
			prevSecrets = secrets
			log.Info().Msg("initialize state done")
		}

		// Check if configmap change
		if !reflect.DeepEqual(configMaps, prevConfigMaps) {
			log.Info().Msg("configmap was change")
			log.Info().Msg("implement configmap logic stub") // TODO: delete pod logic implement
			prevConfigMaps = configMaps
		}

		// Check if secrets change
		if !reflect.DeepEqual(secrets, prevSecrets) {
			log.Info().Msg("secret was change")
			log.Info().Msg("implement secret logic stub") // TODO: delete pod logic implement
			prevSecrets = secrets
		}
		log.Info().Any("configmaps old", prevConfigMaps).Any("configmaps new", configMaps).Msg("watch")
		log.Info().Any("secrets old", prevSecrets).Any("secrets new", secrets).Msg("watch")
	}
}
