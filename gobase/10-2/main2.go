package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type User struct {
	// gin通过binding标签来定义校验规则
	// required 表示该参数必须存在
	Name string `json:"name" binding:"required"`

	// gte=18表示该参数的值必须大于等于18
	Age int `json:"age" binding:"gte=18"`

	// age值等于18
	// Age      int    `json:"age" binding:"eq=18"`

	// age值不等于18
	// Age      int    `json:"age" binding:"ne=18"`

	// email 合法的邮箱格式
	Email string `json:"email" binding:"required,email"`

	// Password字段长度至少为6位且包含123@
	Password string `json:"password" binding:"required,contains=123@,min=6"`

	// eqfield表示跟指定字段相等，同理还有 qfield ,nefield ,gtfield ,gtefield ,ltfield ,ltefield
	// 注意，这里的字段需要跟结构体中的字段名一致，而非是结构体字段的标签名
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`

	// url 合法的url格式
	// URL string `json:"url" binding:"required,url"`

	// 合法的ip格式
	// IP string `json:"ip" binding:"required,ip"`

	// 合法的ipv4格式
	// IPv4 string `json:"ipv4" binding:"required,ipv4"`

	// 合法的ipv6格式
	// IPv6 string `json:"ipv6" binding:"required,ipv6"`

	// []string长度必须大于1，数组中元素string长度必须在2-100之间
	Tags []string `json:"tags" binding:"gt=1,dive,required,min=2,max=100"`

	// 限制key值和value值的长度都必须在2-100之间
	// dive 用于深入到切片、数组或map的内部元素进行验证,对于切片/数组：验证每个元素；对于map：分别验证key和value。
	// keys,min=2,max=100 验证所有key的长度在2-100之间
	// endkeys 结束key验证阶段
	// required,min=2,max=100 验证所有value为必填且长度在2-100之间
	M map[string]string `json:"m" binding:"dive,keys,min=2,max=100,endkeys,required,min=2,max=100"`
	// 默认会递归验证整个结构体
	// 会检查 st 结构体内部的 F1 字段是否满足 required,min=6 验证规则
	Extra1 st `json:"extra1"`
	// 使用了 structonly 标签 只验证结构体本身是否存在，不递归验证内部字段
	// 即使 st 结构体内部的 F1 字段不满足验证规则，也不会影响验证结果
	Extra2 st `json:"extra2" binding:"structonly"`
	// 使用了 - 标签,完全跳过该字段的验证
	// 无论 Extra3 字段是否存在或其内部字段是否满足规则，都不会进行任何验证检查
	Extra3 st `json:"extra3" binding:"-"`

	// 自定义验证器：dateLteNow,时间戳必须小于等于当前时间
	Mtime int64 `json:"mtime" binding:"dateLteNow"`
}

type st struct {
	F1 string `json:"f1" binding:"required,min=6"`
}

func main() {
	router := gin.Default()
	// 注册验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 绑定第一个参数是验证的函数第二个参数是自定义的验证函数
		err := v.RegisterValidation("dateLteNow", dateLteNow)
		if err != nil {
			fmt.Println("register validation failed", err)
		}
	}
	router.POST("register", Register)
	err := router.Run(":80")
	if err != nil {
		fmt.Println("server start failed", err)
	}
}
func Register(c *gin.Context) {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		fmt.Println("param check failed")
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	// 验证 存储操作省略.....
	fmt.Println("register success")
	c.JSON(http.StatusOK, "successful")
}

// 自定义验证函数，通过反射获取字段值并进行校验
func dateLteNow(fileLevel validator.FieldLevel) bool {
	// 获取字段值
	t := fileLevel.Field().Int()
	fmt.Println("t:", t)
	if t == 0 {
		return false
	}
	// 与当前时间对比
	if time.Now().Unix()-t < 0 {
		return false
	}
	return true
}
