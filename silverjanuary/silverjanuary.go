package main

import (
	"C"
	"context"
	"fmt"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/cdata"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type CargoConfig struct {
	Host     string `yaml:"db_host"`
	Port     string `yaml:"db_port"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	Name     string `yaml:"db_name"`
}

//export arrowToDB
func arrowToDB(scptr uintptr, rbptr uintptr, createStatement *C.char, tablename *C.char, partition int64) {

	start := time.Now()

	data, err1 := os.ReadFile("config.yaml")

	CheckError(err1)

	var config CargoConfig

	err2 := config.Parse(data)

	CheckError(err2)

	// Received Create Table statement from pd
	goCreateStatement := C.GoString(createStatement)

	// Received pointer to the arrow schema and array
	schema := cdata.SchemaFromPtr(scptr)
	arr := cdata.ArrayFromPtr(rbptr)

	// Import arrow files
	rec, err := cdata.ImportCRecordBatch(arr, schema)

	// make sure we call the release callback when we're done
	defer rec.Release()

	// String for pgx -> postgres connection
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	// Execute create statement
	_, err = conn.Exec(ctx, goCreateStatement)

	// Slices of slices for pgx COPY
	table := make([][]interface{}, rec.NumRows())
	numRow := int(rec.NumRows())
	numCol := int(rec.NumCols())

	// Get the columnNames from Arrow Schema
	columNames := make([]string, rec.NumCols())

	for i := 0; i < numCol; i++ {
		columNames[i] = rec.ColumnName(i)
	}

	// Get the number of partitions
	numPartitions := int(partition)
	goTablename := C.GoString(tablename)

	// Create WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numPartitions)

	// Insert into postgres with pgx COPY in parallel
	for i := 0; i < numPartitions; i++ {
		go func(partition int) {
			// Create context
			ctx := context.Background()
			// Connect to postgres
			conn, err := pgx.Connect(ctx, connString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
				os.Exit(1)
			}
			// Defer close connection
			defer conn.Close(ctx)

			// Get the lower and upper bound of the partition
			lowerBound := partition * (numRow / numPartitions)
			upperBound := (partition + 1) * (numRow / numPartitions)
			if partition == numPartitions-1 {
				upperBound = numRow
			}
			// Reconstruct Arrow table and take values out of arrow array

			for _, v := range rec.Columns() {
				for i := lowerBound; i < upperBound; i++ {
					switch v.(type) {
					case *array.String:
						table[i] = append(table[i], v.(*array.String).Value(i))
					case *array.Boolean:
						table[i] = append(table[i], v.(*array.Boolean).Value(i))
					case *array.Float32:
						table[i] = append(table[i], v.(*array.Float32).Value(i))
					case *array.Float64:
						table[i] = append(table[i], v.(*array.Float64).Value(i))
					case *array.Date32:
						table[i] = append(table[i], v.(*array.Date32).Value(i).FormattedString())
					case *array.Date64:
						table[i] = append(table[i], v.(*array.Date64).Value(i).FormattedString())
					case *array.Int32:
						table[i] = append(table[i], v.(*array.Int32).Value(i))
					case *array.Int64:
						table[i] = append(table[i], v.(*array.Int64).Value(i))
					}
				}
			}

			// Copy data to postgres
			_, err = conn.CopyFrom(ctx, pgx.Identifier{goTablename}, columNames, pgx.CopyFromRows(table[lowerBound:upperBound]))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to copy data to postgres: %v\n", err)
				os.Exit(1)
			}
			// Done with this partition
			wg.Done()
		}(i)
	}
	// Wait for all goroutines to finish
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Partition: %v Time taken: %s\n", numPartitions, elapsed)
}

func main() {}

func (c *CargoConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
