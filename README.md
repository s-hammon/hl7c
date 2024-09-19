# hl7c

hl7c is a CLI tool which can read a configuration file and generate Go code containing objects (types) and entities (models) which represent various business logic within the healthcare space. One could then unmarshal a "JSON-ified" HL7v2 message into these models using the configured JSON tags.

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

```
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
                    "value": "^~\&"
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
                            "value": "1
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

With hl7c, one could easily configure a YAML file which defines certain business logic which are expected from various HL7 message types, and even within a message type