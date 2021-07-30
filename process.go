package filacct

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

const dateFmt = "2006-01-02"

type balance struct {
	Available  string `json:"availableBalance"`
	Balance    string `json:"balance"`
	Height     int    `json:"height"`
	Pledged    string `json:"sectorPledgeBalance"`
	Timestramp int64  `json:"timestamp"`
	Vesting    string `json:"vestingFunds"`
}

type Balance struct {
	Available string `json:"available,omitempty"`
	Pledged   string `json:"pledged,omitempty"`
	Locked    string `json:"locked,omitempty"`
}

type transfer struct {
	From      string `json:"from,omitempty"`
	Height    int    `json:"height,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	To        string `json:"to,omitempty"`
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
}

type transfers struct {
	TotalCount int        `json:"totalCount"`
	Transfers  []transfer `json:"transfers"`
	Types      []string   `json:"types"`
}

type Assets struct {
	Transferred string `json:"transferred,omitempty"`
}

type blocks struct {
	Blocks     []block `json:"blocks"`
	TotalCount int     `json:"totalCount"`
}

type block struct {
	CID       string `json:"cid"`
	Height    int    `json:"height"`
	Reward    string `json:"reward"`
	Timestamp int64  `json:"timestamp"`
	WinCount  int    `json:"winCount"`
}

type Blocks struct {
	Count  int    `json:"count,omitempty"`
	Reward string `json:"reward,omitempty"`
}

type Fees struct {
	MinerFee     string `json:"miner_fee,omitempty"`
	BurnFee      string `json:"burn_fee,omitempty"`
	WindowPoSt   string `json:"window_post,omitempty"`
	PreCommit    string `json:"pre_commit,omitempty"`
	ProveCommit  string `json:"prove_commit,omitempty"`
	MinerPenalty string `json:"miner_penalty,omitempty"`
	Other        string `json:"other,omitempty"`
}

type Penalty struct {
	Value string
}

type Result struct {
	Balance `json:"balance,omitempty"`
	Assets  `json:"assets,omitempty"`
	Blocks  `json:"blocks,omitempty"`
	Fees    `json:"fees,omitempty"`
	Penalty `json:"penalty,omitempty"`
}

type resultDay struct {
	Day             int
	CountPreCom     int
	CountProveCom   int
	CountSubmitPost int
	CountOther      int
	BurnPreCom      uint64
	BurnProveCom    uint64
	BurnSubmitPost  uint64
	BurnOther       uint64
	MinerFee        uint64
	GasUsed         uint64
}

type Query struct {
	StartDate, EndDate time.Time
	MinerID            string
}

func QueryMiner(q Query) (Result, error) {
	m := &Miner{Address: q.MinerID}

	// turn time.Time to time.Unix
	start := q.StartDate.Unix()
	end := q.EndDate.Unix()

	balance, err := m.GetBalance()
	if err != nil {
		return Result{}, err
	}

	assets, err := m.GetAssets(start, end)
	if err != nil {
		return Result{}, err
	}

	penalties, err := m.GetPenalties(start, end)
	if err != nil {
		return Result{}, err
	}

	blocks, err := m.GetBlocks(start, end)
	if err != nil {
		return Result{}, err
	}

	fees, err := m.GetFees(start, end)
	if err != nil {
		return Result{}, err
	}

	res := Result{
		Balance: balance,
		Assets:  assets,
		Blocks:  blocks,
		Fees:    fees,
		Penalty: penalties,
	}

	return res, nil
}

func (m *Miner) GetBalance() (Balance, error) {
	// read balance json file
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/balance-stats.json", m.Address))
	if err != nil {
		return Balance{}, err
	}

	b := []balance{}
	err = json.Unmarshal(content, &b)
	if err != nil {
		return Balance{}, err
	}
	sort.SliceStable(b, func(i, j int) bool { return b[i].Height > b[j].Height })

	available, err := strconv.ParseFloat(b[0].Available, 64)
	if err != nil {
		return Balance{}, err
	}
	pledged, err := strconv.ParseFloat(b[0].Pledged, 64)
	if err != nil {
		return Balance{}, err
	}
	locked, err := strconv.ParseFloat(b[0].Vesting, 64)
	if err != nil {
		return Balance{}, err
	}
	// get latest balance element
	return Balance{Available: FilFloat(available), Pledged: FilFloat(pledged), Locked: FilFloat(locked)}, nil
}

func (m *Miner) GetPenalties(start, end int64) (Penalty, error) {
	// transfers to the burn account that have no message id are penalties
	// get transfer amounts to determine whether any FIL has been sent from the miner
	filename, err := m.getLatestJsonFile("transfers")
	if err != nil {
		return Penalty{}, err
	}
	acontent, err := ioutil.ReadFile(fmt.Sprintf("%s/transfers/%s.json", m.Address, filename))
	if err != nil {
		return Penalty{}, err
	}

	transf := []transfer{}
	err = json.Unmarshal(acontent, &transf)
	if err != nil {
		return Penalty{}, err
	}
	sort.SliceStable(transf, func(i, j int) bool { return transf[i].Height < transf[j].Height })

	subset := []transfer{}
	for _, t := range transf {
		if t.Timestamp >= start && t.Timestamp <= end {
			subset = append(subset, t)
		}
	}

	var amount float64
	for _, s := range subset {
		if s.Message == "" {
			a, _ := strconv.ParseFloat(s.Value, 64)
			amount += a
		}
	}

	return Penalty{Value: fmt.Sprintf("%v", amount)}, nil
}

func (m *Miner) GetAssets(start, end int64) (Assets, error) {
	filename, err := m.getLatestJsonFile("transfers")
	if err != nil {
		return Assets{}, err
	}
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/transfers/%s.json", m.Address, filename))
	if err != nil {
		return Assets{}, err
	}

	transf := []transfer{}
	err = json.Unmarshal(content, &transf)
	if err != nil {
		return Assets{}, err
	}
	sort.SliceStable(transf, func(i, j int) bool { return transf[i].Height < transf[j].Height })

	subset := []transfer{}
	for _, t := range transf {
		if t.Timestamp >= start && t.Timestamp <= end {
			subset = append(subset, t)
		}
	}
	var amount float64
	for _, s := range subset {
		a, _ := strconv.ParseFloat(s.Value, 64)
		amount += a
	}

	asset := Assets{Transferred: FilFloat(amount)}

	return asset, nil
}

func (m *Miner) GetBlocks(start, end int64) (Blocks, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/blocks.json", m.Address))
	if err != nil {
		return Blocks{}, err
	}

	blks := blocks{}
	err = json.Unmarshal(content, &blks)
	if err != nil {
		return Blocks{}, err
	}
	sort.SliceStable(blks.Blocks, func(i, j int) bool { return blks.Blocks[i].Height < blks.Blocks[j].Height })

	filtered := []block{}
	for _, b := range blks.Blocks {
		if b.Timestamp >= start && b.Timestamp <= end {
			filtered = append(filtered, b)
		}
	}
	var reward float64
	for _, f := range filtered {
		r, _ := strconv.ParseFloat(f.Reward, 64)
		reward += r
	}

	blk := Blocks{
		Count:  len(filtered),
		Reward: FilFloat(reward),
	}

	return blk, nil
}

func (m *Miner) GetFees(start, end int64) (Fees, error) {
	messages, err := m.MarshalAllMsgs()
	if err != nil {
		return Fees{}, err
	}
	sort.SliceStable(messages, func(i, j int) bool { return messages[i].Timestamp > messages[j].Timestamp })

	subset := []MsgCut{}
	for _, m := range messages {
		if m.Timestamp >= start && m.Timestamp <= end {
			subset = append(subset, m)
		}
	}

	fees := struct {
		WindowPoSt   float64
		PreCommit    float64
		ProveCommit  float64
		Other        float64
		MinerFee     float64
		MinerPenalty float64
	}{}
	for _, msg := range subset {
		if m.Address == msg.To {
			burn, _ := strconv.ParseFloat(msg.Fee.BaseFeeBurn, 64)
			oeb, _ := strconv.ParseFloat(msg.Fee.OverEstimationBurn, 64)
			mfee, _ := strconv.ParseFloat(msg.Fee.MinerTip, 64)
			mpen, _ := strconv.ParseFloat(msg.Fee.MinerPenalty, 64)
			// gas = m.Receipt.GasUsed
			switch msg.MethodNumber {
			case 5:
				fees.WindowPoSt += burn + oeb
			case 6:
				fees.PreCommit += burn + oeb
			case 7:
				fees.ProveCommit += burn + oeb
			default:
				fees.Other += burn + oeb
			}
			fees.MinerFee += mfee
			fees.MinerPenalty += mpen
		}
	}

	fee := Fees{
		WindowPoSt:   FilFloat(fees.WindowPoSt),
		PreCommit:    FilFloat(fees.PreCommit),
		ProveCommit:  FilFloat(fees.ProveCommit),
		MinerFee:     FilFloat(fees.MinerFee),
		MinerPenalty: FilFloat(fees.MinerPenalty),
		BurnFee:      FilFloat(fees.Other),
	}

	return fee, nil
}

func ProcessDaysResults(addresses []string, dayRaster map[int][]MsgCut) ([]resultDay, error) {
	resultingDays := make([]resultDay, 0)

	// go through all the msgs for each day
	for day, msgs := range dayRaster {
		resultDay := resultDay{}
		resultDay.Day = day
		for _, msg := range msgs {
			for _, address := range addresses {
				if address == msg.To {
					// if address == msg.From {
					burn, _ := strconv.ParseUint(msg.Fee.BaseFeeBurn, 10, 64)
					oeb, _ := strconv.ParseUint(msg.Fee.OverEstimationBurn, 10, 64)
					mfee, _ := strconv.ParseUint(msg.Fee.MinerTip, 10, 64)
					resultDay.GasUsed = uint64(msg.Receipt.GasUsed)
					switch msg.MethodNumber {
					case 5:
						resultDay.BurnSubmitPost = resultDay.BurnSubmitPost + burn + oeb
						resultDay.CountSubmitPost++
					case 6:
						resultDay.BurnPreCom = resultDay.BurnPreCom + burn + oeb
						resultDay.CountPreCom++
					case 7:
						resultDay.BurnProveCom = resultDay.BurnProveCom + burn + oeb
						resultDay.CountProveCom++
					default:
						resultDay.BurnOther = resultDay.BurnOther + burn + oeb
						resultDay.CountOther++
					}
					resultDay.MinerFee = resultDay.MinerFee + mfee
				}
			}
		}
		resultingDays = append(resultingDays, resultDay)
	}
	sort.SliceStable(resultingDays, func(i, j int) bool {
		return resultingDays[i].Day < resultingDays[j].Day
	})

	return resultingDays, nil
}

func WriteResultsToCSV(resultingDays []resultDay) {
	// write results to csv file
	t := time.Now()
	fileName := "burn-report_" + t.Format("2006-01-02_15-04-05") + ".csv"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)
	// tableHead := "Filecoin-Day,Start-Epoch,End-Epoch,SumbitWindowedPostMessages,SubmitPreCommitMessages,SubmitProveCommitMessages,OtherMessages,SumbitWindowedPostBurn,SubmitPreCommitBurn,SubmitProveCommitBurn,OtherBurn,MinerFees,TotalBurn"
	tableHead := "Filecoin-Day,Start-Epoch,End-Epoch,SumbitWindowedPostBurn,SubmitPreCommitBurn,SubmitProveCommitBurn,OtherBurn,MinerFees,TotalBurn,GasUsed"
	_, _ = datawriter.WriteString(tableHead + "\n")
	for _, day := range resultingDays {
		burnTotal := day.BurnSubmitPost + day.BurnPreCom + day.BurnProveCom + day.BurnOther + day.MinerFee
		sEpoch := day.Day * 2880
		eEpoch := sEpoch + 2879
		dayStr := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
			fmt.Sprintf("%d", day.Day),
			fmt.Sprintf("%d", sEpoch),
			fmt.Sprintf("%d", eEpoch),
			fmt.Sprintf("%.18f", float64(day.BurnSubmitPost)*0.000000000000000001),
			fmt.Sprintf("%.18f", float64(day.BurnPreCom)*0.000000000000000001),
			fmt.Sprintf("%.18f", float64(day.BurnProveCom)*0.000000000000000001),
			fmt.Sprintf("%.18f", float64(day.BurnOther)*0.000000000000000001),
			fmt.Sprintf("%.18f", float64(day.MinerFee)*0.000000000000000001),
			fmt.Sprintf("%.18f", float64(burnTotal)*0.000000000000000001),
			fmt.Sprintf("%d", day.GasUsed),
		)
		_, _ = datawriter.WriteString(dayStr + "\n")
	}
	datawriter.Flush()
	file.Close()
}

func FilFloat(v float64) string {
	return fmt.Sprintf("%.18f", v*0.000000000000000001)
}
