package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"we_book/internal/integration/startup"
	"we_book/internal/repository/dao"
	ijwt "we_book/internal/web/jwt"
)

// TestArticleSuite 测试套件
func TestArticleSuite(t *testing.T) {
	suite.Run(t, &ArticleSuite{})
}

// 测试套件
type ArticleSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Date T      `json:"data"`
}

// SetupSuite 在测试套件开始之前执行
func (as *ArticleSuite) SetupSuite() {
	// 在测试套件开始之前执行
	as.server = gin.Default()
	as.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 1,
		})
	})
	as.db = startup.InitDB()
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRouters(as.server)
}

// TearDownSuite 在测试套件结束之后执行
func (as *ArticleSuite) TearDownSuite() {
	as.db.Exec("TRUNCATE TABLE articles")
}

func (as *ArticleSuite) TestEdit() {
	t := as.T()
	testCases := []struct {
		name string

		// 集成测试之前的数据
		before func(t *testing.T)

		// 集成测试校验的数据
		after func(t *testing.T)

		// 预期的输入
		art Article

		// 预期输出的错误码
		wantedErrCode int

		// 预期输出的结果集
		wantedRes Result[int64]
	}{
		{
			name: "test edit",
			before: func(t *testing.T) {

			},

			after: func(t *testing.T) {
				// 从数据库中直接查
				var art dao.Article
				err := as.db.Where("title = ?", "This is a test title").First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Ctime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "This is a test title",
					Content:  "This is a test content",
					AuthorId: 1,
				}, art)
			},
			art: Article{
				Title:   "This is a test title",
				Content: "This is a test content",
			},
			wantedErrCode: 200,
			wantedRes: Result[int64]{
				Date: 1,
				Msg:  "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			as.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantedErrCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantedRes, webRes)
			tc.after(t)
		})
	}
}
