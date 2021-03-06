package bios

import "math"

func getPeerIndexesToMeshWith(total, myPos int) map[int]bool {
	list := map[int]bool{}
	firstNeighbour := (myPos + 1) % total

	list[firstNeighbour] = true

	for i := 0; i < numConnectionsRequired(total); i++ {
		nextNeighbour := int(math.Pow(2.0, float64(i+2))) + myPos
		if nextNeighbour >= total {
			nextNeighbour = (nextNeighbour) % total
		}
		list[nextNeighbour] = true
	}
	return list
}

func numConnectionsRequired(numberOfNodes int) int {
	return int(math.Ceil(math.Sqrt(float64(numberOfNodes))))
}

func (b *BIOS) computeMyMeshP2PAddresses() []string {
	otherPeers := []string{}
	otherPeersMap := map[string]bool{}
	if len(b.MyPeers) > 0 {
		myFirstPeer := b.MyPeers[0]
		myPosition := -1
		for idx, peer := range b.ShuffledProducers {
			if myFirstPeer.AccountName() == peer.AccountName() {
				myPosition = idx
			}
		}
		if myPosition != -1 {
			peerIDs := getPeerIndexesToMeshWith(len(b.ShuffledProducers), myPosition)
			for idx, peer := range b.ShuffledProducers {
				p2pAddr := peer.Discovery.EOSIOP2P
				if peerIDs[idx] && !otherPeersMap[p2pAddr] {
					otherPeers = append(otherPeers, p2pAddr)
					otherPeersMap[p2pAddr] = true
				}
			}
		}
	}
	return otherPeers
}
