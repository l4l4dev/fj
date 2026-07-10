package config

import (
	"fmt"
	"net/url"
	"strings"
)

func (configuration Configuration) Validate() error {
	if len(configuration.Instances) == 0 {
		return fmt.Errorf("configuration must contain at least one instance")
	}

	names := make(map[string]struct{}, len(configuration.Instances))
	for index, instance := range configuration.Instances {
		if strings.TrimSpace(instance.Name) == "" {
			return fmt.Errorf("instance %d: name is required", index+1)
		}
		if _, exists := names[instance.Name]; exists {
			return fmt.Errorf("instance %q: name must be unique", instance.Name)
		}
		names[instance.Name] = struct{}{}

		if err := validateEndpoint(instance.Name, instance.Endpoint); err != nil {
			return err
		}
		if strings.TrimSpace(string(instance.Credential)) == "" {
			return fmt.Errorf("instance %q: credential reference is required", instance.Name)
		}
	}

	return nil
}

func validateEndpoint(instanceName string, endpoint Endpoint) error {
	rawEndpoint := strings.TrimSpace(string(endpoint))
	if rawEndpoint == "" {
		return fmt.Errorf("instance %q: endpoint is required", instanceName)
	}

	parsedEndpoint, err := url.Parse(rawEndpoint)
	if err != nil || parsedEndpoint.Host == "" || (parsedEndpoint.Scheme != "http" && parsedEndpoint.Scheme != "https") {
		return fmt.Errorf("instance %q: endpoint must be an absolute HTTP or HTTPS URL", instanceName)
	}
	if parsedEndpoint.User != nil {
		return fmt.Errorf("instance %q: endpoint must not contain credentials", instanceName)
	}

	return nil
}
