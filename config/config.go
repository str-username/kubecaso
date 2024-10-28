package config

import (
    "github.com/go-playground/validator/v10"
    "github.com/rs/zerolog/log"
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    Watch struct {
        Any   string `yaml:"any" validate:"required"`
        Label string `yaml:"label" validate:"required"`
    }
}

func (c Config) Validate() error {
    return validator.New().Struct(c)
}

type Valid interface {
    Validate() error
}

func NewConfig(configPath string) (*Config, error) {
    config := &Config{}
    configFile, err := os.Open(configPath)
    
    if err != nil {
        return nil, err
    }
    
    defer configFile.Close()
    
    configNewDecoder := yaml.NewDecoder(configFile)
    
    if err = configNewDecoder.Decode(config); err != nil {
        return nil, err
    }
    
    var validConfig Valid = config
    
    if err = validConfig.Validate(); err != nil {
        return nil, err
    }
    
    log.Info().Str("kubecaso", "config").Msg("create config")
    return config, nil
}
