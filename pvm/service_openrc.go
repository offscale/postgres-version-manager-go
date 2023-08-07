package pvm

import (
	"fmt"
	"os"
	"strings"
)

func openRcInstall(args *Args) error {
	// https://git.alpinelinux.org/aports/plain/main/postgresql/postgresql.confd?h=3.13-stable
	postgresqlConfd :=
		fmt.Sprintf(`
# Which port and socket to bind PostgreSQL.
# This may be overriden in postgresql.conf.
port="%d"

# How long to wait for server to start in seconds.
#start_timeout=10

# Number of seconds to wait for clients to disconnect from the server before
# shutting down. Set to zero to disable this timeout.
#nice_timeout=60

# Timeout in seconds for rude quit - forecfully disconnect clients from server
# and shut down. This is performed after nice_timeout exceeded. Terminated
# client connections have their open transactions rolled back.
# Set "rude_quit=no" to disable.
#rude_quit="yes"
#rude_timeout=30

# Timeout in seconds for force quit - if the server still fails to shutdown,
# you can force it to quit and a recover-run will execute on the next startup.
# Set "force_quit=yes" to enable.
#force_quit="no"
#force_timeout="2"

# Extra options to run postmaster with, e.g.:
#   -N is the maximal number of client connections
#   -B is the number of shared buffers (has to be at least 2x the value for -N)
# Please read man postgres(1) for more options. Many of these options can be
# set directly in the configuration file.
#pg_opts="-N 512 -B 1024"

# Pass extra environment variables. If you have to export environment variables
# for the database process, this can be done here.
# Don't forget to escape quotes.
#env_vars="PGPASSFILE=\"/path/to/.pgpass\""

# Location of postmaster.log. Default is $data_dir/postmaster.log.
logfile="/var/log/postgresql/postmaster.log"

# Automatically set up a new database if missing on startup.
#auto_setup="yes"


##############################################################################
#
# The following values should NOT be arbitrarily changed!
#
# The initscript uses these variables to inform PostgreSQL where to find
# its data directory and configuration files.

# Where the data directory is located/to be created.
data_dir="%s"

# Location of configuration files. Default is $data_dir.
conf_dir="%s"

# Additional options to pass to initdb.
# See man initdb(1) for available options.
#initdb_opts="--locale=en_US.UTF-8"
`, args.Port, args.DataPath, args.DataPath)
	var f0 *os.File
	var err error

	if f0, err = os.OpenFile(args.InstallService.OpenRc.ConfigInstallPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f0)
	if _, err = f0.WriteString(postgresqlConfd); err != nil {
		return err
	}

	// https://git.alpinelinux.org/aports/tree/main/postgresql/postgresql.initd?h=3.13-stable
	var postgresqlInitd strings.Builder
	var f1 *os.File
	postgresqlInitd.WriteString(fmt.Sprintf(`#!/sbin/openrc-run

description="PostgreSQL @VERSION@ server"

extra_started_commands="stop_fast stop_force stop_smart reload reload_force"
description_stop_fast="Stop using Fast Shutdown mode (SIGINT)"
description_stop_force="Stop using Immediate Shutdown mode (SIGQUIT)"
description_stop_smart="Stop using Smart Shutdown mode (SIGTERM)"
description_reload="Reload configuration"
description_reload_force="Reload configuration and restart if needed"

extra_stopped_commands="setup"
description_setup="Initialize a new PostgreSQL cluster"

# Note: Uppercase variables are here for backward compatibility.

: ${user:=${PGUSER:-"%s"}}
: ${group:=${PGGROUP:-"%s"}}`, args.InstallService.OpenRc.User, args.InstallService.OpenRc.Group))
	postgresqlInitd.WriteString(`
: ${auto_setup:=${AUTO_SETUP:-"yes"}}
: ${start_timeout:=${START_TIMEOUT:-10}}
: ${nice_timeout:=${NICE_TIMEOUT:-60}}
: ${rude_quit:=${RUDE_QUIT:-"yes"}}
: ${rude_timeout:=${RUDE_TIMEOUT:-30}}
: ${force_quit:=${FORCE_QUIT:-"no"}}
: ${force_timeout:=${FORCE_TIMEOUT:-2}}`)
	postgresqlInitd.WriteString(fmt.Sprintf(`
: ${data_dir:=${PGDATA:-"%s"}}
: ${conf_dir:=$data_dir}
: ${env_vars:=${PG_EXTRA_ENV:-}}
: ${initdb_opts:=${PG_INITDB_OPTS:-}}
: ${logfile:="$data_dir/postmaster.log"}
: ${pg_opts:=${PGOPTS:-}}
: ${port:=${PGPORT:-%d}}

command="%s/postgres"`, args.DataPath, args.Port, args.BinariesPath))
	postgresqlInitd.WriteString(`
conffile="$conf_dir/postgresql.conf"
pidfile="$data_dir/postmaster.pid"
start_stop_daemon_args="
	--user $user
	--group $group
	--pidfile $pidfile
	--wait 100"

depend() {
	use net
	after firewall

	if [ "$(get_config log_destination)" = "syslog" ]; then
		use logger
	fi
}

start_pre() {
	check_deprecated_var WAIT_FOR_START start_timeout
	check_deprecated_var WAIT_FOR_DISCONNECT nice_timeout
	check_deprecated_var WAIT_FOR_CLEANUP rude_timeout
	check_deprecated_var WAIT_FOR_QUIT force_timeout

	if [ ! -d "$data_dir/base" ]; then
		if yesno "$auto_setup"; then
			setup || return 1
		else
			eerror "Database not found at: $data_dir"
			eerror "Please make sure that 'data_dir' points to the right path."
			eerror "You can run '/etc/init.d/postgresql setup' to setup a new database cluster."
			return 1
		fi
	fi

	local socket_dirs=$(get_config "unix_socket_directories" "/run/postgresql")
	local port=$(get_config "port" "$port")

	start_stop_daemon_args="$start_stop_daemon_args --env PGPORT=$port"

	(
		# Set the proper permission for the socket paths and create them if
		# then don't exist.
		set -f; IFS=","
		for dir in $socket_dirs; do
			if [ -e "${dir%/}/.s.PGSQL.$port" ]; then
				eerror "Socket conflict. A server is already listening on:"
				eerror "    ${dir%/}/.s.PGSQL.$port"
				eerror "Hint: Change 'port' to listen on a different socket."
				return 1
			elif [ "${dir%/}" != "/tmp" ]; then
				checkpath -d -m 1775 -o $user:$group "$dir"
			fi
		done
	)
}

start() {
	local retval

	ebegin "Starting PostgreSQL"

	local var; for var in $env_vars; do
		start_stop_daemon_args="$start_stop_daemon_args --env $var"
	done

	rm -f "$pidfile"
	start-stop-daemon --start \
		$start_stop_daemon_args \
		--exec /usr/bin/pg_ctl \
		-- start \
			--silent \
			-w --timeout="$start_timeout" \
			--log="$logfile" \
			--pgdata="$conf_dir" \
			-o "--data-directory=$data_dir $pg_opts"
	retval=$?

	if [ $retval -ne 0 ]; then
		eerror "Check the log for a possible explanation of the above error:"
		eerror "    $logfile"
	fi
	eend $retval
}

stop() {
	local retry="SIGTERM/$nice_timeout"

	yesno "$rude_quit" \
		&& retry="$retry/SIGINT/$rude_timeout" \
		|| rude_timeout=0

	yesno "$force_quit" \
		&& retry="$retry/SIGQUIT/$force_timeout" \
		|| force_timeout=0

	local seconds=$(( $nice_timeout + $rude_timeout + $force_timeout ))

	ebegin "Stopping PostgreSQL (this can take up to $seconds seconds)"

	start-stop-daemon --stop \
		--exec "$command" \
		--retry "$retry" \
		--progress \
		--pidfile "$pidfile"
	eend $?
}

stop_smart() {
	_stop SIGTERM "smart shutdown"
}

stop_fast() {
	_stop SIGINT "fast shutdown"
}

stop_force() {
	_stop SIGQUIT "immediate shutdown"
}

_stop() {
	ebegin "Stopping PostgreSQL ($2)"

	start-stop-daemon --stop \
		--exec "$command" \
		--signal "$1" \
		--pidfile "$pidfile" \
		&& mark_service_stopped "$RC_SVCNAME"
	eend $?
}

reload() {
	ebegin "Reloading PostgreSQL configuration"

	start-stop-daemon --signal HUP --pidfile "$pidfile" && check_config_errors
	local retval=$?

	is_pending_restart || true

	eend $retval
}

reload_force() {
	ebegin "Reloading PostgreSQL configuration"

	start-stop-daemon --signal HUP --pidfile "$pidfile" && check_config_errors
	local retval=$?

	if [ $retval -eq 0 ] && is_pending_restart; then
		rc-service --nodeps "$RC_SVCNAME" restart
		retval=$?
	fi
	eend $retval
}

setup() {
	local bkpdir

	ebegin "Creating a new PostgreSQL database cluster"

	if [ -d "$data_dir/base" ]; then
		eend 1 "$data_dir/base already exists!"; return 1
	fi

	# If data_dir exists, backup configs.
	if [ -d "$data_dir" ]; then
		bkpdir="$(mktemp -d)"
		find "$data_dir" -type f -name "*.conf" -maxdepth 1 \
			-exec mv -v {} "$bkpdir"/ \;
		rm -rf "$data_dir"/*
	fi

	install -d -m 0700 -o $user -g $group "$data_dir"
	install -d -m 0750 -o $user -g $group "$conf_dir"

	cd "$data_dir"  # to avoid the: could not change directory to "/root"
	su $user -c "/usr/bin/initdb $initdb_opts --pgdata $data_dir"
	local retval=$?

	if [ -d "$bkpdir" ]; then
		# Move backuped configs back.
		mv -v "$bkpdir"/* "$data_dir"/
		rm -rf "$bkpdir"
	fi

	if [ "${data_dir%/}" != "${conf_dir%/}" ]; then
		# Move configs from data_dir to conf_dir and symlink them to data_dir.
		local name newname
		for name in postgresql.conf pg_hba.conf pg_ident.conf; do
			newname="$name"
			[ ! -e "$conf_dir"/$name ] || newname="$name.new"

			mv "$data_dir"/$name "$conf_dir"/$newname
			ln -s "$conf_dir"/$name "$data_dir"/$name
		done
	fi

	eend $retval
} 


get_config() {
	local name="$1"
	local default="${2:-}"

	if [ ! -f "$conffile" ]; then
		printf '%s\n' "$default"
		return 1
	fi
	sed -En "/^\s*${name}\b/{                      # find line starting with the name
		  s/^\s*${name}\s*=?\s*([^#]+).*/\1/;  # capture the value
		  s/\s*$//;                            # trim trailing whitespaces
		  s/^['\"](.*)['\"]$/\1/;              # remove delimiting quotes
		  p
		}" "$conffile" \
		| grep . || printf '%s\n' "$default"
}

check_config_errors() {
	local out; out=$(psql_command "
		select
		  sourcefile || ': line ' || sourceline || ': ' || error ||
		    case when name is not null
		    then ': ' || name || ' = ''' || setting || ''''
		    else ''
		    end
		from pg_file_settings
		where error is not null
		  and name not in (select name from pg_settings where pending_restart = true);
		")
	if [ $? -eq 0 ] && [ "$out" ]; then
		eerror 'Configuration file contains errors:'
		printf '%s\n' "$out" | while read line; do
			eerror "  $line"
		done
		return 1
	fi
}

is_pending_restart() {
	local out; out=$(psql_command "select name from pg_settings where pending_restart = true;")

	if [ $? -eq 0 ] && [ "$out" ]; then
		ewarn 'PostgreSQL must be restarted to apply changes in the following parameters:'
		local line; for line in $out; do
			ewarn "  $line"
		done
		return 0
	fi
	return 1
}

check_deprecated_var() {
	local old_name="$1"
	local new_name="$2"

	if [ -n "$(getval "$old_name")" ]; then
		ewarn "Variable '$old_name' has been removed, please use '$new_name' instead."
	fi
}

getval() {
	eval "printf '%s\n' \"\$$1\""
}

psql_command() {
	su $user -c "psql --no-psqlrc --no-align --tuples-only -q -c \"$1\""
}`)

	if f1, err = os.OpenFile(args.InstallService.OpenRc.ServiceInstallPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f1)
	_, err = f1.WriteString(postgresqlInitd.String())

	return err
}
