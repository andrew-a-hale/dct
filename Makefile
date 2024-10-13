art:
	go run main.go art ;

diff:
	go run main.go diff -m ./test/resources/metrics.json a ./test/resources/left.csv ./test/resources/right.csv ; 

peek:
	go run main.go peek ./test/resources/left.csv ; 
