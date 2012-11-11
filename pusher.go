package main


import "fmt"
import "../../../usr/lib/go/src/pkg/github.com/kr/beanstalk"
import "time"


type DataRow struct {
	x float64
	y float64
}

func main() {

	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	var tube = beanstalk.Tube{c, "unsorted"}
	id, err := tube.Put([]byte{DataRow{4,0},
		DataRow{2,0},
		DataRow{3,0},
		DataRow{6,0},DataRow{1,0},DataRow{5,0}}, 1, 0, 120*time.Second)
	fmt.Println("Id ", id)

	if err != nil {
		fmt.Println(err)
	}
}