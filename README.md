postgres-version-manager-go
===========================
[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech)
[![License](https://img.shields.io/badge/license-Apache--2.0%20OR%20MIT%20OR%20CC0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![goreleaser](https://github.com/offscale/postgres-version-manager-go/actions/workflows/release.yml/badge.svg)](https://github.com/offscale/postgres-version-manager-go/actions/workflows/release.yml)

Like [`rvm`](https://rvm.io)/[`nvm`](https://github.com/nvm-sh/nvm) but for PostgreSQL. Installs any of 135+ different versions of PostgreSQL for your system.

## Development

Tested on Go 1.20, `git clone` this repo then `go build ./cmd`

## Usage

    PostgreSQL version manager
    Usage: pvm-go [--postgres-version POSTGRES-VERSION] [--port PORT] [--database DATABASE] [--username USERNAME] [--password PASSWORD] [--versionmanagerroot VERSIONMANAGERROOT] [--runtime-path RUNTIME-PATH] [--data-path DATA-PATH] [--binary-path BINARY-PATH] [--logs-path LOGS-PATH] [--locale LOCALE] [--binary-repository-url BINARY-REPOSITORY-URL] [--no-remote] <command> [<args>]

    Commands:
    env                    Print out database connection string
    start                  Start specified PostgreSQL server
    stop                   Stop specific (running) PostgreSQL server
    install                Install specified PostgreSQL version
    ls                     List what versions of PostgreSQL are installed
    ls-remote              List what versions of PostgreSQL are available
    get-data-path          Get data path, i.e., where pg_hba and postgres.conf are for specified PostgreSQL version
    install-service        Install service (daemon), e.g., systemd

#### Global options

Common to all subcommands

    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: PGPORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --versionmanagerroot VERSIONMANAGERROOT [default: $HOME/postgres-version-manager, env: VERSION_MANAGER_ROOT]
    --runtime-path RUNTIME-PATH [default: $HOME/postgres-version-manager/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: $HOME/postgres-version-manager/latest/data, env: PGDATA]
    --binary-path BINARY-PATH [default: $HOME/postgres-version-manager/latest, env: BINARY_PATH]
    --logs-path LOGS-PATH [default: $HOME/postgres-version-manager/latest/logs, env: LOGS_PATH]
    --locale LOCALE [default: en_US.UTF-8, env: LC_ALL]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote            Disable HTTPS calls for everything except 'install' [default: false]
    --help, -h             display this help and exit

### `env`

Print out database connection string

    Usage: pvm-go env

### `start`

Start specified PostgreSQL server

    Usage: pvm-go start [--no-install] [POSTGRES_VERSION]
    
    Positional arguments:
    POSTGRES_VERSION
    
    Options:
    --no-install           Inverts default of installing nonexistent version [default: false]

### `stop`

Stop specific (running) PostgreSQL server

    Usage: pvm-go stop [POSTGRES_VERSION]
    
    Positional arguments:
    POSTGRES_VERSION

### `install`

Install specified PostgreSQL version
    
    Usage: pvm-go install [POSTGRES_VERSION]
    
    Positional arguments:
    POSTGRES_VERSION

### `ls`

List what versions of PostgreSQL are installed

    Usage: pvm-go ls

### `ls-remote`

List what versions of PostgreSQL are available

    Usage: pvm-go ls-remote

### `get-data-path`

Get data path, i.e., where pg_hba and postgres.conf are for specified PostgreSQL version

    Usage: pvm-go get-data-path

### `install-service`

Install service (daemon), e.g., systemd

    Usage: pvm-go install-service

    Commands:
    systemd                Install systemd service

#### `systemd`

Install systemd service

    Usage: pvm-go install-service systemd [--service-install-path SERVICE-INSTALL-PATH]
    
    Options:
    --service-install-path SERVICE-INSTALL-PATH [default: /etc/systemd/system/postgresql.service]

---

## License

Licensed under any of:

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or <https://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](LICENSE-MIT) or <https://opensource.org/licenses/MIT>)
- CC0 license ([LICENSE-CC0](LICENSE-CC0) or <https://creativecommons.org/publicdomain/zero/1.0/legalcode>)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
licensed as above, without any additional terms or conditions.
