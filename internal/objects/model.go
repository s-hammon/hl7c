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

type CX struct {
	Id                 string `json:"1"`
	CheckDigit         string `json:"2"`
	CheckDigitScheme   string `json:"3"`
	AssigningAuthority string `json:"4"`
	IdentifierTypeCode string `json:"5"`
	AssigningFacility  string `json:"6"`
}

type Patient struct {
	Base
	Mrn  CX        `json:"PID.3"`
	Name string    `json:"PID.5"`
	Dob  time.Time `json:"PID.7"`
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
	m.Dob, _ = time.Parse("20060102", aux.Dob)

	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}
