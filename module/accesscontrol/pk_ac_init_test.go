package accesscontrol

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/common/v2/helper"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/stretchr/testify/require"
)

const (
	TestPK1 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCwu3sB7+5LJMtqsP2Y3sn5eW8b
r/dv1d9WXa6p0UEsUoVY4bgDZ471X9e0htnZUWcuvI5B0wHkoaJiKhUxSk5AJ8OY
5IvFI0OqQS07IMqIj3/u3iERVluuawA5IUjPFmCiubJ/Pb/JZCpFDQbZb209h7Vs
OkD1v94WlNJN8sm7qQIDAQAB
-----END PUBLIC KEY-----`
	TestSK1 = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCwu3sB7+5LJMtqsP2Y3sn5eW8br/dv1d9WXa6p0UEsUoVY4bgD
Z471X9e0htnZUWcuvI5B0wHkoaJiKhUxSk5AJ8OY5IvFI0OqQS07IMqIj3/u3iER
VluuawA5IUjPFmCiubJ/Pb/JZCpFDQbZb209h7VsOkD1v94WlNJN8sm7qQIDAQAB
AoGAXKZcjR5wSTqH3W3d9KdPMRcFNXmheSKhC9DfAS2vQgIc4AStCDPhESfmmEBd
snznX+v/k+h/xJEr5NR0+bsfm8hbNFBGDysEZt9tzD5YItafs/ePl+PDsAZVv8d/
q3CDwisgZ3oJzPIw+CqRSm2WeL3umQy1RdQXdJ2RsdnRD+UCQQDrB3U0Wefk+Vid
unYpXd+tAXN2cB7GxLsKYAGNz+1AyFLDnbG/K35tihdy0kaa1L/4K9h39xfDnCeK
/c3AfhSHAkEAwIBn4qLuuSQaRwqTGii0IV7Y8Asj/yyUTN/dOo3SB2804nFyNPHX
qfnhUhjRZhxDFLPH4KwKMMOMkilyz3LqTwJAS0yvY19uqXCt0JL96pD16dLuMEMJ
yTschd1uggXdCIVl5uBuI0aHEgdNLe9qyY5iFtvNVdonlfdAwApC0mpSnwJABmT9
jnC9H1dMrCl0w3ywpx8gc7DbDEHt1zPkhGprnKWcCx2bnpieAl5zlqeOZSbxL4Hd
VOBCImaMh9pqnuuBTwJATglQw9U1JbQ10Xq85MdJAVwbd1DgXmM1NAF4V+aeVpFf
V2+2cdpXmxG5Hrp9dy/oR8AGLPQ/JFeUYDxjEU5GNw==
-----END RSA PRIVATE KEY-----`
	TestNodeId1 = "QmQXjPB4DS8fNxsbqWzozSfwRiBbDDZg3t5qTxeb7R8BV5"

	TestPK2 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcGCdkB7vruKPcGXIneql989b3
pGYfE7icZBHdNrCY6QqDqhFytlLL/zVTpDJVopOQRAonEyZtbhFFPoMntIwRR7Uw
oePfJmPw/woQ7oLJtVWW1a9am15ZIMro/dCzywi1Z6qlU+Vf1uOR7YzPF1OeIymv
0SXG2DMny9m1OJ40ZwIDAQAB
-----END PUBLIC KEY-----`
	TestSK2 = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDcGCdkB7vruKPcGXIneql989b3pGYfE7icZBHdNrCY6QqDqhFy
tlLL/zVTpDJVopOQRAonEyZtbhFFPoMntIwRR7UwoePfJmPw/woQ7oLJtVWW1a9a
m15ZIMro/dCzywi1Z6qlU+Vf1uOR7YzPF1OeIymv0SXG2DMny9m1OJ40ZwIDAQAB
AoGBAI2q3m/8qnEoABEEL/5Jbh+sfIoaP8FxKDtCDl2dfj5ugl4Ncf2sbc7xDpov
7lZAt0r9AKv2H54AYw13F2TPSfgD0Vw4YdAgHehG5HgzUO/pRkIy80ARywm7KCM9
RvMBXlbkVsBP/Pwc7KXxDVvlg9GFyDzzb3MgoqsZ56/FsvaBAkEA8p6hvyUTRXE8
YhLnmISXt+uT3D7mZuxgVJ0vBn0sIumCMXtYSHO8tWmlU0XysNy9H0LiTsGc+cSm
FFH1Hb8xWQJBAOg7gP11lMN9mjjHlrf1VviR/r5F9H1lU+MIsW/oxLrOHwwny9B5
FsNui9unDjG93R8ggt3tqeSDgX619YZkG78CQCDjiCGVMQuU0g6paWOvdbGk6aJN
lIYXPOe7dwh2J2mEJfX3Nnx70/TzoUmsjb2T7r8yHeN3M4RYN/tBMO0bYeECQFLm
ovZX2gIrPTmVriz/LMvROjHsQQneeSKrwMOlQU06NYUeU7iY8VJUjSKdMQj6sQvi
jDTzGVnUxA5aoEoYRHsCQG9bw6jejRyQ0CM3Xv1gurAy089LyDF04mn7GuHWFj6B
0NBQdGsuTBTxd0uXa/kMEwd3F782IEL2h1tDCBjhYAM=
-----END RSA PRIVATE KEY-----`
	TestNodeId2 = "QmRUuqP9WkNmHv2NR8P9RUyBKSdjHuz4uu79hjgM4rWri4"

	TestPK3 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4YX5dJBVlHNh6LKAd2dTeW9yx
NWbAPGzgiCF8tOmRuJ+jfBt2VY8GkBxMwPwsEf8mp/4l1NXRp/jSJ69SRwU06rb4
ujl7NzkeBwzc7cFbFxoqoKeI0nmIJcn7tI2YCSX7HtoddnsyLZNCFFKDRB/Db4f1
5Z5rDZWMAMlxeffoGQIDAQAB
-----END PUBLIC KEY-----`
	TestSK3 = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC4YX5dJBVlHNh6LKAd2dTeW9yxNWbAPGzgiCF8tOmRuJ+jfBt2
VY8GkBxMwPwsEf8mp/4l1NXRp/jSJ69SRwU06rb4ujl7NzkeBwzc7cFbFxoqoKeI
0nmIJcn7tI2YCSX7HtoddnsyLZNCFFKDRB/Db4f15Z5rDZWMAMlxeffoGQIDAQAB
AoGAL8eb7lkGblBeTLK5v2KOhhy6APX8rX47HKhKPT3IdSmpvLzRhQXA7Yt0ufMc
pfL38rV/55/S1OS5VwRPq3uZ/l44084zvBMLVrZRpGd4Map4QLNfsXhgNWE+qLne
Bj4HoNSsyzMmzUccp6KKs3FaeTwtLA3cU4qO3xV7Piq0vrECQQDyd/rr42p5Fj3R
Imxc29dSAX+Qph9DV2hDMc4U49/OVsku8yufGQagDWjvAt9jLTu41QtwM7c20asK
iJ8pEn8tAkEAwqukx3RQhhH3cdRFKqQq1sd+IgNMhpdboZhf5vs4+Z3gGKroQVKj
NfHL1SDCtDlxppgFQnI7vzaoJKFzVg2AHQJAQ/xwVwQFLr6VxrYoPEFINq5E3oI1
8ePoUC7+4cyjTG/5KTj12j5iJS6dZacgi+Z7AHB8LJHTpYNUujdkqVeOYQJBAITf
3dhad0Ab8Vcb+Z4Scj8p6dlTgR95Ho1dUVB697fB4B1WQrObsVV31paCBwQ3FXEN
4MEq8cchioF+RhhdnK0CQQCxsKE9Q6Cc+u1LzzZ3mwIJnjSmnpXAsGbP40kCEqHH
Aw0Q6CM4gQkL4N0LgqqmyWglIkfCscVOTA50OGYEioxj
-----END RSA PRIVATE KEY-----`
	TestNodeId3 = "Qmd8o58EHnsfBbDikRra4XNsCmArXjXLSdZdkYwDcdsUvQ"

	TestPK4 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDjToONx0tUBD3lvXscgXkR81r8
c1doBh02g0tjseWeOehROM0hIVYHeY/6CoBtmdXSIlc0ITI/dm4wLw47egkd6gQO
ph2EeuFarzTcXJjztBZElDXEQpppj4E5vx7Qp21uc4uAv/x5o6ORy0tL90bo/n/z
bUzV0TbkOzV4pib2swIDAQAB
-----END PUBLIC KEY-----`
	TestSK4 = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDjToONx0tUBD3lvXscgXkR81r8c1doBh02g0tjseWeOehROM0h
IVYHeY/6CoBtmdXSIlc0ITI/dm4wLw47egkd6gQOph2EeuFarzTcXJjztBZElDXE
Qpppj4E5vx7Qp21uc4uAv/x5o6ORy0tL90bo/n/zbUzV0TbkOzV4pib2swIDAQAB
AoGBAN2eWkMsUTRsIlFRSawETB+FTmuepVTFyUux/RoJg5+eQ/SU1eL8Vp1ZF1gp
TwgNGd0UIEOyLgSUGmCeMFkq5aDPKxwI356YdvHB7SdEJMy5I7fYEKJatEQ3FOtM
F5Il1NYPgtrwKeAAs9mSdH2Lp7/KH+UAmhH5wSPZ3MRbk5oRAkEA+iaGYvFTn1pP
PgEbun53ZNGYirVK8qqYe8lDTLYPtM9Q/HM5OP1Fojd06e+6mBtFX9A8Jm8rnmJZ
l4CgF0K86QJBAOifPTzbGC51daI0k7hBPZCdm3KBQu91gx9G07wpEPV5a6RCho18
8ugO3qGCvCHKXmrAh81sEVVPLCjPZ17h5TsCQH2uI3DMrOXwOsX9SpAtgBEQWWK/
aVN4sLnoyb5d7pA6ZQchYQun/HdfA4eRoZ9QfE+CUOZCjpi58ydyQXzOVBkCQB1j
wQDnTW7ROEN+EQu+cmDLCNC2tBY86owRDr8/EP1ykb73CLjniGj5N/d/5PT/9F3Y
ZU/2z1nP3uxpB85dC/ECQCHuQLatz4sIVuHYQNl2nRiu+1eNfgpB6wGtIh3O782/
Wk0+loBRjk4pBj9djnlephMX9gHptNqyYvJtM3LMJm0=
-----END RSA PRIVATE KEY-----`
	TestNodeId4 = "QmQ8nYaAMm5DdMzf3GaY2NPkmGneqmRyaSJDNRaFwuoxwV"

	TestPK5 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ2H7PWt3l8CF+zvQuRssnVAyY
xyB4tZ6KpXPpb9Y7fzK4mju8pBBFyErBSNM5uqgGYrnmf9Plh3MZvwSGd5ZVHZ3b
f9qx4cttzzAMsHcjra6AJ44qo5jxX0bVIiyErqEqGvYvAqUEb7Ye0kFKs/Z2tq2e
o79CwhOI0TLHOpby8wIDAQAB
-----END PUBLIC KEY-----`
	TestSK5 = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDJ2H7PWt3l8CF+zvQuRssnVAyYxyB4tZ6KpXPpb9Y7fzK4mju8
pBBFyErBSNM5uqgGYrnmf9Plh3MZvwSGd5ZVHZ3bf9qx4cttzzAMsHcjra6AJ44q
o5jxX0bVIiyErqEqGvYvAqUEb7Ye0kFKs/Z2tq2eo79CwhOI0TLHOpby8wIDAQAB
AoGAbLtbVIg2kO9Sm+UQVP194qm8P3DFZUExLq8CSfYdCd/zis5K78vRmEXVP1nj
r22Fpir4ydqCY1sb/fqQjX9OU4ZhwmvWAnHehjPTuzGgkwnKcUyDnWROX6B2wFOc
d9AaUxgZm4Dek7jXImfiKkb0YZxuMXiBZZmdqEhbxlM4QcECQQDwJH+egh6fKXMU
AbUmKuYDlH/D6m52FTozkQU9wWa5xW+/u+XJOVuX8gaaa9S15sX2gDqpB7v0HkO2
BnAqX2qlAkEA1yya9Hx/A/l1Z/rk0JRBiLSxYvdBmjLg43+9mhKCoSI+1adRyri4
26HJcnMAZa15hNWFHN0/QX3BNG+zIhgrtwJAGwJP5DkITqhvzAFBKZDLm/14vUVB
tUA/8orOBxsYfa5qGit89bvgxF8xRO751pelDktvzZEUH6nDvdZNiUaADQJAZsxq
o08vJ2jwjGKzGmsZ/APHk25pKxAPnOCUZp1dRzojJtOvIdiqiFN8+G60y97a5XlV
BPs2k0VPHowW2r0NdQJBAIVCNla0HUB1+1qtsI1qJHqzbI5eXTjEmKJ3FHxYvC9U
E0pyQvClu6PjjyAWRNMSGms4slKFgGffMJf4o1nkB1k=
-----END RSA PRIVATE KEY-----`
	TestNodeId5 = "QmXVYiWKFyQoSSoD6xptx5gCEmQTusnJ9veepeSvkXLTUF"

	TestPK6 = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDEyJWuNjX1BNfsbPtmgzcbfFzb
tHTbOH6I9svBaOAL8B1VyPK+RGIHtZgU5D4EOyWmFTwuXIQ98XGAkQHCJm75mKcG
Th/a3Q0HAf9guR8q5MT+mp0yckHuieMKQD30RafxoZi9H3LBKyGR2KMkNLjdhL3o
Vj+jFNeHf0gOkiQifwIDAQAB
-----END PUBLIC KEY-----`
	TestSK6 = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDEyJWuNjX1BNfsbPtmgzcbfFzbtHTbOH6I9svBaOAL8B1VyPK+
RGIHtZgU5D4EOyWmFTwuXIQ98XGAkQHCJm75mKcGTh/a3Q0HAf9guR8q5MT+mp0y
ckHuieMKQD30RafxoZi9H3LBKyGR2KMkNLjdhL3oVj+jFNeHf0gOkiQifwIDAQAB
AoGARc4hyrrQSSp+rg+63pKNaeKjzgwlp95ShKOHhAR/9bwnq9asxXHclH+Gg2Kz
3SxeHpxJzOhkwNR1PvYxeX3Iv4HKtdrruY9w15i6xNCBClIj/mK4A5IygSnbz6dC
S6gs6KBVpiLmY2IP4Uh28NcKW4eFe233xoKh6S5WHLUJDIECQQD2P9bmSfOJ5GTn
mrjBF0In8MKqq+Nu5miJ14gO3XrTMK8YkoE1V6qufvYBVD5YVHZySLam7pMXAxIO
zSy1fPm1AkEAzJNS7vQokTz+Ud25SYqFVT/g82VYxGRLuOYNWNJCXIY6G86p7/VQ
iRvJ6QM2DF41MMA3ox/bX4VllG4TP2474wJBANCGZOukSehGETCTM8rHcE00MxSl
9DUwRewcKOo1oVH/kvai8WmDcFTNzHJ5rUXNWHQUoR+hPcup3PvNwQN67lUCQAod
HnSBzZ+gjFIvzAE+v+i/B7gAwqqy6qtxdCd3/Z/lYuoNBYm/bwPYQ9spNXrXDXoj
hpyh7o6CYcs8xebU5FECQQCH0ZcTXrcX452gG+WUUL+0M0uDy4tF7KOkO89Aisqm
9Pr4zSK60hTtx3dAuD83HhezMjBxRqLHjLk0QcuNgHeq
-----END RSA PRIVATE KEY-----`
	TestNodeId6 = "QmboYJpeHZqvKPhmeUKQhydP2FKBzeqDARhZvRmbdGQt7p"

	TestPK = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQOJSqXyNB+Q62cXT+lx4TGNDU
Ast/pGwRFzPo+Ofef7lafQqu60gbkq4spQYqEgWyd7xr5tEw3tnQr3VEnSaQu2nS
gJDcT4yol0brUV0b2im9PNA45Q8QT+cZVILPLf3jJZtIxBFLps9q2Js65Xc8P314
UGClc2AZd8w7G7MLlwIDAQAB
-----END PUBLIC KEY-----`
)
const (
	testPermissionedKeyAuthType = "permissionedWithKey"
)

var testPKChainConfig = &config.ChainConfig{
	ChainId:  testChainId,
	Version:  testVersion,
	AuthType: testPermissionedKeyAuthType,
	Sequence: 0,
	Crypto: &config.CryptoConfig{
		Hash: testHashType,
	},
	Block: nil,
	Core:  nil,
	Consensus: &config.ConsensusConfig{
		Type: 0,
		Nodes: []*config.OrgConfig{{
			OrgId:  testOrg1,
			NodeId: nil,
		}, {
			OrgId:  testOrg2,
			NodeId: nil,
		}, {
			OrgId:  testOrg3,
			NodeId: nil,
		}, {
			OrgId:  testOrg4,
			NodeId: nil,
		},
		},
		ExtConfig: nil,
	},
	TrustRoots: []*config.TrustRootConfig{
		{
			OrgId: testOrg1,
			Root:  []string{testCAOrg1},
		},
		{
			OrgId: testOrg2,
			Root:  []string{testCAOrg2},
		},
		{
			OrgId: testOrg3,
			Root:  []string{testCAOrg3},
		},
		{
			OrgId: testOrg4,
			Root:  []string{testCAOrg4},
		},
	},
	TrustMembers: []*config.TrustMemberConfig{
		{
			OrgId:      testOrg5,
			Role:       "admin",
			MemberInfo: testTrustMember1,
		},
		{
			OrgId:      testOrg5,
			Role:       "admin",
			MemberInfo: testTrustMember2,
		},
	},
}

func TestGetNodeIdFromPK(t *testing.T) {
	var nodeId string
	pk, err := asym.PublicKeyFromPEM([]byte(TestPK))
	require.Nil(t, err)
	nodeId, err = helper.CreateLibp2pPeerIdWithPublicKey(pk)
	require.Nil(t, err)
	fmt.Println("nodeId:", nodeId)
}
