package main

import (
	"fmt"
	"github.com/intelsdi-x/snap-plugin-collector-mysql/mysqlplugin"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
)

func main() {
	p := mysqlplugin.New()

	cfg1 := plugin.NewPluginConfigType()
	cfg2 := cdata.NewNode()
	cfg1.AddItem("mysql_connection_string", ctypes.ConfigValueStr{Value: "root:r00tme@tcp(localhost:3306)/"})
	cfg2.AddItem("mysql_connection_string", ctypes.ConfigValueStr{Value: "root:r00tme@tcp(localhost:3306)/"})

	cfg1.AddItem("mysql_use_innodb", ctypes.ConfigValueBool{Value: true})
	cfg2.AddItem("mysql_use_innodb", ctypes.ConfigValueBool{Value: true})

	_ = cfg2

	mts1, err := p.GetMetricTypes(cfg1)
	mts2, err := p.CollectMetrics(mts1)

	_ = err
	for _, x := range mts2 {
		fmt.Println(x.Namespace(), x.Data())
	}

	//fmt.Printf("%v\n%#v\n", err, mts2)
}

/*
func main() {
	a, b := stats.New("root:r00tme@tcp(localhost:3306)/")
	fmt.Println(*a, b)

	x := mysqlplugin.Fee()
	x.StatsSource = a
	x.UseInnodb = true

	//q, e := x.Discover()
	//fmt.Printf("%v\n%#v\n", e, q)
	//x, e := a.GetSlaveStatus() //a.GetStatus(true)
	//fmt.Printf("%v\n%#v\n", e, x)

	x.Collect(map[int]bool{0: true})
	time.Sleep(3 * time.Second)

	q, e := x.Collect(map[int]bool{0: true, 3: true})

	fmt.Printf("%v\n%#v\n", e, q)
}
*/
