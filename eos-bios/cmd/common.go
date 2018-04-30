package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	bios "github.com/eoscanada/eos-bios"
	shell "github.com/ipfs/go-ipfs-api"
)

func fetchNetwork(ipfs *bios.IPFS) (*bios.Network, error) {
	net := bios.NewNetwork(cachePath, myDiscoveryFile, ipfs)

	net.ForceFetch = !noDiscovery

	if err := net.TraverseGraph(); err != nil {
		return nil, fmt.Errorf("fetch-all error: %s", err)
	}

	if err := net.VerifyGraph(); err != nil {
		return nil, fmt.Errorf("graph inconsistent: %s", err)
	}

	if err := net.CalculateWeights(); err != nil {
		return nil, fmt.Errorf("error calculating weights: %s", err)
	}

	return net, nil
}

func ipfsClient() (*shell.IdOutput, *shell.Shell) {
	apiFile, err := ioutil.ReadFile(ipfsAPIFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading ipfs api file: %s\n", err)
		os.Exit(1)
	}

	ipfsClient := shell.NewShell(string(apiFile))

	fmt.Printf("Pinging ipfs node... ")
	info, err := ipfsClient.ID()
	if err != nil {
		fmt.Println("failed")
		fmt.Fprintf(os.Stderr, "error reaching ipfs api: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	return info, ipfsClient
}