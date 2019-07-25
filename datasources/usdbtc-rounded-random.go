package datasources

import (
	//"math"
	"math/rand"
	"time"
	"github.com/gertjaap/dlcoracle/gcfg"

)

type UsdBtcRoundedRandom struct {
}

func (ds *UsdBtcRoundedRandom) Id() uint64 {
	return 1
}

func (ds *UsdBtcRoundedRandom) Name() string {
	return "US Dollar"
}

func (ds *UsdBtcRoundedRandom) Description() string {
	return "Publishes the value of USD denominated in 1/100000000th units of BTC (satoshi) in multitudes of 100"
}

func (ds *UsdBtcRoundedRandom) Interval() uint64 {
	return gcfg.Interval
}

func (ds *UsdBtcRoundedRandom) Value() (uint64, error) {
	//satoshiValue := uint64(math.Floor(float64(random(100, 150))+0.5)) * 100
	//return satoshiValue, nil
	return gcfg.ValueToPublish, nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
