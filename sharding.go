package sharding

import (
	"gorm.io/sharding"
)

type Sharding struct {
	tables                []any
	shardingKey           string
	numberOfShards        uint
	primaryKeyGenerator   int
	primaryKeyGeneratorFn func(int64) int64
	shardingAlgorithm     func(any) (string, error)
}

func New(shardingKey string, numberOfShards uint, opts ...Option) *Sharding {
	s := &Sharding{
		shardingKey:         shardingKey,
		numberOfShards:      numberOfShards,
		primaryKeyGenerator: sharding.PKCustom,
		primaryKeyGeneratorFn: func(id int64) int64 {
			return 0
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(*Sharding)

func WithPrimaryKeyGenerator(pk int) Option {
	return func(s *Sharding) {
		s.primaryKeyGenerator = pk
	}
}

func WithPrimaryKeyGeneratorFn(fn func(int64) int64) Option {
	return func(s *Sharding) {
		s.primaryKeyGeneratorFn = fn
	}
}

func WithTable(tables ...any) Option {
	return func(s *Sharding) {
		s.tables = tables
	}
}

// WithShardingAlgorithm 根据分表列的值来指定分表表名的后缀
func WithShardingAlgorithm(fn func(columnValue any) (suffix string, err error)) Option {
	return func(s *Sharding) {
		s.shardingAlgorithm = fn
	}
}

func (s *Sharding) Register() *sharding.Sharding {
	return sharding.Register(sharding.Config{
		ShardingKey:           s.shardingKey,
		NumberOfShards:        s.numberOfShards,
		PrimaryKeyGenerator:   s.primaryKeyGenerator,
		PrimaryKeyGeneratorFn: s.primaryKeyGeneratorFn,
		ShardingAlgorithm:     s.shardingAlgorithm,
	}, s.tables...)
}
