package etcdutil_test

import (
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	etcdutil "github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd/util"
)

// Current version does not provide Nested struct, only flat struct is supported.
type Nested struct {
	ID     string `etcd:"id"`
	Detail struct {
		Name string `etcd:"name"`
		Age  int    `etcd:"age"`
	} `etcd:"detail"`
}

func (n *Nested) Name() string {
	return "nested"
}

var timestamp = time.Now().String()

type group struct {
	str   map[string]string
	model etcdutil.Model
	data  etcdutil.Model
}

var cases []struct {
	name  string
	whole map[string]string
	group []group
} = []struct {
	name  string
	whole map[string]string
	group []group
}{
	{
		name: "Worker",
		whole: map[string]string{
			"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d":          "",
			"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/platform": "kis",
			"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/status":   "waiting",
			"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/lease_id":  "0",
			"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/timestamp":  timestamp,
			"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886":          "",
			"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/platform": "kis",
			"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/status":   "active",
			"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/lease_id":  "0",
			"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/timestamp":  timestamp,
		},
		group: []group{
			{
				str: map[string]string{
					"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d":          "",
					"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/platform": "kis",
					"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/status":   "waiting",
					"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/lease_id":  "0",
					"/worker/9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d/timestamp":  timestamp,
				},
				model: &etcd.Worker{},
				data: &etcd.Worker{
					ID:       "9cf226f7-4ee8-4a5c-9d2f-6d7c74f6727d",
					Platform: "kis",
					Status:   "waiting",
					LeaseID:  "0",
					Timestamp: timestamp,
				},
			},
			{
				str: map[string]string{
					"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886":          "",
					"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/platform": "kis",
					"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/status":   "active",
					"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/lease_id":  "0",
					"/worker/b9992d7b-a926-483a-84f8-bbc05dee7886/timestamp":  timestamp,
				},
				model: &etcd.Worker{},
				data: &etcd.Worker{
					ID:       "b9992d7b-a926-483a-84f8-bbc05dee7886",
					Platform: "kis",
					Status:   "active",
					LeaseID:  "0",
					Timestamp: timestamp,
				},
			},
		},
	},
	{
		name: "Product",
		whole: map[string]string{
			"/product/test.goboolean.kor":          "",
			"/product/test.goboolean.kor/platform": "kis",
			"/product/test.goboolean.kor/symbol":   "goboolean",
			"/product/test.goboolean.kor/locale":   "kor",
			"/product/test.goboolean.kor/market":   "stock",
			"/product/test.goboolean.eng":          "",
			"/product/test.goboolean.eng/platform": "polygon",
			"/product/test.goboolean.eng/symbol":   "gofalse",
			"/product/test.goboolean.eng/locale":   "eng",
			"/product/test.goboolean.eng/market":   "stock",
		},
		group: []group{
			{
				str: map[string]string{
					"/product/test.goboolean.kor":          "",
					"/product/test.goboolean.kor/platform": "kis",
					"/product/test.goboolean.kor/symbol":   "goboolean",
					"/product/test.goboolean.kor/locale":   "kor",
					"/product/test.goboolean.kor/market":   "stock",
				},
				model: &etcd.Product{},
				data: &etcd.Product{
					ID:       "test.goboolean.kor",
					Platform: "kis",
					Symbol:   "goboolean",
					Locale:   "kor",
					Market:   "stock",
				},
			},
			{
				str: map[string]string{
					"/product/test.goboolean.eng":          "",
					"/product/test.goboolean.eng/platform": "polygon",
					"/product/test.goboolean.eng/symbol":   "gofalse",
					"/product/test.goboolean.eng/locale":   "eng",
					"/product/test.goboolean.eng/market":   "stock",
				},
				model: &etcd.Product{},
				data: &etcd.Product{
					ID:       "test.goboolean.eng",
					Platform: "polygon",
					Symbol:   "gofalse",
					Locale:   "eng",
					Market:   "stock",
				},
			},
		},
	},
}