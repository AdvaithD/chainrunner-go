package bot

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/sync/errgroup"
)

var (
	enable_pprof  = flag.Bool("pprof", false, "pprof profiling")
	repl_mode = flag.Bool(
		"repl_mode", false, "repl mode to inspect DB",
	)
)

const (
	ESTIMATE_GAS_FAIL    = "estimate_gas_fail.json"
	FOUND_ARBS           = "found_arbs.json"
	REASONS_NOT_WORTH_IT = "not_worth.json"
)

type clients struct {
	newBlock      *ethclient.Client
	uniswap       *ethclient.Client
	sushiswap     *ethclient.Client
	otherwise     *ethclient.Client
}

type ReportMessage struct {
	Error    error
	Role     string
	When     time.Time
	ArgsUsed interface{}
	Extra    interface{}
}

type report_logs struct {
	found_arb         chan *ReportMessage
	estimate_fail_log chan *ReportMessage
	not_worth_it_log  chan *ReportMessage
}

type Bot struct {
	db                  *leveldb.DB
	clients 			clients
	shutdown            chan struct{}
	log_update_incoming report_logs
}

func NewBot(db_path, client_path string) (*Bot, error) {
	client, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	newBlockClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	uniswapClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	sushiSwapClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	db, err := leveldb.RecoverFile(db_path, nil)

	if err != nil {
		return nil, err
	}

	return &Bot{
		clients: clients{
			newBlock: newBlockClient,
			uniswap: uniswapClient,
			sushiswap: sushiSwapClient,
			otherwise: client,
		},
		log_update_incoming: report_logs{
			make(chan *ReportMessage),
			make(chan *ReportMessage),
			make(chan *ReportMessage),
		},
		shutdown: make(chan struct{}),
		db: db,
	}, nil
}

func(b *Bot) CloseResources() error {
	b.clients.newBlock.Close()
	b.clients.otherwise.Close()
	b.clients.uniswap.Close()
	b.clients.sushiswap.Close()

	return b.db.Close()
}

func (b *Bot) KickoffFailureLogs(file_used string, failures chan *ReportMessage) {
	f, err := os.OpenFile(file_used, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Println("come back to this")
		return
	}
	
	for {
		select {
		case <-b.shutdown:
			f.Close()
			return
		case msg := <-failures:
			s, err := json.MarshalIndent(msg, "", "\t")
			if err != nil {
				log.Println("come back to this")
				return
			}

			if _, err := f.Write(s); err != nil {
				log.Println("some error on writing to "+file_used, err)
				f.Close()
				return
			}

			f.WriteString("\n")
		}
	}
}

func(b *Bot) Run() (e error) {
	var g errgroup.Group

	go b.KickoffFailureLogs(ESTIMATE_GAS_FAIL, b.log_update_incoming.estimate_fail_log)
	go b.KickoffFailureLogs(REASONS_NOT_WORTH_IT, b.log_update_incoming.not_worth_it_log)

	// g.Go(func() error {
	// 	return nil
	// })

	g.Go(func() error {
		interrupt := make(chan os.Signal)
		defer signal.Stop(interrupt)
		defer close(interrupt)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		close(b.shutdown)
		time.Sleep(time.Second * 2)
		return nil
	})

	return g.Wait()
}