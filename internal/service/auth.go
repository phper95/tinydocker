package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phper95/tinydocker/internal/model"
	"github.com/phper95/tinydocker/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务结构体，负责处理用户认证、授权和JWT令牌管理
type AuthService struct {
	dataDir        string                       // 用户数据存储目录路径
	jwtSecret      []byte                       // JWT签名密钥
	permissions    map[string]*model.Permission // 权限映射表，key为权限ID
	roles          map[string]*model.Role       // 角色映射表，key为角色ID
	tokenBlacklist map[string]int64             // Token黑名单，存储已注销的Token及其过期时间
	mu             sync.RWMutex                 // 读写锁，保护tokenBlacklist的并发访问
}

// AuthDataDir 认证数据存储的默认目录路径常量
const AuthDataDir = "/var/lib/tinydocker/auth"

// NewAuthService 创建并初始化认证服务实例
// 参数 jwtSecret: JWT签名密钥
// 返回值 *AuthService: 初始化完成的认证服务实例
func NewAuthService(jwtSecret []byte) *AuthService {
	service := &AuthService{
		dataDir:        AuthDataDir,
		jwtSecret:      jwtSecret,
		permissions:    make(map[string]*model.Permission), // 初始化权限映射表
		roles:          make(map[string]*model.Role),       // 初始化角色映射表
		tokenBlacklist: make(map[string]int64),             // 初始化Token黑名单
	}

	// 初始化默认权限和角色
	service.initDefaultPermissions()
	service.initDefaultRoles()

	return service
}

