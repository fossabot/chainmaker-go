module chainmaker.org/chainmaker-go/net

go 1.15

require (
	chainmaker.org/chainmaker-go/localconf v0.0.0
	chainmaker.org/chainmaker-go/logger v0.0.0
	chainmaker.org/chainmaker-go/utils v0.0.0
	chainmaker.org/chainmaker/common v0.0.0-20210812042900-40fd24729b4a
	chainmaker.org/chainmaker/pb-go v0.0.0-20210825102713-0125b30c15d4
	chainmaker.org/chainmaker/protocol v0.0.0-20210817020238-7ad0d408ae23
	github.com/gogo/protobuf v1.3.2
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-circuit v0.3.1
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-discovery v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/libp2p/go-libp2p-pubsub v0.3.5
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/stretchr/testify v1.7.0
	github.com/tjfoc/gmsm v1.3.2
	github.com/tjfoc/gmtls v1.2.1
)

replace (
	chainmaker.org/chainmaker-go/localconf => ./../conf/localconf
	chainmaker.org/chainmaker-go/logger => ./../logger

	chainmaker.org/chainmaker-go/utils => ../utils
	github.com/libp2p/go-libp2p => ./p2p/libp2p
	github.com/libp2p/go-libp2p-core => ./p2p/libp2pcore
	github.com/libp2p/go-libp2p-pubsub => ./p2p/libp2ppubsub
)
