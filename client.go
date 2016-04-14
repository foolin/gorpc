package gorpc

import (
	"net/http"
	"bytes"
	"fmt"
	"time"
	"errors"
	"io/ioutil"
	"encoding/json"
)

type Client struct{
	BaseUrl string
	Secret string
}

func NewClient(baseUrl, secret string) *Client {
	return &Client{baseUrl, secret}
}


func (this *Client) Call(name string, args interface{}, reply interface{}) error  {
	//request header
	reqParams, err := json.Marshal(args)
	if err != nil {
		return err
	}
	reqTimestmap := fmt.Sprintf("%v", time.Now().Unix())
	reqSign := "rpc"
	if this.Secret != ""{
		reqSign = makeSign(reqTimestmap + string(reqParams), this.Secret)
	}
	url := this.BaseUrl
	//request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqParams))
	if err != nil {
		return err
	}
	req.Header.Set("sign", reqSign)
	req.Header.Set("timestamp", reqTimestmap)
	req.Header.Set("action", name)
	//client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//response header
	msg := resp.Header.Get("msg")
	if msg != "" {
		return errors.New(msg)
	}
	//如果不需要返回数据
	if reply == nil{
		return nil
	}
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteBody, reply)
	return err
}

