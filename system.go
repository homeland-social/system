package main

import (
	"net/http"
	"system/network"
	"system/update"

	"github.com/gin-gonic/gin"
)

func getInterfaces(c *gin.Context) {
	ifaces, err := network.ParseInterfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open network configuration"})
		return
	}

	iface_list := make([]network.Interface, 0, len(ifaces))
	for _, v := range ifaces {
		iface_list = append(iface_list, *v)
	}

	c.JSON(http.StatusOK, iface_list)
}

func setInterfaces(c *gin.Context) {

}

func getInterface(c *gin.Context) {
	name := c.Param("name")

	ifaces, err := network.ParseInterfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open network configuration"})
		return
	}

	iface := ifaces[name]
	if iface == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid interface name"})
		return
	}

	c.JSON(http.StatusOK, iface)
}

func setInterface(c *gin.Context) {

}

func getUpdates(c *gin.Context) {
	update.CheckUpdates()
}

func main() {
	r := gin.Default()

	r.GET("/network", getInterfaces)
	r.GET("/network/:name", getInterface)

	r.POST("/network", setInterfaces)
	r.POST("/network/:name", setInterface)

	r.GET("/update/available", getUpdates)

	r.Run()
}