// Login 用户登录验证
// 参数 req: 登录请求参数，包含用户名和密码
// 返回值 *model.LoginResponse, error: 登录响应信息和可能的错误
func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// 根据用户名查找用户信息
	user, err := s.getUserByUsername(req.Username)
	if err != nil {
		logger.Error("查找用户失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	// 使用bcrypt验证密码是否匹配
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Error("密码验证失败: %v", err)
		return nil, fmt.Errorf("密码错误")
	}

	// 检查用户账户是否处于激活状态
	if !user.IsActive {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 更新用户最后登录时间戳
	now := time.Now()
	user.LastLoginAt = &now
	err = s.saveUser(user)
	if err != nil {
		logger.Error("更新用户最后登录时间失败: %v", err)
	}

	// 为用户生成JWT访问令牌
	token, expires, err := s.generateJWTToken(&user.UserBase)
	if err != nil {
		logger.Error("生成 Token 失败: %v", err)
		return nil, fmt.Errorf("生成 Token 失败: %v", err)
	}

	// 构造并返回登录响应
	return &model.LoginResponse{
		Token:   token,          // JWT访问令牌
		User:    &user.UserBase, // 用户基本信息
		Expires: expires,        // 令牌过期时间戳
	}, nil
}

// Register 用户注册
// 参数 req: 用户注册请求参数
// 返回值 *model.User, error: 注册成功的用户信息和可能的错误
func (s *AuthService) Register(req *model.CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已被占用
	if _, err := s.getUserByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已被占用
	if _, err := s.getUserByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 使用bcrypt对用户密码进行加密处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密码加密失败: %v", err)
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建新的用户对象
	user := &model.User{
		UserBase: model.UserBase{
			ID:        generateUserID(),         // 生成唯一用户ID
			Username:  req.Username,             // 用户名
			Email:     req.Email,                // 邮箱
			Roles:     req.Roles,                // 用户角色列表
			CreatedAt: time.Now(),               // 创建时间
			UpdatedAt: time.Now(),               // 更新时间
			IsActive:  true,                     // 默认激活状态
			Metadata:  make(map[string]string)}, // 用户元数据
	}
	user.Password = string(hashedPassword) // 设置加密后的密码

	// 根据用户角色获取对应的权限列表
	user.Permissions = s.getUserPermissions(user)

	// 将新用户信息持久化存储
	if err := s.saveUser(user); err != nil {
		logger.Error("保存用户失败: %v", err)
		return nil, fmt.Errorf("保存用户失败: %v", err)
	}

	return user, nil
}

// ValidateToken 验证JWT令牌的有效性
// 参数 tokenString: 待验证的JWT令牌字符串
// 返回值 *model.AuthContext, error: 认证上下文信息和可能的错误
func (s *AuthService) ValidateToken(tokenString string) (*model.AuthContext, error) {
	// 首先检查令牌是否在黑名单中（已注销）
	if s.isTokenRevoked(tokenString) {
		logger.Error("Token 已被注销")
		return nil, fmt.Errorf("Token 已被注销")
	}

	// 解析并验证JWT令牌
	// 第二个参数是一个回调函数，用于提供验证令牌所需的密钥和执行额外验证
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为HMAC（防止JWT算法混淆攻击）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		logger.Error("Token 解析失败: %v", err)
		return nil, fmt.Errorf("Token 解析失败: %v", err)
	}

	// 验证令牌是否有效
	if !token.Valid {
		return nil, fmt.Errorf("Token 无效")
	}

	// 提取令牌中的声明信息
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Token 声明无效")
	}

	// 从声明中获取用户ID
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("Token 中缺少用户 ID")
	}

	// 根据用户ID获取完整的用户信息
	user, err := s.getUserByID(userID)
	if err != nil {
		logger.Error("获取用户信息失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	// 检查用户账户是否处于激活状态
	if !user.IsActive {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 构造并返回认证上下文
	return &model.AuthContext{
		UserID:      user.ID,          // 用户ID
		Username:    user.Username,    // 用户名
		Roles:       user.Roles,       // 用户角色列表
		Permissions: user.Permissions, // 用户权限列表
	}, nil
}

// CheckPermission 检查用户是否具有指定资源和操作的权限
// 参数 ctx: 认证上下文，包含用户信息
// 参数 resource: 资源名称（如containers、images等）
// 参数 action: 操作名称（如create、delete等）
// 返回值 bool: 是否具有权限
func (s *AuthService) CheckPermission(ctx *model.AuthContext, resource, action string) bool {
	// 首先检查用户的直接权限（支持通配符匹配）
	for _, perm := range ctx.Permissions {
		if s.permissionMatch(perm, resource, action) {
			return true
		}
	}

	// 然后检查用户角色对应的权限（支持通配符匹配）
	for _, roleName := range ctx.Roles {
		if role, exists := s.roles[roleName]; exists {
			for _, perm := range role.Permissions {
				if s.permissionMatch(perm, resource, action) {
					return true
				}
			}
		}
	}

	// 如果都没有匹配的权限，则返回false
	return false
}

// generateJWTToken 生成JWT访问令牌
// 参数 user: 用户基本信息
// 返回值 string, int64, error: 令牌字符串、过期时间戳和可能的错误
func (s *AuthService) generateJWTToken(user *model.UserBase) (string, int64, error) {
	// 设置令牌24小时后过期
	expires := time.Now().Add(24 * time.Hour).Unix()

	// 构造JWT声明信息
	claims := jwt.MapClaims{
		"user_id":  user.ID,           // 用户ID
		"username": user.Username,     // 用户名
		"roles":    user.Roles,        // 用户角色列表
		"exp":      expires,           // 过期时间
		"iat":      time.Now().Unix(), // 签发时间
	}

	// 创建JWT令牌对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对令牌进行签名
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		logger.Error("Token 签名失败: %v", err)
		return "", 0, err
	}

	return tokenString, expires, nil
}

// getUserByUsername 根据用户名查找用户
// 参数 username: 用户名
// 返回值 *model.User, error: 用户对象和可能的错误
func (s *AuthService) getUserByUsername(username string) (*model.User, error) {
	// 获取所有用户列表
	users, err := s.getAllUsers()
	if err != nil {
		logger.Error("获取所有用户失败: %v", err)
		return nil, err
	}

	// 遍历查找匹配的用户名
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("用户不存在")
}

// getUserByEmail 根据邮箱查找用户
// 参数 email: 用户邮箱
// 返回值 *model.User, error: 用户对象和可能的错误
func (s *AuthService) getUserByEmail(email string) (*model.User, error) {
	// 获取所有用户列表
	users, err := s.getAllUsers()
	if err != nil {
		logger.Error("获取所有用户失败: %v", err)
		return nil, err
	}

	// 遍历查找匹配的邮箱
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("用户不存在")
}

// getUserByID 根据用户ID查找用户
// 参数 id: 用户ID
// 返回值 *model.User, error: 用户对象和可能的错误
func (s *AuthService) getUserByID(id string) (*model.User, error) {
	// 获取所有用户列表
	users, err := s.getAllUsers()
	if err != nil {
		logger.Error("获取所有用户失败: %v", err)
		return nil, err
	}

	// 遍历查找匹配的用户ID
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, fmt.Errorf("用户不存在")
}

// getAllUsers 获取所有用户列表
// 返回值 []*model.User, error: 用户对象列表和可能的错误
func (s *AuthService) getAllUsers() ([]*model.User, error) {
	// 构造用户数据目录路径
	usersDir := filepath.Join(s.dataDir, "users")

	// 检查用户目录是否存在
	if _, err := os.Stat(usersDir); os.IsNotExist(err) {
		return []*model.User{}, nil
	}

	// 读取用户目录下的所有条目
	entries, err := os.ReadDir(usersDir)
	if err != nil {
		logger.Error("读取用户目录失败: %v", err)
		return nil, err
	}

	// 遍历目录条目，加载所有用户数据
	var users []*model.User
	for _, entry := range entries {
		if entry.IsDir() {
			user, err := s.loadUser(entry.Name())
			if err != nil {
				logger.Error("加载用户失败: %v", err)
				continue
			}
			users = append(users, user)
		}
	}

	return users, nil
}

// loadUser 从文件加载指定用户的数据
// 参数 id: 用户ID
// 返回值 *model.User, error: 用户对象和可能的错误
func (s *AuthService) loadUser(id string) (*model.User, error) {
	// 构造用户配置文件路径
	configPath := filepath.Join(s.dataDir, "users", id, "config.json")

	// 读取配置文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("读取用户配置失败: %v", err)
		return nil, err
	}

	// 解析JSON格式的用户数据
	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		logger.Error("解析用户配置失败: %v", err)
		return nil, err
	}

	return &user, nil
}

