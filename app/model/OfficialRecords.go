package model

//OfficialRecords used to hold data from records
type OfficialRecords struct {
	ParcelID string `json:"parcelid"`
	Records  []Parcel
}

//Parcel used to hold each parcel record
type Parcel struct {
	ParcelID          string `json:"parcelid"`
	FirstDirectName   string `json:"firstdirectname"`
	FirstInDirectName string `json:"firstindirectname"`
	BookType          string `json:"booktype"`
	BookPage          string `json:"bookpage"`
	DateRecorded      string `json:"daterecorded"`
	DocType           string `json:"doctype"`
	InstrumentNumber  string `json:"instrumentnumber"`
	Legal             string `json:"legal"`
}
