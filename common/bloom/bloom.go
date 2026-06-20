package bloom

import (
	"context"
	"go-zero-ecommerce/common/errx"

	"github.com/bits-and-blooms/bitset"
	"github.com/redis/go-redis/v9"
	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	m      uint
	k      uint
	bitset *bitset.BitSet
}

func NewBloomFilter(m uint, k uint) *BloomFilter {
	return &BloomFilter{
		m:      m,
		k:      k,
		bitset: bitset.New(m),
	}
}

func (bf *BloomFilter) Add(data []byte) {
	h1, h2 := hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + uint32(i)*h2) % uint32(bf.m)
		bf.bitset.Set(uint(idx))
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	h1, h2 := hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + uint32(i)*h2) % uint32(bf.m)
		if !bf.bitset.Test(uint(idx)) {
			return false
		}
	}
	return true
}

func hash(data []byte) (uint32, uint32) {
	h := murmur3.New128()
	h.Write(data)
	s1, s2 := h.Sum128()
	return uint32(s1), uint32(s2)
}

type RedisBloomFilter struct {
	key string
	m   uint
	k   uint
	rdb interface {
		SetBit(ctx context.Context, key string, offset int64, value int) (int64, error)
		GetBit(ctx context.Context, key string, offset int64) (int64, error)
	}
}

func NewRedisBloomFilter(rdb interface{}, key string, m uint, k uint) *RedisBloomFilter {
	return &RedisBloomFilter{
		key: key,
		m:   m,
		k:   k,
		rdb: rdb,
	}
}

func (bf *RedisBloomFilter) Add(ctx context.Context, data []byte) error {
	h1, h2 := hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := int64((h1 + uint32(i)*h2) % uint32(bf.m))
		_, err := bf.rdb.SetBit(ctx, bf.key, idx, 1)
		if err != nil {
			return errx.ErrInternalServer
		}
	}
	return nil
}

func (bf *RedisBloomFilter) Contains(ctx context.Context, data []byte) (bool, error) {
	h1, h2 := hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := int64((h1 + uint32(i)*h2) % uint32(bf.m))
		bit, err := bf.rdb.GetBit(ctx, bf.key, idx)
		if err != nil {
			return false, errx.ErrInternalServer
		}
		if bit == 0 {
			return false, nil
		}
	}
	return true, nil
}
