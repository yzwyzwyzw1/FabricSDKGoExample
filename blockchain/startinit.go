package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	//"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt/opts"
)

const ChaincodeVersion  = "1.0"

//该结构体中定义了Fabric的启动参数
type FabricSetup struct {
	ConfigFile        string            // SDK配置文件所在路径
	Instantiated      bool		        // 是否初始化
	OrgAdmin          string
	OrgName           string

	ChannelID          string
	ChannelConfigPath  string
	OrdererOrgName     string

	OrgResMgmt         *resmgmt.Client

	ChaincodePath      string
	ChaincodeGoPath    string
	ChaincodeID        string

	UserName           string

}

//实例化SDK对象
func InstantiateSdk(configfile string,instantiated bool) (*fabsdk.FabricSDK,error) {
	fmt.Println("Instantiating SDK...")

	if instantiated {
		return nil,fmt.Errorf("SDK has been initallized")
	}

	// init sdk config file
	sdk,err := fabsdk.New(config.FromFile(configfile))
	if err != nil {
		fmt.Printf("Instantiate Fabric SDK failed: %v", err)
	}
	fmt.Println("Instantiate Fabric SDK successed!")
	return sdk,nil
}

// create channel and join peers
func CreateChannel(sdk *fabsdk.FabricSDK,fab *FabricSetup) error {
	fmt.Println("creating channel...")
   clientContext := sdk.Context(fabsdk.WithUser(fab.OrgAdmin), fabsdk.WithOrg(fab.OrgName))   //ORgAdmin   OrgName
   if clientContext == nil {
   	return fmt.Errorf("create client context failed!")
   }
	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("For resMgmtClient context --->create resMgmtClient failed: %v", err)
	}

	//  New creates a new Client instance
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(fab.OrgName))
	if err != nil {
		return fmt.Errorf("For OrgName ---> Create Org MSP client instance failed!: %v", err)
	}
	//  obtain signing identity
	adminIdentity, err := mspClient.GetSigningIdentity(fab.OrgAdmin)
	if err != nil {
		return fmt.Errorf(" obtain signing identity failed: %v", err)//获取指定id的签名标识失败
	}

	//SaveChannelRequest holds parameters for save channel request
	channelReq := resmgmt.SaveChannelRequest{ChannelID:fab.ChannelID,ChannelConfigPath:fab.ChannelConfigPath,SigningIdentities:[]msp.SigningIdentity{adminIdentity}}//ChannelID  ChannelConfigpath
	// save channel response with transaction ID
	_,err = resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(fab.OrdererOrgName))  //OrdererOrgName
    if err != nil {
    	return fmt.Errorf("create channel failed:%v",err)
	}
    fmt.Println("create channel succeed!")


    fab.OrgResMgmt = resMgmtClient
	// allows for peers to join existing channel with optional custom options (specific peers, filtered peers). If peer(s) are not specified in options it will default to all peers that belong to client's MSP.
	err = fab.OrgResMgmt.JoinChannel(fab.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(fab.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("Peers join channel failed: %v", err)
	}

	fmt.Println("Peers join channel succeed!")

	return nil

   }

func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK,fab *FabricSetup) (*channel.Client,error) {
	fmt.Println("InstallAndInstantiateCC...")

	//create chaincode package
	ccPkg, err := gopackager.NewCCPackage(fab.ChaincodePath, fab.ChaincodeGoPath) //ChaincodePath  ChaincodeGoPath
	if err != nil {
		return nil, fmt.Errorf("create chaincode package failed!: %v", err)
	}


	// contains install chaincode request parameters
	installCCReq := resmgmt.InstallCCRequest{Name: fab.ChaincodeID, Path: fab.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg} //ChaincodeID
	// allows administrators to install chaincode onto the filesystem of a peer
	_, err = fab.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("install chaincode failed: %v", err)
	}
	fmt.Println("install chaincode succeed!")


	fmt.Println("Instantiate chaincode...")
	//  returns a policy that requires one valid
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.example.com"})

	//save instantiate chaincode request information
	var Args=[][]byte{[]byte("init")}
	instantiateCCReq := resmgmt.InstantiateCCRequest{Name: fab.ChaincodeID, Path: fab.ChaincodePath, Version: ChaincodeVersion, Args:Args, Policy: ccPolicy} //ChaincodeID ChaincodePath

	// instantiates chaincode with optional custom options (specific peers, filtered peers, timeout). If peer(s) are not specified
	_, err = fab.OrgResMgmt.InstantiateCC(fab.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//_, err = fab.OrgResMgmt.InstantiateCC(fab.ChannelID, instantiateCCReq)

	if err != nil {
		return nil, fmt.Errorf("Instantiate chaincode failed: %v", err)
	}

	fmt.Println("Instantiate chaincode succeed!")



	clientChannelContext := sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.UserName), fabsdk.WithOrg(fab.OrgName))  //UserName   OrgName
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("create application channel client failed: %v", err)
	}



	fmt.Println("create application channel client succeed!")//通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务

	return channelClient, nil


}
