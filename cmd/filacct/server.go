package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
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

type ResultCache struct {
	sync.RWMutex
	data map[string]map[string]*filacct.Result
}

func NewResultCache() *ResultCache {
	data := make(map[string]map[string]*filacct.Result)
	return &ResultCache{data: data}
}

func (rc *ResultCache) Get(address, start, end string) (map[string]*filacct.Result, bool) {
	rc.RLock()
	defer rc.RUnlock()
	key := strings.Join([]string{address, start, end}, ":")
	res, ok := rc.data[key]
	return res, ok
}

func (rc *ResultCache) Put(address, start, end string, res map[string]*filacct.Result) {
	rc.Lock()
	defer rc.Unlock()
	key := strings.Join([]string{address, start, end}, ":")
	rc.data[key] = res
}

type Default struct {
	Start string
	End   string
}

func (rc *ResultCache) account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFS(content, "template/form.gtpl")
		if err != nil {
			log.Print(err)
		}
		today := time.Now()
		end := strings.Join([]string{string(today.Year()), string(today.Month()), string(today.Day())}, "-")
		data := &Default{
			Start: "2021-10-01",
			End:   end,
		}
		t.Execute(w, data)
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
		address := r.Form["miner-id"][0]

		q := filacct.Query{Address: address, StartDate: start, EndDate: end}

		qresults := map[string]*filacct.Result{}
		if res, ok := rc.Get(address, string(start.Unix()), string(end.Unix())); ok {
			qresults = res
		} else {
			qresults, err = filacct.QueryAddress(q)
			if err != nil {
				log.Print(err)
			}
			rc.Put(address, string(start.Unix()), string(end.Unix()), qresults)
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
