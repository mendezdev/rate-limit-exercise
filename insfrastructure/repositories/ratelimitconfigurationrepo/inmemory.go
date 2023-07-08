package ratelimitconfigurationrepo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mendezdev/rate-limit-example/core/domain"
	"github.com/mendezdev/rate-limit-example/core/ports"
)

const (
	path = "./mock_store/rate_limit_configuration_repo.json"
)

var (
	byName map[string]domain.RateLimitConfiguration
)

type inmemory struct {
	byName map[string]domain.RateLimitConfiguration
}

func NewInMemory() ports.RateLimitConfigurationRepository {
	return &inmemory{
		byName: byName,
	}
}

// GetRateLimitConfiguration implements ports.RateLimitConfigurationRepository.
func (mem *inmemory) GetRateLimitConfiguration(configType string) (*domain.RateLimitConfiguration, error) {
	config, ok := mem.byName[configType]
	if !ok {
		return nil, nil
	}
	return &config, nil
}

func init() {
	fmt.Println("initializing rate_limit_configuration inmemory db...")

	byName = make(map[string]domain.RateLimitConfiguration)

	var configurations []domain.RateLimitConfiguration
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("error reading json file for initialize rate limit configuration db: %s", err.Error()))
	}

	jsonErr := json.Unmarshal(data, &configurations)
	if jsonErr != nil {
		panic(fmt.Errorf("error unmarshaling json data: %s", jsonErr.Error()))
	}

	for _, c := range configurations {
		byName[c.Name] = c
	}

	fmt.Println("inmemory rate limit configuration db initilized.")
}
