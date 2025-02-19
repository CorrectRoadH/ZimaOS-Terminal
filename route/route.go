package route

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/external"
	"github.com/IceWhaleTech/CasaOS-Common/utils/jwt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "embed"

	"github.com/CorrectRoadH/ZimaOS-Terminal/codegen"
	"github.com/CorrectRoadH/ZimaOS-Terminal/config"
	"github.com/CorrectRoadH/ZimaOS-Terminal/service"
)

type TerminalRouter struct{}

var (
	_swagger *openapi3.T

	APIPath string
)

func init() {
	swagger, err := codegen.GetSwagger()
	if err != nil {
		panic(err)
	}

	_swagger = swagger

	u, err := url.Parse(_swagger.Servers[0].URL)
	if err != nil {
		panic(err)
	}

	APIPath = strings.TrimRight(u.Path, "/")
}

func GetRouter() http.Handler {
	hello := NewHelloWorld()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.POST, echo.GET, echo.OPTIONS, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderXCSRFToken, echo.HeaderContentType, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderAccessControlAllowMethods, echo.HeaderConnection, echo.HeaderOrigin, echo.HeaderXRequestedWith},
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders},
		MaxAge:           172800,
		AllowCredentials: true,
	}))

	e.Use(middleware.Gzip())

	e.Use(middleware.Logger())

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return c.RealIP() == "::1" || c.RealIP() == "127.0.0.1"
		},
		ParseTokenFunc: func(token string, c echo.Context) (interface{}, error) {
			valid, claims, err := jwt.Validate(token, func() (*ecdsa.PublicKey, error) { return external.GetPublicKey(config.CommonInfo.RuntimePath) })
			if err != nil || !valid {
				return nil, echo.ErrUnauthorized
			}

			c.Request().Header.Set("user_id", strconv.Itoa(claims.ID))

			return claims, nil
		},
		TokenLookupFuncs: []middleware.ValuesExtractor{
			func(c echo.Context) ([]string, error) {
				return []string{c.Request().Header.Get(echo.HeaderAuthorization)}, nil
			},
		},
	}))

	codegen.RegisterHandlersWithBaseURL(e, hello, APIPath)

	return e
}

func NewHelloWorld() codegen.ServerInterface {
	return &TerminalRouter{}
}

func (h *TerminalRouter) Ping(ctx echo.Context) error {
	ping := service.HelloWorld.Ping()

	return ctx.JSON(http.StatusOK, codegen.PongOK{
		Data: &ping,
	})
}

func (h *TerminalRouter) OpenTerminal(ctx echo.Context) error {
	// 随机生成一个20000-60000之间的端口号
	port := rand.Intn(40000) + 20000

	// 执行ttyd命令
	cmd := exec.Command("ttyd", "-p", fmt.Sprintf("%d", port), "-W", "-o", "env", "SHELL=/bin/bash", "/usr/bin/btop")
	if err := cmd.Start(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "无法启动ttyd")
	}

	// 返回端口号给客户端
	return ctx.JSON(http.StatusOK, codegen.OpenTerminalOK{
		Port: &port,
	})
}
