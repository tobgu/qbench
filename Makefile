pybench:
	py.test --benchmark-columns="Mean,Median,StdDev,Rounds" --benchmark-sort=Name

gobench:
	go test -bench=.