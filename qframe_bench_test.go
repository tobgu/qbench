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

/*
BenchmarkQFrame_ReadCsv-2            	       5	 202199966 ns/op	164318040 B/op	    1500 allocs/op
BenchmarkQFrame_WriteJsonRecords-2   	       5	 224061889 ns/op	69792524 B/op	      74 allocs/op
BenchmarkGota_ReadCSV-2              	       2	 771431210 ns/op	228591928 B/op	 3686954 allocs/op
BenchmarkGota_WriteJsonRecords-2     	       1	2298617520 ns/op	482345720 B/op	 5827950 allocs/op
*/
