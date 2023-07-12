package pvm

var versionsFromMaven []string = nil

// ConfigStruct originally from github.com/fergusstrange/embedded-postgres@v1.20.0/config.go
type ConfigStruct struct {
	PostgresVersion     string `arg:"--postgres-version,env:POSTGRES_VERSION" default:"latest"`
	Port                uint32 `arg:"-p,env:PGPORT" default:"5432"`
	Database            string `arg:"-d,env:POSTGRES_DATABASE" default:"database"`
	Username            string `arg:"-u,env:POSTGRES_USERNAME" default:"username"`
	Password            string `arg:"env:POSTGRES_PASSWORD" default:"password"`
	VersionManagerRoot  string `arg:"env:VERSION_MANAGER_ROOT"`
	RuntimePath         string `arg:"--runtime-path,env:RUNTIME_PATH"`
	DataPath            string `arg:"--data-path,env:PGDATA"`
	BinariesPath        string `arg:"--binary-path,env:BINARY_PATH"`
	LogsPath            string `arg:"--logs-path,env:LOGS_PATH"`
	Locale              string `arg:"--locale,env:LC_ALL" default:"en_US.UTF-8"`
	BinaryRepositoryURL string `arg:"--binary-repository-url,env:BINARY_REPOSITORY_URL" default:"https://repo1.maven.org/maven2"`
}

type EnvCmd struct{}

type GetDataPathCmd struct{}

type InstallCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type InstallServiceForSystemdCmd struct {
	ServiceInstallPath string `arg:"--service-install-path" default:"/etc/systemd/system/postgresql.service"`
}

type InstallServiceCmd struct {
	Systemd *InstallServiceForSystemdCmd `arg:"subcommand:systemd" help:"Install systemd service"`
}

type LsCmd struct{}

type LsRemoteCmd struct{}

type StartCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
	NoInstall       bool   `arg:"--no-install" default:"false" help:"Inverts default of installing nonexistent version"`
}

type StopCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type Args struct {
	ConfigStruct
	NoRemote       bool               `arg:"--no-remote" default:"false" help:"Disable HTTPS calls for everything except 'install'"`
	Env            *EnvCmd            `arg:"subcommand:env" help:"Print out database connection string"`
	GetDataPath    *GetDataPathCmd    `arg:"subcommand:get-data-path" help:"Get data path, i.e., where pg_hba and postgres.conf are for specified PostgreSQL version"`
	Install        *InstallCmd        `arg:"subcommand:install" help:"Install specified PostgreSQL version"`
	InstallService *InstallServiceCmd `arg:"subcommand:install-service" help:"Install service (daemon), e.g., systemd"`
	Ls             *LsCmd             `arg:"subcommand:ls" help:"List what versions of PostgreSQL are installed"`
	LsRemote       *LsRemoteCmd       `arg:"subcommand:ls-remote" help:"List what versions of PostgreSQL are available"`
	Start          *StartCmd          `arg:"subcommand:start" help:"Start specified PostgreSQL server"`
	Stop           *StopCmd           `arg:"subcommand:stop" help:"Stop specific (running) PostgreSQL server"`
}

func (Args) Description() string {
	return "PostgreSQL version manager"
}

func (Args) Version() string {
	return "pvm 0.0.9"
}
