package main

import (
	"fmt"
	"os"
	"sim"
	"time"
)

type SimInterfaceStruct struct {
}

var g_cmdclientip [5000]string
var g_imsi [5000]string

func (sim *SimInterfaceStruct) SimClientConnectNotify(clientid int, clientip string, clientport int) {
	fmt.Printf("SimClientConnectNotify clientid=%d clientip=%s clientport=%d\n", clientid, clientip, clientport)

}

func (sim *SimInterfaceStruct) SimClientDisConnectNotify(clientid int, clientip string, clientport int) {
	fmt.Printf("SimClientDisConnectNotify clientid=%d clientip=%s clientport=%d\n", clientid, clientip, clientport)

}

func (sim *SimInterfaceStruct) SimState(slot int, imsi string, state int, simtype int, clientip string) {
	fmt.Printf("SimState slot=%d imsi=%s state=%d simtype=%d clientip=%s\n", slot, imsi, state, simtype, clientip)
	g_imsi[slot] = imsi
	g_cmdclientip[slot] = clientip

}

func (sim *SimInterfaceStruct) AuthResult(slot int, imsi string, result int, sres string, kc string, ik string, clientip string) {
	fmt.Printf("AuthResult slot=%d imsi=%s result=%d sres=%s kc=%s ik=%s clientip=%s\n", slot, imsi, result, sres, kc, ik, clientip)

}

var randstr string = "b93011f0874d02e5a2a77035966"
var autnstr string = "8ec48bf37bc7ff16fbb23a833cd"
var count int = 0

func getrandautn() (ran string, autn string) {
	count++
	count = count % 10000
	ran = fmt.Sprintf("%s%05d", randstr, count)
	autn = fmt.Sprintf("%s%05d", autnstr, count)
	return
}

func main() {

	simctl := SimInterfaceStruct{}

	ret := sim.InitSimControl("210.32.1.229", 8888, 8889, &simctl)
	if ret != 0 {
		fmt.Println("sim.InitSimControl error ret=", ret)
		os.Exit(1)
	}

	//var imsi string = "460078571900404"
	//var imsi string = "460005562309046"
	//sim.Auth(7,imsi,randstr,"",g_cmdclientip)
	//var imsi string = "460008832529783"

	var ran, autn string
	fmt.Println(ran, autn)

	for {
		time.Sleep(2 * time.Second)

		slot := 1
		ran, autn = getrandautn()
		sim.Auth(slot, g_imsi[slot], ran, autn, g_cmdclientip[slot])
		/*
			slot = 2
			ran,autn = getrandautn()
			sim.Auth(slot,g_imsi[slot],ran,autn,g_cmdclientip[slot])

			slot = 9
			ran,autn = getrandautn()
			sim.Auth(slot,g_imsi[slot],ran,autn,g_cmdclientip[slot])
		*/
	}
}
