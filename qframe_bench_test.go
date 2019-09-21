package qframe_test

import (
	"bytes"
	qf "github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/groupby"
	"github.com/tobgu/qframe/types"
	"os"
	"testing"
)

const rowCount = 73861
const columnCount = 23
const dataFileName = "recipeData-utf8.csv"

func intRange(size int) []int {
	result := make([]int, size)
	for i := range result {
		result[i] = i
	}
	return result
}

func qframeReadCsv() (qf.QFrame, error) {
	f, err := os.Open(dataFileName)
	if err != nil {
		return qf.QFrame{}, err
	}

	defer f.Close()
	frame := qf.ReadCSV(f)
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
	// fmt.Println(f.Slice(0,2))
	// return
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		err := f.ToJSON(buf)
		if err != nil {
			b.Fatalf("Unexpected JSON error: %s", f.Err)
		}
	    // Previous length: 33363821, check diff
		if buf.Len() != 33407314 {
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

func F(comparator interface{}, column string, arg interface{}) qf.Filter {
	return qf.Filter{Comparator: comparator, Column: column, Arg: arg}
}

func BenchmarkQFrame_Filter(b *testing.B) {
	f, _ := qframeReadCsv()

	benchmarks := []struct {
		name          string
		filter        qf.FilterClause
		expectedCount int
	}{
		{"Float gt", F(">", "Size(L)", 21.0), 26823},
		{"Float custom gt", F(func(f float64) bool { return f > 21.0 }, "Size(L)", nil), 26823},
		{"Combine or", qf.Or(F(">", "Size(L)", 21.0), F(">", "StyleID", 100)), 39818},
		{"Combine and", qf.And(F(">", "Size(L)", 21.0), F(">", "StyleID", 100)), 7280},
		{"String eq", F("=", "Style", "Cream Ale"), 830},
		{"String like case sensitive", F("like", "Name", "%Ale%"), 9118},
		{"String like case insensitive", F("ilike", "Name", "%ale%"), 11912},
		{"String regex case sensitive", F("like", "Name", ".*Ale.*"), 9118},
		{"String regex case insensitive", F("ilike", "Name", ".*ale.*"), 11912},
		{"Integer in", F("in", "StyleID", intRange(100)), 53514},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := f.Filter(bc.filter)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}

				if newF.Len() != bc.expectedCount {
					b.Fatalf("Unexpected frame length: %d", newF.Len())
				}
			}
		})
	}
}

func BenchmarkQFrame_Eval(b *testing.B) {
	f, _ := qframeReadCsv()

	benchmarks := []struct {
		name   string
		expr   qf.Expression
		dstCol string
	}{
		{"Float abs", qf.Expr("abs", types.ColumnName("BoilSize")), "dstCol"},
		{"Add columns", qf.Expr("+", types.ColumnName("OG"), types.ColumnName("FG")), "dstCol"},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := f.Eval(bc.dstCol, bc.expr)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}

				if newF.Len() != f.Len() {
					b.Fatalf("Unexpected frame length: %d", newF.Len())
				}

				if len(newF.ColumnNames()) != len(f.ColumnNames())+1 {
					b.Fatalf("Unexpected column count: %d", len(newF.ColumnNames()))
				}
			}
		})
	}
}

func Mean(column []float64) float64 {
	sum := 0.0
	for _, val := range column {
		sum += val
	}
	return sum / float64(len(column))
}

func BenchmarkQFrame_Aggregate(b *testing.B) {
	f, _ := qframeReadCsv()

	benchmarks := []struct {
		name          string
		groupColumns  []string
		aggregations  []qf.Aggregation
		expectedCount int
	}{
		// Note that there is a difference of +1 rows here compared to the pandas benchmark. This is due to the fact that pandas does not
		// group by N/A while qframe does (in this case it actually reads it as the string "N/A")
		{"Single col string single float mean", []string{"Style"}, []qf.Aggregation{{Fn: Mean, Column: "OG"}}, 176},
		{"Single col integer single float mean", []string{"StyleID"}, []qf.Aggregation{{Fn: Mean, Column: "OG"}}, 176},
		{"Double col string single float mean", []string{"Style", "Color"}, []qf.Aggregation{{Fn: Mean, Column: "OG"}}, 39557},
		{"Single col string double float mean", []string{"Style"}, []qf.Aggregation{{Fn: Mean, Column: "OG"}, {Fn: Mean, Column: "FG"}}, 176},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := f.GroupBy(groupby.Columns(bc.groupColumns...)).Aggregate(bc.aggregations...)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}

				if newF.Len() != bc.expectedCount {
					b.Fatalf("Unexpected frame length: %d", newF.Len())
				}
			}
		})
	}
}

