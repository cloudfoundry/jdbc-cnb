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

package main

import (
	"testing"

	"github.com/buildpacks/libbuildpack/v2/buildplan"
	"github.com/buildpacks/libbuildpack/v2/detect"
	"github.com/cloudfoundry/jdbc-cnb/jdbc"
	"github.com/cloudfoundry/libcfbuildpack/v2/services"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("fails without service", func() {
			g.Expect(d(f.Detect)).To(gomega.Equal(detect.FailStatusCode))
		})

		it("passes with mariadb service", func() {
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jdbc.MariaDBDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: jdbc.MariaDBDependency},
				},
			}))
		})

		it("passes with mysql service", func() {
			f.AddService("mysql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jdbc.MariaDBDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: jdbc.MariaDBDependency},
				},
			}))
		})

		it("passes with postgresql service", func() {
			f.AddService("postgresql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(gomega.Equal(buildplan.Plans{
				Plan: buildplan.Plan{
					Provides: []buildplan.Provided{
						{Name: jdbc.PostgreSQLDependency},
					},
					Requires: []buildplan.Required{
						{Name: "jvm-application"},
						{Name: jdbc.PostgreSQLDependency},
					},
				},
			}))
		})

		it("passes with mariadb and postgresql services", func() {
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})
			f.AddService("postgresql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jdbc.MariaDBDependency},
					{Name: jdbc.PostgreSQLDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: jdbc.MariaDBDependency},
					{Name: jdbc.PostgreSQLDependency},
				},
			}))
		})

		it("passes with mysql and postgresql services", func() {
			f.AddService("mysql", services.Credentials{"test-key": "test-value"})
			f.AddService("postgresql", services.Credentials{"test-key": "test-value"})

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jdbc.MariaDBDependency},
					{Name: jdbc.PostgreSQLDependency},
				},
				Requires: []buildplan.Required{
					{Name: "jvm-application"},
					{Name: jdbc.MariaDBDependency},
					{Name: jdbc.PostgreSQLDependency},
				},
			}))
		})

		it("fails with mariadb service and a matching jar available", func() {
			f.AddService("mariadb", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "mariadb-java-client-1.2.3.jar")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.FailStatusCode))
		})

		it("fails with mysql service and a matching jar available", func() {
			f.AddService("mysql", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "mysql-connector-java-1.2.3.jar")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.FailStatusCode))
		})

		it("fails with postgres service and a matching jar available", func() {
			f.AddService("postgres", services.Credentials{"test-key": "test-value"})
			test.TouchFile(t, f.Detect.Application.Root, "subdir", "postgresql-1.2.3.jar")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.FailStatusCode))
		})

	}, spec.Report(report.Terminal{}))
}
