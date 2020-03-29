This repository contains a number of basic benchmarks that compare
the performance of [QFrame](https://github.com/tobgu/qframe) with
[Pandas](https://pandas.pydata.org) and [Gota](https://github.com/kniren/gota) (where applicable).

As always with benchmarks: Take these results with a grain of salt and
benchmark your own use cases to get a proper feel for the performance.

## About
The benchmarks in this repository have mostly been constructed with
the use case of [qocache](https://github.com/tobgu/qocache) in mind.
All optimizations done on QFrame so far have also targeted these use
cases.

If you have ideas for new benchmarks or improvements to the existing
ones please don't hesitate to open an issue.

All benchmarks use the [Brewer's Friend Beer Recipes](https://www.kaggle.com/jtrofe/beer-recipes)
dataset. It contains ~75000 lines and is also checked in to this repository.

Most of the benchmark operations are nonsensical from a beer brewers
perspective.

## Environment
The results presented here are executed in a VirtualBox with two cores assigned
to it and 4 Gb of base memory. The host OS is Windows 10 and the guest OS
Ubuntu Linux. Running the benchmarks in a VM like this may have some negative
effects on repeatability and consistency. For most of the results
here it does not matter too much though. These benchmarks are pretty
high level and the difference in performance is so great between the
different implementations that +- a couple of percent does should not
matter for the interpretation of them.

It would be really interesting if someone would like to run this on
a modern bare metal machine for comparison!

Processor:
```
$ cat /proc/cpuinfo
...
model name	: Intel(R) Core(TM) i7-3517U CPU @ 1.90GHz
cache size	: 4096 KB
...
```

OS:
```
$ lsb_release -a
No LSB modules are available.
Distributor ID:	Ubuntu
Description:	Ubuntu 14.04 LTS
Release:	14.04
Codename:	trusty
```

Go:
```
$ go version
go version go1.10 linux/amd64
```

Pandas:
```
>>> pd.show_versions()
INSTALLED VERSIONS
------------------
commit: None
python: 3.6.3.final.0
python-bits: 64
OS: Linux
OS-release: 3.13.0-83-generic
machine: x86_64
processor: x86_64
byteorder: little
LC_ALL: None
LANG: en_US.UTF-8
LOCALE: en_US.UTF-8

pandas: 0.22.0
pytest: 3.5.0
pip: 9.0.3
setuptools: 39.0.1
Cython: None
numpy: 1.14.2
scipy: None
pyarrow: None
xarray: None
IPython: None
sphinx: None
patsy: None
dateutil: 2.7.2
pytz: 2018.4
blosc: None
bottleneck: None
tables: None
numexpr: None
feather: None
matplotlib: None
openpyxl: None
xlrd: None
xlwt: None
xlsxwriter: None
lxml: None
bs4: None
html5lib: None
sqlalchemy: None
pymysql: None
psycopg2: None
jinja2: None
s3fs: None
fastparquet: None
pandas_gbq: None
pandas_datareader: None
```

## Results
### QFrame
```
BenchmarkQFrame_ReadCsv-2            	                                    5	 207966949 ns/op	164317979 B/op	    1501 allocs/op
BenchmarkQFrame_WriteJsonRecords-2   	                                    5	 230018909 ns/op	69792521 B/op	      74 allocs/op
BenchmarkQFrame_Sort/UserId_-_Int-2  	                                  100	  10306055 ns/op	  303152 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Name_-__string-2         	                           20	  86066995 ns/op	  303184 B/op	       3 allocs/op
BenchmarkQFrame_Sort/Multi_column-2           	                           10	 156379617 ns/op	  303344 B/op	       5 allocs/op
BenchmarkQFrame_Filter/Float_gt-2             	                         2000	    925995 ns/op	  196608 B/op	       2 allocs/op
BenchmarkQFrame_Filter/Float_custom_gt-2    	                         2000	   1130335 ns/op	  196608 B/op	       2 allocs/op
BenchmarkQFrame_Filter/Combine_or-2           	                         1000	   1352508 ns/op	  245936 B/op	       4 allocs/op
BenchmarkQFrame_Filter/Combine_and-2          	                         1000	   1197479 ns/op	  256800 B/op	       6 allocs/op
BenchmarkQFrame_Filter/String_eq-2            	                         2000	   1076656 ns/op	   85376 B/op	       2 allocs/op
BenchmarkQFrame_Filter/String_like_case_sensitive-2         	          500	   3777651 ns/op	  122896 B/op	       3 allocs/op
BenchmarkQFrame_Filter/String_like_case_insensitive-2       	          100	  14342672 ns/op	  131704 B/op	      16 allocs/op
BenchmarkQFrame_Filter/String_regex_case_sensitive-2        	           20	  67469026 ns/op	  164712 B/op	      58 allocs/op
BenchmarkQFrame_Filter/String_regex_case_insensitive-2      	           20	  88342431 ns/op	  172972 B/op	      61 allocs/op
BenchmarkQFrame_Filter/Integer_in-2                         	          500	   3179964 ns/op	  304811 B/op	      10 allocs/op
BenchmarkQFrame_Eval/Float_abs-2         	                             2000	    637177 ns/op	  612885 B/op	      41 allocs/op
BenchmarkQFrame_Eval/Add_columns-2       	                             2000	    780004 ns/op	  612833 B/op	      40 allocs/op

// These are currently 10% - 20% slower than the Pandas equivalents
BenchmarkQFrame_Aggregate/Single_col_string_single_float_mean-2    	      100	  12083567 ns/op	 2475952 B/op	    1396 allocs/op
BenchmarkQFrame_Aggregate/Single_col_integer_single_float_mean-2   	      200	   6040112 ns/op	 2535856 B/op	    1389 allocs/op
BenchmarkQFrame_Aggregate/Double_col_string_single_float_mean-2    	       30	  34800359 ns/op	15473961 B/op	   48456 allocs/op
BenchmarkQFrame_Aggregate/Single_col_string_double_float_mean-2    	      100	  12443931 ns/op	 2627584 B/op	    1405 allocs/op
```

### Gota
```
BenchmarkGota_ReadCSV-2                                            	       2	 758721612 ns/op	228591928 B/op	 3686954 allocs/op
BenchmarkGota_WriteJsonRecords-2                                   	       1	2771840823 ns/op	482439320 B/op	 5828275 allocs/op
BenchmarkGota_Sort/UserId_-_Int-2                                  	      30	  53656268 ns/op	42841668 B/op	     131 allocs/op
BenchmarkGota_Sort/Name_-__string-2                                	      10	 152335582 ns/op	48951630 B/op	     113 allocs/op
BenchmarkGota_Sort/Multi_column-2                                  	       5	 285486561 ns/op	78037472 B/op	     241 allocs/op
BenchmarkGota_Filter/Float_gt-2                                    	      50	  32328655 ns/op	38116730 B/op	     609 allocs/op
BenchmarkGota_Filter/Combine_or-2                                  	      20	  51720372 ns/op	60522417 B/op	     663 allocs/op
BenchmarkGota_Filter/Combine_and-2                                 	      30	  42981669 ns/op	48312308 B/op	    1103 allocs/op
BenchmarkGota_Filter/String_eq-2                                   	     200	   9020487 ns/op	 1087112 B/op	     310 allocs/op
BenchmarkGota_Filter/Integer_in-2                                  	      10	 189430304 ns/op	77769508 B/op	     779 allocs/op
```

### Pandas
```
Name (time in ms)                                                         Mean              Median            StdDev        Rounds
----------------------------------------------------------------------------------------------------------------------------------
test_aggregation[double col string single float mean-<lambda>-39065]     30.4766 (5.90)     30.0730 (6.37)    1.6894 (1.62)     31
test_aggregation[single col int single float mean-<lambda>-176]           5.1619 (1.0)       4.7210 (1.0)     3.3494 (3.22)    129
test_aggregation[single col string double float mean-<lambda>-175]       11.0780 (2.15)     10.8296 (2.29)    1.0403 (1.0)      71
test_aggregation[single col string single float mean-<lambda>-175]       10.8325 (2.10)     10.5123 (2.23)    1.2459 (1.20)     73
test_filter[combine and-<lambda>-7280]                                    4.6996 (1.0)       4.5539 (1.0)     0.4474 (1.0)     192
test_filter[combine or-<lambda>-39818]                                   10.2313 (2.18)     10.0739 (2.21)    0.6902 (1.54)     89
test_filter[contains case insensitive-<lambda>-11912]                   100.4382 (21.37)    99.3378 (21.81)   3.4898 (7.80)     11
test_filter[contains case sensitive-<lambda>-9118]                       59.0800 (12.57)    58.2704 (12.80)   2.0572 (4.60)     18
test_filter[integer in-<lambda>-53514]                                   11.6831 (2.49)     11.4531 (2.51)    0.9151 (2.05)     85
test_filter[regex case insensitive-<lambda>-11912]                      304.9066 (64.88)   306.2898 (67.26)   2.9705 (6.64)      5
test_filter[regex case sensitive-<lambda>-9118]                         126.5499 (26.93)   125.9932 (27.67)   3.2897 (7.35)      8
test_filter[single float-<lambda>-26823]                                  7.3997 (1.57)      7.1886 (1.58)    0.8638 (1.93)    113
test_filter[string eq-<lambda>-830]                                      10.1308 (2.16)      9.9458 (2.18)    1.0313 (2.31)     92
test_read_csv                                                           416.8440 (81.83)   416.2372 (84.13)   2.7442 (5.98)      5
test_sort[columns0]                                                      22.8061 (4.48)     22.5379 (4.56)    1.0598 (2.31)     42
test_sort[columns1]                                                     142.2413 (27.92)   143.0750 (28.92)   2.4375 (5.31)      8
test_sort[columns2]                                                     184.1537 (36.15)   183.4041 (37.07)   5.8602 (12.77)     6
test_write_json_records                                                 317.4334 (62.31)   318.0416 (64.28)   4.9470 (10.78)     5
test_eval[float abs-destCol = abs(BoilSize)]                             14.4224 (1.0)      14.1708 (1.0)     2.0511 (1.45)     29
test_eval[float add-destCol = OG + FG]                                   15.8143 (1.10)     15.8394 (1.12)    1.4177 (1.0)      41
----------------------------------------------------------------------------------------------------------------------------------
```

### Summary
Overall QFrame performs well ahead of Pandas in most benchmarks. The
only benchmarks that it's currently slower at are those containing
some sort of grouping. Here Pandas beats QFrame by 10 - 20% in runtime.

Compared to Gota QFrame is much faster in all benchmarks.

## Install Python benchmarks
```
virtualenv -p python3.6 pvenv
./pvenv/bin/activate
pip install -r requirements.txt
```

## Run Python benchmarks
```
make pybench
```

## Run Go benchmarks
Go dep is used for dependency management so it needs to be installed.

```
dep ensure
make gobench
```
