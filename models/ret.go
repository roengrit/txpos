package models

import (
	"time"
)

//Provinces Provinces
type Provinces struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//RetModel _
type RetModel struct {
	RetOK      bool
	RetCount   int64
	RetData    string
	ID         int64
	Name       string
	Del        string
	Title      string
	Alert      string
	XSRF       string
	FlagAction string
	ListData   interface{}
	Data1      interface{}
	Data2      interface{}
	Data3      interface{}
}
