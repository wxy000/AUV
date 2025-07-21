package handlers

import (
	"AUV/common/response"
	"AUV/config"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHitokoto(c *gin.Context) {
	// 读取文本文件
	content, err := os.ReadFile(config.Cfg.HitokotoFile)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "无法读取数据文件")
		return
	}

	// 解析请求参数
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 || count > 10 {
		response.Fail(c, http.StatusBadRequest, "参数count应为1-10之间的整数")
		return
	}

	// 处理数据
	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 {
		response.Fail(c, http.StatusNotFound, "没有可用的数据")
		return
	}

	// 随机选择
	results := make([]string, 0, count)
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		idx := localRand.Intn(len(lines))
		results = append(results, strings.TrimSpace(lines[idx]))
	}

	response.Success(c, gin.H{
		"data":  results,
		"count": len(results),
	})
}
