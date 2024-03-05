package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// getHostname returns the host name reported by the kernel.
func getHostname(c *gin.Context) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"hostname": name})
}

// getHealthStatus returns the health status of your API.
func getHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// ping quick check to verify API status.
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func getAppHostname(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return string(body), err
}

func app1(c *gin.Context) {
	url := "https://poc-app-2-cloudhub-eu-west-1.dev.hydra.sophos.com/hostname"
	resp, err := getAppHostname(url)

	if err != nil {
		fmt.Println("Error:", err)
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func app2(c *gin.Context) {
	url := "https://poc-app-2-cloudhub-eu-west-1.dev.hydra.sophos.com/hostname"
	resp, err := getAppHostname(url)

	if err != nil {
		fmt.Println("Error:", err)
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func app3(c *gin.Context) {
	url := "https://poc-app-3-cloudhub-eu-west-1.dev.hydra.sophos.com/hostname"
	resp, err := getAppHostname(url)

	if err != nil {
		fmt.Println("Error:", err)
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func mainRouter() http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/hostname", getHostname)
	engine.GET("/ping", ping)
	engine.GET("/app1", app1)
	engine.GET("/app2", app2)
	engine.GET("/app3", app3)
	return engine
}

func healthRouter() http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/health", getHealthStatus)
	return engine
}

func main() {
	mainServer := &http.Server{
		Addr:         ":8080",
		Handler:      mainRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	healthServer := &http.Server{
		Addr:         ":8081",
		Handler:      healthRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := mainServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := healthServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
