package service

import (
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"server/repository"
)

func (s *service) VisitorP2P(key string) (P2PConfig, error) {
	db, _, _ := repository.Get("")
	p2p, err := db.GostClientP2P.Preload(db.GostClientP2P.Node).Where(db.GostClientP2P.VKey.Eq(key)).First()
	if err != nil {
		return P2PConfig{}, err
	}

	auth, err := db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(p2p.Code)).First()
	if auth == nil {
		return P2PConfig{}, err
	}

	serverAddr, serverPort := p2p.Node.GetAddress()
	result := P2PConfig{
		Common: v1.ClientCommonConfig{
			Auth: v1.AuthClientConfig{
				Token: p2p.NodeCode,
			},
			ServerAddr: serverAddr,
			ServerPort: serverPort,
			Transport: v1.ClientTransportConfig{
				Protocol: p2p.Node.Protocol,
			},
			Metadatas: map[string]string{
				"user":     auth.User,
				"password": auth.Password,
			},
			LoginFailExit: new(bool),
		},
		STCP: v1.STCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: p2p.Code + "_visitorp2pstcp",
				Type: "stcp",
				Transport: v1.VisitorTransport{
					UseEncryption:  true,
					UseCompression: true,
				},
				SecretKey:  p2p.VKey,
				ServerName: p2p.Code + "_p2pstcp",
			},
		},
		XTCP: v1.XTCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: p2p.Code + "_visitorp2pxtcp",
				Type: "xtcp",
				Transport: v1.VisitorTransport{
					UseEncryption:  true,
					UseCompression: true,
				},
				SecretKey:  p2p.VKey,
				ServerName: p2p.Code + "_p2pxtcp",
			},
			KeepTunnelOpen:    true,
			MaxRetriesAnHour:  60,
			MinRetryInterval:  60,
			FallbackTo:        p2p.Code + "_visitorp2pstcp",
			FallbackTimeoutMs: 1500,
		},
	}

	// 判断是否需要中继
	if p2p.Forward != 1 || p2p.Node.P2PDisableForward == 1 {
		result.STCP = v1.STCPVisitorConfig{}
		result.XTCP.FallbackTo = ""
	}
	return result, nil
}
