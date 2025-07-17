package main

import "time"

type Lot struct {
	Id          string
	Status      string
	Name        string
	Description string
	CreateDate  time.Time
	BidEndTime  time.Time
	Image       []byte
	IsStopped   bool
	HasAppeals  bool
	IsAnnulled  bool
}
