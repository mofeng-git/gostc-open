package service

import (
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"server/repository"
)

func (s *service) VisitorTunnel(key string) (TunnelConfig, error) {
	db, _, _ := repository.Get("")
	tunnel, err := db.GostClientTunnel.Preload(db.GostClientTunnel.Node).Where(db.GostClientTunnel.VKey.Eq(key)).First()
	if err != nil {
		return TunnelConfig{}, err
	}

	auth, err := db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(tunnel.Code)).First()
	if auth == nil {
		return TunnelConfig{}, err
	}

	serverAddr, serverPort := tunnel.Node.GetAddress()
	return TunnelConfig{
		Common: v1.ClientCommonConfig{
			Auth: v1.AuthClientConfig{
				Token: tunnel.NodeCode,
			},
			ServerAddr: serverAddr,
			ServerPort: serverPort,
			Transport: v1.ClientTransportConfig{
				Protocol: tunnel.Node.Protocol,
			},
			Metadatas: map[string]string{
				"user":     auth.User,
				"password": auth.Password,
			},
			LoginFailExit: new(bool),
		},
		STCP: v1.STCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: tunnel.Code + "_visitorstcp",
				Type: "stcp",
				Transport: v1.VisitorTransport{
					UseEncryption:  true,
					UseCompression: true,
				},
				SecretKey:  tunnel.VKey + "_stcp",
				ServerName: tunnel.Code + "_stcp",
			},
		},
		SUDP: v1.SUDPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: tunnel.Code + "_visitorsudp",
				Type: "sudp",
				Transport: v1.VisitorTransport{
					UseEncryption:  true,
					UseCompression: true,
				},
				SecretKey:  tunnel.VKey + "_sudp",
				ServerName: tunnel.Code + "_sudp",
			},
		},
	}, nil
}
