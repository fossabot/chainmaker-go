module chainmaker.org/chainmaker-go/net

go 1.15

require (
	chainmaker.org/chainmaker/chainmaker-net-common v0.0.6-0.20210929043521-02e40bf96300
	chainmaker.org/chainmaker/chainmaker-net-libp2p v0.0.11-0.20210929043636-4e46c072735d
	chainmaker.org/chainmaker/chainmaker-net-liquid v0.0.8-0.20210929043651-03cf1a4650f8
	chainmaker.org/chainmaker/common/v2 v2.0.1-0.20210928092334-f8be4fb05660
	chainmaker.org/chainmaker/logger/v2 v2.0.0-20210907134457-53647922a89d
	chainmaker.org/chainmaker/pb-go/v2 v2.0.1-0.20210901132412-435b75070bf2
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20210928092254-cfa32191bfce
	github.com/gogo/protobuf v1.3.2
	github.com/stretchr/testify v1.7.0
)

replace github.com/libp2p/go-libp2p-core => chainmaker.org/chainmaker/libp2p-core v0.0.2
