package autoprovisioning

import (
	"context"
	"errors"
	"os"

	publicAPI "github.com/Escape-Technologies/cli/pkg/api"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/google/uuid"
)

type Autoprovisioner struct {
	client       *publicAPI.Client
	repeaterName string
	locationId   uuid.UUID
}

func NewAutoprovisioner() (*Autoprovisioner, error) {
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
	res, err := a.client.GetV1Locations(ctx)
	if err != nil {
		return "", err
	}
	locations, err := publicAPI.ParseGetV1LocationsResponse(res)
	if err != nil {
		return "", err
	}
	if locations.JSON200 == nil {
		return "", errors.New("no locations found")
	}

	for _, location := range *locations.JSON200 {
		if location.Name == a.repeaterName {
			a.locationId = location.Id
			logger.Info("Repeater found in location %s", a.repeaterName)
			return a.locationId.String(), nil
		}
	}
	logger.Info("Repeater not found in location, creating it")

	// Create the repeater
	res, err = a.client.PostV1Locations(ctx, publicAPI.PostV1LocationsJSONRequestBody{
		Name: a.repeaterName,
	})
	if err != nil {
		return "", err
	}
	location, err := publicAPI.ParsePostV1LocationsResponse(res)
	if err != nil {
		return "", err
	}
	if location.JSON200 == nil {
		return "", errors.New("no location created")
	}
	a.locationId = location.JSON200.Id
	logger.Info("New repeater created")
	return a.locationId.String(), nil
}

// Create a kubernetes integration if it doesn't exist
func (a *Autoprovisioner) CreateIntegration(ctx context.Context) error {
	logger.Debug("Looking up for integration bound to repeater %s", a.repeaterName)
	if a.locationId == uuid.Nil {
		_, err := a.getId(ctx)
		if err != nil {
			return err
		}
	}
	res, err := a.client.GetV1Integrations(ctx)
	if err != nil {
		return err
	}
	integrations, err := publicAPI.ParseGetV1IntegrationsResponse(res)
	if err != nil {
		return err
	}
	if integrations.JSON200 == nil {
		return errors.New("no integrations found")
	}

	// for _, integration := range *integrations.JSON200 {
	// 	if integration.Kind == "KUBERNETES" && integration.LocationId != nil && *integration.LocationId == a.locationId {
	// 		logger.Debug("Integration found, nothing to do")
	// 		return nil
	// 	}
	// }

	// id, err := uuid.Parse(a.locationId)
	// if err != nil {
	// 	return err
	// }
	// // Create the integration
	// logger.Info("KUBERNETES integration bound to repeater %s not found, creating it", a.repeaterName)
	// res, err = a.client.PostV1Integrations(ctx, publicAPI.PostV1IntegrationsJSONRequestBody{
	// 	LocationId: id,
	// 	Name:       a.repeaterName,
	// })
	// if err != nil {
	// 	return err
	// }
	return nil
}
