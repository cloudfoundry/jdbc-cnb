/*
 * Copyright 2018-2019 the original author or authors.
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
	"fmt"
	"os"
	"regexp"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jdbc-cnb/jdbc"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/helper"
)

var mp = regexp.MustCompile(".*(mariadb-java-client|mysql-connector-java)-[\\d].*\\.jar")

var pp = regexp.MustCompile(".*postgresql-([\\d].*)\\.jar")

func main() {
	detect, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err)
		os.Exit(101)
	}

	if code, err := d(detect); err != nil {
		detect.Logger.TerminalError(detect.Buildpack, err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func d(detect detect.Detect) (int, error) {
	m := detect.Services.HasService("mariadb") || detect.Services.HasService("mysql")
	mf, err := helper.HasFile(detect.Application.Root, mp)
	if err != nil {
		return detect.Error(102), err
	}

	p := detect.Services.HasService("postgres")
	pf, err := helper.HasFile(detect.Application.Root, pp)
	if err != nil {
		return detect.Error(102), err
	}

	if !(m && !mf) && !(p && !pf) {
		return detect.Fail(), nil
	}

	q := buildplan.Plan{
		Requires: []buildplan.Required{
			{Name: jvmapplication.Dependency},
		},
	}

	if m && !mf {
		q.Provides = append(q.Provides, buildplan.Provided{Name: jdbc.MariaDBDependency})
		q.Requires = append(q.Requires, buildplan.Required{Name: jdbc.MariaDBDependency})
	}

	if p && !pf {
		q.Provides = append(q.Provides, buildplan.Provided{Name: jdbc.PostgreSQLDependency})
		q.Requires = append(q.Requires, buildplan.Required{Name: jdbc.PostgreSQLDependency})
	}

	return detect.Pass(q)
}
