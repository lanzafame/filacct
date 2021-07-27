package filacct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Miner struct {
	Address string
}

type msgStub struct {
	Cid       string `json:"cid"`
	Height    int    `json:"height"`
	Timestamp int    `json:"timestamp"`
}

type msgList struct {
	TotalCount int       `json:"totalCount"`
	Messages   []msgStub `json:"messages"`
}

type MsgFee struct {
	BaseFeeBurn        string `json:"baseFeeBurn"`
	OverEstimationBurn string `json:"overEstimationBurn"`
	MinerTip           string `json:"minerTip"`
	MinerPenalty       string `json:"minerPenalty"`
	Refund             string `json:"refund"`
}

type receipt struct {
	GasUsed int64 `json:"gasUsed"`
}

type MsgCut struct {
	Cid          string  `json:"cid"`
	Height       int     `json:"height"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Fee          MsgFee  `json:"fee"`
	Receipt      receipt `json:"receipt"`
	MethodNumber int     `json:"methodNumber"`
	Timestamp    int64   `json:"timestamp"`
}

func FetchAddress(address string) error {
	m := &Miner{Address: address}
	err := m.fetchMessages()
	if err != nil {
		return err
	}

	err = m.fetchBalance()
	if err != nil {
		return err
	}

	err = m.fetchBlocks()
	if err != nil {
		return err
	}

	err = m.fetchTransfers()
	if err != nil {
		return err
	}

	return nil
}

func (m *Miner) fetchMessages() error {
	initMsgList, err := m.fetchMsgLists(fmt.Sprintf("https://filfox.info/api/v1/address/%s/messages?pageSize=%d&page=%d", m.Address, 1, 0))
	if err != nil {
		return err
	}

	prevMsgStub, err := m.getLatestStoredMsg()
	if err != nil {
		return err
	}

	if initMsgList.Messages[0].Height <= prevMsgStub.Height {
		log.Print("no new msgs found")
		return nil
	}

	pageCalls := initMsgList.TotalCount / 100
	msgStubs := make([]msgStub, 0)
	for i := 0; i <= pageCalls; i++ {
		msgList, err := m.fetchMsgLists(fmt.Sprintf("https://filfox.info/api/v1/address/%s/messages?pageSize=%d&page=%d", m.Address, 100, i))
		if err != nil {
			return err
		}
		msgStubs = append(msgStubs, msgList.Messages...)
		if msgStubs[len(msgStubs)-1].Height <= prevMsgStub.Height {
			break
		}
	}

	delta := extractDeltaSlice(msgStubs, prevMsgStub.Height)

	jsonMsgStubs, err := json.Marshal(delta)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/msglist/%d.json", m.Address, time.Now().Unix()), jsonMsgStubs, 0666)
	if err != nil {
		return err
	}

	for i, msg := range delta {
		log.Printf("msg %d/%d", i, len(delta))
		msgCid := msg.Cid
		err := m.fetchMessage("https://filfox.info/api/v1/message/%s", msgCid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Miner) fetchLists(url string) ([]byte, error) {
	log.Println("Fetching ", url)
	filfoxClient := http.Client{
		Timeout: time.Second * 120,
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
	return body, nil
}

func (m *Miner) fetchTransfers() error {
	initTransferList, err := m.fetchTransfersPage(fmt.Sprintf("https://filfox.info/api/v1/address/%s/transfers?pageSize=%d&page=%d", m.Address, 1, 0))
	if err != nil {
		return err
	}

	// prevMsgStub, err := m.getLatestStoredMsg()
	// if err != nil {
	// 	return err
	// }

	// if initMsgList.Messages[0].Height <= prevMsgStub.Height {
	// 	log.Print("no new msgs found")
	// 	return nil
	// }

	pageCalls := initTransferList.TotalCount / 100
	ts := []transfer{}
	for i := 0; i <= pageCalls; i++ {
		transferList, err := m.fetchTransfersPage(fmt.Sprintf("https://filfox.info/api/v1/address/%s/transfers?pageSize=%d&page=%d", m.Address, 100, i))
		if err != nil {
			return err
		}
		ts = append(ts, transferList.Transfers...)
		// if ts[len(ts)-1].Height <= prevMsgStub.Height {
		// 	break
		// }
	}

	jsonTransfers, err := json.Marshal(ts)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/transfers/%d.json", m.Address, time.Now().Unix()), jsonTransfers, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (m *Miner) fetchTransfersPage(url string) (transfers, error) {
	body, err := m.fetchLists(url)
	if err != nil {
		return transfers{}, err
	}

	tList := transfers{}
	err = json.Unmarshal(body, &tList)
	if err != nil {
		return transfers{}, err
	}
	return tList, nil
}

func (m *Miner) fetchMsgLists(url string) (msgList, error) {
	body, err := m.fetchLists(url)
	if err != nil {
		return msgList{}, err
	}

	msgs := msgList{}
	err = json.Unmarshal(body, &msgs)
	if err != nil {
		return msgList{}, err
	}
	return msgs, nil
}

func (m *Miner) fetchMessage(urlfmt string, msgCid string) error {
	url := fmt.Sprintf(urlfmt, msgCid)

	log.Println("Fetching ", url)
	filfoxClient := http.Client{
		Timeout: time.Second * 240,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, getErr := filfoxClient.Do(req)
	if getErr != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	backoff(res)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/messages/%s.json", m.Address, msgCid), body, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (m *Miner) fetch(minerID string, thing string) error {
	url := fmt.Sprintf("https://filfox.info/api/v1/address/%s/%s", m.Address, thing)

	log.Println("Fetching ", url)
	filfoxClient := http.Client{
		Timeout: time.Second * 240,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, getErr := filfoxClient.Do(req)
	if getErr != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	backoff(res)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", m.Address, thing), body, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (m *Miner) fetchBlocks() error {
	return m.fetch(m.Address, "blocks")
}

func (m *Miner) fetchBalance() error {
	return m.fetch(m.Address, "balance-stats")
}

func readJSONFiles(dir string) ([][]byte, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	contents := [][]byte{}
	for _, file := range files {
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (m *Miner) MarshalAllMsgs() ([]MsgCut, error) {
	jsonMsgs, err := readJSONFiles(fmt.Sprintf("%s/messages", m.Address))
	if err != nil {
		return nil, err
	}

	msgCuts := []MsgCut{}
	for _, jm := range jsonMsgs {
		ffMsgCut := MsgCut{}
		err = json.Unmarshal(jm, &ffMsgCut)
		if err != nil {
			return nil, err
		}
		msgCuts = append(msgCuts, ffMsgCut)
	}
	return msgCuts, nil
}

func backoff(res *http.Response) {
	defaultBackoff, _ := time.ParseDuration("20s")
	if i, err := strconv.Atoi(res.Header.Get("x-ratelimit-remaining")); err == nil {
		if i < 5 {
			sleep, err := strconv.Atoi(res.Header.Get("x-ratelimit-reset"))
			if err != nil {
				log.Printf("backing off for %s", defaultBackoff)
				time.Sleep(defaultBackoff) // overly cautious in the case where we don't get the reset value
			}
			sleepDur, err := time.ParseDuration(fmt.Sprintf("%ds", sleep))
			if err != nil {
				log.Printf("backing off for %s", defaultBackoff)
				time.Sleep(defaultBackoff) // overly cautious in the case where we don't get the reset value
			}
			log.Printf("backing off for %s", sleepDur)
			time.Sleep(sleepDur)
		}
	}
}

func (m *Miner) getLatestJsonFile(dir string) (string, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", m.Address, dir))
	if err != nil {
		return "", err
	}

	if len(files) <= 0 {
		return "", nil
	}

	names := []string{}
	for _, f := range files {
		names = append(names, strings.TrimRight(f.Name(), ".json"))
	}

	sort.SliceStable(names, func(i int, j int) bool {
		return names[i] > names[j]
	})

	return names[0], nil
}

func (m *Miner) getLatestContents(dir string) ([]byte, error) {
	latest, err := m.getLatestJsonFile(dir)
	if err != nil {
		return nil, err
	}

	if latest == "" {
		return nil, nil
	}

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s.json", m.Address, dir, latest))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (m *Miner) getLatestStoredMsg() (msgStub, error) {
	// read json file with greatest timestamp
	msgStubs, err := m.getLatestContents("msglist")
	if err != nil {
		return msgStub{}, err
	}

	if msgStubs == nil {
		return msgStub{
			Height: 0,
		}, nil
	}

	// unmarshal json
	msgs := []msgStub{}
	err = json.Unmarshal(msgStubs, &msgs)
	if err != nil {
		return msgStub{}, err
	}

	// sort slice as a precaution by height
	sort.SliceStable(msgs, func(i, j int) bool {
		return msgs[i].Height > msgs[j].Height
	})

	// return the first element of the slice
	return msgs[0], nil
}

func extractDeltaSlice(msgs []msgStub, height int) []msgStub {
	subset := []msgStub{}
	for _, m := range msgs {
		if m.Height > height {
			subset = append(subset, m)
		}
	}

	return subset
}
