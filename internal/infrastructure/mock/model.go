package mock

import (
	"fmt"
	"time"
	"math/rand"
)


type Trade struct {
	Symbol    string
	Price     float64
	Amount    int
	Timestamp time.Time
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randstr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var cache = make(map[string]struct{})



func Symbol() string {

	var name string

	for {
		name = randstr(5)
		if _, exist := cache[name]; !exist {
			cache[name] = struct{}{}
			break
		}
	}

	return fmt.Sprintf("%s.%s.%s", "MOCK",name,"io")
}