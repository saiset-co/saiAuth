package internal

import (
	"time"
)

type FloodList map[string]*Flood

type Flood struct {
	Count   int
	Expired int64
}

var Flooders = make(FloodList, 10000)

func (is *InternalService) FloodAdd(ip string) {
	if flood := is.FloodGet(ip); flood != nil {
		Flooders[ip].Count++
		Flooders[ip].Expired = time.Now().Add(time.Minute * time.Duration(is.AuthFloodDuration)).Unix()
	} else {
		Flooders[ip] = &Flood{
			Count:   1,
			Expired: time.Now().Add(time.Minute * time.Duration(is.AuthFloodDuration)).Unix(),
		}
	}

}

func (is *InternalService) FloodGet(ip string) *Flood {
	if flood, ok := Flooders[ip]; ok {
		return flood
	}

	return nil
}

func (is *InternalService) isFlooder(ip string) bool {
	if flood := is.FloodGet(ip); flood != nil && flood.Count >= is.AuthFloodLimit && flood.Expired > time.Now().Unix() {
		return true
	}

	return false
}

func (is *InternalService) FloodClear() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				for i, v := range Flooders {
					if v.Expired < time.Now().Unix() {
						delete(Flooders, i)
					}
				}
			}
		}
	}()
}
