package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"net/http"
	"net/url"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/gorilla/websocket"
)

const (
	host = "ops.koreainvestment.com:21000"
	keyIssueURL = "https://openapi.koreainvestment.com:9443/oauth2/Approval"
)

type Client struct {
	conn *websocket.Conn

	approvalKey string

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	ch chan *Trade
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
		ch: make(chan *Trade, buf),
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

	return instance, nil
}


func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	c.cancel()
	close(c.ch)
	c.wg.Wait()
	return nil
}

func (c *Client) Ping() error {
	return nil
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

	var res *getApprovalKeyResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.ApprovalKey == "" {
		return res.ApprovalKey, errors.New("invalid request")
	}

	return res.ApprovalKey, nil
}



func (c *Client) Subscribe(stocks ...string) (<-chan *Trade, error) {

	hangingStockList := make(map[string]struct{})
	for _, stock := range stocks {
		hangingStockList[stock] = struct{}{}
	}
	finCh := make(chan struct{})

	c.wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			if err := ctx.Err(); err != nil {
				return
			}
	
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				continue
			}

			fmt.Println(string(message))

			if stock, ok := parseToSubscriptionResponse(message); ok {
				delete(hangingStockList, stock)
				if len(hangingStockList) == 0 {
					finCh <- struct{}{}
				}
				continue
			}
	
			data, err := UnmarshalTrade(string(message))
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, d := range data {
				c.ch <- d
			}	
		}
	}(c.ctx, &c.wg)

	for _, stock := range stocks {
		if err := c.subscribeProduct(stock); err != nil {
			return nil, err
		}
	}

	<- finCh
	return c.ch, nil
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

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	jsonStr := string(data)
	fmt.Println(jsonStr)

	if err := c.conn.WriteJSON(req); err != nil {
		return err
	}
	return nil
}