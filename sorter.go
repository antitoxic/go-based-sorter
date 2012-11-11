package main


import "fmt"
import "../../../usr/lib/go/src/pkg/github.com/kr/beanstalk"
import "time"
import "sort"


type DataRow struct {
	x float
	y float
}

type Dataset[]*DataRow

func (s Dataset) Len() int	{ return len(s) }
func (s Dataset) Swap(i, j int)	{ s[i], s[j] = s[j], s[i] }

// Decending implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Dataset value.
type Decending struct{ Dataset }

func (s Decending) Less(i, j int) bool	{ return s.Dataset[i].x < s.Dataset[j].x }


func main() {

	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	tubeSet := beanstalk.NewTubeSet(c, "unosrted")
	id, body, err := tubeSet.Reserve(30 * time.Second)
	fmt.Println("Id ", id)
	fmt.Println(body)
	dataSet := Dataset{body}
	c.Delete(id)
	sort.Sort(Decending(dataSet))
	fmt.Println(dataSet)

	if err != nil {
		fmt.Println(err)
	}
}