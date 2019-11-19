package main

import (
	"errors"
	"flag"
	"fmt"
)

type config struct {
	help bool
	kind string
	size int
	rand float64
	seed int64
}

const (
	helpDefault = false

	kindDefault = kindAVL
	kindAVL     = "avl"
	kindRB      = "rb"

	sizeDefault = 10000000
	sizeMin     = 1
	sizeMax     = 99999999

	randDefault = randMax
	randMin     = 0
	randMax     = 1

	seedDefault = seedMin
	seedMin     = 1
)

var (
	errInvalidKind = errors.New("invalid kind")
	errInvalidSize = errors.New("invalid size")
	errInvalidRand = errors.New("invalid rand")
	errInvalidSeed = errors.New("invalid seed")
)

func parse() (*config, error) {
	var conf config

	flag.BoolVar(&conf.help, "h", helpDefault, fmt.Sprintf("show help messages (default %t)", helpDefault))
	flag.StringVar(&conf.kind, "k", kindDefault, fmt.Sprintf("kind of tree: %q | %q", kindAVL, kindRB))
	flag.IntVar(&conf.size, "n", sizeDefault, fmt.Sprintf("size of tree: [%v, %v]", sizeMin, sizeMax))
	flag.Float64Var(&conf.rand, "r", randDefault, fmt.Sprintf("randomness: [%v, %v]", randMin, randMax))
	flag.Int64Var(&conf.seed, "s", seedDefault, fmt.Sprintf("seed: should not be smaller than %v", seedMin))
	flag.Parse()

	if conf.kind != kindAVL && conf.kind != kindRB {
		return nil, errInvalidKind
	}

	if conf.size < sizeMin || sizeMax < conf.size {
		return nil, errInvalidSize
	}

	if conf.rand < randMin || randMax < conf.rand {
		return nil, errInvalidRand
	}

	if conf.seed < seedMin {
		return nil, errInvalidSeed
	}

	return &conf, nil
}
