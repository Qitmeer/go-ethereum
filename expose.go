// Copyright (c) 2017-2022 The qitmeer developers

package ethereum

import (
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/internal/web3ext"
	"github.com/urfave/cli/v2"
	"math/big"
)

var Modules = web3ext.Modules

func GlobalBig(ctx *cli.Context, name string) *big.Int {
	return flags.GlobalBig(ctx, name)
}
