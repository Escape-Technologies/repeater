package autoprovisioning

import (
	"context"
	"errors"
	"os"

	publicAPI "github.com/Escape-Technologies/cli/pkg/api"
	"github.com/Escape-Technologies/cli/pkg/log"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Autoprovisioner struct {
	client        *publicAPI.ClientWithResponses
	repeaterName  string
	locationId    uuid.UUID
	integrationId uuid.UUID
}

func NewAutoprovisioner() (*Autoprovisioner, error) {
	log.SetLevel(logrus.TraceLevel)
	client, err := publicAPI.NewAPIClient()
	if err != nil {
		return nil, err
	}
	repeaterName := os.Getenv("ESCAPE_REPEATER_NAME")
	if repeaterName == "" {
		return nil, errors.New("ESCAPE_REPEATER_NAME is not set")
	}
	return &Autoprovisioner{client: client, repeaterName: repeaterName,
		locationId: uuid.Nil,
	}, nil
}

// Get the repeater ID from locations using the public API
// If the repeater is not found, create it
func (a *Autoprovisioner) GetId(ctx context.Context) (string, error) {
	if a.locationId != uuid.Nil {
		return a.locationId.String(), nil
	}
	return a.getId(ctx)
}

func (a *Autoprovisioner) getId(ctx context.Context) (string, error) {
	logger.Info("Looking up for repeater %s", a.repeaterName)
	locations, err := a.client.GetV1LocationsWithResponse(ctx)
	if err != nil {
		return "", err
	}
	if locations.JSON200 == nil {
		return "", errors.New("no locations found")
	}

	for _, location := range *locations.JSON200 {
		if location.Name == a.repeaterName && location.Type == "REPEATER" {
			a.locationId = location.Id
			logger.Info("Repeater found in location %s", a.repeaterName)
			return a.locationId.String(), nil
		}
	}
	logger.Info("Repeater not found in location, creating it")

	// Create the repeater
	location, err := a.client.PostV1LocationsWithResponse(ctx, publicAPI.PostV1LocationsJSONRequestBody{
		Name: a.repeaterName,
	})
	if err != nil {
		return "", err
	}
	if location.JSON200 == nil {
		return "", errors.New("no location created")
	}
	a.locationId = location.JSON200.Id
	logger.Info("New repeater created with name %s", a.repeaterName)
	return a.locationId.String(), nil
}

// Create a kubernetes integration if it doesn't exist
func (a *Autoprovisioner) CreateIntegration(ctx context.Context) error {
	if a.integrationId != uuid.Nil {
		return nil
	}

	logger.Debug("Looking up for integration bound to repeater %s", a.repeaterName)
	if a.locationId == uuid.Nil {
		_, err := a.getId(ctx)
		if err != nil {
			return err
		}
	}
	integrations, err := a.client.GetV1IntegrationsKubernetesWithResponse(ctx)
	if err != nil {
		return err
	}
	if integrations.JSON200 == nil {
		return errors.New("no integrations found")
	}

	for _, integration := range *integrations.JSON200 {
		if integration.LocationId != nil && *integration.LocationId == a.locationId {
			logger.Debug("Integration found, nothing to do")
			return nil

		}
	}

	// Create the integration
	logger.Info("No kubernetes integration is bound to repeater %s, creating it", a.repeaterName)
	integration, err := a.client.PostV1IntegrationsKubernetesWithResponse(ctx, publicAPI.PostV1IntegrationsKubernetesJSONRequestBody{
		LocationId: &a.locationId,
		Name:       a.repeaterName,
	})
	if err != nil {
		return err
	}
	if integration.JSON200 == nil {
		return errors.New("no integration created")
	}
	a.integrationId = integration.JSON200.Id
	logger.Info("Kubernetes integration created with id %s", a.integrationId)
	return nil
}
