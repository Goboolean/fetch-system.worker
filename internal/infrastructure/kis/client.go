package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"net/http"
	"net/url"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/gorilla/websocket"
)

const (
	host = "ops.koreainvestment.com:21000"
	approvalKeyIssueURL = "https://openapi.koreainvestment.com:9443/oauth2/Approval"
	checkHolidayURL     = "https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/chk-holiday"
	tokenIssueURL       = "https://openapivts.koreainvestment.com:29443/oauth2/tokenP"
)

const (
	defaultDataBufSize = 10000
	defaulSubBufSize  = 3000
	defaultMsgBufSize = 100
	defaultErrBufSize = 100
)

type Client struct {
	conn *websocket.Conn

	approvalKey string
	accessKey   string
	appKey 	    string
	setretKey   string

	mode        string

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	dataCh chan *Trade
	msgCh chan struct{}
	subCh chan string
	errCh chan error
}

func New(c *resolver.ConfigMap) (*Client, error) {

	u := url.URL{
		Scheme: "ws", Host: host,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	buf, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Client{
		conn:       conn,
		ctx:        ctx,
		cancel:     cancel,
		dataCh:     make(chan *Trade, buf),
		msgCh:      make(chan struct{}, defaultMsgBufSize),
		subCh:      make(chan string, defaulSubBufSize),
		errCh:      make(chan error, defaultErrBufSize),
	}

	appkey, err := c.GetStringKey("APPKEY")
	if err != nil {
		return nil, err
	}

	secretkey, err := c.GetStringKey("SECRET")
	if err != nil {
		return nil, err
	}

	approvalKey, err := instance.GetApprovalKey(ctx, appkey, secretkey)
	if err != nil {
		return nil, err
	}

	instance.approvalKey = approvalKey
	instance.appKey = appkey
	instance.setretKey = secretkey

	mode, flag, err := c.GetStringKeyOptional("MODE")
	if err != nil {
		return nil, err
	}

	if flag {
		switch mode {
		case "PRODUCTION":
			instance.accessKey, err = instance.IssueAccessToken(ctx, appkey, secretkey)
			if err != nil {
				return nil, err
			}
		default:
			break
		}
	}

	instance.wg.Add(1)
	go instance.runReader(instance.ctx, &instance.wg)

	return instance, nil
}


func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	c.cancel()
	c.wg.Wait()
	close(c.dataCh)
	close(c.msgCh)
	close(c.subCh)
	close(c.errCh)
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	select {
	case <- ctx.Done():
		return ctx.Err()
	case <- c.msgCh:
		return nil
	}
}


func (c *Client) tryVacatingMsgCh() {
	if len(c.msgCh) == defaultMsgBufSize {
		for {
			select {
			case <- c.msgCh:
				continue
			default:
				break
			}
		}
	}
}

func (c *Client) runReader(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {

		if err := ctx.Err(); err != nil {
			return
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			time.Sleep(time.Millisecond * 1)
			continue
		}

		if isPingpongMsg(message) {
			c.tryVacatingMsgCh()
			c.msgCh <- struct{}{}
			continue
		}

		if stock, ok := tryParsingToSubResp(message); ok {
			c.subCh <- stock
			continue
		}

		data, err := parseTrade(string(message))
		if err != nil {
			c.errCh <- err
			continue
		}

		for _, d := range data {
			c.dataCh <- d
		}	
	}
}


