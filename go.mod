module chainrunner

go 1.15

require (
	github.com/ALTree/bigfloat v0.0.0-20201218142103-4a33235224ec
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/ethereum/go-ethereum v0.0.0-local
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/exp v0.0.0-20220209042442-160e291fcf24 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gonum.org/v1/gonum v0.0.0-unpublished
)

replace (
	github.com/ethereum/go-ethereum v0.0.0-local => /home/bot/geth-storageat/
	gonum.org/v1/gonum v0.0.0-unpublished => /home/bot/gonum-fork/
)
