package kis

import (
	"os"
	"testing"
)

func TestGetApprovalKey(t *testing.T) {

	var client *Client

	t.Run("GetApprovalKey (case:fail)", func(t *testing.T) {

		_, err := client.GetApprovalKey("", "")
		if err == nil {
			t.Errorf("GetApprovalKey() = %v, want error", err)
			return
		}
	})

	t.Run("GetApprovalKey (case:success)", func(t *testing.T) {

		key, err := client.GetApprovalKey(os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		if err != nil {
			t.Errorf("GetApprovalKey() = %v", err)
			return
		}

		if key == "" {
			t.Errorf("GetApprovalKey() = %v, want not empty", key)
			return
		}
	})
}
