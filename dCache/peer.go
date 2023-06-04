package dCache

import pb "github.com/DistributedCache/dCache/protobuf"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool) // 根据相应的key选择相应的节点
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error // 从对应的Groups中查找缓存值
}
