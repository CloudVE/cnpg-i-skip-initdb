package lifecycle

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/pluginhelper/decoder"
	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/pluginhelper/object"
	"github.com/cloudnative-pg/cnpg-i/pkg/lifecycle"
	"github.com/cloudnative-pg/machinery/pkg/log"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/leonardoce/cnpg-i-skip-initdb/internal/utils"
)

// Implementation is the implementation of the lifecycle handler
type Implementation struct {
	lifecycle.UnimplementedOperatorLifecycleServer
}

// GetCapabilities exposes the lifecycle capabilities
func (impl Implementation) GetCapabilities(
	_ context.Context,
	_ *lifecycle.OperatorLifecycleCapabilitiesRequest,
) (*lifecycle.OperatorLifecycleCapabilitiesResponse, error) {
	return &lifecycle.OperatorLifecycleCapabilitiesResponse{
		LifecycleCapabilities: []*lifecycle.OperatorLifecycleCapabilities{
			{
				Group: "batch",
				Kind:  "Job",
				OperationTypes: []*lifecycle.OperatorOperationType{
					{
						Type: lifecycle.OperatorOperationType_TYPE_CREATE,
					},
				},
			},
		},
	}, nil
}

// LifecycleHook is called when creating Kubernetes services
func (impl Implementation) LifecycleHook(
	ctx context.Context,
	request *lifecycle.OperatorLifecycleRequest,
) (*lifecycle.OperatorLifecycleResponse, error) {
	logger := log.FromContext(ctx).WithName("cnpg_i_skip_initdb")

	kind, err := utils.GetKind(request.GetObjectDefinition())
	if err != nil {
		logger.Error(err, "Failed to get object kind")
		return nil, err
	}

	operation := request.GetOperationType().GetType().Enum()
	if operation == nil {
		logger.Error(errors.New("no operation set"), "No operation type in request")
		return nil, errors.New("no operation set")
	}

	logger.Info("LifecycleHook called", "kind", kind, "operation", operation.String())

	//nolint: gocritic
	switch kind {
	case "Job":
		switch *operation {
		case lifecycle.OperatorOperationType_TYPE_CREATE:
			logger.Info("Intercepting Job CREATE operation")
			return impl.reconcileJob(ctx, request)
		}
	}

	logger.Info("No action taken for this object", "kind", kind, "operation", operation.String())
	return &lifecycle.OperatorLifecycleResponse{}, nil
}

// reconcileJob is called when intercepting Job creation
func (impl Implementation) reconcileJob(
	ctx context.Context,
	request *lifecycle.OperatorLifecycleRequest,
) (*lifecycle.OperatorLifecycleResponse, error) {
	logger := log.FromContext(ctx).WithName("cnpg_i_skip_initdb")

	cluster, err := decoder.DecodeClusterLenient(request.GetClusterDefinition())
	if err != nil {
		logger.Error(err, "Failed to decode cluster definition")
		return nil, err
	}

	logger.Info("Decoded cluster", "clusterName", cluster.Name, "namespace", cluster.Namespace)

	var job batchv1.Job
	if err := decoder.DecodeObjectLenient(request.GetObjectDefinition(), &job); err != nil {
		logger.Error(err, "Failed to decode job definition")
		return nil, err
	}

	jobRole := job.Labels["cnpg.io/jobRole"]
	logger.Info("Intercepted Job creation",
		"jobName", job.Name,
		"jobRole", jobRole,
		"clusterName", cluster.Name)

	// Check if this is an initdb job by label OR by name pattern
	// CNPG may not always set the jobRole label, but initdb jobs are named *-initdb
	isInitdbJob := jobRole == "initdb" || strings.HasSuffix(job.Name, "-initdb")

	if !isInitdbJob {
		logger.Info("Job is not an initdb job, skipping", "jobRole", jobRole, "jobName", job.Name)
		return &lifecycle.OperatorLifecycleResponse{}, nil
	}

	logger.Info("INITDB JOB DETECTED - Replacing with no-op",
		"jobName", job.Name,
		"clusterName", cluster.Name,
		"originalImage", cluster.Status.Image)

	// Definitely not so elegant, but this container has no side effects
	mutatedJob := job.DeepCopy()
	mutatedJob.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:  "skip",
			Image: cluster.Status.Image,
			Command: []string{
				"bash",
				"-c",
				"sleep 1",
			},
		},
	}

	patch, err := object.CreatePatch(mutatedJob, &job)
	if err != nil {
		logger.Error(err, "Failed to create patch for job mutation")
		return nil, err
	}

	logger.Info("Generated JSON patch for initdb skip",
		"patchSize", len(patch),
		"patch", string(patch))

	logger.Info("SUCCESS: Initdb job replaced with no-op sleep command")

	return &lifecycle.OperatorLifecycleResponse{
		JsonPatch: patch,
	}, nil
}
