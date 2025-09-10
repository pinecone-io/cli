package models

import (
	"context"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/cache"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// CachedModel is a simplified version of pinecone.ModelInfo for caching
type CachedModel struct {
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
func GetModels(ctx context.Context, useCache bool) ([]pinecone.ModelInfo, error) {
	if useCache {
		return getModelsWithCache(ctx)
	}

	// When not using cache, fetch from API and update cache
	models, err := getModelsFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache with fresh data
	cachedModels := convertModelInfoToCached(models)
	cache.Cache.Set("models", cachedModels, 24*time.Hour)
	return models, nil
}

// getModelsWithCache tries cache first, then API if not found
func getModelsWithCache(ctx context.Context) ([]pinecone.ModelInfo, error) {
	// Try to get from cache first
	cached, found, err := cache.GetCached[[]CachedModel]("models")
	if found && err == nil {
		// Convert cached models to pinecone.ModelInfo
		models := convertCachedToModelInfo(*cached)
		return models, nil
	}

	// Fetch from API if not in cache
	models, err := getModelsFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to cached models and cache
	cachedModels := convertModelInfoToCached(models)
	cache.CacheWithTTL("models", cachedModels, 24*time.Hour)
	return models, nil
}

// getModelsFromAPI fetches models directly from the API
func getModelsFromAPI(ctx context.Context) ([]pinecone.ModelInfo, error) {
	pc := sdk.NewPineconeClient()
	embed := "embed"
	embedModels, err := pc.Inference.ListModels(ctx, &pinecone.ListModelsParams{Type: &embed})
	if err != nil {
		return nil, err
	}

	if embedModels == nil || embedModels.Models == nil {
		return []pinecone.ModelInfo{}, nil
	}

	return *embedModels.Models, nil
}

// convertCachedToModelInfo converts CachedModel to pinecone.ModelInfo
func convertCachedToModelInfo(cached []CachedModel) []pinecone.ModelInfo {
	models := make([]pinecone.ModelInfo, len(cached))
	for i, cachedModel := range cached {
		models[i] = pinecone.ModelInfo{
			Model:               cachedModel.Model,
			Type:                cachedModel.Type,
			VectorType:          cachedModel.VectorType,
			DefaultDimension:    cachedModel.DefaultDimension,
			ProviderName:        cachedModel.ProviderName,
			ShortDescription:    cachedModel.ShortDescription,
			MaxBatchSize:        cachedModel.MaxBatchSize,
			MaxSequenceLength:   cachedModel.MaxSequenceLength,
			Modality:            cachedModel.Modality,
			SupportedDimensions: cachedModel.SupportedDimensions,
			SupportedMetrics:    cachedModel.SupportedMetrics,
		}
	}
	return models
}

// convertModelInfoToCached converts pinecone.ModelInfo to CachedModel
func convertModelInfoToCached(models []pinecone.ModelInfo) []CachedModel {
	cached := make([]CachedModel, len(models))
	for i, model := range models {
		cached[i] = CachedModel{
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
	return cached
}
