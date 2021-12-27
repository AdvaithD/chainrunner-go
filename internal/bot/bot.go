package bot

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type clients struct {
	newBlock      *ethclient.Client
	uniswap       *ethclient.Client
	sushiswap     *ethclient.Client
	otherwise *ethclient.Client
}

type report_msg struct {
	Error    error
	Role     string
	When     time.Time
	ArgsUsed interface{}
	Extra    interface{}
}

type report_logs struct {
	estimate_fail_log chan *report_msg
	not_worth_it_log  chan *report_msg
	clients           clients
}

type Bot struct {
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

	return &Bot{
		clients: clients{
			newBlock: newBlockClient,
			uniswap: uniswapClient,
			sushiswap: sushiSwapClient,
			otherwise: client,
		},
		log_update_incoming: report_logs{
			make(chan *report_msg), make(chan *report_msg),
		},
		shutdown: make(chan struct{}),
	}, nil
}