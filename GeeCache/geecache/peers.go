package geecache

import pb "geecache/geecachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a spcefic key
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer
type PeerGetter interface {
	//Get(group string, key string) ([]byte, error)
	Get(in *pb.Request) (*pb.Response, error)
}
