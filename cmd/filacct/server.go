package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/lanzafame/filacct"
)

const dateFmt = "2006-01-02"

type Result struct {
	MinerID     string
	StartDate   string
	EndDate     string
	Available   string
	Pledged     string
	Locked      string
	Transferred string
	Penalty     string
	Sent        string
	MinerFee    string
	BlocksWon   int
	Burn        string
	BurnFee     string
	WindowPoSt  string
	PreCommit   string
	ProveCommit string
	FILWon      string
}

type OwnerResult struct {
	Owner  Result
	Miners []Result
}

func account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFS(content, "template/form.gtpl")
		if err != nil {
			log.Print(err)
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		log.Printf("req: ip: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
		log.Printf("miner-id: %s, start-date: %s, end-date: %s", r.Form["miner-id"], r.Form["start-date"], r.Form["end-date"])
		start, err := time.Parse(dateFmt, r.Form["start-date"][0])
		if err != nil {
			log.Print(err)
		}
		end, err := time.Parse(dateFmt, r.Form["end-date"][0])
		if err != nil {
			log.Print(err)
		}

		q := filacct.Query{Address: r.Form["miner-id"][0], StartDate: start, EndDate: end}

		qresults, err := filacct.QueryAddress(q)
		if err != nil {
			log.Print(err)
		}

		if len(qresults) == 1 {
			var res Result
			for id, result := range qresults {
				res = Result{
					MinerID:     id,
					StartDate:   q.StartDate.Local().String(),
					EndDate:     q.EndDate.Local().String(),
					Available:   result.Available,
					Pledged:     result.Pledged,
					Locked:      result.Locked,
					Penalty:     result.Penalty.Value,
					Sent:        result.Sent.Value,
					Transferred: result.Transferred,
					MinerFee:    result.MinerFee,
					BlocksWon:   result.Count,
					Burn:        result.Burn.Value,
					BurnFee:     result.BurnFee,
					WindowPoSt:  result.WindowPoSt,
					PreCommit:   result.PreCommit,
					ProveCommit: result.ProveCommit,
					FILWon:      result.Reward,
				}
			}
			t, err := template.ParseFS(content, "template/results.gtpl")
			if err != nil {
				log.Print(err)
			}
			t.Execute(w, res)
		} else if len(qresults) > 1 {
			var results []Result
			var owner Result
			for id, result := range qresults {
				res := Result{
					MinerID:     id,
					StartDate:   q.StartDate.Local().String(),
					EndDate:     q.EndDate.Local().String(),
					Available:   result.Available,
					Pledged:     result.Pledged,
					Locked:      result.Locked,
					Penalty:     result.Penalty.Value,
					Sent:        result.Sent.Value,
					Transferred: result.Transferred,
					MinerFee:    result.MinerFee,
					BlocksWon:   result.Count,
					Burn:        result.Burn.Value,
					BurnFee:     result.BurnFee,
					WindowPoSt:  result.WindowPoSt,
					PreCommit:   result.PreCommit,
					ProveCommit: result.ProveCommit,
					FILWon:      result.Reward,
				}
				if id != q.Address {
					results = append(results, res)
				} else {
					owner = res
				}
			}
			ownerresults := &OwnerResult{owner, results}
			sort.SliceStable(ownerresults.Miners, func(i, j int) bool { return ownerresults.Miners[i].MinerID < ownerresults.Miners[j].MinerID })
			t, err := template.ParseFS(content, "template/ownerresults.gtpl")
			if err != nil {
				log.Print(err)
			}
			t.Execute(w, ownerresults)
		}
	}
}
