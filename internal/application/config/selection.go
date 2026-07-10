package config

import "fmt"

func (configuration Configuration) SelectInstance(explicitName string) (Instance, error) {
	if err := configuration.Validate(); err != nil {
		return Instance{}, err
	}

	if explicitName != "" {
		for _, instance := range configuration.Instances {
			if instance.Name == explicitName {
				return instance, nil
			}
		}
		return Instance{}, fmt.Errorf("instance profile %q not found", explicitName)
	}

	if len(configuration.Instances) == 1 {
		return configuration.Instances[0], nil
	}

	return Instance{}, fmt.Errorf("instance selection is ambiguous: specify one of %d configured profiles", len(configuration.Instances))
}
