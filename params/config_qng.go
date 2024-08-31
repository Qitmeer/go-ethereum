package params

import (
	"fmt"
	"math/big"
)

const (
	Meer   = "meer"
	Amana  = "amana"
	Flana  = "flana"
	Mizana = "mizana"
)

type MeerChainConfig struct {
	ChainID   *big.Int // chainId identifies the current chain and is used for replay protection
	Name      string
	Type      string
	Consensus string // TODO:reserve
}

func (mcc MeerChainConfig) String() string {
	return fmt.Sprintf("Name:%s,ID:%s,Type:%s", mcc.Name, mcc.ChainID.String(), mcc.Type)
}

var (
	QngMainnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(813),
		Name:    "qng",
		Type:    Meer,
	}
	QngTestnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8131),
		Name:    "qng-test",
		Type:    Meer,
	}
	QngMixnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8132),
		Name:    "qng-mix",
		Type:    Meer,
	}
	QngPrivnetChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8133),
		Name:    "qng-priv",
		Type:    Meer,
	}
	AmanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8134),
		Name:    "amana",
		Type:    Amana,
	}
	FlanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8135),
		Name:    "flana",
		Type:    Flana,
	}
	MizanaChainConfig = &MeerChainConfig{
		ChainID: big.NewInt(8136),
		Name:    "mizana",
		Type:    Mizana,
	}

	meerChainConfigs = []*MeerChainConfig{
		QngMainnetChainConfig,
		QngTestnetChainConfig,
		QngMixnetChainConfig,
		QngPrivnetChainConfig,
		AmanaChainConfig,
		FlanaChainConfig,
		MizanaChainConfig,
	}
)

func init() {
	for _, c := range meerChainConfigs {
		NetworkNames[c.ChainID.String()] = c.Name
	}
}

func IsQngNetwork(chainID *big.Int) bool {
	return IsNetwork(chainID, Meer)
}

func IsAmanaNetwork(chainID *big.Int) bool {
	return IsNetwork(chainID, Amana)
}

func IsFlanaNetwork(chainID *big.Int) bool {
	return IsNetwork(chainID, Flana)
}

func IsMizanaNetwork(chainID *big.Int) bool {
	return IsNetwork(chainID, Mizana)
}

func QngEIPsBanner(banner string, c *ChainConfig) string {
	banner += "\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Pre-Merge hard forks (block based):\n"
	banner += fmt.Sprintf(" - Homestead:                   #%-8v (EIP-606 Requires:EIP-2, EIP-7, EIP-8)\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    #%-8v (EIP-779 Requires:EIP-606)\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): #%-8v (EIP-608 Requires:EIP-150, EIP-779)\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): #%-8v (EIP-607 Requires:EIP-155, EIP-160, EIP-161, EIP-170, EIP-608)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): #%-8v (EIP-607 Requires:EIP-155, EIP-160, EIP-161, EIP-170, EIP-608)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Byzantium:                   #%-8v (EIP-609 Requires:EIP-100, EIP-140, EIP-196, EIP-197, EIP-198, EIP-211, EIP-214, EIP-607, EIP-649, EIP-658)\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              #%-8v (EIP-1013 Requires:EIP-145, EIP-609, EIP-1014, EIP-1052, EIP-1234, EIP-1283)\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  #%-8v (EIP-1716 Requires:EIP-1013, EIP-1283)\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    #%-8v (EIP-1679 Requires:EIP-152, EIP-1108, EIP-1344, EIP-1716, EIP-1884, EIP-2028, EIP-2200)\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                #%-8v (EIP-2387 Requires:EIP-1679, EIP-2384)\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Berlin:                      #%-8v (EIP-2565 EIP-2929 EIP-2718 EIP-2930 Requires:EIP-198)\n", c.BerlinBlock)
	banner += fmt.Sprintf(" - London:                      #%-8v (EIP-1559 EIP-3198 EIP-3529 EIP-3541 EIP-3554 Requires:EIP-2718, EIP-2930, EIP-2200, EIP-2929, EIP-2930)\n", c.LondonBlock)
	if c.ArrowGlacierBlock != nil {
		banner += fmt.Sprintf(" - Arrow Glacier:               #%-8v (EIP-4345)\n", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		banner += fmt.Sprintf(" - Gray Glacier:                #%-8v (EIP-5133)\n", c.GrayGlacierBlock)
	}
	banner += "\n"

	// Add a special section for the merge as it's non-obvious
	if c.TerminalTotalDifficulty == nil {
		banner += "The Merge is not yet available for this network!\n"
		banner += " - Hard-fork specification: EIP-3675 EIP-4399 Requires:EIP-2124\n"
	} else {
		banner += "Merge configured:\n"
		banner += " - Hard-fork specification: EIP-3675 EIP-4399 Requires:EIP-2124\n"
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
		banner += fmt.Sprintf(" - Shanghai:                    @%-10v (EIP-3651 EIP-3855 EIP-3860 EIP-4895 EIP-6049 Requires:EIP-2929, EIP-170)\n", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		banner += fmt.Sprintf(" - Cancun:                      @%-10v\n", *c.CancunTime)
	}
	if c.PragueTime != nil {
		banner += fmt.Sprintf(" - Prague:                      @%-10v\n", *c.PragueTime)
	}

	return banner
}

func AddMeerChainConfig(cfg *MeerChainConfig) error {
	for _, c := range meerChainConfigs {
		if c.ChainID.Cmp(cfg.ChainID) == 0 {
			return fmt.Errorf("It already exists:%s, but now add %s", c.String(), cfg.String())
		}
	}
	meerChainConfigs = append(meerChainConfigs, cfg)
	NetworkNames[cfg.ChainID.String()] = cfg.Name
	return nil
}

func IsNetwork(chainID *big.Int, mtype string) bool {
	for _, c := range meerChainConfigs {
		if chainID.Cmp(c.ChainID) == 0 && c.Type == mtype {
			return true
		}
	}
	return false
}
