// Copyright The OpenTelemetry Authors
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

package k8s // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter/internal/metadata/internal/k8s"

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter/internal/metadata/provider"

	"go.uber.org/zap"
)

var _ provider.HostnameProvider = (*Provider)(nil)

type Provider struct {
	logger              *zap.Logger
	nodeNameProvider    nodeNameProvider
	clusterNameProvider provider.ClusterNameProvider
}

// Hostname returns the Kubernetes node name followed by the cluster name if available.
func (p *Provider) Hostname(ctx context.Context) (string, error) {
	nodeName, err := p.nodeNameProvider.NodeName(ctx)
	if err != nil {
		return "", fmt.Errorf("node name not available: %w", err)
	}

	clusterName, err := p.clusterNameProvider.ClusterName(ctx)
	if err != nil {
		p.logger.Debug("failed to get valid cluster name", zap.Error(err))
		return nodeName, nil
	}

	return fmt.Sprintf("%s-%s", nodeName, clusterName), nil
}

// NewProvider creates a new Kubernetes hostname provider.
func NewProvider(logger *zap.Logger, clusterProvider provider.ClusterNameProvider) (*Provider, error) {
	return &Provider{
		logger:              logger,
		nodeNameProvider:    newNodeNameProvider(),
		clusterNameProvider: clusterProvider,
	}, nil
}
