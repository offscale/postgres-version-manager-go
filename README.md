postgres-version-manager-go
===========================
[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech)
[![License](https://img.shields.io/badge/license-Apache--2.0%20OR%20MIT%20OR%20CC0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![goreleaser](https://github.com/offscale/postgres-version-manager-go/actions/workflows/release.yml/badge.svg)](https://github.com/offscale/postgres-version-manager-go/actions/workflows/release.yml)

Like [`rvm`](https://rvm.io)/[`nvm`](https://github.com/nvm-sh/nvm) but for PostgreSQL. Installs any of 135+ different versions of PostgreSQL for your system.

## Development

Tested on Go 1.20, `git clone` this repo then `go build .`

## Usage

    PostgreSQL version manager
    Usage: postgres-version-manager-go [--postgres-version POSTGRES-VERSION] [--port PORT] [--database DATABASE] [--username USERNAME] [--password PASSWORD] [--runtime-path RUNTIME-PATH] [--data-path DATA-PATH] [--binary-path BINARY-PATH] [--locale LOCALE] [--binary-repository-url BINARY-REPOSITORY-URL] [--no-remote] <command> [<args>]
    
    Options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit
    
    Commands:
    env
    start
    stop
    install
    ls-remote

### `env`

    PostgreSQL version manager
    Usage: postgres-version-manager-go env
    
    Global options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit

### `start`

    PostgreSQL version manager
    Usage: postgres-version-manager-go start [--pid PID]
    
    Options:
    --pid PID              If PID provided and exists, will stop that process. [default: -1]
    
    Global options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit

### `stop`

    PostgreSQL version manager
    Usage: postgres-version-manager-go stop [--pid PID]
    
    Options:
    --pid PID [default: -1]
    
    Global options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit

### `install`

    PostgreSQL version manager
    Usage: postgres-version-manager-go install [POSTGRESVERSION]
    
    Positional arguments:
    POSTGRESVERSION
    
    Global options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit

### `ls-remote`

    PostgreSQL version manager
    Usage: postgres-version-manager-go ls-remote
    
    Global options:
    --postgres-version POSTGRES-VERSION [default: latest, env: POSTGRES_VERSION]
    --port PORT, -p PORT [default: 5432, env: POSTGRES_PORT]
    --database DATABASE, -d DATABASE [default: database, env: POSTGRES_DATABASE]
    --username USERNAME, -u USERNAME [default: username, env: POSTGRES_USERNAME]
    --password PASSWORD [default: password, env: POSTGRES_PASSWORD]
    --runtime-path RUNTIME-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/run, env: RUNTIME_PATH]
    --data-path DATA-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest/data, env: DATA_PATH]
    --binary-path BINARY-PATH [default: /home/samuel/.config/postgres-version-manager-go/latest, env: BINARY_PATH]
    --locale LOCALE [default: en_US.UTF-8,env:LANGUAGE]
    --binary-repository-url BINARY-REPOSITORY-URL [default: https://repo1.maven.org/maven2, env: BINARY_REPOSITORY_URL]
    --no-remote [default: false]
    --help, -h             display this help and exit

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
