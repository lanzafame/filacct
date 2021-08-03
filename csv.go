package filacct

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

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
