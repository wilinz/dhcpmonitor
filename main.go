package main

import (
	_ "encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

// DhcpLease represents a DHCP lease information
type DhcpLease struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Expires  string `json:"expires"`
}

// getDhcpLeases gets the DHCP lease information from uci
func getDhcpLeases() ([]DhcpLease, error) {
	cmd := exec.Command("sh", "-c", "cat /tmp/dhcp.leases")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var leases []DhcpLease
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		lease := DhcpLease{
			Expires:  fields[0],
			MAC:      fields[1],
			IP:       fields[2],
			Hostname: fields[3],
		}
		leases = append(leases, lease)
	}
	return leases, nil
}

func main() {
	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password123", // Replace with your desired username and password
	}))

	authorized.GET("/api/dhcp", func(c *gin.Context) {
		leases, err := getDhcpLeases()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DHCP leases"})
			return
		}
		c.JSON(http.StatusOK, leases)
	})

	r.Run(":8081") // Start the server on port 8080
}
