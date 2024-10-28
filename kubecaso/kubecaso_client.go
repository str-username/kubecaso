package kubecaso

import (
    "kube-caso/config"
)

// KubeCaso save kubernetes client and config
type KubeCaso struct {
    cfg *config.Config
    cli *KubernetesClient
}

// NewKubeCasoClient creat new client
func NewKubeCasoClient(configPath string) (*KubeCaso, error) {
    cfg, err := config.NewConfig(configPath)
    if err != nil {
        return nil, err
    }
    
    cli, err := NewKubernetesClient()
    if err != nil {
        return nil, err
    }
    
    return &KubeCaso{
        cfg: cfg,
        cli: cli,
    }, nil
}
