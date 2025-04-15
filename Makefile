art:
	go run main.go art ;

diff:
	go run main.go diff -m test/resources/metrics.json a test/resources/left.csv test/resources/right.csv ; 

peek:
	go run main.go peek test/resources/left.csv ; 

chart:
	go run main.go chart -w 50 test/resources/left.csv 1 count ;
	go run main.go chart -w 23 test/resources/right.csv 1 sum ;
	go run main.go chart -w 10 test/resources/right.csv 1 max ;
	go run main.go chart -w 5 test/resources/chart.csv 1 count_distinct ;
	go run main.go chart test/resources/chart.csv 1 count ;

gen:
	go run main.go gen -n 2 -s test/resources/generator-schema.json ;
