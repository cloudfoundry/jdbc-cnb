/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the License);
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

package main

import (
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/detect"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/services"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("fails without service", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})

			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

		it("fails without jvm-application", func() {
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

		it("passes with mariadb service", func() {
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})

			g.Expect(d(f.Detect)).To(Equal(detect.PassStatusCode))
		})

		it("passes with mysql service", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			f.AddService("mysql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(Equal(detect.PassStatusCode))
		})

		it("passes with postgresql service", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			f.AddService("postgresql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(Equal(detect.PassStatusCode))
		})

		it("fails with mariadb service and a matching jar available", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "mariadb-java-client-1.2.3.jar")

			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

		it("fails with mysql service and a matching jar available", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			f.AddService("mysql", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "mysql-connector-java-1.2.3.jar")

			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

		it("fails with postgres service and a matching jar available", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			f.AddService("postgres", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "subdir", "postgresql-1.2.3.jar")

			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

	}, spec.Report(report.Terminal{}))
}
