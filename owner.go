package filacct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Owner is a struct for the json response from the filfox api
/*
{
    "actor": "account",
    "address": "f3vkfate35djn77rgehygngt3swjezzfg6nxleiaqein53do4d4podenpgr6sv74eln47qjaevriy3hqwl472q",
    "balance": "97465975019593953399",
    "createHeight": 585487,
    "createTimestamp": 1615871010,
    "id": "f0409359",
    "lastSeenHeight": 985832,
    "lastSeenTimestamp": 1627881360,
    "messageCount": 4163,
    "ownedMiners": [
        "f0410938",
        "f0410939",
        "f0410941"
    ],
    "robust": "f3vkfate35djn77rgehygngt3swjezzfg6nxleiaqein53do4d4podenpgr6sv74eln47qjaevriy3hqwl472q",
    "timestamp": 1627881390,
    "workerMiners": [
        "f0410938",
        "f0410939",
        "f0410941"
    ]
}
*/
type Owner struct {
	Actor             string   `json:"actor"`
	Address           string   `json:"address"`
	Balance           string   `json:"balance"`
	CreateHeight      int      `json:"createHeight"`
	CreateTime        int      `json:"createTimestamp"`
	ID                string   `json:"id"`
	LastSeenHeight    int      `json:"lastSeenHeight"`
	LastSeenTimestamp int      `json:"lastSeenTimestamp"`
	MessageCount      int      `json:"messageCount"`
	OwnedMiners       []string `json:"ownedMiners"`
	RobustAddress     string   `json:"robust"`
	Timestamp         int      `json:"timestamp"`
	WorkerMiners      []string `json:"workerMiners"`
}

// QueryAddress determines whether a given address is an owner address or miner address
// and returns the apropriate result
func QueryAddress(q Query) (map[string]*Result, error) {
	// if owner address, query owner directory for which miner addresses need to be processed
	if ok, err := IsOwnerAddress(q.Address); ok && err == nil {
		return QueryOwner(q)
	} else if err != nil {
		return nil, err
	}

	// just a miner address
	res, err := QueryMiner(q)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*Result)
	result[q.Address] = res
	return result, nil
}

func IsOwnerAddress(address string) (bool, error) {
	files, err := ioutil.ReadDir("owners")
	if err != nil {
		return false, err
	}

	if len(files) <= 0 {
		return false, nil
	}

	for _, f := range files {
		// shortcut for when the user provides the owner ID
		if strings.TrimRight(f.Name(), ".json") == address {
			return true, nil
		}
		//FIXME: allow user to provide long/robust address
		// // if the user provided the long/robust address, check the contents of the owner file
		// rawOwner, err := ioutil.ReadFile(fmt.Sprintf("owners/%s", f))
		// if err != nil {
		// 	return false, err
		// }
		// var owner Owner
		// err = json.Unmarshal(rawOwner, &owner)
		// if err != nil {
		// 	return false, err
		// }
		// if address == owner.Address {
		// 	return true, nil
		// }
	}

	return false, nil
}

// QueryOwner queries filfox api to retrieve the associated miner
// addresses of provided owner address
func QueryOwner(q Query) (map[string]*Result, error) {
	rawOwner, err := ioutil.ReadFile(fmt.Sprintf("owners/%s.json", q.Address))
	if err != nil {
		return nil, err
	}
	var owner Owner
	err = json.Unmarshal(rawOwner, &owner)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Result)
	for _, miner := range owner.OwnedMiners {
		res, err := QueryMiner(Query{Address: miner, StartDate: q.StartDate, EndDate: q.EndDate})
		if err != nil {
			return nil, err
		}
		result[miner] = res
	}

	result[q.Address] = SummariseResults(result)

	return result, nil
}

// SummariseResults takes a slice of Results and sums the fields
// returning a single Result
func SummariseResults(results map[string]*Result) *Result {
	var res Result
	for _, r := range results {
		res.Balance.Add(r.Balance)
		res.Assets.Add(r.Assets)
		res.Blocks.Add(r.Blocks)
		res.Fees.Add(r.Fees)
		res.Penalty.Add(r.Penalty)
	}
	return &res
}

func (a *Penalty) Add(b Penalty) {
	value, _ := strconv.ParseFloat(a.Value, 64)
	adder, _ := strconv.ParseFloat(b.Value, 64)
	value += adder
	a.Value = strconv.FormatFloat(value, 'f', -1, 64)
}

