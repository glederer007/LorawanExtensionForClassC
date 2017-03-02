package main

import (
	"bufio"
	"os"
  	//import the Paho Go MQTT library
  	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
        "crypto/tls"
        "io/ioutil"
        "bytes"
        "github.com/brocaar/lorawan"
	"encoding/hex"
	"strconv"
	"strings"
//	"encoding/base64"
//	"net/url"
)

// const LoraServer="https://192.168.0.100:8888"
// const MQTTServer="tcp://192.168.0.100:1883"

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
  fmt.Printf("TOPIC: %s\n", msg.Topic())
  fmt.Printf("MSG: %s\n", msg.Payload())
}

type ApiGetNodeSessionResponse struct {
	DevAddr string `json:"devAddr"`
	AppEUI string `json:"appEUI"`
	DevEUI string `json:"devEUI"`
	AppSKey string `json:"appSKey"`
	NwkSKey string `json:"nwkSKey"`
	FCntUp uint32 `json:"fCntUp"`
	FCntDown uint32 `json:"fCntDown"`
	RxDelay int `json:"rxDelay"`
	Rx1DROffset int `json:"rx1DROffset"`
	CFList []int `json:"cFList"`
	RxWindow string `json:"rxWindow"`
	Rx2DR int `json:"rx2DR"`
	RelaxFCnt bool `json:"relaxFCnt"`
	AdrInterval int `json:"adrInterval"`
	InstallationMargin int `json:"installationMargin"`
	NbTrans int `json:"nbTrans"`
	TxPower int `json:"txPower"`
}

type ApiGetNodes struct {
	Result []struct {
		AdrInterval int `json:"adrInterval"`
		AppEUI string `json:"appEUI"`
		AppKey string `json:"appKey"`
		ChannelListID string `json:"channelListID"`
		DevEUI string `json:"devEUI"`
		InstallationMargin int `json:"installationMargin"`
		Name string `json:"name"`
		RelaxFCnt bool `json:"relaxFCnt"`
		Rx1DROffset int `json:"rx1DROffset"`
		Rx2DR int `json:"rx2DR"`
		RxDelay int `json:"rxDelay"`
		RxWindow string `json:"rxWindow"`
	} `json:"result"`
	TotalCount string `json:"totalCount"`
}

func main() {

//input Lora server and MQTT server addresses 

        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter Lora server address (https://192.168.1.100:8888): ")
        sel, _ := reader.ReadString('\n')
//      loraServer := strings.TrimRight(sel,"\n")   // for linux
        loraServer := strings.TrimRight(sel,"\r\n")  // for windows
        if (loraServer =="")  {
                fmt.Println("using default address: https://192.168.1.100:8888")
                loraServer = "https://192.168.1.100:8888"
        }
        fmt.Println("selected Lora Server: ",loraServer)

        reader = bufio.NewReader(os.Stdin)
        fmt.Print("Enter MQTT server address (tcp://192.168.1.100:1883): ")
        sel, _ = reader.ReadString('\n')
//      mQTTServer := strings.TrimRight(sel,"\n")   // for linux
        mQTTServer := strings.TrimRight(sel,"\r\n")  // for windows
        if (mQTTServer =="")  {
                fmt.Println("using default address: tcp://192.168.1.100:1883")
                mQTTServer = "tcp://192.168.1.100:1883"
        }
        fmt.Println("selected MQTT Server: ",mQTTServer)


//nodes start

  
  url := fmt.Sprintf(loraServer+"/api/node?limit=20")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

  tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()
  	var recordOfNodes ApiGetNodes

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&recordOfNodes); err != nil {
		log.Println(err)
	}
	nodesTotalCount , _ := strconv.Atoi(recordOfNodes.TotalCount)
	fmt.Println("TotalCount = ", nodesTotalCount)
	for i := 0; i < nodesTotalCount; i++ {
		fmt.Println("(",i+1,")", " DevEUI = ", recordOfNodes.Result[i].DevEUI)
	}
	fmt.Println("( 0 ) is exit")

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter selection: ")
	sel, _ = reader.ReadString('\n')
	selectedNode, _ := strconv.Atoi(strings.TrimRight(sel,"\r\n")) // for windows
//	selectedNode, _ := strconv.Atoi(strings.TrimRight(sel,"\n")) // for linux
	if selectedNode == 0 {
		fmt.Println("no selection")
		os.Exit(0)
	}
//	fmt.Println(selectedNode)
	selectedNode=selectedNode-1 //index starts with 0


//nodes end

//get data from selected node start

	url = fmt.Sprintf(loraServer+"/api/nodeSession/"+recordOfNodes.Result[selectedNode].DevEUI)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

  tr = &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client = &http.Client{Transport: tr}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()
	var record2 ApiGetNodeSessionResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record2); err != nil {
		log.Println(err)
	}
        if record2.DevEUI == "" {
                fmt.Println("There is no session for selected DevEUI")
                os.Exit(0)
        }
 
	fmt.Println("devEUI = ", record2.DevEUI)
	fmt.Println("devAddr = ", record2.DevAddr)
	fmt.Println("nwkSKey = ", record2.NwkSKey)
	
//get data from selecetd node end

//input port and data start

        reader = bufio.NewReader(os.Stdin)
        fmt.Print("Enter port number (1-200): ")
        sel, _ = reader.ReadString('\n')
