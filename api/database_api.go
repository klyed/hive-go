package api

import (
	"encoding/json"

	"github.com/klyed/hive-go/transports"
	_ "github.com/klyed/hive-go/types"
)

var APInameCondenser = "condenser_api"
var APInameDB = "database_api"

//GetConfig api request get_config
func (api *API) GetConfig() (*Config, error) {
	var resp Config
	err := api.call(APInameCondenser, "get_config", transports.EmptyParams, &resp)
	return &resp, err
}

//GetDynamicGlobalProperties api request get_dynamic_global_properties
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call(APInameDB, "get_dynamic_global_properties", transports.EmptyParams, &resp)
	return &resp, err
}

//GetBlock api request get_block
func (api *API) GetBlock(blockNum uint32) (*Block, error) {
	var resp Block
	err := api.call(APInameCondenser, "get_block", []uint32{blockNum}, &resp)
	resp.Number = blockNum
	return &resp, err
}

//GetBlockHeader api request get_block_header
func (api *API) GetBlockHeader(blockNum uint32) (*BlockHeader, error) {
	var resp BlockHeader
	err := api.call(APInameCondenser, "get_block_header", []uint32{blockNum}, &resp)
	resp.Number = blockNum
	return &resp, err
}

// Set callback to invoke as soon as a new block is applied
func (api *API) SetBlockAppliedCallback(notice func(header *BlockHeader, error error)) (err error) {
	err = api.setCallback(APInameCondenser, "set_block_applied_callback", func(raw json.RawMessage) {
		var header []BlockHeader
		if err := json.Unmarshal(raw, &header); err != nil {
			notice(nil, err)
		}
		for _, b := range header {
			notice(&b, nil)
		}
	})
	return
}
