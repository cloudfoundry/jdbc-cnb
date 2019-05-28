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
	"fmt"
	"os"
	"reflect"
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

	if err := detect.BuildPlan.Init(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Build Plan: %s\n", err)
		os.Exit(101)
	}

	if code, err := d(detect); err != nil {
		detect.Logger.Info(err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func d(detect detect.Detect) (int, error) {
	bp := buildplan.BuildPlan{}

	if _, ok := detect.BuildPlan[jvmapplication.Dependency]; ok {
		if detect.Services.HasService("mariadb") || detect.Services.HasService("mysql") {
			if ok, err := helper.HasFile(detect.Application.Root, mp); err != nil {
				return detect.Error(102), err
			} else if !ok {
				bp[jdbc.MariaDBDependency] = detect.BuildPlan[jdbc.MariaDBDependency]
			}
		}

		if detect.Services.HasService("postgres") {
			if ok, err := helper.HasFile(detect.Application.Root, pp); err != nil {
				return detect.Error(102), err
			} else if !ok {
				bp[jdbc.PostgreSQLDependency] = detect.BuildPlan[jdbc.PostgreSQLDependency]
			}
		}
	}

	if reflect.DeepEqual(bp, buildplan.BuildPlan{}) {
		return detect.Fail(), nil
	}

	return detect.Pass(bp)
}
