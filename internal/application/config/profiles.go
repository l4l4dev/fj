package config

import "fmt"

type Profile struct {
	Name     string
	Endpoint Endpoint
}

func (configuration Configuration) ListProfiles() ([]Profile, error) {
	if err := configuration.Validate(); err != nil {
		return nil, err
	}

	profiles := make([]Profile, len(configuration.Instances))
	for index, instance := range configuration.Instances {
		profiles[index] = profileFromInstance(instance)
	}

	return profiles, nil
}

func (configuration Configuration) InspectProfile(name string) (Profile, error) {
	if err := configuration.Validate(); err != nil {
		return Profile{}, err
	}

	for _, instance := range configuration.Instances {
		if instance.Name == name {
			return profileFromInstance(instance), nil
		}
	}

	return Profile{}, fmt.Errorf("profile %q not found", name)
}

func profileFromInstance(instance Instance) Profile {
	return Profile{Name: instance.Name, Endpoint: instance.Endpoint}
}
