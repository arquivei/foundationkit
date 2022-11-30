// This is a simple example that shows how to use avroutil package.
// This can be run with `go run main.go avroschemas.go`
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/arquivei/foundationkit/avroutil"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/schemaregistry/implschemaregistry"
)

// exampleStruct is a struct that contains all fields related to SchemaExample
type exampleStruct struct {
	Field1 int64  `avro:"Field1"`
	Field2 string `avro:"Field2"`
}

func main() {
	// Just some basic logger initialization
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create an example data
	exampleData := exampleStruct{
		Field1: 1,
		Field2: "2",
	}

	// Setup schema registry mock to avoid external connection
	mockSchemaRegistry := implschemaregistry.MustNewMock(map[schemaregistry.ID]string{
		1: schemaExample,
	})

	// Create a new encoder based in schemaExample
	encoder, err := avroutil.NewEncoder(context.Background(), mockSchemaRegistry, avroSubjectschemaExample, schemaExample)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to initialize new encoder")
	}

	// Create a new decoder with mock schema registry
	decoder := avroutil.NewDecoder(mockSchemaRegistry)

	// Let's encode example data
	encodedData, err := encoder.Encode(exampleData)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to encode data")
	}

	// Let's decode the output encoded data
	var decodedData exampleStruct
	err = decoder.Decode(context.Background(), encodedData, &decodedData)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to decode data")
	}

	fmt.Printf("Example data %+v\n", exampleData)
	fmt.Printf("Decoded data %+v\n", decodedData)

	// Let's just compare the example data and output decode data for fun.
	if exampleData == decodedData {
		log.Info().Msg("Encoded and Decoded data with success")
	}
}
