package models

import (
	"context"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/cache"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// ModelInfo is our CLI's model representation
type ModelInfo struct {
	Model               string                  `json:"model"`
	Type                string                  `json:"type"`
	VectorType          *string                 `json:"vector_type"`
	DefaultDimension    *int32                  `json:"default_dimension"`
	ProviderName        *string                 `json:"provider_name"`
	ShortDescription    string                  `json:"short_description"`
	MaxBatchSize        *int32                  `json:"max_batch_size"`
	MaxSequenceLength   *int32                  `json:"max_sequence_length"`
	Modality            *string                 `json:"modality"`
	SupportedDimensions *[]int32                `json:"supported_dimensions"`
	SupportedMetrics    *[]pinecone.IndexMetric `json:"supported_metrics"`
}

// GetModels fetches models from API or cache
func GetModels(ctx context.Context, useCache bool) ([]ModelInfo, error) {
	if useCache {
		return getModelsWithCache(ctx)
	}

	// When not using cache, fetch from API and update cache
	models, err := getModelsFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache with fresh data
	cache.Cache.Set("models", models, 24*time.Hour)
	return models, nil
}

// getModelsWithCache tries cache first, then API if not found
func getModelsWithCache(ctx context.Context) ([]ModelInfo, error) {
	// Try to get from cache first
	cached, found, err := cache.GetCached[[]ModelInfo]("models")
	if found && err == nil {
		return *cached, nil
	}

	// Fetch from API if not in cache
	models, err := getModelsFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the models
	cache.CacheWithTTL("models", models, 24*time.Hour)
	return models, nil
}

// getModelsFromAPI fetches models directly from the API
func getModelsFromAPI(ctx context.Context) ([]ModelInfo, error) {
	pc := sdk.NewPineconeClient()
	embed := "embed"
	embedModels, err := pc.Inference.ListModels(ctx, &pinecone.ListModelsParams{Type: &embed})
	if err != nil {
		return nil, err
	}

	if embedModels == nil || embedModels.Models == nil {
		return []ModelInfo{}, nil
	}

	// Convert pinecone.ModelInfo to our ModelInfo
	models := make([]ModelInfo, len(*embedModels.Models))
	for i, model := range *embedModels.Models {
		models[i] = ModelInfo{
			Model:               model.Model,
			Type:                model.Type,
			VectorType:          model.VectorType,
			DefaultDimension:    model.DefaultDimension,
			ProviderName:        model.ProviderName,
			ShortDescription:    model.ShortDescription,
			MaxBatchSize:        model.MaxBatchSize,
			MaxSequenceLength:   model.MaxSequenceLength,
			Modality:            model.Modality,
			SupportedDimensions: model.SupportedDimensions,
			SupportedMetrics:    model.SupportedMetrics,
		}
	}

	return models, nil
}
