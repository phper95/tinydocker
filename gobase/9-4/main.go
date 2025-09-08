package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"go.etcd.io/bbolt"
)

const TableName = "mytable"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 打开/创建BoltDB数据库
	db, err := bbolt.Open("mydb.db", 0600, nil)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	// 重要：在函数退出前关闭数据库，避免文件句柄泄漏
	defer db.Close()

	log.Println("数据库初始化成功")
	// 后续读写操作...
	err = db.Update(func(tx *bbolt.Tx) error {
		// 创建Bucket（若不存在）
		bkt, err := tx.CreateBucketIfNotExists([]byte(TableName))
		if err != nil {
			return fmt.Errorf("创建Bucket失败: %v", err)
		}

		// 写入数据：key和value均为[]byte
		// 字符串类型
		err = bkt.Put([]byte("name"), []byte("张三"))
		if err != nil {
			return err
		}
		// 数字类型（需用strconv转换为字符串）
		age := 28
		err = bkt.Put([]byte("age"), []byte(strconv.Itoa(age)))
		if err != nil {
			return err
		}
		// JSON类型（复杂结构需序列化为[]byte）
		user := map[string]interface{}{"gender": "male", "city": "Beijing"}
		userJSON, _ := json.Marshal(user)
		err = bkt.Put([]byte("detail"), userJSON)
		if err != nil {
			return err
		}

		return nil // 事务无错误，BoltDB自动提交
	})
	if err != nil {
		log.Fatalf("写入数据失败: %v", err)
	}
	log.Println("数据写入成功")

	// 读取数据
	err = db.View(func(tx *bbolt.Tx) error {
		// 获取Bucket
		bkt := tx.Bucket([]byte(TableName))
		if bkt == nil {
			return fmt.Errorf("Bucket 'user_info' 不存在")
		}

		// 读取单个key
		name := bkt.Get([]byte("name"))
		if name == nil {
			return fmt.Errorf("key 'name' 不存在")
		}
		log.Printf("姓名: %s\n", name) // 输出：姓名: 张三

		// 读取数字类型（需转换）
		ageBytes := bkt.Get([]byte("age"))
		age, _ := strconv.Atoi(string(ageBytes))
		log.Printf("年龄: %d\n", age) // 输出：年龄: 28

		// 读取JSON类型（需反序列化）
		detailBytes := bkt.Get([]byte("detail"))
		var detail map[string]interface{}
		err = json.Unmarshal(detailBytes, &detail)
		if err != nil {
			log.Println("JSON反序列化失败: %v", err)
		}
		log.Printf("详细信息: %v\n", detail) // 输出：详细信息: map[city:Beijing gender:male]

		return nil // 只读事务无需提交，自动结束
	})
	if err != nil {
		log.Fatalf("读取数据失败: %v", err)
	}
}
