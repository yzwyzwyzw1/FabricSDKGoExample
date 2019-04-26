package main

import (
	"fmt"
	"github.com/blockchaintest.com/FabricSDKGoExample/blockchain"
	"github.com/blockchaintest.com/FabricSDKGoExample/service"
)


//CHANNELCONFIGPATH= os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/kongyixueyuan/fixtures/artifacts/channel.tx"
const (
  CONFIGFILE="config.yaml"
  ORGADMIN="Admin"
	//ORGADMIN="Org1MSP."
  ORGNAME="Org1"
  CHANNELID="mychannel"
  CHANNELCONFIGPATH= "/home/yzw/GoSpace/gopath/src/github.com/blockchaintest.com/FabricSDKGoExample/fixtures/channel-artifacts/mychannel.tx"
  ORDERERORGNAME="orderer.example.com"
  CHAINCODEPATH="github.com/blockchaintest.com/FabricSDKGoExample/chaincode"
  CHAINCODEGOPATH="/home/yzw/GoSpace/gopath"
  CHAINCODEID="mycc"
  USERNAME="User1"
)


func FabricSetupInit() *blockchain.FabricSetup{
	f := blockchain.FabricSetup{
		ConfigFile:CONFIGFILE,
		Instantiated:false,
		OrgAdmin:ORGADMIN,
		OrgName:ORGNAME,
		ChannelID:CHANNELID,
		ChannelConfigPath:CHANNELCONFIGPATH,
		ChaincodePath:CHAINCODEPATH,
		ChaincodeGoPath:CHAINCODEGOPATH,
		ChaincodeID:CHAINCODEID,
		UserName:USERNAME,
		OrdererOrgName:ORDERERORGNAME,

	}
	return &f
}


func main() {
	fab := FabricSetupInit()

    //Instantiate SDK  objection
	sdk,err := blockchain.InstantiateSdk(fab.ConfigFile,fab.Instantiated)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sdk.Close()

	err = blockchain.CreateChannel(sdk,fab)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := blockchain.InstallAndInstantiateCC(sdk, fab)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)


	//-------------------------------------//




  serviceSetup := service.FabricSetupService{ChannelClient:channelClient,Fab:fab}

	fmt.Println("save...")
	msg,err := serviceSetup.Save("a","1000","b","200")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

	fmt.Println("query...")
	msg,err = serviceSetup.Query("a")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

	fmt.Println("transfer...")
    msg,err = serviceSetup.Transfar("a","b","100")
    if err != nil {
    	fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

	fmt.Println("query...")
	msg,err = serviceSetup.Query("a")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

	fmt.Println("query...")
	msg,err = serviceSetup.Query("b")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

}
