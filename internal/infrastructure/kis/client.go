package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	keyIssueURL = "https://openapi.koreainvestment.com:9443/oauth2/Approval"
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

	approvalKey, err := instance.GetApprovalKey(appkey, secretkey)
	if err != nil {
		return nil, err
	}

	instance.approvalKey = approvalKey

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

		if tryParsingToPingMsg(message) {
			c.tryVacatingMsgCh()
			c.msgCh <- struct{}{}
			continue
		}

		if stock, ok := tryParsingToSubResp(message); ok {
			c.subCh <- stock
			continue
		}

		data, err := UnmarshalTrade(string(message))
		if err != nil {
			c.errCh <- err
			continue
		}

		for _, d := range data {
			c.dataCh <- d
		}	
	}
}


func (c *Client) GetApprovalKey(Appkey string, Secretkey string) (string, error) {
	data := &getApprovalKeyReqeust{
		GrantType: "client_credentials",
		AppKey:    Appkey,
		SecretKey: Secretkey,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := http.Post(keyIssueURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res getApprovalKeyResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.ApprovalKey == "" {
		return res.ApprovalKey, errors.New("invalid request")
	}

	return res.ApprovalKey, nil
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
			Input: RequestInputJson{
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