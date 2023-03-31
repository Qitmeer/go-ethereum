package params

import "math/big"

type MeerChainConfig struct {
	ChainID *big.Int // chainId identifies the current chain and is used for replay protection
}

var (
	QngMainnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(813),
	}
	QngTestnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8131),
	}
	QngMixnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8132),
	}
	QngPrivnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8133),
	}

	AmanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8134),
	}
	AmanaTestnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81341),
	}
	AmanaMixnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81342),
	}
	AmanaPrivnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81343),
	}

	FlanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8135),
	}
	FlanaTestnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81351),
	}
	FlanaMixnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81352),
	}
	FlanaPrivnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81353),
	}

	MizanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8136),
	}
	MizanaTestnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81361),
	}
	MizanaMixnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81362),
	}
	MizanaPrivnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(81363),
	}
)

func init() {
	NetworkNames[QngMainnetChainConfig.ChainID.String()] = "qng"
	NetworkNames[QngTestnetChainConfig.ChainID.String()] = "qng-test"
	NetworkNames[QngMixnetChainConfig.ChainID.String()] = "qng-mix"
	NetworkNames[QngPrivnetChainConfig.ChainID.String()] = "qng-priv"

	NetworkNames[AmanaChainConfig.ChainID.String()] = "amana"
	NetworkNames[AmanaTestnetChainConfig.ChainID.String()] = "amana-test"
	NetworkNames[AmanaMixnetChainConfig.ChainID.String()] = "amana-mix"
	NetworkNames[AmanaPrivnetChainConfig.ChainID.String()] = "amana-priv"

	NetworkNames[FlanaChainConfig.ChainID.String()] = "flana"
	NetworkNames[FlanaTestnetChainConfig.ChainID.String()] = "flana-test"
	NetworkNames[FlanaMixnetChainConfig.ChainID.String()] = "flana-mix"
	NetworkNames[FlanaPrivnetChainConfig.ChainID.String()] = "flana-priv"

	NetworkNames[MizanaChainConfig.ChainID.String()] = "mizana"
	NetworkNames[MizanaTestnetChainConfig.ChainID.String()] = "mizana-test"
	NetworkNames[MizanaMixnetChainConfig.ChainID.String()] = "mizana-mix"
	NetworkNames[MizanaPrivnetChainConfig.ChainID.String()] = "mizana-priv"
}

func IsQngNetwork(chainID *big.Int) bool {
	if chainID == QngMainnetChainConfig.ChainID ||
		chainID == QngTestnetChainConfig.ChainID ||
		chainID == QngMixnetChainConfig.ChainID ||
		chainID == QngPrivnetChainConfig.ChainID {
		return true
	}
	return false
}

func IsAmanaNetwork(chainID *big.Int) bool {
	if chainID == AmanaChainConfig.ChainID ||
		chainID == AmanaTestnetChainConfig.ChainID ||
		chainID == AmanaMixnetChainConfig.ChainID ||
		chainID == AmanaPrivnetChainConfig.ChainID {
		return true
	}
	return false
}

func IsFlanaNetwork(chainID *big.Int) bool {
	if chainID == FlanaChainConfig.ChainID ||
		chainID == FlanaTestnetChainConfig.ChainID ||
		chainID == FlanaMixnetChainConfig.ChainID ||
		chainID == FlanaPrivnetChainConfig.ChainID {
		return true
	}
	return false
}

func IsMizanaNetwork(chainID *big.Int) bool {
	if chainID == MizanaChainConfig.ChainID ||
		chainID == MizanaTestnetChainConfig.ChainID ||
		chainID == MizanaMixnetChainConfig.ChainID ||
		chainID == MizanaPrivnetChainConfig.ChainID {
		return true
	}
	return false
}