// saveUser 将用户数据保存到文件
// 参数 user: 用户对象
// 返回值 error: 可能的错误
func (s *AuthService) saveUser(user *model.User) error {
	// 创建用户数据目录
	userDir := filepath.Join(s.dataDir, "users", user.ID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		logger.Error("创建用户目录失败: %v", err)
		return err
	}

	// 构造配置文件路径
	configPath := filepath.Join(userDir, "config.json")

	// 将用户对象序列化为JSON格式
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		logger.Error("序列化用户数据失败: %v", err)
		return err
	}

	// 写入配置文件
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		logger.Error("写入用户配置文件失败: %v", err)
		return err
	}

	return nil
}

// getUserPermissions 根据用户角色获取用户权限列表
// 参数 user: 用户对象
// 返回值 []string: 权限ID列表
func (s *AuthService) getUserPermissions(user *model.User) []string {
	var permissions []string

	// 遍历用户的所有角色，收集对应权限
	for _, roleName := range user.Roles {
		if role, exists := s.roles[roleName]; exists {
			permissions = append(permissions, role.Permissions...)
		}
	}

	return permissions
}

// initDefaultPermissions 初始化默认权限
func (s *AuthService) initDefaultPermissions() {
	// 定义系统默认权限列表
	permissions := []*model.Permission{
		{ID: "containers:list", Name: "列出容器", Description: "查看容器列表", Resource: "containers", Action: "list"},
		{ID: "containers:create", Name: "创建容器", Description: "创建新容器", Resource: "containers", Action: "create"},
		{ID: "containers:get", Name: "查看容器", Description: "查看容器详情", Resource: "containers", Action: "get"},
		{ID: "containers:update", Name: "更新容器", Description: "更新容器配置", Resource: "containers", Action: "update"},
		{ID: "containers:delete", Name: "删除容器", Description: "删除容器", Resource: "containers", Action: "delete"},
		{ID: "containers:start", Name: "启动容器", Description: "启动容器", Resource: "containers", Action: "start"},
		{ID: "containers:stop", Name: "停止容器", Description: "停止容器", Resource: "containers", Action: "stop"},
		{ID: "images:list", Name: "列出镜像", Description: "查看镜像列表", Resource: "images", Action: "list"},
		{ID: "images:create", Name: "创建镜像", Description: "构建或导入镜像", Resource: "images", Action: "create"},
		{ID: "images:get", Name: "查看镜像", Description: "查看镜像详情", Resource: "images", Action: "get"},
		{ID: "images:delete", Name: "删除镜像", Description: "删除镜像", Resource: "images", Action: "delete"},
		{ID: "networks:list", Name: "列出网络", Description: "查看网络列表", Resource: "networks", Action: "list"},
		{ID: "networks:create", Name: "创建网络", Description: "创建新网络", Resource: "networks", Action: "create"},
		{ID: "networks:get", Name: "查看网络", Description: "查看网络详情", Resource: "networks", Action: "get"},
		{ID: "networks:delete", Name: "删除网络", Description: "删除网络", Resource: "networks", Action: "delete"},
		{ID: "system:read", Name: "查看系统信息", Description: "查看系统信息", Resource: "system", Action: "read"},
	}

	// 将权限添加到权限映射表中
	for _, perm := range permissions {
		s.permissions[perm.ID] = perm
	}
}

