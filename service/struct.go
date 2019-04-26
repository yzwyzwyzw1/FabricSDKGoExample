package service

import (
	"github.com/blockchaintest.com/FabricSDKGoExample/blockchain"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type FabricSetupService struct {
	ChannelClient   *channel.Client
	Fab             *blockchain.FabricSetup
}
