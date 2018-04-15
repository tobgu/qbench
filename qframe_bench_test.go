package qframe_test

import (
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

/*
BenchmarkQFrame_ReadCsv-2   	       5	 204340010 ns/op	164317976 B/op	    1501 allocs/op
BenchmarkGota_ReadCSV-2     	       2	 794546955 ns/op	228592016 B/op	 3686955 allocs/op
*/
