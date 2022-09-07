package client

import "net/http"

type FluentBitMethods interface {
	GetMetricData() (*Response, error)
}

type FluentBitClient struct {
	FBHost     string
	FBPort     int
	HTTPClient http.Client
}
type ChunksStorage struct {
	TotalChunks  int `json:"total_chunks,omitempty"`
	MemChunks    int `json:"mem_chunks,omitempty"`
	FsChunks     int `json:"fs_chunks,omitempty"`
	FsChunksUp   int `json:"fs_chunks_up,omitempty"`
	FsChunksDown int `json:"fs_chunks_down,omitempty"`
}
type StorageLayer struct {
	Chunks ChunksStorage `json:"chunks,omitempty"`
}
type InputChunks struct {
	Containers GenericInputStruct `json:"containers,omitempty"`
	Systemd    GenericInputStruct `json:"systemd,omitempty"`
	Audit      GenericInputStruct `json:"audit,omitempty"`
}

type GenericInputStruct struct {
	Status Status      `json:"status,omitempty"`
	Chunks ChunksInput `json:"chunks,omitempty"`
}

type Status struct {
	Overlimit bool   `json:"overlimit,omitempty"`
	MemSize   string `json:"mem_size,omitempty"`
	MemLimit  string `json:"mem_limit,omitempty"`
}

type ChunksInput struct {
	Total    int    `json:"total,omitempty"`
	Up       int    `json:"up,omitempty"`
	Down     int    `json:"down,omitempty"`
	Busy     int    `json:"busy,omitempty"`
	BusySize string `json:"busy_size,omitempty"`
}

type Response struct {
	StorageLayer StorageLayer `json:"storage_layer,omitempty"`
	InputChunks  InputChunks  `json:"input_chunks,omitempty"`
}
