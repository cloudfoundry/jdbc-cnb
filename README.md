# `jdbc-cnb`
The Cloud Foundry JDBC Buildpack is a Cloud Native Buildpack V3 that provides JDBC drivers to applications.

This buildpack is designed to work in collaboration with other buildpacks.

## Detection
The detection phase passes if:

* A service is bound with a payload containing a `binding_name`, `instance_name`, `label`, or `tag` containing either `mariadb` or `mysql` as a substring and build plan contains `jvm-application`
* and the application does not contain either mariadb-java-client-_version_.jar or mysql-connector-java-_version_.jar
    * Contributes `mariadb-jdbc` to the build plan

* A service is bound with a payload containing a `binding_name`, `instance_name`, `label`, or `tag` containing `postgres` as a substring and build plan contains `jvm-application`
* and the application does not contain postgresql-_version_.jar
    * Contributes `postgresql-jdbc` to the build plan

## Build
If the build plan contains

* `mariadb-jdbc`
  * Contributes the MariaDB JDBC Driver to a layer marked `launch`
  * Adds the MariaDB JDBC Driver to the classpath

* `postgresql-jdbc`
  * Contributes the PostgreSQL JDBC Driver to a layer marked `launch`
  * Adds the PostgreSQL JDBC Driver to the classpath

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0
