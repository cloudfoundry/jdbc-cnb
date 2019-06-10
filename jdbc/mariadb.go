/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// MariaDBDependency indicates that a JVM application should be run with MariaDB JDBC enabled.
const MariaDBDependency = "mariadb-jdbc"

// MariaDB represents a MariaDB contribution by the buildpack.
type MariaDB struct {
	layer layers.DependencyLayer
}

// Contribute makes the contribution to launch.
func (m MariaDB) Contribute() error {
	return m.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Copying to %s", layer.Root)

		destination := filepath.Join(layer.Root, layer.ArtifactName())

		if err := helper.CopyFile(artifact, destination); err != nil {
			return err
		}

		return layer.AppendPathLaunchEnv("CLASSPATH", ":%s", destination)
	}, layers.Launch)
}

// NewMariaDB creates a new MariaDB instance.
func NewMariaDB(build build.Build) (MariaDB, bool, error) {
	bp, ok := build.BuildPlan[MariaDBDependency]
	if !ok {
		return MariaDB{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return MariaDB{}, false, err
	}

	dep, err := deps.Best(MariaDBDependency, bp.Version, build.Stack)
	if err != nil {
		return MariaDB{}, false, err
	}

	return MariaDB{build.Layers.DependencyLayer(dep)}, true, nil
}
