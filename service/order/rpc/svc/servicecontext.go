package svc

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-ecommerce/common/kafkaq"
	"go-zero-ecommerce/common/snowflake"
	"go-zero-ecommerce/service/order/rpc/config"
	"go-zero-ecommerce/service/order/rpc/model"
	stock2 "go-zero-ecommerce/service/stock/rpc/stock"
	"go-zero-ecommerce/service/product/rpc/product"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type OrderKafkaMsg struct {
	OrderNo       string            `json:"orderNo"`
	UserId        int64             `json:"userId"`
	Items         []OrderItemKafkaMsg `json:"items"`
	TotalAmount   int64             `json:"totalAmount"`
	PayAmount     int64             `json:"payAmount"`
	ReceiverName    string          `json:"receiverName"`
	ReceiverPhone   string          `json:"receiverPhone"`
	ReceiverAddress string          `json:"receiverAddress"`
	Remark        string            `json:"remark"`
}

type OrderItemKafkaMsg struct {
	ProductId int64 `json:"productId"`
	Num       int32 `json:"num"`
	Price     int64 `json:"price"`
}

type ServiceContext struct {
	Config        config.Config
	OrderModel    model.OrderModel
	OrderItemModel model.OrderItemModel
	Rdb           *redis.Client
	SnowWorker    *snowflake.Worker
	KafkaProducer *kafkaq.Producer
	KafkaConsumer *kafkaq.Consumer
	StockRpc      stock2.StockClient
	ProductRpc    product.ProductClient
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

	worker, err := snowflake.NewWorker(1)
	if err != nil {
		panic(err)
	}

	var kafkaProducer *kafkaq.Producer
	var kafkaConsumer *kafkaq.Consumer
	if c.Kafka.Brokers != "" {
		brokers := strings.Split(c.Kafka.Brokers, ",")
		kafkaProducer = kafkaq.NewProducer(brokers)
		kafkaConsumer = kafkaq.NewConsumer(brokers, c.Kafka.Topic, c.Kafka.Group)

		go startOrderConsumer(kafkaConsumer, db, worker)
	}

	stockRpc := stock2.NewStockClient(zrpc.MustNewClient(c.StockRpc).Conn())
	productRpc := product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn())

	return &ServiceContext{
		Config:        c,
		OrderModel:    model.NewOrderModel(db),
		OrderItemModel: model.NewOrderItemModel(db),
		Rdb:           rdb,
		SnowWorker:    worker,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
		StockRpc:      stockRpc,
		ProductRpc:    productRpc,
	}
}

func startOrderConsumer(consumer *kafkaq.Consumer, db *gorm.DB, worker *snowflake.Worker) {
	logx.Info("Order Kafka consumer started")
	err := consumer.Consume(context.Background(), func(ctx context.Context, key string, value []byte) error {
		var msg OrderKafkaMsg
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorf("Order consumer unmarshal error: %v", err)
			return nil
		}
		logx.Infof("Order consumer received: orderNo=%s", msg.OrderNo)
		processOrderFromKafka(ctx, db, worker, &msg)
		return nil
	})
	if err != nil {
		logx.Errorf("Order consumer error: %v", err)
	}
}

func processOrderFromKafka(ctx context.Context, db *gorm.DB, worker *snowflake.Worker, msg *OrderKafkaMsg) {
	orderModel := model.NewOrderModel(db)

	statusKey := "order_status:" + msg.OrderNo
	exists, _ := redis.NewClient(&redis.Options{Addr: "redis:6379"}).Exists(ctx, statusKey).Result()
	if exists > 0 {
		logx.Infof("Order already processed: %s", msg.OrderNo)
		return
	}

	order := &model.Order{
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
		items = append(items, &model.OrderItem{
			ProductId: it.ProductId, Num: it.Num, Price: it.Price, TotalPrice: it.Price * int64(it.Num),
		})
	}

	err := orderModel.Insert(ctx, order, items)
	if err != nil {
		logx.Errorf("Order insert error: orderNo=%s, err=%v", msg.OrderNo, err)
		return
	}

	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	rdb.Set(ctx, statusKey, "1", 24*3600*1e9)
	logx.Infof("Order created from kafka: orderNo=%s, id=%d", msg.OrderNo, order.Id)
}
