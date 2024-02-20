package middlw

import (
	"encoding/json"
	"gitlab.com/ptflp/gopubsub/queue"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Email string
	Phone string
}

type Limit struct {
	req  int
	queM queue.MessageQueuer
}

func NewLimit(r int, que queue.MessageQueuer) *Limit {

	return &Limit{req: r, queM: que}
}
func (l *Limit) LimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(l.req)), l.req)
	data := &Data{
		Email: "@example",
		Phone: "891234567890",
	}
	jsData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !limiter.Allow() {
			log.Println("месс сенд ту брокер")
			err = l.queM.Publish("limitreq", jsData)
			if err != nil {
				log.Fatal(err)
			}

			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
