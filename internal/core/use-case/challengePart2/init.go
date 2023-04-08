package challengePart2

import (
	"encoding/json"
	"fmt"
)

type ChangeRecord struct {
	StockCode string
	Price     int64
	Quantity  int64
}

type IndexMember struct {
	StockCode string
	IndexCode string
}

type Summary struct {
	StockCode string   `json:"stock_code"`
	IndexCode []string `json:"index_code"`
	Open      int64    `json:"open"`
	High      int64    `json:"high"`
	Low       int64    `json:"low"`
	Close     int64    `json:"close"`
	Prev      int64    `json:"prev"`
}

var result = map[string]Summary{}

// DO NOT INCLUDE THIS ORIGINAL FUNCTION INTO YOUR SUBMITTED SOLUTION! AS THIS
// WILL CAUSE MISCALCULATION WHEN STOCKBIT'S TEAM REVIEWS YOUR SUBMISSION, THIS
// ACTION WILL NEGATIVELY IMPACT YOUR SCORE.
func ohlc(x []string, w []ChangeRecord, p []IndexMember) map[string]Summary {
	for _, y := range x {
		found, not := result[y]
		if not {
			found = Summary{}
		}
		found.StockCode = y
		for _, u := range w {
			if u.StockCode == y {
				if u.Quantity == 0 {
					found.Prev = u.Price
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				} else if u.Quantity > 0 && result[y].Open == 0 {
					found.Open = u.Price
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				} else {
					found.Close = u.Price
					if found.High < u.Price {
						found.High = u.Price
					}
					if found.Low > u.Price {
						found.Low = u.Price
					}
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				}
			} else {
				fmt.Println("done")
				fmt.Println("price updated")
				result[y] = found
			}
		}
		for _, i := range p {
			if i.StockCode == y {
				found.IndexCode = append(found.IndexCode, i.IndexCode)
				fmt.Println("index updated")
				result[y] = found
			} else {
				fmt.Println("index updated")
				result[y] = found
			}
		}
	}

	return result
}

func main() {
	x := []string{}
	w := []ChangeRecord{}
	p := []IndexMember{}

	r := ohlc(x, w, p)
	for _, v := range r {
		jss, _ := json.Marshal(v)
		fmt.Println("summary: ", string(jss))
	}
}