/*
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

func BenchmarkGota_Filter(b *testing.B) {
	f, _ := gotaReadCsv()

	benchmarks := []struct {
		name          string
		filter        func(df.DataFrame) df.DataFrame
		expectedCount int
	}{
		{
			name: "Float gt",
			filter: func(f df.DataFrame) df.DataFrame {
				return f.Filter(df.F{Colname: "Size(L)", Comparator: ">", Comparando: 21.0})
			},
			expectedCount: 26823},
		{
			name: "Combine or",
			filter: func(f df.DataFrame) df.DataFrame {
				return f.Filter(df.F{Colname: "Size(L)", Comparator: ">", Comparando: 21.0}, df.F{Colname: "StyleID", Comparator: ">", Comparando: 100})
			},
			expectedCount: 39818},
		{
			name: "Combine and",
			filter: func(f df.DataFrame) df.DataFrame {
				return f.Filter(df.F{Colname: "Size(L)", Comparator: ">", Comparando: 21.0}).Filter(df.F{Colname: "StyleID", Comparator: ">", Comparando: 100})
			},
			expectedCount: 7280},
		{
			name: "String eq",
			filter: func(f df.DataFrame) df.DataFrame {
				return f.Filter(df.F{Colname: "Style", Comparator: "==", Comparando: "Cream Ale"})
			},
			expectedCount: 830},
		{
			name: "Integer in",
			filter: func(f df.DataFrame) df.DataFrame {
				return f.Filter(df.F{Colname: "StyleID", Comparator: "in", Comparando: intRange(100)})
			},
			expectedCount: 53514},
	}

	for _, bc := range benchmarks {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newF := bc.filter(f)
				if newF.Err != nil {
					b.Fatalf("Unexpected error: %s", newF.Err.Error())
				}

				l, _ := newF.Dims()
				if l != bc.expectedCount {
					b.Fatalf("Unexpected frame length: %d", l)
				}
			}
		})
	}
}
*/

/*
QFrame
------
BenchmarkQFrame_ReadCsv-2            	       5	 207966949 ns/op	164317979 B/op	    1501 allocs/op
BenchmarkQFrame_WriteJsonRecords-2   	       5	 230018909 ns/op	69792521 B/op	      74 allocs/op
BenchmarkQFrame_Sort/UserId_-_Int-2  	     100	  10306055 ns/op	  303152 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Name_-__string-2         	      20	  86066995 ns/op	  303184 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Multi_column-2           	      10	 156379617 ns/op	  303344 B/op	       5 allocs/op
BenchmarkQFrame_Filter/Float_gt-2             	    2000	    925995 ns/op	  196608 B/op	       2 allocs/op
BenchmarkQFrame_Filter/Float_custom_gt-2    	    2000	   1130335 ns/op	  196608 B/op	       2 allocs/op
BenchmarkQFrame_Filter/Combine_or-2           	    1000	   1352508 ns/op	  245936 B/op	       4 allocs/op
BenchmarkQFrame_Filter/Combine_and-2          	    1000	   1197479 ns/op	  256800 B/op	       6 allocs/op
BenchmarkQFrame_Filter/String_eq-2            	    2000	   1076656 ns/op	   85376 B/op	       2 allocs/op
BenchmarkQFrame_Filter/String_like_case_sensitive-2         	     500	   3777651 ns/op	  122896 B/op	       3 allocs/op
BenchmarkQFrame_Filter/String_like_case_insensitive-2       	     100	  14342672 ns/op	  131704 B/op	      16 allocs/op
BenchmarkQFrame_Filter/String_regex_case_sensitive-2        	      20	  67469026 ns/op	  164712 B/op	      58 allocs/op
BenchmarkQFrame_Filter/String_regex_case_insensitive-2      	      20	  88342431 ns/op	  172972 B/op	      61 allocs/op
BenchmarkQFrame_Filter/Integer_in-2                         	     500	   3179964 ns/op	  304811 B/op	      10 allocs/op
BenchmarkQFrame_Eval/Float_abs-2         	    2000	    637177 ns/op	  612885 B/op	      41 allocs/op
BenchmarkQFrame_Eval/Add_columns-2       	    2000	    780004 ns/op	  612833 B/op	      40 allocs/op

// These are currently 10% - 20% slower than the Pandas equivalents
BenchmarkQFrame_Aggregate/Single_col_string_single_float_mean-2         	     100	  12083567 ns/op	 2475952 B/op	    1396 allocs/op
BenchmarkQFrame_Aggregate/Single_col_integer_single_float_mean-2        	     200	   6040112 ns/op	 2535856 B/op	    1389 allocs/op
BenchmarkQFrame_Aggregate/Double_col_string_single_float_mean-2         	      30	  34800359 ns/op	15473961 B/op	   48456 allocs/op
BenchmarkQFrame_Aggregate/Single_col_string_double_float_mean-2         	     100	  12443931 ns/op	 2627584 B/op	    1405 allocs/op

Gota
----
BenchmarkGota_ReadCSV-2                                     	       2	 758721612 ns/op	228591928 B/op	 3686954 allocs/op
BenchmarkGota_WriteJsonRecords-2                            	       1	2771840823 ns/op	482439320 B/op	 5828275 allocs/op
BenchmarkGota_Sort/UserId_-_Int-2                           	      30	  53656268 ns/op	42841668 B/op	     131 allocs/op
BenchmarkGota_Sort/Name_-__string-2                         	      10	 152335582 ns/op	48951630 B/op	     113 allocs/op
BenchmarkGota_Sort/Multi_column-2                           	       5	 285486561 ns/op	78037472 B/op	     241 allocs/op
BenchmarkGota_Filter/Float_gt-2                             	      50	  32328655 ns/op	38116730 B/op	     609 allocs/op
BenchmarkGota_Filter/Combine_or-2                           	      20	  51720372 ns/op	60522417 B/op	     663 allocs/op
BenchmarkGota_Filter/Combine_and-2                          	      30	  42981669 ns/op	48312308 B/op	    1103 allocs/op
BenchmarkGota_Filter/String_eq-2                            	     200	   9020487 ns/op	 1087112 B/op	     310 allocs/op
BenchmarkGota_Filter/Integer_in-2                           	      10	 189430304 ns/op	77769508 B/op	     779 allocs/op
*/
