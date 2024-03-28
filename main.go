package main

import "cronJob/cmd"

func main() {
	cmd.Execute()
	//client := g.Client()
	//
	//client.ContentJson()
	//r, err := client.Post(gctx.GetInitCtx(), "http://httpbin.org/delay/1")
	//if err != nil {
	//	panic(err)
	//}
	//defer r.Close()
	//
	//fmt.Println(r.ReadAllString())
}
