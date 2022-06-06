/*
Copyright 2022 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"

	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/modules"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	log "github.com/sirupsen/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type DeckhouseConfigValidator struct {
	ModulesDir     string
	GlobalHooksDir string
}

func NewDeckhouseConfigValidator() *DeckhouseConfigValidator {
	return &DeckhouseConfigValidator{}
}

func (c *DeckhouseConfigValidator) Init() (err error) {
	log.Infof("Load OpenAPI schemas")
	return modules.Registry().Init(c.GlobalHooksDir, c.ModulesDir)
}

func (c *DeckhouseConfigValidator) Validate(_ context.Context, review *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhvalidating.ValidatorResult, error) {
	if review.Operation == kwhmodel.OperationDelete && review.Name == "global" {
		return rejectResult("deleting DeckhouseConfig/global is not allowed")
	}

	untypedCfg, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("expect DeckhouseConfig as unstructured, got %T", obj)
	}

	if untypedCfg.GetKind() != "DeckhouseConfig" {
		return nil, fmt.Errorf("expect DeckhouseConfig, got %s", untypedCfg.GetKind())
	}

	var cfg d8config_v1.DeckhouseConfig
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(untypedCfg.UnstructuredContent(), &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Name != "global" && !modules.Registry().HasModule(cfg.Name) {
		return allowResult(fmt.Sprintf("module name '%s' is unknown for deckhouse", cfg.Name))
	}

	if cfg.Spec.ConfigVersion == "" {
		return rejectResult("spec.configVersion is required")
	}

	err = d8config.ValidateValues(&cfg)
	if err != nil {
		return rejectResult(fmt.Sprintf("validate: %v", err))
	}

	return allowResult("")
}

func allowResult(warnMsg string) (*kwhvalidating.ValidatorResult, error) {
	var warnings []string
	if warnMsg != "" {
		warnings = []string{warnMsg}
	}
	return &kwhvalidating.ValidatorResult{
		Valid:    true,
		Warnings: warnings,
	}, nil
}

func rejectResult(msg string) (*kwhvalidating.ValidatorResult, error) {
	return &kwhvalidating.ValidatorResult{
		Valid:   false,
		Message: msg,
	}, nil
}
