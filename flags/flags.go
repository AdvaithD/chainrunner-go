package flags

import "flag"

var (
	DEFAULT_CLIENT = flag.String(
		"client_dial", "ws://157.90.35.22:8545", "could be websocket or IPC",
	)

	ENABLE_PPROF  = flag.Bool("pprof", false, "pprof profiling")

	DB_PATH             = flag.String(
		"db_path", "chainrunner-test-db", "where is the db",
	)
)
