// Copyright (c) 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhooks

import (
	"context"
	"net/http"

	"emperror.dev/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/go-logr/logr"

	clusterregistrycontrollerapiv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
)

// ClusterValidator validates cluster CRs of the cluster registry.
type ClusterValidator struct {
	// logger is the log interface to use inside the validator.
	logger logr.Logger

	// manager is responsible for handling the communication between the
	// validator and the Kubernetes API server.
	manager ctrl.Manager

	// decoder is responsible for decoding the webhook request into structured
	// data.
	decoder *admission.Decoder
}

// NewClusterValidator instantiates a cluster CR validator set to work on the
// specified Kubernetes API server.
func NewClusterValidator(logger logr.Logger, manager ctrl.Manager) *ClusterValidator {
	return &ClusterValidator{
		logger:  logger,
		manager: manager,
		decoder: nil,
	}
}

// Handle handles the validator's admission requests and determines whether the
// specified request can be allowed.
func (validator *ClusterValidator) Handle(ctx context.Context, request admission.Request) admission.Response {
	validator.logger.Info("validating cluster CR", "request", request)

	newClusterCR := &clusterregistrycontrollerapiv1alpha1.Cluster{}

	err := validator.decoder.Decode(request, newClusterCR)
	if err != nil {
		err = errors.Wrap(err, "decoding admission request as new cluster CR failed")

		validator.logger.Error(err, "validating cluster CR failed", "request", request)

		return admission.Errored(http.StatusBadRequest, err)
	}

	k8sClient := validator.manager.GetClient()

	oldClusterCRs := &clusterregistrycontrollerapiv1alpha1.ClusterList{}
	err = k8sClient.List(ctx, oldClusterCRs)
	if err != nil {
		err = errors.Wrap(err, "listing existing cluster CRs failed")

		validator.logger.Error(err, "validating cluster CR failed")

		return admission.Errored(http.StatusConflict, err)
	}

	err = validator.validateUniqueLocalCluster(ctx, newClusterCR, oldClusterCRs)
	if err != nil {
		err = errors.Wrap(err, "validating unique local cluster CR failed")

		validator.logger.Error(err, "validating cluster CR failed", "newClusterCR", newClusterCR)

		return admission.Errored(http.StatusConflict, err)
	}

	validator.logger.Info("validating cluster CR succeeded", "request", request)

	return admission.Allowed("")
}

// validateUniqueLocalCluster returns an error if the specified new cluster CR
// has a matching local clusterID among the old cluster CRs, nil otherwise.
func (validator *ClusterValidator) validateUniqueLocalCluster(
	ctx context.Context,
	newClusterCR *clusterregistrycontrollerapiv1alpha1.Cluster,
	oldClusterCRs *clusterregistrycontrollerapiv1alpha1.ClusterList,
) error {
	// Note: the new cluster CR usually doesn't have a status.type set yet, so
	// we only return on explicit peer cluster status.type, we examine new
	// cluster CRs with local or nil status.type further.
	if newClusterCR.Status.Type == clusterregistrycontrollerapiv1alpha1.ClusterTypePeer {
		validator.logger.Info("skipping peer cluster unique local validation", "newClusterCR", newClusterCR)

		return nil // Note: for peer cluster CRs this validation is a no-operation.
	}

	for _, oldClusterCR := range oldClusterCRs.Items {
		if oldClusterCR.Status.Type == clusterregistrycontrollerapiv1alpha1.ClusterTypeLocal &&
			oldClusterCR.Spec.ClusterID == newClusterCR.Spec.ClusterID &&
			oldClusterCR.Name != newClusterCR.Name { // Note: in case a cluster CR is updated in place it should not be rejected because of the older version of itself.
			err := errors.Errorf(
				"a local cluster CR with name %s already exists with the same clusterID %s",
				oldClusterCR.ObjectMeta.Name,
				newClusterCR.Spec.ClusterID,
			)

			validator.logger.Error(
				err,
				"validating unique local cluster CR failed",
				"newClusterCR", newClusterCR,
				"oldClusterCR", oldClusterCR,
			)

			return err
		}
	}

	validator.logger.Info(
		"validating new cluster CR unique local cluster ID succeeded",
		"newClusterCRClusterID", newClusterCR.Spec.ClusterID,
	)

	return nil
}

// InjectDecoder sets the cluster CR decoder object.
func (validator *ClusterValidator) InjectDecoder(decoder *admission.Decoder) error {
	validator.decoder = decoder

	return nil
}
