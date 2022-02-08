package flags

import "flag"

var (
	// Note: WS client here - ws://157.90.35.22:8545
	DEFAULT_CLIENT = flag.String(
		"client_dial", "/home/bot/.ethereum/geth.ipc", "could be websocket or IPC",
	)

	ENABLE_PPROF = flag.Bool("pprof", false, "pprof profiling")

	DB_PATH = flag.String(
		"db_path", "chainrunner-test-db", "where is the db",
	)

	DRY_RUN = flag.Bool("dry_run", false, "dry run ")
)
