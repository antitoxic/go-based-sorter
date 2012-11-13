package main


import (
	beanstalk "github.com/iwanbk/gobeanstalk"
	"log"
	"encoding/json"
)

type DataRow struct {
	X float64
	Y float64
	Xi float64
	Yi float64
}

type Dataset struct {
	Spearman float64
	Significance float64
	N int
	Data []DataRow
}

func main() {
	// Connect to beanstalk
	beanstalkConn, err := beanstalk.Dial("127.0.0.1:11300")
	if err != nil {
		log.Fatal(err)
	}

	// Pick a tube 
	err = beanstalkConn.Use("unsorted")
	if err != nil {
		log.Fatal(err)
	}

	data := []DataRow{
		{4,0,0,0},
		{2,0,0,0},
		{3,0,0,0},
		{6,0,0,0},
		{1,0,0,0},
		{5,0,0,0},
	}

	// Serialize
	serializedJob, err := json.Marshal(Dataset{2,7,3,data})
	if err != nil {
		log.Fatal(err)
	}

	// Store job in queue
	_, err = beanstalkConn.Put(serializedJob, 0, 0, 120)
	if err != nil {
		log.Fatal(err)
	}
}