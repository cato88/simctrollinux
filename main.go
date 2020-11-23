package main

import (
	"fmt"
	"os"
	"sim"
	"time"
)



type SimInterfaceStruct struct {
}

var g_cmdclientip string

func (sim *SimInterfaceStruct) SimClientConnectNotify  (clientid int,clientip string,clientport int){
	fmt.Printf("SimClientConnectNotify clientid=%d clientip=%s clientport=%d\n",clientid,clientip,clientport)
	g_cmdclientip = clientip

}

func (sim *SimInterfaceStruct) SimClientDisConnectNotify (clientid int,clientip string,clientport int){
	fmt.Printf("SimClientDisConnectNotify clientid=%d clientip=%s clientport=%d\n",clientid,clientip,clientport)

}

func (sim *SimInterfaceStruct) SimState (slot int,imsi string,state int,simtype int,clientip string){
	fmt.Printf("SimState slot=%d imsi=%s state=%d simtype=%d clientip=%s\n",slot,imsi,state,simtype,clientip)

}

func (sim *SimInterfaceStruct) AuthResult(slot int,imsi string,result int,sres string,kc string,ik string,clientip string){
	fmt.Printf("AuthResult slot=%d imsi=%s result=%d sres=%s kc=%s ik=%s clientip=%s\n",slot,imsi,result,sres,kc,ik,clientip)

}


func main() {

	simctl := SimInterfaceStruct{}

	ret :=sim.InitSimControl("210.32.1.229", 8888, 8889,&simctl)
	if ret !=0 {
		fmt.Println("sim.InitSimControl error ret=",ret);
		os.Exit(1)
	}


	for {
		time.Sleep(30*time.Second)

		var randstr string = "b93011f0874d02e5a2a7703596631dcb"
		//var imsi string = "460078571900404"
		//var imsi string = "460005562309046"
		//sim.Auth(7,imsi,randstr,"",g_cmdclientip)
		var imsi string = "460008832529783"
		sim.Auth(10,imsi,randstr,"",g_cmdclientip)

	}
}
