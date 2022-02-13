package flags

import "flag"

var (
	// Note: WS client here - ws://157.90.35.22:8545
	DEFAULT_CLIENT = flag.String(
		"client_dial", "/home/bot/.ethereum/geth.ipc", "could be websocket or IPC",
	)

	BOR_CLIENT = flag.String(
		"bor_ipc", "/home/bot/.bor/data/bor.ipc", "could be websocket or IPC",
	)

	ENABLE_PPROF = flag.Bool("pprof", false, "pprof profiling")

	DB_PATH = flag.String(
		"db_path", "chainrunner-test-db", "where is the db",
	)

	DRY_RUN = flag.Bool("dry", false, "dry run ")

	CPU_PROFILE = flag.String("cpuprofile", "", "write cpu profile to `file`")
	MEM_PROFILE = flag.String("memprofile", "", "write memory profile to `file`")

	CONFIG_PATH = flag.String("config", "./config.yml", "path to config file")
)
