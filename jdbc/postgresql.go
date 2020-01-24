/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jdbc

import (
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/helper"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
)

// PostgreSQLDependency indicates that a JVM application should be run with PostgreSQL JDBC enabled.
const PostgreSQLDependency = "postgresql-jdbc"

// PostgreSQL represents a PostgreSQL contribution by the buildpack.
type PostgreSQL struct {
	layer layers.DependencyLayer
}

// Contribute makes the contribution to launch.
func (p PostgreSQL) Contribute() error {
	return p.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.Body("Copying to %s", layer.Root)

		destination := filepath.Join(layer.Root, layer.ArtifactName())

		if err := helper.CopyFile(artifact, destination); err != nil {
			return err
		}

		return layer.PrependPathLaunchEnv("CLASSPATH", "%s", destination)
	}, layers.Launch)
}

// NewPostgreSQL creates a new PostgreSQL instance.
func NewPostgreSQL(build build.Build) (PostgreSQL, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(PostgreSQLDependency)
	if err != nil {
		return PostgreSQL{}, false, err
	} else if !ok {
		return PostgreSQL{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return PostgreSQL{}, false, err
	}

	dep, err := deps.Best(PostgreSQLDependency, p.Version, build.Stack)
	if err != nil {
		return PostgreSQL{}, false, err
	}

	return PostgreSQL{build.Layers.DependencyLayer(dep)}, true, nil
}
