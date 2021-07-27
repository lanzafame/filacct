package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lanzafame/filacct"
	cli "github.com/urfave/cli/v2"
)

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
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() <= 0 {
			log.Println("please provide at least one miner address, i.e. f0410941")
		}
		addresses := cctx.Args().Slice()
		for _, address := range addresses {
			err := filacct.FetchAddress(address)
			if err != nil {
				return err
			}
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

type Result struct {
	MinerID     string
	StartDate   string
	EndDate     string
	Available   string
	Pledged     string
	Locked      string
	Transferred string
	MinerFee    string
	BlocksWon   int
	BurnFee     string
	FILWon      string
}

func account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("form.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		log.Printf("req: ip: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
		log.Printf("miner-id: %s, start-date: %s, end-date: %s", r.Form["miner-id"], r.Form["start-date"], r.Form["end-date"])
		const dateFmt = "2006-01-02"
		start, err := time.Parse(dateFmt, r.Form["start-date"][0])
		if err != nil {
			log.Print(err)
		}
		end, err := time.Parse(dateFmt, r.Form["end-date"][0])
		if err != nil {
			log.Print(err)
		}

		q := filacct.Query{MinerID: r.Form["miner-id"][0], StartDate: start, EndDate: end}

		results, err := filacct.QueryMiner(q)
		if err != nil {
			log.Print(err)
		}

		res := &Result{
			MinerID:     q.MinerID,
			StartDate:   q.StartDate.Local().String(),
			EndDate:     q.EndDate.Local().String(),
			Available:   results.Available,
			Pledged:     results.Pledged,
			Locked:      results.Locked,
			Transferred: results.Transferred,
			MinerFee:    results.MinerFee,
			BlocksWon:   results.Count,
			BurnFee:     results.BurnFee,
			FILWon:      results.Reward,
		}
		t, err := template.ParseFiles("results.gtpl")
		if err != nil {
			log.Print(err)
		}
		t.Execute(w, res)
	}
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
		err := http.ListenAndServe("0.0.0.0:80", nil) // setting listening port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		return nil
	},
}

var initCmd = &cli.Command{
	Name: "init",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() <= 0 {
			log.Println("please provide at least one miner address, i.e. f0410941")
		}
		addresses := cctx.Args().Slice()
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
	},
}
