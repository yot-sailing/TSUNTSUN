package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/yot-sailing/TSUNTSUN/body"
	"github.com/yot-sailing/TSUNTSUN/interfaces/controllers"
	authMiddleware "github.com/yot-sailing/TSUNTSUN/middleware"
)

func Init() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	userController := controllers.NewUserController(NewSqlHandler())
	tsundokuController := controllers.NewTsundokuController(NewSqlHandler())
	tagController := controllers.NewTagController(NewSqlHandler())
	tsundokuTagController := controllers.NewTsundokuTagController(NewSqlHandler())

	// Middleware
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: logFormat(),
		Output: os.Stdout,
	})
	e.Use(logger)
	e.Use(middleware.Recover())

	// 接続テスト
	e.GET("/api/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "This is test!")
	})

	// LINE
	// ログイン
	e.POST("/api/line_login", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}

		userExcludeLine := body.UesrExcludeLine{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		return c.JSON(http.StatusOK, userExcludeLine)
	})

	// ログアウト
	e.POST("/api/line_logout", func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")

		revokeRequestBody := &body.RevokeRequestBody{
			ClientID:      os.Getenv("CHANNEL_ID"),
			ClientSercret: os.Getenv("CHANNEL_SECRET"),
			AccessToken:   accessToken[7:],
		}

		revokeJsonString, err := json.Marshal(revokeRequestBody)
		if err != nil {
			fmt.Println(err)
		}

		endpoint := "https://api.line.me/oauth2/v2.1/revoke"
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(revokeJsonString))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		// // 成功していたら空
		// byteArray, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		return c.String(resp.StatusCode, "logout")
	})

	// ユーザー全取得
	e.GET("/api/users", func(c echo.Context) error {
		users := userController.GetUser()
		c.Bind(&users)
		return c.JSON(http.StatusOK, users)
	})

	// ユーザー作成
	e.POST("/api/users", func(c echo.Context) error {
		userController.Create(c)
		return c.String(http.StatusOK, "created")
	})

	// ユーザー削除
	e.DELETE("/api/users", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}
		userController.Delete(user.ID)
		return c.String(http.StatusOK, "deleted")
	})

	// 積読全取得
	e.GET("api/tsundokus", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}

		tsundokus := tsundokuController.GetTsundoku(user.ID)
		c.Bind(&tsundokus)
		for i, tsundoku := range tsundokus {
			tsundokuTags := tsundokuTagController.GetTsundokuTagsByTsundokuIDandUserID(tsundoku.ID, user.ID)
			var tagIDs []int
			for _, tsundokuTag := range tsundokuTags {
				tagIDs = append(tagIDs, tsundokuTag.TagID)
			}
			tsundokus[i].Tags = tagController.GetTags(tagIDs)
		}
		return c.JSON(http.StatusOK, tsundokus)
	})

	// 積読追加
	e.POST("api/tsundokus", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}
		tsundokuController.CreateTsundoku(c, user.ID)
		return c.String(http.StatusOK, "created tsundoku")
	})

	// 積読削除
	// TODO:ユーザーが管理しているかの判定
	e.DELETE("api/tsundokus/:tsundokuID", func(c echo.Context) error {
		str_tsundokuID := c.Param("tsundokuID")
		tsundokuID, err := strconv.Atoi(str_tsundokuID)
		if err != nil {
			fmt.Println(err)
		}
		tsundokuController.Delete(tsundokuID)
		return c.String(http.StatusOK, "deleted tsundoku")
	})

	// ある時間以内に読める本を取得
	e.GET("api/time/:time", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}
		total_min, _ := strconv.Atoi(c.Param("time"))
		tsundokus := tsundokuController.GetFreeTsundoku(c, user.ID, total_min)
		c.Bind(&tsundokus)
		return c.JSON(http.StatusOK, tsundokus)
	})

	// ユーザーが管理するタグ全取得
	e.GET("api/tags", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}

		// TsundokuTagテーブルのユーザーの管理下のものを取得
		tsundokuTags := tsundokuTagController.GetTsundokuTags(user.ID)
		var tagIDs []int
		for _, tsundokuTag := range tsundokuTags {
			tagIDs = append(tagIDs, tsundokuTag.TagID)
		}
		// tagIDからtagを取得
		tags := tagController.GetTags(tagIDs)
		// c.Bind(&tags)
		return c.JSON(http.StatusOK, tags)
	})

	// ユーザーが管理する積読についているタグ全取得
	e.GET("api/tsundokus/:tsundokuID/tags", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}

		str_tsundokuID := c.Param("tsundokuID")
		// intに変換
		tsundokuID, err := strconv.Atoi(str_tsundokuID)
		if err != nil {
			return err
		}

		tsundokuTags := tsundokuTagController.GetTsundokuTagsByTsundokuIDandUserID(tsundokuID, user.ID)
		var tagIDs []int
		for _, tsundokuTag := range tsundokuTags {
			tagIDs = append(tagIDs, tsundokuTag.TagID)
		}
		// tagIDからtagを取得
		tags := tagController.GetTags(tagIDs)

		return c.JSON(http.StatusOK, tags)
	})

	// タグ追加
	e.POST("api/tsundokus/:tsundokuID/tags", func(c echo.Context) error {
		user, err := authMiddleware.AuthUser(c.Request().Header.Get("Authorization"), userController)
		if err != nil {
			return err
		}

		str_tsundokuID := c.Param("tsundokuID")
		// intに変換
		tsundokuID, err := strconv.Atoi(str_tsundokuID)
		if err != nil {
			return err
		}

		// Tagsテーブルにレコードを追加
		tagID := tagController.CreateTag(c, user.ID)
		// TsundokuTagsテーブルにレコードを追加
		tsundokuTagController.CreateTsundokuTag(c, tsundokuID, user.ID, tagID)

		return c.JSON(http.StatusOK, "created tag")
	})

	// タグ削除
	// TODO:ユーザーが管理しているかの判定
	e.DELETE("api/tsundokus/:tsundokuID/tags/:tagID", func(c echo.Context) error {
		str_tagID := c.Param("tagID")
		tagID, err := strconv.Atoi(str_tagID)
		if err != nil {
			fmt.Println(err)
		}
		tagController.Delete(tagID)
		return c.String(http.StatusOK, "deleted tag")
	})

	port := os.Getenv("PORT")
	// start server
	e.Logger.Fatal(e.Start(":" + port))
}

func logFormat() string {
	// Refer to https://github.com/tkuchiki/alp
	var format string
	format += "time:${time_rfc3339}\t"
	format += "host:${remote_ip}\t"
	format += "forwardedfor:${header:x-forwarded-for}\t"
	format += "req:-\t"
	format += "status:${status}\t"
	format += "method:${method}\t"
	format += "uri:${uri}\t"
	format += "size:${bytes_out}\t"
	format += "referer:${referer}\t"
	format += "ua:${user_agent}\t"
	format += "reqtime_ns:${latency}\t"
	format += "cache:-\t"
	format += "runtime:-\t"
	format += "apptime:-\t"
	format += "vhost:${host}\t"
	format += "reqtime_human:${latency_human}\t"
	format += "x-request-id:${id}\t"
	format += "host:${host}\n"

	return format
}
