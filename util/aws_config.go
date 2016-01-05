package util

// The Configuration struct encapsulates the parsed data from config.json.
// Use the Initialize function at program start to parse. Reference variables
// from util.Config to get values. 

/* Note that the config.json file must be in the same directory as the executable.*/

import (
	"encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
    SessionRegion string
	Vpc vpc
	PubSubnet subnet
	PvtSubnet subnet
	InetGateway gateway
	NatGateway gateway
	NatEip gateway
}

type subnet struct {
	AvailabilityZone string
	CidrBlock string
	DryRun	bool
	VpcId string
	TagName string
}

type gateway struct {
	Name string
}

type vpc struct {
	InstanceTenancy string
	TagName string
	DryRun	bool
	CidrBlock string
}

var Config Configuration

// Read values from config file
func Initialize() {
	fmt.Println("Initializing config object...")
    file, _ := os.Open("config.json")
    decoder := json.NewDecoder(file)
    err := decoder.Decode(&Config)
    if err != nil {
        fmt.Println("error:", err)
    }	
}