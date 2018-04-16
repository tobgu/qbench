package qframe_test

import (
	"bytes"
	df "github.com/kniren/gota/dataframe"
	qf "github.com/tobgu/qframe"
	"os"
	"testing"
)

const rowCount = 73861
const columnCount = 23
const dataFileName = "recipeData-utf8.csv"

func qframeReadCsv() (qf.QFrame, error) {
	f, err := os.Open(dataFileName)
	if err != nil {
		return qf.QFrame{}, err
	}

	defer f.Close()
	frame := qf.ReadCsv(f)
	return frame, frame.Err
}

func BenchmarkQFrame_ReadCsv(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f, err := qframeReadCsv()
		if err != nil {
			b.Fatalf("Unexpected CSV error: %s", err.Error())
		}

		if f.Len() != rowCount {
			b.Fatalf("Unexpected row count: %d", f.Len())
		}

		if len(f.ColumnNames()) != columnCount {
			b.Fatalf("Unexpected column count: %d", len(f.ColumnNames()))
		}
	}
}

func BenchmarkQFrame_WriteJsonRecords(b *testing.B) {
	b.ReportAllocs()
	f, _ := qframeReadCsv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		f.ToJson(buf, "records")
		if f.Err != nil {
			b.Fatalf("Unexpected JSON error: %s", f.Err)
		}

		if buf.Len() != 33363821 {
			b.Fatalf("Unexpected JSON length: %d", buf.Len())
		}
	}
}

func BenchmarkQFrame_Sort(b *testing.B) {
	f, _ := qframeReadCsv()

	benchmarks := []struct {
		name   string
		orders []qf.Order
	}{
		{"UserId - Int", []qf.Order{{Column: "UserId"}}},
		{"Name -  string", []qf.Order{{Column: "Name"}}},
		{"Multi column", []qf.Order{{Column: "Style"}, {Column: "Name"}, {Column: "BrewMethod"}}},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := f.Sort(bc.orders...)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}

				if newF.Len() != f.Len() {
					b.Fatalf("Unexpected frame length: %d", newF.Len())
				}
			}

		})
	}
}

func gotaReadCsv() (df.DataFrame, error) {
	f, err := os.Open(dataFileName)
	if err != nil {
		return df.DataFrame{}, err
	}

	defer f.Close()
	frame := df.ReadCSV(f)
	return frame, frame.Err
}

func BenchmarkGota_ReadCSV(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f, err := gotaReadCsv()
		if err != nil {
			b.Fatalf("Unexpected CSV error: %s", err.Error())
		}

		xDim, yDim := f.Dims()
		if xDim != rowCount {
			b.Fatalf("Unexpected row count: %d", xDim)
		}

		if yDim != columnCount {
			b.Fatalf("Unexpected column count: %d", yDim)
		}
	}
}

func BenchmarkGota_WriteJsonRecords(b *testing.B) {
	b.ReportAllocs()
	f, _ := gotaReadCsv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		f.WriteJSON(buf)
		if f.Err != nil {
			b.Fatalf("Unexpected JSON error: %s", f.Err)
		}

		if buf.Len() != 33409725 {
			b.Fatalf("Unexpected JSON length: %d", buf.Len())
		}
	}
}

func BenchmarkGota_Sort(b *testing.B) {
	f, _ := gotaReadCsv()

	benchmarks := []struct {
		name   string
		orders []df.Order
	}{
		{"UserId - Int", []df.Order{{Colname: "UserId"}}},
		{"Name -  string", []df.Order{{Colname: "Name"}}},
		{"Multi column", []df.Order{{Colname: "Style"}, {Colname: "Name"}, {Colname: "BrewMethod"}}},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := f.Arrange(bc.orders...)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}
			}
		})
	}
}

/*
BenchmarkQFrame_ReadCsv-2            	       5	 202199966 ns/op	164318040 B/op	    1500 allocs/op
BenchmarkQFrame_WriteJsonRecords-2   	       5	 224061889 ns/op	69792524 B/op	      74 allocs/op
BenchmarkGota_ReadCSV-2              	       2	 771431210 ns/op	228591928 B/op	 3686954 allocs/op
BenchmarkGota_WriteJsonRecords-2     	       1	2298617520 ns/op	482345720 B/op	 5827950 allocs/op

BenchmarkQFrame_Sort/UserId_-_Int-2         	     100	  10522800 ns/op	  303152 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Name_-__string-2       	      20	  92802526 ns/op	  303184 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Multi_column-2         	      10	 174326684 ns/op	  303344 B/op	       5 allocs/op
BenchmarkGota_Sort/UserId_-_Int-2           	      30	  54875433 ns/op	42841667 B/op	     131 allocs/op
BenchmarkGota_Sort/Name_-__string-2         	      10	 159895352 ns/op	48951627 B/op	     113 allocs/op
BenchmarkGota_Sort/Multi_column-2           	       5	 289134941 ns/op	78037472 B/op	     241 allocs/op
*/
