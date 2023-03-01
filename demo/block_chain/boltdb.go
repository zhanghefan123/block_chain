package block_chain

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type BlockChainDao struct {
	// 创建的数据库
	DB *bolt.DB
	// 数据库的名称
	DBName string
	// bucketName
	BucketName string
}

func (dao *BlockChainDao) WriteBlock(hash []byte, block *Block) {
	err := dao.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dao.BucketName))
		// 如果不存在bucket，就创建一个
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte(dao.BucketName))
			if err != nil {
				return err
			}
			// 将区块序列化之后，存储到数据库中
			err = bucket.Put(hash, block.Serialize())
			if err != nil {
				return err
			}
			// 将最后一个区块的哈希进行更新
			err = bucket.Put([]byte("lastHash"), hash)
			if err != nil {
				return err
			}
		} else {
			// 将区块序列化之后，存储到数据库中
			err := bucket.Put(hash, block.Serialize())
			if err != nil {
				return err
			}
			// 将最后一个区块的哈希进行更新
			err = bucket.Put([]byte("lastHash"), hash)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
}

// GetBlockByHash 根据哈希值获取区块
func (dao *BlockChainDao) GetBlockByHash(hash []byte) *Block {
	var resultBlock *Block
	err := dao.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dao.BucketName))
		if bucket == nil {
			fmt.Printf("bucket is nil")
		}
		// 进行相应的键的查找
		value := bucket.Get(hash)
		resultBlock = Deserialize(value)
		return nil
	})
	if err != nil {
		panic(err)
	} else {
		return resultBlock
	}
}

// CreateBlockChainDao 创建一个区块链数据库对象
func CreateBlockChainDao(tempDBName string, bucketName string, readOnly bool) *BlockChainDao {
	tempdb, err := bolt.Open(tempDBName, 0600, &bolt.Options{ReadOnly: readOnly})
	if err != nil {
		fmt.Printf("open db failed, err: %v", err)
	}
	return &BlockChainDao{
		DB:         tempdb,
		DBName:     tempDBName,
		BucketName: bucketName,
	}
}

func (dao *BlockChainDao) GetLastHash() []byte {
	// 进行lasthash的获取
	var result []byte
	err := dao.DB.View(func(tx *bolt.Tx) error {
		result = tx.Bucket([]byte(dao.BucketName)).Get([]byte("lastHash"))
		return nil
	})
	if err != nil {
		return nil
	} else {
		return result
	}
}