func (a *Fees) Add(b Fees) {
	minerfee, _ := strconv.ParseFloat(a.MinerFee, 64)
	adder, _ := strconv.ParseFloat(b.MinerFee, 64)
	minerfee += adder
	a.MinerFee = strconv.FormatFloat(minerfee, 'f', -1, 64)

	burnfee, _ := strconv.ParseFloat(a.BurnFee, 64)
	adder, _ = strconv.ParseFloat(b.BurnFee, 64)
	burnfee += adder
	a.BurnFee = strconv.FormatFloat(burnfee, 'f', -1, 64)

	window, _ := strconv.ParseFloat(a.WindowPoSt, 64)
	adder, _ = strconv.ParseFloat(b.WindowPoSt, 64)
	window += adder
	a.WindowPoSt = strconv.FormatFloat(window, 'f', -1, 64)

	precommit, _ := strconv.ParseFloat(a.PreCommit, 64)
	adder, _ = strconv.ParseFloat(b.PreCommit, 64)
	precommit += adder
	a.PreCommit = strconv.FormatFloat(precommit, 'f', -1, 64)

	provecommit, _ := strconv.ParseFloat(a.ProveCommit, 64)
	adder, _ = strconv.ParseFloat(b.ProveCommit, 64)
	provecommit += adder
	a.ProveCommit = strconv.FormatFloat(provecommit, 'f', -1, 64)

	minerpenalty, _ := strconv.ParseFloat(a.MinerPenalty, 64)
	adder, _ = strconv.ParseFloat(b.MinerPenalty, 64)
	minerpenalty += adder
	a.MinerPenalty = strconv.FormatFloat(minerpenalty, 'f', -1, 64)

	other, _ := strconv.ParseFloat(a.Other, 64)
	adder, _ = strconv.ParseFloat(b.Other, 64)
	other += adder
	a.Other = strconv.FormatFloat(other, 'f', -1, 64)
}

func (a *Blocks) Add(b Blocks) {
	a.Count += b.Count
	reward, _ := strconv.ParseFloat(a.Reward, 64)
	adder, _ := strconv.ParseFloat(b.Reward, 64)
	reward += adder
	a.Reward = strconv.FormatFloat(reward, 'f', -1, 64)
}

func (a *Assets) Add(b Assets) {
	transferred, _ := strconv.ParseFloat(a.Transferred, 64)
	adder, _ := strconv.ParseFloat(b.Transferred, 64)
	transferred += adder
	a.Transferred = strconv.FormatFloat(transferred, 'f', -1, 64)
}

func (a *Balance) Add(b Balance) {
	balance, _ := strconv.ParseFloat(a.Available, 64)
	adder, _ := strconv.ParseFloat(b.Available, 64)
	balance += adder
	a.Available = strconv.FormatFloat(balance, 'f', -1, 64)

	pledged, _ := strconv.ParseFloat(a.Pledged, 64)
	padder, _ := strconv.ParseFloat(b.Pledged, 64)
	pledged += padder
	a.Pledged = strconv.FormatFloat(pledged, 'f', -1, 64)

	locked, _ := strconv.ParseFloat(a.Locked, 64)
	ladder, _ := strconv.ParseFloat(b.Locked, 64)
	locked += ladder
	a.Locked = strconv.FormatFloat(locked, 'f', -1, 64)
}

// FetchOwner fetches the owner from the filfox api and returns
// the list of associated miners or an error
func FetchOwner(ownerAddress string) ([]string, error) {
	url := fmt.Sprintf("https://filfox.info/api/v1/address/%s", ownerAddress)
	filfoxClient := http.Client{
		Timeout: time.Second * 240,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, getErr := filfoxClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	backoff(res)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	var owner Owner
	err = json.Unmarshal(body, &owner)
	if err != nil {
		return nil, err
	}

	// store under the owner id, instead of owner address
	err = ioutil.WriteFile(fmt.Sprintf("owners/%s.json", owner.ID), body, 0666)
	if err != nil {
		return nil, err
	}
	return owner.OwnedMiners, nil
}

// OwnerMiners returns a slice of all the owned miners
func OwnerMiners(ownerAddr string) ([]string, error) {
	rawOwner, err := ioutil.ReadFile(fmt.Sprintf("owners/%s.json", ownerAddr))
	if err != nil {
		return nil, err
	}
	var owner Owner
	err = json.Unmarshal(rawOwner, &owner)
	if err != nil {
		return nil, err
	}
	var miners []string
	miners = append(miners, owner.OwnedMiners...)
	return miners, nil
}
