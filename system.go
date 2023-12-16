package main

import (
	"net/http"
	"system/network"
	"system/power"
	"system/update"
	"system/wireless"

	"github.com/gin-gonic/gin"
)

func getInterfaces(c *gin.Context) {
	ifaces, err := network.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, ifaces.Values())
}

func setInterfaces(c *gin.Context) {
	var iface_list []*network.Interface

	c.BindJSON(&iface_list)
	ifaces := network.FromList(iface_list)

	err := ifaces.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, ifaces)
}

func getInterface(c *gin.Context) {
	name := c.Param("name")

	ifaces, err := network.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	iface := ifaces.Get(name)
	if iface == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid interface name"})
		return
	}

	c.JSON(http.StatusOK, iface)
}

func setInterface(c *gin.Context) {
	name := c.Param("name")
	var iface network.Interface
	c.BindJSON(&iface)
	ifaces, err := network.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if iface.Name != name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name does not match"})
		return
	}

	var status int
	if ifaces.Has(name) {
		status = http.StatusOK
	} else {
		status = http.StatusCreated
	}

	ifaces.Add(&iface)
	err = ifaces.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(status, iface)
}

func getUpdates(c *gin.Context) {
	update.CheckUpdates()
}

func doUpdate(c *gin.Context) {

}

func getWireless(c *gin.Context) {
	wifi, err := wireless.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, wifi)
}

func setWireless(c *gin.Context) {
}

func doShutdown(c *gin.Context) {
	err := power.Shutdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "shutting down"})
}

func doRestart(c *gin.Context) {
	err := power.Restart()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "rebooting"})
}

func main() {
	r := gin.Default()

	r.GET("/network", getInterfaces)
	r.GET("/network/:name", getInterface)

	r.PUT("/network", setInterfaces)
	r.PUT("/network/:name", setInterface)

	r.GET("/update/available", getUpdates)
	r.POST("/update", doUpdate)

	r.GET("/wireless", getWireless)
	r.PUT("/wireless", setWireless)

	r.POST("/power/shutdown", doShutdown)
	r.POST("/power/restart", doRestart)

	r.Run()
}
