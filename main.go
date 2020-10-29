package main

import (
	"fmt"
	"sim"
	"time"
)



type SimInterfaceStruct struct {
}


func (sim *SimInterfaceStruct) SimClientConnectNotify  (clientid int,clientip string,clientport int){
	fmt.Printf("SimClientConnectNotify clientid=%d clientip=%s clientport=%d\n",clientid,clientip,clientport)

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

	sim.InitSimControl("210.32.1.229", 8888, 8889,&simctl)


	for {
		var randstr string = "b93011f0874d02e5a2a7703596631dcb"
		var imsi string = "460028671303247"
		sim.Auth(1,imsi,randstr,"","210.32.1.229")
		time.Sleep(10000)
	}

}
