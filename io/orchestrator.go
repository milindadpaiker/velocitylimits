package io

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	validate "github.com/milind/velocitylimit/validate"
	"github.com/pkg/errors"
)

//Process processes orchestrates flow. Read from input call validator and then sink to output
func Process(ctx context.Context, input Ingester, output Sink, validator *validate.Validator) {

	in := make(chan string)
	out := make(chan string)
	var wg sync.WaitGroup
	ct := context.Background()
	ioctx, iocancel := context.WithCancel(ct)
	wg.Add(1)
	go input.Read(ioctx, in, &wg)
	go output.Write(ioctx, out, &wg)
	defer func() {
		iocancel()
		wg.Wait()
	}()
	for {
		select {
		case data := <-in:

			if data == "" {
				//In pipe has been closed. ingestion completed.
				//close outpie and signal shutdown
				close(out)
				iocancel()
				return
			}
			deposit, err := preProcess(data)
			if err != nil {
				log.Println(err)
				continue
			}
			result, err := validator.Process(deposit)
			if err != nil {
				log.Printf("Fund %s not processed. Error: %v\n", data, err)
				continue
			}
			fundStatus, err := postProcess(result)
			if err != nil {
				log.Println(err)
				continue
			}
			out <- fundStatus

		case <-ctx.Done():
			//Ctrl + C hit, shut gracefully
			iocancel()
			return
		}
	}

}

type incomingFund struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	LoadAmount string    `json:"load_amount"`
	Time       time.Time `json:"time"`
}

type fundStatus struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
	//to be removed
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

func preProcess(fund string) (*validate.Deposit, error) {
	inFund := &incomingFund{}

	//make sure incoming fund is a valid json
	err := json.Unmarshal([]byte(fund), inFund)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse deposit")
	}

	// Remove the $ and convert to float64
	inFund.LoadAmount = strings.Replace(inFund.LoadAmount, "$", "", 1)
	loadAmnt, err := strconv.ParseFloat(inFund.LoadAmount, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Could not conver load amount")
	}

	//load id must be numeric
	loadID, err := strconv.Atoi(inFund.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert id")
	}
	//customer id must be numeric
	custID, err := strconv.Atoi(inFund.CustomerID)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert customerID")
	}
	a := &validate.Deposit{
		ID:         loadID,
		CustomerID: custID,
		LoadAmount: loadAmnt,
		Time:       inFund.Time,
	}
	return a, nil
}

func postProcess(fundRslt *validate.DepositStatus) (string, error) {

	tmp := fundStatus{
		ID:         strconv.Itoa(fundRslt.ID),
		CustomerID: strconv.Itoa(fundRslt.CustomerID),
		Accepted:   fundRslt.Accepted,
		//to be removed
		Amount: fundRslt.Amount,
		Time:   fundRslt.Time.UTC(),
	}
	result, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
