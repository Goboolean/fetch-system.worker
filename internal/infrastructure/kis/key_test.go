package kis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kis"
	"github.com/stretchr/testify/assert"
)

func TestGetApprovalKey(t *testing.T) {

	var client = new(kis.Client)

	t.Run("GetApprovalKey (case:fail)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
		defer cancel()

		_, err := client.GetApprovalKey(ctx, "", "")
		assert.Error(t, err)
	})

	t.Run("GetApprovalKey (case:success)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
		defer cancel()

		key, err := client.GetApprovalKey(ctx, os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key)
	})

	t.Run("GetMultipleApprovalKey", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
		defer cancel()

		key1, err := client.GetApprovalKey(ctx, os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key1)

		key2, err := client.GetApprovalKey(ctx, os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, key2)

		assert.NotEqual(t, key1, key2)
	})
}



func TestIssueToken(t *testing.T) {

	t.Skip("Skip this test on development mode since issuing token is limited to 1 token per day.")

	var client = new(kis.Client)

	t.Run("IssueToken (case:fail)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
		defer cancel()

		_, err := client.IssueAccessToken(ctx, "", "")
		assert.Error(t, err)
	})

	t.Run("IssueToken (case:success)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
		defer cancel()

		token, err := client.IssueAccessToken(ctx, os.Getenv("KIS_APPKEY"), os.Getenv("KIS_SECRET"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}