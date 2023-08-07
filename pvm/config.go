package pvm

var versionsFromMaven []string = nil

// ConfigStruct originally from github.com/fergusstrange/embedded-postgres@v1.20.0/config.go
type ConfigStruct struct {
	PostgresVersion     string `arg:"--postgres-version,env:POSTGRES_VERSION" default:"latest" json:"-"`
	Port                uint32 `arg:"-p,env:PGPORT" default:"5432"`
	Database            string `arg:"-d,env:POSTGRES_DATABASE" default:"database"`
	Username            string `arg:"-u,env:POSTGRES_USERNAME" default:"username"`
	Password            string `arg:"env:POSTGRES_PASSWORD" default:"password"`
	VersionManagerRoot  string `arg:"--version-manager-root,env:VERSION_MANAGER_ROOT" placeholder:"VERSION_MANAGER_ROOT"`
	RuntimePath         string `arg:"--runtime-path,env:RUNTIME_PATH"`
	DataPath            string `arg:"--data-path,env:PGDATA"`
	BinariesPath        string `arg:"--binary-path,env:BINARY_PATH"`
	LogsPath            string `arg:"--logs-path,env:LOGS_PATH"`
	Locale              string `arg:"--locale,env:LC_ALL" default:"en_US.UTF-8"`
	BinaryRepositoryURL string `arg:"--binary-repository-url,env:BINARY_REPOSITORY_URL" default:"https://repo1.maven.org/maven2"`
}

type ConfigStructs []ConfigStruct

type DownloadCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type EnvCmd struct{}

type GetPathCmd struct {
	DirectoryToFind string `arg:"positional" help:"bin|data|log|runtime"`
}

type InstallCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type InstallServiceForOpenRcCmd struct {
	Group              string `arg:"--group" default:"postgres"`
	ConfigInstallPath  string `arg:"--service-install-path" default:"/etc/conf.d/postgresql"`
	ServiceInstallPath string `arg:"--service-install-path" default:"/etc/init.d/postgresql"`
	User               string `arg:"--user" default:"postgres"`
}

type InstallServiceForSystemdCmd struct {
	Group              string `arg:"--group" default:"postgres"`
	ServiceInstallPath string `arg:"--service-install-path" default:"/etc/systemd/system/postgresql.service"`
	User               string `arg:"--user" default:"postgres"`
}

type InstallServiceForWindowsCmd struct {
	Name        string `arg:"--service-name" default:"PostgreSQL"`
	Description string `arg:"--service-description" default:"open-source relational database management system"`
}

type InstallServiceCmd struct {
	OpenRc         *InstallServiceForOpenRcCmd  `arg:"subcommand:openrc" help:"Install OpenRC service"`
	Systemd        *InstallServiceForSystemdCmd `arg:"subcommand:systemd" help:"Install systemd service"`
	WindowsService *InstallServiceForWindowsCmd `arg:"subcommand:windows-service" help:"Install Windows Service"`
}

type LsCmd struct{}

type LsRemoteCmd struct{}

type PingCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type ReloadCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type StartCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
	NoInstall       bool   `arg:"--no-install" default:"false" help:"Inverts default of installing nonexistent version"`
}

type StopCmd struct {
	PostgresVersion string `arg:"positional" placeholder:"POSTGRES_VERSION" default:""`
}

type UriCmd struct{}

type Args struct {
	ConfigStruct

	ConfigFile    string `arg:"-c,--config" help:"Config filepath to use"`
	NoConfigRead  bool   `arg:"--no-config-read" default:"false" help:"Do not read the config file"`
	NoConfigWrite bool   `arg:"--no-config-write" default:"false" help:"Do not write to config file"`
	NoRemote      bool   `arg:"--no-remote" default:"false" help:"Disable HTTPS calls for everything except 'install'"`

	Download       *DownloadCmd       `arg:"subcommand:download" help:"Download specified PostgreSQL version"`
	Env            *EnvCmd            `arg:"subcommand:env" help:"Print out associated environment variables"`
	GetPath        *GetPathCmd        `arg:"subcommand:get-path" help:"One of: bin, data, log, runtime"`
	Install        *InstallCmd        `arg:"subcommand:install" help:"Install specified PostgreSQL version"`
	InstallService *InstallServiceCmd `arg:"subcommand:install-service" help:"Install service (daemon), e.g., systemd"`
	Ls             *LsCmd             `arg:"subcommand:ls" help:"List what versions of PostgreSQL are installed"`
	LsRemote       *LsRemoteCmd       `arg:"subcommand:ls-remote" help:"List what versions of PostgreSQL are available"`
	Ping           *PingCmd           `arg:"subcommand:ping" help:"Confirm server is online and auth works"`
	Reload         *ReloadCmd         `arg:"subcommand:reload" help:"Reload specified PostgreSQL server"`
	Start          *StartCmd          `arg:"subcommand:start" help:"Start specified PostgreSQL server"`
	Stop           *StopCmd           `arg:"subcommand:stop" help:"Stop specific (running) PostgreSQL server"`
	Uri            *UriCmd            `arg:"subcommand:uri" help:"Print out database connection string"`
}

func (Args) Description() string {
	return "PostgreSQL version manager"
}

func (Args) Version() string {
	return "pvm 0.0.21"
}
