package main

import (
	"context"
	"fmt"
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/csv"
	"github.com/jackc/pgx/v5"
	"os"
	"sync"
	"time"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "admin"
	password = "password"
	dbname   = "january"
)

var files = []string{
	"./january/people-100.csv",
	"./january/people-1000.csv",
	"./january/people-10000.csv",
	"./january/people-100000.csv",
}

var runTime []int

var numPartitionList = []int{16, 32, 64}

var numRows = map[string]int{
	"./january/people-100.csv":    100,
	"./january/people-1000.csv":   1000,
	"./january/people-10000.csv":  10000,
	"./january/people-100000.csv": 100000,
}

func main() {
	for _, numPartitions := range numPartitionList {
		for _, file := range files {
			// Open the CSV file
			f, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			schema := arrow.NewSchema([]arrow.Field{
				{Name: "index", Type: arrow.PrimitiveTypes.Int64},
				{Name: "userid", Type: arrow.BinaryTypes.String},
				{Name: "firstname", Type: arrow.BinaryTypes.String},
				{Name: "lastname", Type: arrow.BinaryTypes.String},
				{Name: "sex", Type: arrow.BinaryTypes.String},
				{Name: "email", Type: arrow.BinaryTypes.String},
				{Name: "phone", Type: arrow.BinaryTypes.String},
				{Name: "dateofbirth", Type: arrow.FixedWidthTypes.Date32},
				{Name: "jobtitle", Type: arrow.BinaryTypes.String},
			}, nil)

			numRow := numRows[file]
			// Read the CSV data into a struct or map
			r := csv.NewReader(
				f, schema,
				csv.WithComma(','),
				csv.WithHeader(true),
				csv.WithChunk(numRow),
			)
			defer r.Release()
			table := make([][]interface{}, numRow)

			for r.Next() {
				rec := r.Record()

				connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

				// Connect to the database
				conn, err := pgx.Connect(context.Background(), connString)
				if err != nil {
					fmt.Println(err)
				}
				defer conn.Close(context.Background())

				// Get the columnNames from Arrow Schema
				columNames := make([]string, 9)
				for i := 0; i < 9; i++ {
					columNames[i] = schema.Field(i).Name
				}

				// Insert into postgres with pgx COPY
				tableName := "people"

				start := time.Now()
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
						_, err = conn.CopyFrom(ctx, pgx.Identifier{tableName}, columNames, pgx.CopyFromRows(table[lowerBound:upperBound]))
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
				fmt.Printf("File: %s . Partition: %v Time taken: %s\n", file, numPartitions, elapsed)
				defer rec.Release()
			}
		}
	}
}
