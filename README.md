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
The benchmarks are all executed in a VirtualBox with two cores assigned
to it and 4 Gb of base memory. I'm aware of the limitations on repeatability,
etc. that running the benchmarks in a VM means. For most of the results
here it does not matter too much though. In most cases the difference
in performance is so great between the different implementations
that +- a couple of percent does not matter.

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
Can be found as comments in the end of `qframe_bench_test.go` and
`test_pandas_bench.py`. I hope to find the time to make a nice table
comparing the different data frames side by side some day.

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
```
make gobench
```