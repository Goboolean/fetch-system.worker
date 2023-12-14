package kis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetApprovalKey(t *testing.T) {

	var client *Client

	t.Run("GetApprovalKey (case:fail)", func(t *testing.T) {

		_, err := client.GetApprovalKey("", "")
		assert.Error(t, err)
	})

	t.Run("GetApprovalKey (case:success)", func(t *testing.T) {

		key, err := client.GetApprovalKey(os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key)
	})

	t.Run("GetMultipleApprovalKey", func(t *testing.T) {

		key1, err := client.GetApprovalKey(os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key1)

		key2, err := client.GetApprovalKey(os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key2)

		assert.NotEqual(t, key1, key2)
	})
}
