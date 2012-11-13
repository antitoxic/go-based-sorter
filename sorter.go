package main

import (
	beanstalk "github.com/iwanbk/gobeanstalk"
	"fmt"
	"sort"
	"encoding/json"
	"log"
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

func (s Dataset) Len() int	{ return len(s.Data) }
func (s Dataset) Swap(i, j int)	{ s.Data[i], s.Data[j] = s.Data[j], s.Data[i] }

// Decending implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Dataset value.
type xDecending struct{ Dataset }
func (s xDecending) Less(i, j int) bool	{ return s.Data[i].X < s.Data[j].X }


func main() {
	// Connect to beanstalk
	beanstalkConn, err := beanstalk.Dial("127.0.0.1:11300")
	if err != nil {
		log.Fatal(err)
	}

	// Pick a tube 
	_, err = beanstalkConn.Watch("unsorted")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch a job
	job, err := beanstalkConn.Reserve()
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize 
	var jobData Dataset
	err = json.Unmarshal(job.Body, &jobData)
	if err != nil {
		log.Fatal(err)
	}

	// Sort 
	toSort := xDecending{jobData}
	sort.Sort(toSort)
	sorted := toSort.Dataset

	// Put back sorted data into a different beanstalk queue
	err = beanstalkConn.Use("sortedX")
	if err != nil {
		log.Fatal(err)
	}
	serializedJob, err := json.Marshal(sorted)
	if err != nil {
		log.Fatal(err)
	}
	_, err = beanstalkConn.Put(serializedJob, 0, 0, 120);
	if err != nil {
		log.Fatal(err)
	}
	
	// Delete successfully handled job
	beanstalkConn.Delete(job.Id)

	if err != nil {
		fmt.Println(err)
	}
}