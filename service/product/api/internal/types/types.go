package types

type CreateProductReq struct {
	CategoryId    int64  `json:"categoryId"`
	Name          string `json:"name"`
	Subtitle      string `json:"subtitle"`
	MainImage     string `json:"mainImage"`
	SubImages     string `json:"subImages"`
	Detail        string `json:"detail"`
	Spec          string `json:"spec"`
	Price         int64  `json:"price"`
	OriginalPrice int64  `json:"originalPrice"`
	Stock         int32  `json:"stock"`
	Brand         string `json:"brand"`
}

type CreateProductResp struct {
	Id         int64  `json:"id"`
	CategoryId int64  `json:"categoryId"`
	Name       string `json:"name"`
	Price      int64  `json:"price"`
}

type ProductInfoReq struct {
	Id int64 `form:"id"`
}

type ProductInfoResp struct {
	Id            int64  `json:"id"`
	CategoryId    int64  `json:"categoryId"`
	Name          string `json:"name"`
	Subtitle      string `json:"subtitle"`
	MainImage     string `json:"mainImage"`
	SubImages     string `json:"subImages"`
	Detail        string `json:"detail"`
	Spec          string `json:"spec"`
	Price         int64  `json:"price"`
	OriginalPrice int64  `json:"originalPrice"`
	Stock         int32  `json:"stock"`
	Sales         int32  `json:"sales"`
	Status        int32  `json:"status"`
	Brand         string `json:"brand"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type ListProductReq struct {
	Page       int32  `form:"page"`
	PageSize   int32  `form:"pageSize"`
	CategoryId int64  `form:"categoryId"`
	Keyword    string `form:"keyword"`
}

type ProductItem struct {
	Id            int64  `json:"id"`
	CategoryId    int64  `json:"categoryId"`
	Name          string `json:"name"`
	Subtitle      string `json:"subtitle"`
	MainImage     string `json:"mainImage"`
	Price         int64  `json:"price"`
	OriginalPrice int64  `json:"originalPrice"`
	Sales         int32  `json:"sales"`
	Brand         string `json:"brand"`
}

type ListProductResp struct {
	Products []ProductItem `json:"products"`
	Total    int64         `json:"total"`
}

type UpdateProductReq struct {
	Id            int64  `json:"id"`
	CategoryId    int64  `json:"categoryId"`
	Name          string `json:"name"`
	Subtitle      string `json:"subtitle"`
	MainImage     string `json:"mainImage"`
	Price         int64  `json:"price"`
	OriginalPrice int64  `json:"originalPrice"`
	Status        int32  `json:"status"`
	Brand         string `json:"brand"`
}

type UpdateProductResp struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status int32  `json:"status"`
}

type ListCategoryReq struct {
	ParentId int64 `form:"parentId"`
}

type CategoryItem struct {
	Id        int64  `json:"id"`
	ParentId  int64  `json:"parentId"`
	Name      string `json:"name"`
	SortOrder int32  `json:"sortOrder"`
}

type ListCategoryResp struct {
	Categories []CategoryItem `json:"categories"`
}

type SearchProductReq struct {
	Keyword    string `form:"keyword"`
	Page       int32  `form:"page"`
	PageSize   int32  `form:"pageSize"`
	CategoryId int64  `form:"categoryId"`
}

type SearchProductResp struct {
	Products []ProductItem `json:"products"`
	Total    int64         `json:"total"`
}
