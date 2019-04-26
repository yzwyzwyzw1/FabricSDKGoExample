package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)



func (f *FabricSetupService)Save(a string,an string,b string,bn string )(string,error) {

	var Args = [][]byte{[]byte(a),[]byte(an),[]byte(b),[]byte(bn)}
	req:= channel.Request{ChaincodeID:f.Fab.ChaincodeID,Fcn:"save",Args:Args}
	response,err := f.ChannelClient.Execute(req)//执行请求消息
	if err != nil {
		return "",fmt.Errorf("save failed:%v",err)
	}
	return string(response.TransactionID),nil
}


//Transaction makes payment of X units from A to B
func (f *FabricSetupService)Transfar(a string,b string,num string)(string,error) {


	var Args = [][]byte{[]byte(a),[]byte(b),[]byte(num)}
	req:= channel.Request{ChaincodeID:f.Fab.ChaincodeID,Fcn:"transfer",Args:Args}
	response,err := f.ChannelClient.Execute(req)//执行请求消息
	if err != nil {
		return "",fmt.Errorf("transfar failed:%v",err)
	}
	return string(response.TransactionID),nil
}

func (f *FabricSetupService)Query(name string)(string,error) {


	var Args = [][]byte{[]byte(name)}
	req:= channel.Request{ChaincodeID:f.Fab.ChaincodeID,Fcn:"query",Args:Args}

	response,err := f.ChannelClient.Query(req)//执行请求消息
	if err != nil {
		return "",fmt.Errorf("query failed:%v",err)
	}
	return string(response.Payload),nil
}