func (c *Client) GetApprovalKey(ctx context.Context, Appkey string, Secretkey string) (string, error) {

	jsonData, err := json.Marshal(getApprovalKeyReqeust{
		GrantType: "client_credentials",
		AppKey:    Appkey,
		SecretKey: Secretkey,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", approvalKeyIssueURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err	
	}

	req.Header.Set("content-type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(string(body))
	}

	var res getApprovalKeyResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.ApprovalKey == "" {
		return res.ApprovalKey, fmt.Errorf("approval key is empty")
	}

	return res.ApprovalKey, nil
}


func (c *Client) AccessTokenAvailable() bool {
	return c.accessKey != ""
}


func (c *Client) issueAccessToken(ctx context.Context, appkey string, appsecret string) (string, bool, error) {

	var retryable bool = false

	jsonData, err := json.Marshal(AccessKeyRequestBodyJson{
		GrantType: "client_credentials",
		AppKey:    appkey,
		AppSecret: appsecret,
	})
	if err != nil {
		return "", retryable, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenIssueURL, bytes.NewBuffer(jsonData))
	req.Header.Set("content-type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", retryable, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", retryable, err
	}

	if resp.StatusCode != http.StatusOK {
		var res AccessKeyErrorResponseJson
		if err := json.Unmarshal(body, &res); err != nil {
			return "", retryable, err
		}

		if res.ErrorCode == "EGW00133" {
			retryable = true
		}
		return "", retryable, fmt.Errorf(string(body))
	}

	var res AccessKeyResponseBodyJson
	if err := json.Unmarshal(body, &res); err != nil {
		return "", retryable, err
	}

	if res.AccessToken == "" {
		return "", retryable, fmt.Errorf("access token is empty")
	}

	return res.AccessToken, retryable, nil
}

func (c *Client) IssueAccessToken(ctx context.Context, appkey string, appsecret string) (string, error){
	token, retryable, err := c.issueAccessToken(ctx, appkey, appsecret)
	if err != nil {
		if retryable {
			select {
			case <- ctx.Done():
				return "", errors.Join(ctx.Err(), err)
			case <- time.After(time.Second * 10):
				return c.IssueAccessToken(ctx, appkey, appsecret)
			}
		}
	}

	return token, err
}



func (c *Client) Subscribe(ctx context.Context, stocks ...string) (<-chan *Trade, error) {

	hangingStockList := make(map[string]struct{})
	for _, stock := range stocks {
		hangingStockList[stock] = struct{}{}
	}

	for _, stock := range stocks {
		if err := c.subscribeProduct(stock); err != nil {
			return nil, err
		}
	}

	for {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case stock := <- c.subCh:
			delete(hangingStockList, stock)
			if len(hangingStockList) == 0 {
				return c.dataCh, nil
			}
		}
	}
}


func (c *Client) subscribeProduct(symbol string) error {
	req := &RequestJson{
		Header: HeaderJson{
			ApprovalKey: c.approvalKey,
			Custtype:    custtype,
			TrType:      tr_type_subscribe,
			ContentType: "utf-8",
		},
		Body: RequestBodyJson{
			Input: RequestInput{
				TrId:  "H0STCNT0", // "HDFSCNT0",
				TrKey: symbol,
			},
		},
	}

	if err := c.conn.WriteJSON(req); err != nil {
		return err
	}

	return nil
}



func (c *Client) IsMarketOn(ctx context.Context) (bool, error) {

	kst := time.FixedZone("KST", 9*60*60)
	now := time.Now().In(kst)

	date := now.Format("20060102")
	url := fmt.Sprintf("%s?BASS_DT=%s&CTX_AREA_NK=&CTX_AREA_FK=", checkHolidayURL, date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.accessKey))
	req.Header.Set("appkey", c.appKey)
	req.Header.Set("appsecret", c.setretKey)
	req.Header.Set("tr_id", "CTCA0903R")
	req.Header.Set("custtype", custtype)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.Join(fmt.Errorf(string(body)), err) 
	}

	var res CheckHolidayResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return false, err
	}

	if len(res.Output) == 0 {
		return false, fmt.Errorf("output is empty")
	}

	if res.Output[0].BzdyYn == "N" || res.Output[0].TrDayYn == "N" || res.Output[0].OpndYn == "N" {
		return false, nil
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, kst)
    end := time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, kst)

	return (now.After(start) && now.Before(end)), nil
}