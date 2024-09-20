package objects

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Base struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CM_MSH struct {
	MessageType  string `json:"1"`
	TriggerEvent string `json:"2"`
}

type CX struct {
	Id                 string `json:"1"`
	CheckDigit         string `json:"2"`
	CheckDigitScheme   string `json:"3"`
	AssigningAuthority string `json:"4"`
	IdentifierTypeCode string `json:"5"`
	AssigningFacility  string `json:"6"`
}

type XPN struct {
	LastName     string `json:"1"`
	FirstName    string `json:"2"`
	MiddleName   string `json:"3"`
	Suffix       string `json:"4"`
	Prefix       string `json:"5"`
	Degree       string `json:"6"`
	NameTypeCode string `json:"7"`
}

type XAD struct {
	StreetAddress    string `json:"1"`
	OtherDesignation string `json:"2"`
	City             string `json:"3"`
	State            string `json:"4"`
	Zip              string `json:"5"`
	Country          string `json:"6"`
	AddressType      string `json:"7"`
}

type Message struct {
	Base
	SendingApp      string    `json:"MSH.3"`
	SendingFacility string    `json:"MSH.4"`
	SentAt          time.Time `json:"MSH.7"`
	MessageType     CM_MSH    `json:"MSH.9"`
	ControlID       string    `json:"MSH.10"`
	ProcessingID    string    `json:"MSH.11"`
	VersionID       string    `json:"MSH.12"`
}

func (m *Message) UnmarshalJSON(b []byte) error {
	type alias Message
	aux := &struct {
		SentAt string `json:"MSH.7"`
		*alias
	}{
		alias: (*alias)(m),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	m.SentAt, _ = time.Parse("20060102150405", aux.SentAt)

	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

type Patient struct {
	Base
	Mrn     CX        `json:"PID.3"`
	Name    XPN       `json:"PID.5"`
	Dob     time.Time `json:"PID.7"`
	Sex     string    `json:"PID.8"`
	Address XAD       `json:"PID.11"`
}

func (m *Patient) UnmarshalJSON(b []byte) error {
	type alias Patient
	aux := &struct {
		Dob string `json:"PID.7"`
		*alias
	}{
		alias: (*alias)(m),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	m.Dob, _ = time.Parse("20060102150405", aux.Dob)

	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}