// initDefaultRoles 初始化默认角色
func (s *AuthService) initDefaultRoles() {
	// 定义系统默认角色及其权限
	roles := []*model.Role{
		{
			ID:          "admin", // 管理员角色
			Name:        "管理员",
			Description: "拥有所有权限",
			Permissions: []string{"*:*"}, // 通配符表示所有权限
		},
		{
			ID:          "developer", // 开发者角色
			Name:        "开发者",
			Description: "可以管理容器和镜像",
			Permissions: []string{
				"containers:*",  // 容器相关所有权限
				"images:*",      // 镜像相关所有权限
				"networks:list", // 网络列表查看权限
				"networks:read", // 网络详情查看权限
				"system:read",   // 系统信息查看权限
			},
		},
		{
			ID:          "operator", // 运维人员角色
			Name:        "运维人员",
			Description: "可以管理容器和网络",
			Permissions: []string{
				"containers:*", // 容器相关所有权限
				"networks:*",   // 网络相关所有权限
				"images:list",  // 镜像列表查看权限
				"images:read",  // 镜像详情查看权限
				"system:read",  // 系统信息查看权限
			},
		},
		{
			ID:          "viewer", // 查看者角色
			Name:        "查看者",
			Description: "只能查看资源",
			Permissions: []string{
				"containers:list", // 容器列表查看权限
				"containers:read", // 容器详情查看权限
				"images:list",     // 镜像列表查看权限
				"images:read",     // 镜像详情查看权限
				"networks:list",   // 网络列表查看权限
				"networks:read",   // 网络详情查看权限
				"system:read",     // 系统信息查看权限
			},
		},
	}

	// 将角色添加到角色映射表中
	for _, role := range roles {
		s.roles[role.ID] = role
	}
}

// generateUserID 生成唯一的用户ID
// 返回值 string: 格式为"user_时间戳纳秒"的唯一ID
func generateUserID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

// Logout 注销当前Token（加入黑名单）
// 参数 tokenString: 要注销的JWT令牌字符串
// 返回值 error: 可能的错误
func (s *AuthService) Logout(tokenString string) error {
	// 检查令牌字符串是否为空
	if strings.TrimSpace(tokenString) == "" {
		return fmt.Errorf("无效的 Token")
	}

	// 默认设置24小时后过期，尽量使用Token自身的exp过期时间
	expTs := time.Now().Add(24 * time.Hour).Unix()

	// 尝试解析令牌以获取真实的过期时间
	if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	}); err == nil {
		// 从令牌声明中提取过期时间
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if expVal, ok := claims["exp"]; ok {
				// 根据不同的数据类型转换过期时间
				switch v := expVal.(type) {
				case float64:
					expTs = int64(v)
				case int64:
					expTs = v
				case json.Number:
					if n, err := v.Int64(); err == nil {
						expTs = n
					} else {
						logger.Error("解析 Token 过期时间失败: %v", err)
					}
				}
			}
		}
	} else {
		logger.Error("解析 Token 失败: %v", err)
	}

	// 将令牌添加到黑名单中
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokenBlacklist[tokenString] = expTs
	return nil
}

// isTokenRevoked 检查令牌是否已被注销（在黑名单中）
// 参数 tokenString: 待检查的JWT令牌字符串
// 返回值 bool: 令牌是否已被注销
func (s *AuthService) isTokenRevoked(tokenString string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查黑名单映射表是否存在
	if s.tokenBlacklist == nil {
		return false
	}

	// 检查令牌是否在黑名单中且未过期
	if exp, ok := s.tokenBlacklist[tokenString]; ok {
		if time.Now().Unix() < exp {
			return true
		}
	}
	return false
}

// permissionMatch 检查权限模式是否匹配指定的资源和操作
// 参数 pattern: 权限模式（如"*:*"、"containers:*"等）
// 参数 resource: 资源名称
// 参数 action: 操作名称
// 返回值 bool: 是否匹配
func (s *AuthService) permissionMatch(pattern, resource, action string) bool {
	// 通配符"*:*"匹配所有权限
	if pattern == "*:*" {
		return true
	}

	// 将模式按":"分割为资源和操作部分
	parts := strings.Split(pattern, ":")
	if len(parts) != 2 {
		return false
	}

	// 分别获取模式中的资源和操作
	pr, pa := parts[0], parts[1]

	// 检查资源是否匹配（"*"表示匹配所有资源）
	if pr != "*" && pr != resource {
		return false
	}

	// 检查操作是否匹配（"*"表示匹配所有操作）
	if pa != "*" && pa != action {
		return false
	}

	return true
}