//      port, _ := strconv.Atoi(strings.TrimRight(sel,"\n"))  // for linux
        port, _ := strconv.Atoi(strings.TrimRight(sel,"\r\n"))  // for windows
        if (port < 1) || (port > 200) {
                fmt.Println("wrong port number")
                os.Exit(0)
        }
        fmt.Println("selected port: ",port)

        reader = bufio.NewReader(os.Stdin)
        fmt.Print("Enter message in bytes (max. 52 bytes): ")
        sel, _ = reader.ReadString('\n')
//      msg, err2 := hex.DecodeString(strings.TrimRight(sel,"\n"))  // for linux
	msg, err2 := hex.DecodeString(strings.TrimRight(sel,"\r\n"))  // for windows

	if err2 != nil  {
                fmt.Println("wrong message")
                os.Exit(0)
        }
        fmt.Println("message: ",msg)


//input port and data end

// prepare data start

var nwkSKey0,appSKey0 [16]byte
var devAddr0 [4]byte
nwkSKey3,_ := hex.DecodeString(record2.NwkSKey)
appSKey3,_ := hex.DecodeString(record2.AppSKey)
devAddr3,_ := hex.DecodeString(record2.DevAddr)

copy(nwkSKey0[:], nwkSKey3[0:16])
copy(appSKey0[:], appSKey3[0:16])
copy(devAddr0[:], devAddr3[0:4])
fPort0 := uint8(port)
//record2.FCntDown = record2.FCntDown + 1
record2.FCntDown ++

fmt.Println("counter down plus", record2.FCntDown)
// prepare data end



//fmt.Println("nwks1: ",nwkSKey1)

//nwkSKey := [16]byte{111, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
//appSKey := [16]byte{161, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
//fPort := uint8(99)


//fmt.Println("nwks0: ",nwkSKey)


phy := lorawan.PHYPayload{
    MHDR: lorawan.MHDR{
        MType: lorawan.UnconfirmedDataDown,
        Major: lorawan.LoRaWANR1,
    },
    MACPayload: &lorawan.MACPayload{
        FHDR: lorawan.FHDR{
            DevAddr: lorawan.DevAddr(devAddr0),
            FCtrl: lorawan.FCtrl{
                ADR:       false,
                ADRACKReq: false,
                ACK:       false,
            },
            FCnt:  record2.FCntDown,
            FOpts: []lorawan.MACCommand{}, // you can leave this out when there is no MAC command to send
        },
        FPort:      &fPort0,
        FRMPayload: []lorawan.Payload{&lorawan.DataPayload{Bytes: msg}},
    },
}

if err := phy.EncryptFRMPayload(appSKey0); err != nil {
    panic(err)
}

if err := phy.SetMIC(nwkSKey0); err != nil {
    panic(err)
}

str, err := phy.MarshalText()
if err != nil {
    panic(err)
}

loraBytes, err := phy.MarshalBinary()
if err != nil {
    panic(err)
}

fmt.Println(string(str))
fmt.Println(loraBytes)

//write back start


	b4, err4 := json.Marshal(record2)
	if err4 != nil {
		fmt.Println("error:", err4)
	}
//	fmt.Println("JSON:", string(b4))
	req4, err4 := http.NewRequest("PUT", url, bytes.NewBuffer(b4))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Accept", "application/json")

    	client = &http.Client{Transport: tr}

	resp, err = client.Do(req4)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

    	defer resp.Body.Close()

    	fmt.Println("response Status:", resp.Status)
    	fmt.Println("response Headers:", resp.Header)
    	body, _ := ioutil.ReadAll(resp.Body)
    	fmt.Println("response Body:", string(body))

//write back ... end

// start mqtt

  //create a ClientOptions struct setting the broker address, clientid, turn
  //off trace output and set the default message handler
  opts := MQTT.NewClientOptions().AddBroker(mQTTServer)
  opts.SetClientID("class_c")
  opts.SetDefaultPublishHandler(f)

  //create and start a client using the above ClientOptions
  c := MQTT.NewClient(opts)
  if token := c.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  //subscribe to the topic /go-mqtt/sample and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := c.Subscribe("gateway/68c90bffffece086/tx", 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }

//create message

//mqttMsg := 
	ts:= string(str) //ts := "YG+jqwcAAQBjAjXuKvdE9g=="
	ms := "{\"phyPayload\": \"%s\",\"txInfo\": {\"codeRate\": \"4/5\",\"dataRate\": {\"bandwidth\": 125,\"modulation\": \"LORA\",\"spreadFactor\": 12 }, \"frequency\": 869525000,\"immediately\": true, \"mac\": \"68c90bffffece086\", \"power\": 14}}"
	mqttMsg :=fmt.Sprintf(ms, ts)
	fmt.Println("mqtt msg: ", mqttMsg)
//	fmt.Println("the end")


  //Publish a message to /go-mqtt/sample at qos 1 and wait for the receipt
  //from the server after sending each message
  
//    text := fmt.Sprintf("this is msg #%d!", 1)
    token := c.Publish("gateway/68c90bffffece086/tx", 0, false, mqttMsg)
    token.Wait()
  

  time.Sleep(3 * time.Second)

  //unsubscribe from /go-mqtt/sample
  if token := c.Unsubscribe("gateway/68c90bffffece086/tx"); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }

  c.Disconnect(250)

// end mqtt

}

