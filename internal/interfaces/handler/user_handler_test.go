package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gva/internal/domain/entity"
	"gva/internal/domain/service"
	"gva/internal/infrastructure/cache"
	"gva/internal/infrastructure/database"
	"gva/internal/infrastructure/repository"
	"gva/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 切换到项目根目录
	if err := os.Chdir("../../../"); err != nil {
		log.Fatalf("Failed to change dir: %v", err)
	}
}

func setupTestRouter() (*gin.Engine, error) {
	// 获取测试数据库连接
	db, err := testutil.GetTestDB()
	if err != nil {
		return nil, err
	}

	// 执行数据库迁移
	if err := database.AutoMigrate(db); err != nil {
		return nil, err
	}

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// 初始化路由
	r := gin.Default()
	userRepo := repository.NewUserRepository(db)
	userCache := cache.NewUserCache(rdb)
	userService := service.NewUserService(userRepo, db, userCache)
	userHandler := NewUserHandler(userService)

	// 注册路由
	r.POST("/api/v1/register", userHandler.Register)
	r.POST("/api/v1/login", userHandler.Login)
	r.PUT("/api/v1/user/profile", userHandler.UpdateProfile)
	r.POST("/api/v1/user/reset-password", userHandler.ResetPassword)
	r.POST("/api/v1/user/avatar", userHandler.UploadAvatar)

	return r, nil
}

func TestUserFlow(t *testing.T) {
	r, err := setupTestRouter()
	assert.NoError(t, err)

	// 1. 测试注册
	t.Run("Register", func(t *testing.T) {
		testUser := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(testUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Register response: %s", w.Body.String())

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "注册成功", response["message"])

		// 验证用户是否真的创建成功
		db, err := testutil.GetTestDB()
		assert.NoError(t, err)
		var user entity.User
		err = db.Where("username = ?", "testuser").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
	})

	// 2. 测试登录成功
	t.Run("Login Success", func(t *testing.T) {
		testUser := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(testUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		t.Logf("Login response: %s", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "登录成功", response["message"])
		assert.NotEmpty(t, response["token"])
		assert.NotNil(t, response["user"])
	})

	// 3. 测试错误密码登录
	t.Run("Login with Wrong Password", func(t *testing.T) {
		testUser := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(testUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "用户名或密码错误", response["error"])
	})
}

func TestExportUsers(t *testing.T) {
	r, err := setupTestRouter()
	assert.NoError(t, err)

	// 1. 先创建一些测试用户
	testUsers := []struct {
		username string
		nickname string
		email    string
		phone    string
	}{
		{"test1", "测试用户1", "test1@example.com", "13800138001"},
		{"test2", "测试用户2", "test2@example.com", "13800138002"},
		{"test3", "测试用户3", "test3@example.com", "13800138003"},
	}

	for _, u := range testUsers {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"username": u.username,
			"password": "password123",
		})
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 2. 登录获取token
	w := httptest.NewRecorder()
	loginBody, _ := json.Marshal(map[string]string{
		"username": "test1",
		"password": "password123",
	})
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	token := loginResp["token"].(string)

	// 3. 测试导出功能
	t.Run("Export Users", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/export", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;filename=users.csv")

		// 验证CSV内容
		lines := strings.Split(w.Body.String(), "\n")
		assert.True(t, len(lines) > 1) // 至少有表头和一行数据

		// 验证表头
		header := lines[0]
		assert.Contains(t, header, "ID")
		assert.Contains(t, header, "用户名")
		assert.Contains(t, header, "昵称")
		assert.Contains(t, header, "邮箱")
		assert.Contains(t, header, "手机号")
		assert.Contains(t, header, "状态")
		assert.Contains(t, header, "创建时间")

		// 验证数据行
		for i, u := range testUsers {
			assert.Contains(t, lines[i+1], u.username)
			assert.Contains(t, lines[i+1], u.nickname)
			assert.Contains(t, lines[i+1], u.email)
			assert.Contains(t, lines[i+1], u.phone)
			assert.Contains(t, lines[i+1], "正常") // 默认状态
		}
	})

	// 4. 测试未认证的情况
	t.Run("Export Users Without Auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/export", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
