# 简介

这是一个Fabric-sdk-go的测试程序例子，重点在configtx和crypto-config文件的书写方法。

# 测试

执行下面的命令就可以直接运行程序
```
make 
```
测试结果如下：
```
save...
c67e86a6a2272748da63c3e567a555ce7b82981ec2b8e024a75a8904d7d3eaa6
query...
1000
transfer...
bf14917453e7b9c6b369096855ccb60c62025b1f49a34328b0619f24352e345d
query...
900
query...
300
```

对于其他的命令的作用，可以查看Makefile文件中的内容获知

# 其他

1.重新生成组织机构关系和密钥证书文件
```
cd fixtures
cryptogen generate  --config ./crypto-config.yaml  --output  crypto-config
export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mychannel
configtxgen  -profile OneOrgOrdererGenesis  -outputBlock     ./channel-artifacts/genesis.block
configtxgen  -profile OneOrgChannel -outputCreateChannelTx   ./channel-artifacts/mychannel.tx      -channelID $CHANNEL_NAME
configtxgen  -profile OneOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg exampleOrg
```
切记：
每次重建以上文件都需要修改 docker-compose-base.yaml 中的ca密钥，再重新启动网络。

2. cli测试
```
make env-clean
make env-up

docker exec -it cli bash
export CHANNEL_NAME=mychannel
peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/mychannel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join -b mychannel.block
peer chaincode install -n mycc -v 1.0 -p github.com/chaincode
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c'{"Args":["init"]}' -P "OR ('org1.example.com.member')"
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["save","a","1000","b","200"]}'
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}' 
```
你也可以通过脚本文件测试
```
make env-clean
make env-up

docker exec -it cli bash
./scripts/startup.sh

```

peer channel list //查看通道
peer chaincode list --instantiated -C mychannel  //查看通道中被安装的链码

# 错误处理



凡是和permission denied有关的，都是和证书与私钥相关，以及节点权限相关
我曾经修改过configtx.yaml文件中的OrdererOrg中的&OrdererOrgPolicies导致permission denied失败。


如果将docker-compose.yaml中的- CORE_PEER_ADDRESSAUTODETECT=true注释掉将会报如下错误：
gRPC Transport Status Code: (4) DeadlineExceeded. Description: context deadline exceeded


Instantiate chaincode failed: sending deploy transaction proposal failed: Transaction processing for endorser [localhost:7151]: Chaincode status Code: (500) UNKNOWN. Description: error starting container: error starting container: API error (404): network fixtures not found
注释掉：CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fixtures


执行命令： configtxgen  -profile OneOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg exampleOrg时报错：
Error on inspectChannelCreateTx: No organization name matching : exampleOrg
我最后发现，我的组织名exampleOrg后面有空格，需要删掉这个空格