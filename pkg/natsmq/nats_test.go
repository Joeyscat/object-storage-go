package natsmq

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeWithReply(t *testing.T) {
	var url = "nats://me.io:4222"

	var subject = "hello"

	subNc, err := nats.Connect(url, nats.Name("sub--------"))
	assert.Nil(t, err)

	subNum := 4
	for i := 0; i < subNum; i++ {
		go func(index int) {
			SubscribeWithReply(subNc, subject, func(msg *nats.Msg) ([]byte, error) {
				var bs []byte
				bs, err = json.Marshal(map[string]string{"addr": fmt.Sprintf("127.0.0.%d", index), "msg": string(msg.Data)})
				return bs, err
			})
		}(i)
	}

	pubNc, err := nats.Connect(url, nats.Name("pub--------"))
	assert.Nil(t, err)

	for i := 0; i < 100; i++ {
		go func(index int) {
			key := fmt.Sprintf("000%d", index)
			replyCount := 4
			r, err := PublichAndWaitForReply(pubNc, subject, []byte(key), time.Millisecond*1000, replyCount)
			assert.Nil(t, err)
			assert.Equal(t, replyCount, len(r))

			rs := make([]string, 0)
			for _, v := range r {
				rs = append(rs, string(v.Data))

				assert.Contains(t, string(v.Data), key)
			}
			t.Logf("key: %s, rs: %v", key, rs)
		}(i)
	}

	time.Sleep(time.Second * 1)
}
