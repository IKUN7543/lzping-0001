package svc

import (
	"context"
	"encoding/json"
	"go-zero-ecommerce/common/kafkaq"
	"go-zero-ecommerce/common/snowflake"
	"go-zero-ecommerce/service/order/rpc/config"
	"go-zero-ecommerce/service/order/rpc/model"
	"go-zero-ecommerce/service/product/rpc/product"
	stock2 "go-zero-ecommerce/service/stock/rpc/stock"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/zrpc"
)

type OrderKafkaMsg struct {
	OrderNo         string              `json:"orderNo"`
	UserId          int64               `json:"userId"`
	Items           []OrderItemKafkaMsg `json:"items"`
	TotalAmount     int64               `json:"totalAmount"`
	PayAmount       int64               `json:"payAmount"`
	ReceiverName    string              `json:"receiverName"`
	ReceiverPhone   string              `json:"receiverPhone"`
	ReceiverAddress string              `json:"receiverAddress"`
	Remark          string              `json:"remark"`
}

type OrderItemKafkaMsg struct {
	ProductId int64 `json:"productId"`
	Num       int32 `json:"num"`
	Price     int64 `json:"price"`
}

type ServiceContext struct {
	Config         config.Config
	OrderModel     model.OrderModel
	OrderItemModel model.OrderItemModel
	Rdb            *redis.Client
	SnowWorker     *snowflake.Worker
	KafkaProducer  *kafkaq.Producer
	KafkaConsumer  *kafkaq.Consumer
	StockRpc       stock2.StockClient
	ProductRpc     product.ProductClient
	cancelFn       context.CancelFunc
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

	worker, err := snowflake.NewWorker(1)
	if err != nil {
		panic(err)
	}

	var kafkaProducer *kafkaq.Producer
	var kafkaConsumer *kafkaq.Consumer
	var cancelFn context.CancelFunc

	if c.Kafka.Brokers != "" {
		brokers := strings.Split(c.Kafka.Brokers, ",")
		kafkaProducer = kafkaq.NewProducer(brokers)
		kafkaConsumer = kafkaq.NewConsumer(brokers, c.Kafka.Topic, c.Kafka.Group)

		svcCtx := &ServiceContext{
			Config:         c,
			OrderModel:     model.NewOrderModel(db),
			OrderItemModel: model.NewOrderItemModel(db),
			Rdb:            rdb,
			SnowWorker:     worker,
			KafkaProducer:  kafkaProducer,
			KafkaConsumer:  kafkaConsumer,
		}

		svcCtx.StockRpc = stock2.NewStockClient(zrpc.MustNewClient(c.StockRpc).Conn())
		svcCtx.ProductRpc = product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn())

		ctx, cancel := context.WithCancel(context.Background())
		cancelFn = cancel
		go startOrderConsumer(ctx, kafkaConsumer, svcCtx)

		return svcCtx
	}

	stockRpc := stock2.NewStockClient(zrpc.MustNewClient(c.StockRpc).Conn())
	productRpc := product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn())

	return &ServiceContext{
		Config:         c,
		OrderModel:     model.NewOrderModel(db),
		OrderItemModel: model.NewOrderItemModel(db),
		Rdb:            rdb,
		SnowWorker:     worker,
		KafkaProducer:  kafkaProducer,
		KafkaConsumer:  kafkaConsumer,
		StockRpc:       stockRpc,
		ProductRpc:     productRpc,
		cancelFn:       cancelFn,
	}
}

func (s *ServiceContext) Close() {
	if s.cancelFn != nil {
		s.cancelFn()
	}
	if s.KafkaConsumer != nil {
		s.KafkaConsumer.Close()
	}
	if s.KafkaProducer != nil {
		s.KafkaProducer.Close()
	}
}

func startOrderConsumer(ctx context.Context, consumer *kafkaq.Consumer, svcCtx *ServiceContext) {
	logx.Info("Order Kafka consumer started, waiting for messages...")
	err := consumer.Consume(ctx, func(consumeCtx context.Context, key string, value []byte) error {
		var msg OrderKafkaMsg
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorf("Order consumer unmarshal error: %v, key=%s", err, key)
			return nil
		}
		logx.Infof("Order consumer received: orderNo=%s, userId=%d, items=%d", msg.OrderNo, msg.UserId, len(msg.Items))
		processOrderFromKafka(consumeCtx, svcCtx, &msg)
		return nil
	})
	if err != nil {
		logx.Errorf("Order consumer stopped with error: %v", err)
	}
}

func processOrderFromKafka(ctx context.Context, svcCtx *ServiceContext, msg *OrderKafkaMsg) {
	statusKey := "order_status:" + msg.OrderNo
	exists, err := svcCtx.Rdb.Exists(ctx, statusKey).Result()
	if err != nil {
		logx.Errorf("Redis exists check error for order %s: %v", msg.OrderNo, err)
	}
	if exists > 0 {
		logx.Infof("Order already processed, skipping: %s", msg.OrderNo)
		return
	}

	productNames := make(map[int64]string)
	productImages := make(map[int64]string)
	if svcCtx.ProductRpc != nil {
		productIds := make([]int64, 0, len(msg.Items))
		for _, it := range msg.Items {
			productIds = append(productIds, it.ProductId)
		}
		prodRes, err := svcCtx.ProductRpc.GetProductList(ctx, &product.GetProductListRequest{Ids: productIds})
		if err == nil && prodRes != nil {
			for _, p := range prodRes.Products {
				productNames[p.Id] = p.Name
				productImages[p.Id] = p.MainImage
			}
		}
	}

	o := &model.Order{
		OrderNo:         msg.OrderNo,
		UserId:          msg.UserId,
		TotalAmount:     msg.TotalAmount,
		PayAmount:       msg.PayAmount,
		Status:          1,
		ReceiverName:    msg.ReceiverName,
		ReceiverPhone:   msg.ReceiverPhone,
		ReceiverAddress: msg.ReceiverAddress,
		Remark:          msg.Remark,
	}

	items := make([]*model.OrderItem, 0, len(msg.Items))
	for _, it := range msg.Items {
		name := productNames[it.ProductId]
		image := productImages[it.ProductId]
		items = append(items, &model.OrderItem{
			ProductId:    it.ProductId,
			Num:          it.Num,
			Price:        it.Price,
			TotalPrice:   it.Price * int64(it.Num),
			ProductName:  name,
			ProductImage: image,
		})
	}

	err = svcCtx.OrderModel.Insert(ctx, o, items)
	if err != nil {
		logx.Errorf("Order insert error: orderNo=%s, err=%v", msg.OrderNo, err)
		return
	}

	svcCtx.Rdb.Set(ctx, statusKey, "1", 24*time.Hour)
	logx.Infof("Order created from kafka: orderNo=%s, id=%d, itemCount=%d", msg.OrderNo, o.Id, len(items))
}
