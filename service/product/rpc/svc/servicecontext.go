package svc

import (
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"go-zero-ecommerce/common/bloom"
	"go-zero-ecommerce/common/cache"
	"go-zero-ecommerce/service/product/rpc/config"
	"go-zero-ecommerce/service/product/rpc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type ServiceContext struct {
	Config        config.Config
	ProductModel  model.ProductModel
	CategoryModel model.CategoryModel
	Rdb           *redis.Client
	ProductCache  *cache.CacheAside
	BloomFilter   *bloom.RedisBloomFilter
	ESClient      *elastic.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       0,
	})

	productCache := cache.NewCacheAside(rdb, "product", 5*time.Minute)
	bloomFilter := bloom.NewRedisBloomFilter(rdb, "bloom:product", 1<<24, 8)

	var esClient *elastic.Client
	if c.ES.Urls != "" {
		esClient, err = elastic.NewClient(elastic.SetURL(c.ES.Urls), elastic.SetSniff(false))
		if err != nil {
			panic(err)
		}
	}

	return &ServiceContext{
		Config:        c,
		ProductModel:  model.NewProductModel(db),
		CategoryModel: model.NewCategoryModel(db),
		Rdb:           rdb,
		ProductCache:  productCache,
		BloomFilter:   bloomFilter,
		ESClient:      esClient,
	}
}
