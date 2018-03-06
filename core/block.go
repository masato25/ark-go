package core

import (
	"net/http"
	"strconv"

	"github.com/kristjank/goark-node/base/model"
)

//GetFullBlocksFromPeer function returns a full list of blocks from current last block on. A random number of blocks is returned,
//due to ddos measures
func (s *ArkClient) GetFullBlocksFromPeer(lastBlockHeight int) (model.BlockResponse, ArkApiResponseError, *http.Response) {
	respData := new(model.BlockResponse)
	respError := new(ArkApiResponseError)

	resp, err := s.sling.New().Get("peer/blocks?lastBlockHeight="+strconv.Itoa(lastBlockHeight)).Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
		respError.ErrorObj = err
	}

	return *respData, *respError, resp
}

//GetPeerHeight function returns node peer height.
func (s *ArkClient) GetPeerHeight() (model.BlockHeightResponse, ArkApiResponseError, *http.Response) {
	respError := new(ArkApiResponseError)
	respData := new(model.BlockHeightResponse)

	resp, err := s.sling.New().Get("api/blocks/getHeight").Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
	}

	return *respData, *respError, resp
}

type BlocksResponse struct {
	Success bool          `json:"success"`
	Blocks  []model.Block `json:"blocks"`
	Count   int           `json:"count,omitempty"`
}

func (s *ArkClient) GetBlocksByGeneratorPublicKey(publicKey string, limit ...int) (BlocksResponse, ArkApiResponseError, *http.Response) {
	respData := new(BlocksResponse)
	respError := new(ArkApiResponseError)
	qLimit := 10
	if len(limit) != 0 {
		qLimit = limit[0]
	}

	resp, err := s.sling.New().Get("api/blocks?limit="+strconv.Itoa(qLimit)+"&generatorPublicKey="+publicKey).Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
		respError.ErrorObj = err
	}

	return *respData, *respError, resp
}

func (s *ArkClient) GetBlocksByHeight(height int) (BlocksResponse, ArkApiResponseError, *http.Response) {
	respData := new(BlocksResponse)
	respError := new(ArkApiResponseError)
	resp, err := s.sling.New().Get("api/blocks?height="+strconv.Itoa(height)).Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
		respError.ErrorObj = err
	}

	return *respData, *respError, resp
}

//PostBlock to selected ARKNetwork
func (s *ArkClient) PostBlock(payload model.BlockReceiveStruct) (model.PostBlockResponse, ArkApiResponseError, *http.Response) {
	respTr := new(model.PostBlockResponse)
	errTr := new(ArkApiResponseError)

	/*var payload transactionPayload
	payload.Transactions = append(payload.Transactions, tx)
	*/
	resp, err := s.sling.New().Post("peer/blocks").BodyJSON(payload).Receive(respTr, errTr)

	if err != nil {
		errTr.ErrorMessage = err.Error()
	}

	return *respTr, *errTr, resp
}
