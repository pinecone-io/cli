package index

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/validation"
	"github.com/spf13/cobra"
)

// ValidateIndexNameArgs validates that exactly one non-empty index name is provided as a positional argument.
// This is the standard validation used across all index commands (create, describe, delete, configure).
func ValidateIndexNameArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("\b" + style.FailMsg("please provide an index name"))
	}
	if len(args) > 1 {
		return errors.New("\b" + style.FailMsg("please provide only one index name"))
	}
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("\b" + style.FailMsg("index name cannot be empty"))
	}
	return nil
}

// CreateOptionsRule creates a new validation rule from a function that takes *CreateOptions
func CreateOptionsRule(fn func(*CreateOptions) string) validation.Rule {
	return func(value interface{}) string {
		config, ok := value.(*CreateOptions)
		if !ok {
			return ""
		}
		return fn(config)
	}
}

// ValidateCreateOptions validates the index creation configuration using the validation framework
func ValidateCreateOptions(config CreateOptions) []string {
	validator := validation.New()

	validator.AddRule(CreateOptionsRule(validateConfigIndexTypeFlags))
	validator.AddRule(CreateOptionsRule(validateConfigHasName))
	validator.AddRule(CreateOptionsRule(validateConfigNameLength))
	validator.AddRule(CreateOptionsRule(validateConfigNameStartsWithAlphanumeric))
	validator.AddRule(CreateOptionsRule(validateConfigNameEndsWithAlphanumeric))
	validator.AddRule(CreateOptionsRule(validateConfigNameCharacters))
	validator.AddRule(CreateOptionsRule(validateConfigServerlessCloud))
	validator.AddRule(CreateOptionsRule(validateConfigServerlessRegion))
	validator.AddRule(CreateOptionsRule(validateConfigPodEnvironment))
	validator.AddRule(CreateOptionsRule(validateConfigPodType))
	validator.AddRule(CreateOptionsRule(validateConfigPodSparseVector))
	validator.AddRule(CreateOptionsRule(validateConfigSparseVectorDimension))
	validator.AddRule(CreateOptionsRule(validateConfigSparseVectorMetric))
	validator.AddRule(CreateOptionsRule(validateConfigDenseVectorDimension))

	return validator.Validate(&config)
}

// validateConfigIndexTypeFlags checks that serverless and pod flags are not both set
func validateConfigIndexTypeFlags(config *CreateOptions) string {
	if config.Serverless.Value && config.Pod.Value {
		return fmt.Sprintf("%s and %s cannot be provided together", style.Code("serverless"), style.Code("pod"))
	}
	return ""
}

// validateConfigHasName checks if the config has a non-empty name
func validateConfigHasName(config *CreateOptions) string {
	if strings.TrimSpace(config.Name.Value) == "" {
		return "index must have a name"
	}
	return ""
}

// validateConfigNameLength checks if the config name is 1-45 characters long
func validateConfigNameLength(config *CreateOptions) string {
	name := strings.TrimSpace(config.Name.Value)
	if len(name) < 1 || len(name) > 45 {
		return "index name must be 1-45 characters long"
	}
	return ""
}

// validateConfigNameStartsWithAlphanumeric checks if the config name starts with an alphanumeric character
func validateConfigNameStartsWithAlphanumeric(config *CreateOptions) string {
	name := strings.TrimSpace(config.Name.Value)
	if len(name) > 0 {
		first := name[0]
		if !((first >= 'a' && first <= 'z') || (first >= '0' && first <= '9')) {
			return "index name must start with an alphanumeric character"
		}
	}
	return ""
}

// validateConfigNameEndsWithAlphanumeric checks if the config name ends with an alphanumeric character
func validateConfigNameEndsWithAlphanumeric(config *CreateOptions) string {
	name := strings.TrimSpace(config.Name.Value)
	if len(name) > 0 {
		last := name[len(name)-1]
		if !((last >= 'a' && last <= 'z') || (last >= '0' && last <= '9')) {
			return "index name must end with an alphanumeric character"
		}
	}
	return ""
}

// validateConfigNameCharacters checks if the config name consists only of lowercase alphanumeric characters or '-'
func validateConfigNameCharacters(config *CreateOptions) string {
	name := strings.TrimSpace(config.Name.Value)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return "index name must consist only of lowercase alphanumeric characters or '-'"
		}
	}
	return ""
}

// validateConfigServerlessCloud checks that cloud is provided for serverless indexes
func validateConfigServerlessCloud(config *CreateOptions) string {
	if config.GetSpec() == IndexSpecServerless && config.Cloud.Value == "" {
		return fmt.Sprintf("%s is required for %s indexes", style.Code("cloud"), style.Code("serverless"))
	}
	return ""
}

// validateConfigServerlessRegion checks that region is provided for serverless indexes
func validateConfigServerlessRegion(config *CreateOptions) string {
	if config.GetSpec() == IndexSpecServerless && config.Region.Value == "" {
		return fmt.Sprintf("%s is required for %s indexes", style.Code("region"), style.Code("serverless"))
	}
	return ""
}

// validateConfigPodEnvironment checks that environment is provided for pod indexes
func validateConfigPodEnvironment(config *CreateOptions) string {
	if config.GetSpec() == IndexSpecPod && config.Environment.Value == "" {
		return fmt.Sprintf("%s is required for %s indexes", style.Code("environment"), style.Code("pod"))
	}
	return ""
}

// validateConfigPodType checks that pod_type is provided for pod indexes
func validateConfigPodType(config *CreateOptions) string {
	if config.GetSpec() == IndexSpecPod && config.PodType.Value == "" {
		return fmt.Sprintf("%s is required for %s indexes", style.Code("pod_type"), style.Code("pod"))
	}
	return ""
}

// validateConfigPodSparseVector checks that pod indexes cannot use sparse vector type
func validateConfigPodSparseVector(config *CreateOptions) string {
	if config.GetSpec() == IndexSpecPod && config.VectorType.Value == "sparse" {
		return fmt.Sprintf("%s vector type is not supported for %s indexes", style.Code("sparse"), style.Code("pod"))
	}
	return ""
}

// validateConfigSparseVectorDimension checks that dimension should not be specified for sparse vector type
func validateConfigSparseVectorDimension(config *CreateOptions) string {
	if config.VectorType.Value == "sparse" && config.Dimension.Value > 0 {
		return fmt.Sprintf("%s should not be specified when vector type is %s", style.Code("dimension"), style.Code("sparse"))
	}
	return ""
}

// validateConfigSparseVectorMetric checks that metric should be 'dotproduct' for sparse vector type
func validateConfigSparseVectorMetric(config *CreateOptions) string {
	if config.VectorType.Value == "sparse" && config.Metric.Value != "" && config.Metric.Value != "dotproduct" {
		return fmt.Sprintf("metric should be %s when vector type is %s", style.Code("dotproduct"), style.Code("sparse"))
	}
	return ""
}

// validateConfigDenseVectorDimension checks that dimension is provided for dense vector indexes
func validateConfigDenseVectorDimension(config *CreateOptions) string {
	// Check if it's a dense vector type (empty string means dense, or explicitly "dense")
	if config.VectorType.Value == "dense" && config.Dimension.Value <= 0 {
		return fmt.Sprintf("%s is required when vector type is %s", style.Code("dimension"), style.Code("dense"))
	}
	return ""
}
