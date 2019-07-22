package restfulapi

type Stream struct {
	Key             string `json:"key"`
	Url             string `json:"Url"`
	StreamId        uint32 `json:"StreamId"`
	VideoTotalBytes uint64 `json:123456`
	VideoSpeed      uint64 `json:123456`
	AudioTotalBytes uint64 `json:123456`
	AudioSpeed      uint64 `json:123456`
}

type Streams struct {
	Publishers []Stream `json:"publishers"`
	Players    []Stream `json:"players"`
}
