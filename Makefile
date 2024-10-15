art:
	go run main.go art ;

diff:
	go run main.go diff -m ./test/resources/metrics.json a ./test/resources/left.csv ./test/resources/right.csv ; 

peek:
	go run main.go peek ./test/resources/left.csv ; 

chart:
	go run main.go chart ./test/resources/left.csv 1 ;
	go run main.go chart ./test/resources/right.csv 1 ;
	go run main.go chart ./test/resources/chart.csv 1 ;
