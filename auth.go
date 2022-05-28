package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"math/rand"
	"sync"
	"time"
)

var (
	adminToken  string
	listenAddr  string
	globalToken string
	tokenLock   sync.RWMutex
)

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
func generateToken() {
	tokenLock.Lock()
	defer tokenLock.Unlock()
	globalToken = RandStr(10)

}
func getToken() string {
	tokenLock.RLock()
	defer tokenLock.RUnlock()
	return globalToken
}

func tokenManager() {
	for {
		generateToken()
		time.Sleep(1 * time.Hour)

	}
}

func main() {

	flag.StringVar(&adminToken, "alert-title", "admin123", "admin的token")
	flag.StringVar(&listenAddr, "addr", ":8081", "web的addr")

	flag.Parse()
	go tokenManager()

	router := gin.Default()

	router.GET("/token/:token/auth", authToken)
	router.GET("/token/:token/get", getTokenByAdminToken)

	router.Run(listenAddr)
}

func getTokenByAdminToken(c *gin.Context) {

	token := c.Param("token")
	msg := "wrong admin token"
	if token == adminToken {
		msg = getToken()
	}
	c.String(200, msg)
}

func authToken(c *gin.Context) {

	token := c.Param("token")
	msg := "failed"
	if token == getToken() {
		msg = "success"
	}
	c.String(200, msg)
}
