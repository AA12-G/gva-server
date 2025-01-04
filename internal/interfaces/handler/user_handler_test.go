package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gva/internal/domain/entity"
	"gva/internal/domain/service"
	"gva/internal/infrastructure/database"
	"gva/internal/infrastructure/repository"
	"gva/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
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

	// 初始化路由
	r := gin.Default()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
	userHandler := NewUserHandler(userService)

	// 注册路由
	r.POST("/api/v1/register", userHandler.Register)
	r.POST("/api/v1/login", userHandler.Login)

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
		db, err := database.GetTestDB()
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
