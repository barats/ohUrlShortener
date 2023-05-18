// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

// APINewAdmin
//
// Add new admin user
func APINewAdmin(ctx *gin.Context) {
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if utils.EmptyString(account) || utils.EmptyString(password) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("用户名或密码不能为空"))
		return
	}

	if len(password) < 8 {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("密码长度最少8位"))
		return
	}

	err := service.NewUser(account, password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, core.ResultJsonSuccess())
}

// APIAdminUpdate
//
// Update password of given admin user
func APIAdminUpdate(ctx *gin.Context) {
	account := ctx.Param("account")
	password := ctx.PostForm("password")

	if utils.EmptyString(account) || utils.EmptyString(password) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("用户名或密码不能为空"))
		return
	}

	if len(password) < 8 {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("密码长度最少8位"))
		return
	}

	err := service.UpdatePassword(strings.TrimSpace(account), password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("修改失败"))
		return
	}

	ctx.JSON(http.StatusOK, core.ResultJsonSuccess())
}

// APIGenShortUrl Generate new short url
func APIGenShortUrl(ctx *gin.Context) {
	url := ctx.PostForm("dest_url")
	memo := ctx.PostForm("memo")
	strOpenType := ctx.PostForm("open_type")
	openType, err := strconv.Atoi(strOpenType)
	if err != nil {
		openType = int(core.OpenInAll)
	}

	if utils.EmptyString(strings.TrimSpace(url)) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("dest_url 不能为空"))
		return
	}

	res, err := service.GenerateShortUrl(strings.TrimSpace(url), strings.TrimSpace(memo), openType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest(err.Error()))
		return
	}

	json := map[string]string{
		"short_url": fmt.Sprintf("%s%s", utils.AppConfig.UrlPrefix, res),
	}
	ctx.JSON(http.StatusOK, core.ResultJsonSuccessWithData(json))
}

// APIUrlInfo Get Short Url Stat Info.
func APIUrlInfo(ctx *gin.Context) {
	url := ctx.Param("url")
	if utils.EmptyString(strings.TrimSpace(url)) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("url 不能为空"))
		return
	}

	stat, err := service.GetShortUrlStats(strings.TrimSpace(url))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, core.ResultJsonError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, core.ResultJsonSuccessWithData(stat))
}

// APIUpdateUrl Enable or Disable Short Url
func APIUpdateUrl(ctx *gin.Context) {
	url := ctx.Param("url")
	enableStr := ctx.PostForm("enable")
	if utils.EmptyString(strings.TrimSpace(url)) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("url 不能为空"))
		return
	}

	enable, err := strconv.ParseBool(enableStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("enable 参数值非法"))
		return
	}

	res, err := service.ChangeState(url, enable)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, core.ResultJsonSuccessWithData(res))
}

// APIDeleteUrl Delete Short Url
func APIDeleteUrl(ctx *gin.Context) {
	url := ctx.Param("url")
	if utils.EmptyString(strings.TrimSpace(url)) {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest("url 不能为空"))
		return
	}
	err := service.DeleteUrlAndAccessLogs(url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, core.ResultJsonSuccess)
}
