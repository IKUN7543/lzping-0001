package types

type GetStockReq struct {
	ProductId int64 `form:"productId"`
}

type GetStockResp struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"productId"`
	Total     int32 `json:"total"`
	Available int32 `json:"available"`
	LockStock int32 `json:"lockStock"`
	Sales     int32 `json:"sales"`
}

type CreateStockReq struct {
	ProductId int64 `json:"productId"`
	Num       int32 `json:"num"`
}

type CreateStockResp struct {
	ProductId int64 `json:"productId"`
	Total     int32 `json:"total"`
	Available int32 `json:"available"`
}

type DeductStockReq struct {
	ProductId int64  `json:"productId"`
	Num       int32  `json:"num"`
	OrderNo   string `json:"orderNo"`
}

type DeductStockResp struct {
	Success bool `json:"success"`
}
