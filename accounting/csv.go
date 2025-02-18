package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date        time.Time
	Description string
	Source      string
	Destination string
	Currency    string
	Price       float32
}

func (trans Transaction) String() string {
	if trans.Description != "" {
		return fmt.Sprintf("%v * \"%v\"\n\t%v  %v %v\n\t%v",
			trans.Date.Format("2006-01-02"),
			trans.Description, trans.Source, trans.Price, trans.Currency, trans.Destination)
	}
	return ""
}

var (
	ErrEmptyString = errors.New("string empty")
	ErrParseDate   = errors.New("time.Parse invalid")
)

var ErrCsvParse *csv.ParseError

func ParseTransaction(fConf *File, data []string) (Transaction, error) {
	if len(data) <= 5 && strings.Join(data[1:], "") == "" {
		return Transaction{}, ErrEmptyString
	}

	date, err := time.Parse(fConf.DateParse, strings.TrimSpace(data[fConf.Date]))
	if err != nil {
		return Transaction{}, ErrParseDate
	}

	price, err := strconv.ParseFloat(strings.TrimSpace(data[fConf.Price]), 32)
	if err != nil {
		return Transaction{}, err
	}
	description := strings.Join(strings.Fields(data[fConf.Description]), " ")

	if fConf.FlipPrice {
		price = price * -1
	}

	return Transaction{
		Date:        date,
		Description: description,
		Price:       float32(price),
		Currency:    "CAD",
	}, nil
}

func ReadCsv(exactRules, partialRules *map[string]string, fConf *File) ([]Transaction, []Transaction) {
	var matched []Transaction
	var unmatched []Transaction

	file, err := os.Open(fConf.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	exactMap := *exactRules
	partialMap := *partialRules

	reader := csv.NewReader(file)
	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("%T ", err)
			panic(err)
		}

		tran, err := ParseTransaction(fConf, data)
		switch err {
		case ErrEmptyString:
		case ErrParseDate:
			continue
		}
		tran.Source = fConf.Account

		// skip empty transactions
		if tran.Price == 0 && tran.Description == "" {
			continue
		}

		if v, ok := exactMap[tran.Description]; ok {
			tran.Destination = v
			matched = append(matched, tran)
		} else {
			partialFound := false
			for inc, a := range partialMap {
				if strings.Contains(tran.Description, inc) {
					tran.Destination = a
					matched = append(matched, tran)
					partialFound = true
				}
			}
			if !partialFound {
				unmatched = append(unmatched, tran)
			}
		}

	}

	return matched, unmatched
}
