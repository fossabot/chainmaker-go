module chainmaker.org/chainmaker-go/txpool

go 1.15

require (
	chainmaker.org/chainmaker-go/monitor v0.0.0
	chainmaker.org/chainmaker/chainconf/v2 v2.0.0-20210913144615-f27c44059848
	chainmaker.org/chainmaker/common/v2 v2.0.1-0.20210916080817-79b5a4160dae
	chainmaker.org/chainmaker/localconf/v2 v2.0.0-20210913154622-9f9774ed7d1b
	chainmaker.org/chainmaker/logger/v2 v2.0.0-20210907133316-af00cea33c97
	chainmaker.org/chainmaker/pb-go/v2 v2.0.1-0.20210916064951-47123db73430
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20210914063622-6f007edc3a98
	chainmaker.org/chainmaker/utils/v2 v2.0.0-20210907033606-84c6c841cbdb
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
)

replace chainmaker.org/chainmaker-go/monitor => ../monitor
