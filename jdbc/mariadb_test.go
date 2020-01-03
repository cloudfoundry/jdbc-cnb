/*
 * Copyright 2019-2020 the original author or authors.
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

package jdbc_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/jdbc-cnb/jdbc"
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMariaDB(t *testing.T) {
	spec.Run(t, "MariaDB", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan does exist", func() {
			f.AddPlan(buildpackplan.Plan{Name: jdbc.MariaDBDependency})
			f.AddDependency(jdbc.MariaDBDependency, filepath.Join("testdata", "stub-mariadb-java-client.jar"))

			_, ok, err := jdbc.NewMariaDB(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := jdbc.NewMariaDB(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("contributes driver", func() {
			f.AddPlan(buildpackplan.Plan{Name: jdbc.MariaDBDependency})
			f.AddDependency(jdbc.MariaDBDependency, filepath.Join("testdata", "stub-mariadb-java-client.jar"))

			y, ok, err := jdbc.NewMariaDB(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())

			g.Expect(y.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("mariadb-jdbc")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "stub-mariadb-java-client.jar")).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HavePrependPathLaunchEnvironment("CLASSPATH", ":%s",
				filepath.Join(layer.Root, "stub-mariadb-java-client.jar")))
		})
	}, spec.Report(report.Terminal{}))
}
