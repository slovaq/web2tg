package main

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
)

func (s Boxs) Len() int {
	return len(s)
}

func (s Boxs) Less(i, j int) bool {
	return s[i].Data < s[j].Data
}

func (s Boxs) Swap(i, j int) {

	s[i], s[j] = s[j], s[i]
}

type Box struct {
	Data int64
}
type Boxs []Box
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}
func (box *Boxs) DBCheck() {

}
func (box *Boxs) checkDateCounter() {
	for {
		select {

		case <-checkDate:
			box.DBCheck()
		}
	}
}
func (box *Boxs) add(item int64) {
	*box = append(*box, Box{Data: item})
}
func (box *Boxs) read(f chan bool) {
	for {
		select {
		case <-f:
			m.Lock()
			sort.Sort(box)
			bx := append(Boxs{}, *box...)
			unixTimeUTC := time.Unix(bx[0].Data, 0)
			unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
			fmc.Printfln("#rbt read> #bbt Time: #gbt %s", unitTimeInRFC3339)
			if len(bx) == 1 {
				*box = Boxs{}
			} else {
				if 1 < len(bx) {
					*box = bx[1:]
				}
			}
			m.Unlock()
		}
	}
}
func (box *Boxs) check(f chan bool) {
	for {
		m.Lock()
		//fmc.Printfln("#rbt check lock")
		sort.Sort(box)
		bx := append(Boxs{}, *box...)
		if 0 < len(bx) {
			dt := time.Now().Local().Unix()
			if (bx[0].Data - dt) < 0 {
				v := true
				f <- v

			}
		}
		m.Unlock()
		time.Sleep(time.Duration(1) * time.Second)
	}
}

var (
	m         sync.Mutex
	d         int
	checkDate chan bool
)

func (box *Boxs) write(f chan bool) {
	for {
		r := rand.New(rand.NewSource(55))
		ir := IntRange{0, 30}
		m.Lock()
		bx := append(Boxs{}, *box...)
		dt := (time.Now().Local().Unix()) + int64(ir.NextRandom(r))
		unixTimeUTC := time.Unix(dt, 0)
		unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
		fmc.Printfln("#rbt write> #gbtdate: #ybt%s", unitTimeInRFC3339)
		bx = append(bx, Box{Data: dt})
		*box = bx
		m.Unlock()
		time.Sleep(time.Duration(5) * time.Second)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())
	box := Boxs{}
	//box.add(2)
	f := make(chan bool)
	go box.check(f)
	go box.read(f)
	//go box.write(f)

	box.checkDateCounter()
}
