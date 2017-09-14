package main

import (
	"flag"
	"fmt"
)

func main()  {
	//把 -flagname 放到 ip 中
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	fmt.Print(*ip)

	//var flagVal int
	//
	//flag.Var(&flagVal, "name", "help message for flagname")
	//
	//fmt.Print(flagVal)

}


