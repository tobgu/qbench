
## Install Python benchmarks
```
virtualenv -p python3.6 pvenv
./pvenv/bin/activate
pip install -r requirements.txt
```

## Run Python benchmarks
```
py.test
```

## Run Go benchmarks
```
go test -bench=. -run=^$
```