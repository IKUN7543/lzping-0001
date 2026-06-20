package types

type OrderItemReq struct {
	ProductId int64 `json:"productId"`
	Num       int32 `json:"num"`
	Price     int64 `json:"price"`
}

type CreateOrderReq struct {
	Items           []OrderItemReq `json:"items"`
	ReceiverName    string         `json:"receiverName"`
	ReceiverPhone   string         `json:"receiverPhone"`
	ReceiverAddress string         `json:"receiverAddress"`
	Remark          string         `json:"remark"`
}

type CreateOrderResp struct {
	OrderNo   string `json:"orderNo"`
	Id        int64  `json:"id"`
	PayAmount int64  `json:"payAmount"`
}

type GetOrderReq struct {
	OrderNo string `form:"orderNo"`
}

type OrderItemResp struct {
	Id           int64  `json:"id"`
	ProductId    int64  `json:"productId"`
	ProductName  string `json:"productName"`
	ProductImage string `json:"productImage"`
	Price        int64  `json:"price"`
	Num          int32  `json:"num"`
	TotalPrice   int64  `json:"totalPrice"`
}

type OrderResp struct {
	Id              int64           `json:"id"`
	OrderNo         string          `json:"orderNo"`
	UserId          int64           `json:"userId"`
	TotalAmount     int64           `json:"totalAmount"`
	PayAmount       int64           `json:"payAmount"`
	FreightAmount   int64           `json:"freightAmount"`
	DiscountAmount  int64           `json:"discountAmount"`
	Status          int32           `json:"status"`
	PayType         int32           `json:"payType"`
	ReceiverName    string          `json:"receiverName"`
	ReceiverPhone   string          `json:"receiverPhone"`
	ReceiverAddress string          `json:"receiverAddress"`
	Remark          string          `json:"remark"`
	Items           []OrderItemResp `json:"items"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

type GetOrderResp struct {
	Order OrderResp `json:"order"`
}

type ListOrderReq struct {
	Page     int32 `form:"page"`
	PageSize int32 `form:"pageSize"`
}

type ListOrderResp struct {
	Orders []OrderResp `json:"orders"`
	Total  int64       `json:"total"`
}

type CancelOrderReq struct {
	OrderNo string `json:"orderNo"`
}

type CancelOrderResp struct {
	Success bool `json:"success"`
}

type PayOrderReq struct {
	OrderNo string `json:"orderNo"`
	PayType int32  `json:"payType"`
}

type PayOrderResp struct {
	Success bool `json:"success"`
}
