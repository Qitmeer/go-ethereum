package params

import (
	"fmt"
	"math/big"
)

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

func QngEIPsBanner(banner string, c *ChainConfig) string {
	banner += "\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Pre-Merge hard forks (block based):\n"
	banner += fmt.Sprintf(" - Homestead:                   #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/homestead.md)\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/dao-fork.md)\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/tangerine-whistle.md)\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Byzantium:                   #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/byzantium.md)\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/constantinople.md)\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/petersburg.md)\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/istanbul.md)\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/muir-glacier.md)\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Berlin:                      #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/berlin.md)\n", c.BerlinBlock)
	banner += fmt.Sprintf(" - London:                      #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/london.md)\n", c.LondonBlock)
	if c.ArrowGlacierBlock != nil {
		banner += fmt.Sprintf(" - Arrow Glacier:               #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/arrow-glacier.md)\n", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		banner += fmt.Sprintf(" - Gray Glacier:                #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/gray-glacier.md)\n", c.GrayGlacierBlock)
	}
	banner += "\n"

	// Add a special section for the merge as it's non-obvious
	if c.TerminalTotalDifficulty == nil {
		banner += "The Merge is not yet available for this network!\n"
		banner += " - Hard-fork specification: https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md\n"
	} else {
		banner += "Merge configured:\n"
		banner += " - Hard-fork specification:    https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md\n"
		banner += fmt.Sprintf(" - Network known to be merged: %v\n", c.TerminalTotalDifficultyPassed)
		banner += fmt.Sprintf(" - Total terminal difficulty:  %v\n", c.TerminalTotalDifficulty)
		if c.MergeNetsplitBlock != nil {
			banner += fmt.Sprintf(" - Merge netsplit block:       #%-8v\n", c.MergeNetsplitBlock)
		}
	}
	banner += "\n"

	// Create a list of forks post-merge
	banner += "Post-Merge hard forks (timestamp based):\n"
	if c.ShanghaiTime != nil {
		banner += fmt.Sprintf(" - Shanghai:                    @%-10v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/shanghai.md)\n", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		banner += fmt.Sprintf(" - Cancun:                      @%-10v\n", *c.CancunTime)
	}
	if c.PragueTime != nil {
		banner += fmt.Sprintf(" - Prague:                      @%-10v\n", *c.PragueTime)
	}

	return banner
}
