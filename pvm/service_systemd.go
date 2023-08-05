package pvm

import (
	"fmt"
	"os"
)

func systemdInstall(args *Args) error {
	systemd :=
		fmt.Sprintf(`[Unit]
Description=PostgreSQL %s database server
After=network.target

[Service]
Type=forking

User=%s
Group=%s

OOMScoreAdjust=-1000
Environment=PG_OOM_ADJUST_FILE=/proc/self/oom_score_adj
Environment=PG_OOM_ADJUST_VALUE=0

Environment=PGSTARTTIMEOUT=270

Environment=PGDATA=%s
Environment=PGPORT=%d


ExecStart=%s/pg_ctl start -D ${PGDATA} -s -w -t ${PGSTARTTIMEOUT}
ExecStop=%s/pg_ctl stop -D ${PGDATA} -s -m fast
ExecReload=%s/pg_ctl reload -D ${PGDATA} -s

TimeoutSec=300

[Install]
WantedBy=multi-user.target
`, args.PostgresVersion, args.InstallService.Systemd.User, args.InstallService.Systemd.Group, args.DataPath, args.Port, args.BinariesPath, args.BinariesPath, args.BinariesPath)
	var f *os.File
	var err error

	if f, err = os.Create(args.InstallService.Systemd.ServiceInstallPath); err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	_, err = f.WriteString(systemd)

	return err
}
