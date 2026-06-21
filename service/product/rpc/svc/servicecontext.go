package svc

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-ecommerce/common/bloom"
	"go-zero-ecommerce/common/cache"
	"go-zero-ecommerce/service/product/rpc/config"
	"go-zero-ecommerce/service/product/rpc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const productIndex = "products"

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

	poolSize := c.Redis.PoolSize
	if poolSize <= 0 {
		poolSize = 100
	}
	minIdle := c.Redis.MinIdleConns
	if minIdle <= 0 {
		minIdle = 10
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Host,
		Password:     c.Redis.Pass,
		DB:           0,
		PoolSize:     poolSize,
		MinIdleConns: minIdle,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		MaxRetries:   3,
	})

	productCache := cache.NewCacheAside(rdb, "product", 5*time.Minute)
	bloomFilter := bloom.NewRedisBloomFilter(rdb, "bloom:product", 1<<24, 8)

	var esClient *elastic.Client
	if c.ES.Urls != "" {
		esClient, err = elastic.NewClient(elastic.SetURL(c.ES.Urls), elastic.SetSniff(false))
		if err != nil {
			panic(err)
		}
		if err := ensureProductIndex(esClient); err != nil {
			logx.Errorf("Failed to ensure ES product index: %v", err)
		} else {
			logx.Infof("ES product index '%s' is ready", productIndex)
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

func ensureProductIndex(client *elastic.Client) error {
	ctx := context.Background()
	exists, err := client.IndexExists(productIndex).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"id":             { "type": "long" },
				"categoryId":     { "type": "long" },
				"name":           { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart" },
				"subtitle":       { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart" },
				"mainImage":      { "type": "keyword" },
				"subImages":      { "type": "keyword" },
				"detail":         { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart" },
				"spec":           { "type": "keyword" },
				"price":          { "type": "long" },
				"originalPrice":  { "type": "long" },
				"stock":          { "type": "integer" },
				"sales":          { "type": "integer" },
				"status":         { "type": "integer" },
				"brand":          { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart", "fields": { "keyword": { "type": "keyword" } } },
				"createdAt":      { "type": "date", "format": "yyyy-MM-dd HH:mm:ss||epoch_millis" },
				"updatedAt":      { "type": "date", "format": "yyyy-MM-dd HH:mm:ss||epoch_millis" }
			}
		},
		"settings": {
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"ik_smart": {
						"type": "custom",
						"tokenizer": "standard"
					},
					"ik_max_word": {
						"type": "custom",
						"tokenizer": "standard"
					}
				}
			}
		}
	}`

	_, err = client.CreateIndex(productIndex).Body(mapping).Do(ctx)
	return err
}
