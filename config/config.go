package config

type IConfig interface {
	App() ConfigApp
	DB() ConfigDB
}

type config struct {
	app *app
	db  *db
}

// App Config
type ConfigApp interface {
	Port() string
	Version() string
}

type app struct {
	port    string
	version string
}

// DB Config
type ConfigDB interface {
	Driver() string
	DSN() string
	MaxOpenConn() int
}

type db struct {
	driver         string
	dsn            string
	maxConnections int
}

// Config Method
func (c *config) App() ConfigApp { return c.app }
func (c *config) DB() ConfigDB   { return c.db }

// App Method
func (a *app) Port() string    { return a.port }
func (a *app) Version() string { return a.version }

// DB Method
func (d *db) DSN() string      { return d.dsn }
func (d *db) Driver() string   { return d.driver }
func (d *db) MaxOpenConn() int { return d.maxConnections }
