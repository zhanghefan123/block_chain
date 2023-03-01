package test

import (
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
)

// 注意全局之中定义的变量，在TestMain之中不能够进行重名
var db *bolt.DB
var testBucket []byte = []byte("testBucket")
var testKey []byte = []byte("testKey")
var testValue []byte = []byte("testValue")
var nameOfDB = "./database/test.db"

// TestInsertRecord 测试将键值对插入到数据库中
func TestInsertRecord(t *testing.T) {
	// 创建表
	err := db.Update(func(tx *bolt.Tx) error {
		// 进行bucket的创建
		bucket, err := tx.CreateBucketIfNotExists(testBucket)
		if err != nil {
			t.Errorf("create bucket failed, err: %v", err)
		}
		// 创建好了bucket之后，就可以进行键值对的插入了
		err = bucket.Put(testKey, testValue)
		return nil
	})
	if err != nil {
		t.Errorf("update db failed, err: %v", err)
	}
}

// TestSearchRecord 测试从数据库中查询键值对
func TestSearchRecord(t *testing.T) {
	// 进入view视图,就不能够进行数据的修改了
	err := db.View(func(tx *bolt.Tx) error {
		// 找到bucket
		bucket := tx.Bucket(testBucket)
		if bucket == nil {
			t.Errorf("bucket is nil")
		}
		// 进行相应的键的查找
		value := bucket.Get(testKey)
		t.Logf("found key: %s, value: %s", testKey, value)
		return nil
	})
	if err != nil {
		t.Errorf("search record failed, err: %v", err)
	}
}

func TestMain(m *testing.M) {
	// 创建或者打开数据库,创建的数据库的相对路径是当前项目的根目录
	// 参数1：数据库文件名
	// 参数2：文件权限
	// 参数3：数据库选项
	tempDb, err := bolt.Open(nameOfDB, 0600, nil)
	if err != nil {
		fmt.Printf("open db failed, err: %v", err)
	} else {
		db = tempDb
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("close db failed, err: %v", err)
		}
	}(db)
	m.Run()
}
