package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lanzafame/filacct"
	cli "github.com/urfave/cli/v2"
)

//go:embed template
var content embed.FS

func main() {
	local := []*cli.Command{
		downloadCmd,
		processCmd,
		serveCmd,
		initCmd,
	}

	app := &cli.App{
		Name:     "filacct",
		Commands: local,
	}
	app.Setup()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var downloadCmd = &cli.Command{
	Name: "download",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "owner",
			Usage: "Owner ID/Address",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() <= 0 && cctx.String("owner") == "" {
			return errors.New("please provide at least one miner ID, i.e. f0410941, or owner ID/address, i.e. --owner f0409359")
		}
		var addresses []string
		var err error
		if cctx.String("owner") != "" {
			log.Println("fetching owner: ", cctx.String("owner"))
			addresses, err = filacct.FetchOwner(cctx.String("owner"))
			if err != nil {
				return err
			}
		}

		if len(addresses) == 0 {
			addresses = cctx.Args().Slice()
		}
		err = initMinerDirs(addresses)
		if err != nil {
			log.Println(err)
		}
		for _, address := range addresses {
			err := filacct.FetchAddress(address)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

var initCmd = &cli.Command{
	Name: "init",
	Action: func(cctx *cli.Context) error {
		/* directory structure for stored data
		owners/
			<owner-id>.json // contains the list of associated miner id's
		<miner-id>/
			messages/
				<message-id>.json // contains the contents of message-id
			msglist/
				<timestamp>.json // contains the list of messages that are new since the last update
			transfers/
				<timestamp>.json // contains the list of transfers that are new since the last update
			balance-stats.json // contains the balance statistics for the entire history of the miner
			blocks.json // contains the list of blocks that the miner has won
		*/
		if cctx.Args().Len() <= 0 {
			return errors.New("please provide at least one miner address, i.e. f0410941")
		}
		err := os.Mkdir("owners", 0777)
		if err != nil {
			log.Println("could not create directory owners")
			log.Println(err)
		}
		addresses := cctx.Args().Slice()
		return initMinerDirs(addresses)
	},
}

// initMinerDirs creates the directories for a slice of miners
func initMinerDirs(addresses []string) error {
	for _, addr := range addresses {
		err := os.Mkdir(addr, 0777)
		if err != nil {
			return err
		}
		err = os.MkdirAll(fmt.Sprintf("%s/messages", addr), 0777)
		if err != nil {
			return err
		}
		err = os.MkdirAll(fmt.Sprintf("%s/msglist", addr), 0777)
		if err != nil {
			return err
		}
		err = os.MkdirAll(fmt.Sprintf("%s/transfers", addr), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

var serveCmd = &cli.Command{
	Name: "serve",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 9090,
			Usage: "Port that the server will use",
		},
	},
	Action: func(cctx *cli.Context) error {
		http.HandleFunc("/", account)

		port := cctx.Int("port")
		log.Printf("serving on 0.0.0.0:%d...", port)
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil) // setting listening port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		return nil
	},
}

var processCmd = &cli.Command{
	Name: "process",
	Action: func(cctx *cli.Context) error {
		fmt.Println("processing")
		if cctx.Args().Len() <= 0 {
			log.Println("please provide at least one miner address, i.e. f0410941")
		}
		addresses := cctx.Args().Slice()
		for _, a := range addresses {
			m := &filacct.Miner{Address: a}
			messages, err := m.MarshalAllMsgs()
			if err != nil {
				return err
			}

			// split messages into days based on epoch
			dayRaster := make(map[int][]filacct.MsgCut)
			for _, message := range messages {
				x := message.Height / 2880
				day, ok := dayRaster[x]
				if ok {
					day = append(day, message)
					dayRaster[x] = day
				} else {
					newDay := make([]filacct.MsgCut, 0)
					newDay = append(newDay, message)
					dayRaster[x] = newDay
				}
			}

			resultingDays, err := filacct.ProcessDaysResults(addresses, dayRaster)
			if err != nil {
				return err
			}

			filacct.WriteResultsToCSV(resultingDays)
		}
		return nil
	},
}
