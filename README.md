# hl7c

hl7c is a CLI tool which can read a configuration file and generate Go code containing value objects (types) and entities (models) which represent various business logic within the healthcare space. One could then unmarshal a "JSON-ified" HL7v2 message into these models using the configured JSON tags.

## Installation

### go install
Requires Go 1.21+

`go install github.com/s-hammon/hl7c/cmd/hl7c`

## HL7 in JSON

A "JSON-ified" HL7 is a JSON representation of an HL7 wherein the segments, fields, components, etc. are nested as key-value pairs based on the hierarchy of HL7 objects:

```
HL7
└── segment
    ├── repeat
    │   └── field
    │       └── component
    │           └── subcomponent
    └── field
        └── component
            └── subcomponent
```

For example, an HL7 could be represented as:

```JSON
{
    "segments" [
        {
            "name": "MSH",
            "value": [
                {
                    "name": "MSH.1",
                    "value": "|",
                },
                {
                    "name": "MSH.2",
                    "value": "^~\\&"
                },
                {
                    "name": "MSH.3",
                    "value": "Imagecast RIS"
                },
                ...
            ],
        },
        {
            "name": "PID",
            "value": [
                {
                    "name": "PID.1",
                    "value": "1"
                },
                {
                    "name": "PID.2",
                    "value": [
                        {
                            "name": "PID.2.1",
                            "value": "0123456789"
                        },
                        {
                            "name": "PID.2.2",
                            "value": "6"
                        },
                        {
                            "name": "PID.2.3",
                            "value": "1"
                        },
                    ]
                },
                ...
            ]
        },
        ...
    ]
}
```

## Why do this?

HL7--even within the new FHIR standard--has **poor semantic interoperability**, despite being (relatively) syntactically rigorious. Although there are conventions in HL7 structure, in truth one could put any kind of data in HL7 format. In fact, even the delimiters could vary between healthcare organizations, although this is quite rare.

This make developing healthcare applications difficult, as one cannot ensuse that each of their healthcare partners store the exact same information in the exact stame spots--that is, not without expensive and/or cumbersome software which serves as an HL7 integration engine. Further, the core business logic within an application does not--and *should* not--change with "different types" of HL7 messaging.

The structure of an HL7 message--specifically, a v2 message--is also archaic. Often, restructuring this format (as described above) is a good first step in digesting the data contained within them. This is where hl7c comes into play.

With hl7c, one can configure the name of the business logic, the data type, and the corresponding HL7 field (represented as a JSON key) all in one YAML file. In the case of "composite" HL7 fields (often things with multiple components like patient names, addresses, etc), one can also define custom types.

For example, the following JSON represents some key pieces of information we want to extract from an `ORM` message:

```JSON
{
    "MSH.1": "|",
    "MSH.2": "^~\\&",
    "MSH.3": "Centricity RIS-IC",
    "MSH.4": "ABC Hospital",
    "MSH.7": "20240308075809",
    "MSH.9": {
        "1": "ORM",
        "2": "O01"
    },
    "MSH.10": "2024030807580934",
    "MSH.11": "P",
    "MSH.12": "2.3",
    "PID.3": {
        "1": "123456",
        "6": "ABC"
    },
    "PID.5" : {
        "1": "DOE",
        "2": "JOHN"
    },
    "PID.7": "19910905",
    "PID.8": "M",
    "PID.11": {
        "1": "123 MAIN ST",
        "3": "NEW YORK",
        "4": "NY",
        "5": "10001"
    }
}
```

With the following YAML, hl7c will generate Go structs which will unmarshal the data in this JSON:

```YAML
meta:
  package: internal/objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: message
    fields:
      - name: sendingApp
        type: string
        tag: MSH.3
      - name: sendingFacility
        type: string
        tag: MSH.4
      - name: sentAt
        type: timestamp
        tag: MSH.7
      - name: messageType
        type: CM_MSH
        tag: MSH.9
      - name: controlID
        type: string
        tag: MSH.10
      - name: processingID
        type: string
        tag: MSH.11
      - name: versionID
        type: string
        tag: MSH.12
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: XPN
        tag: PID.5
      - name: dob
        type: timestamp
        tag: PID.7
      - name: sex
        type: string
        tag: PID.8
      - name: address
        type: XAD
        tag: PID.11

types:
  - name: CM_MSH
    fields:
      - name: messageType
      - name: triggerEvent
  - name: CX
    fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
  - name: XPN
    fields:
      - name: lastName
      - name: firstName
      - name: middleName
      - name: suffix
      - name: prefix
      - name: degree
      - name: nameTypeCode
  - name: XAD
    fields:
      - name: streetAddress
      - name: otherDesignation
      - name: city
      - name: state
      - name: zip
      - name: country
      - name: addressType
```

Running `hl7c generate` in the same directory as this YAML file will create a module called `models.go` at the path `internal/objects`:

```Go
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
```

Each of these models extend a `Base` model with fields that you might expect to use when working with a normalized database. In addition, hl7c will handle formatting and dependencies under the hood--although, you **must** have already initialized a Go project using `go mod init`. Further, currently this only supports `time.Time` and `uuid.UUID` (from `github.com/google/uuid`) types.

## Todo

1. Add support to define inbound HL7 JSON structures
2. Add support to handle deps for other common package types ()
3. Improve README.md
4. Be able to read other types of config files (JSON, XML, etc)