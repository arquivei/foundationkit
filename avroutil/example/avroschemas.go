package main

import (
	// Embedding avro schemas as strings
	_ "embed"

	"github.com/arquivei/foundationkit/schemaregistry"
)

// schemaExample is an avro schema that contains all data related to SchemaExample
//
//go:embed assets/SchemaExample.avsc
var schemaExample string

// avroSubjectschemaExample is the subject that can be used to find the SchemaExample in a schema registry
var avroSubjectschemaExample schemaregistry.Subject = "com.arquivei.subject"
