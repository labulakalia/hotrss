package util

import (
	"testing"

	"github.com/antlabs/pcurl"
)

func TestRemoveSlash(t *testing.T) {
	const cURLData = `'https://xueqiu.com/statuses/hot/listV2.json?since_id=-1&max_id=-1&size=15' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Pragma: no-cache' \
-H 'Cookie: CNZZDATA1256793290=292804978-1604453779-%7C1604972190; Hm_lpvt_1db88642e346389874251b5a1eded6e3=1604973929; Hm_lvt_1db88642e346389874251b5a1eded6e3=1604456626,1604973482,1604973929; u=781604973480446; xq_a_token=db48cfe87b71562f38e03269b22f459d974aa8ae; xq_id_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOi0xLCJpc3MiOiJ1YyIsImV4cCI6MTYwNjk2MzA1MCwiY3RtIjoxNjA0OTczNDY3NTM3LCJjaWQiOiJkOWQwbjRBWnVwIn0.E4GK-vwlTQYx3OgjPrSEqKCORY02uu6l0ezDEvih8DEs0THi__fpRkc113dGSpLWOVTTVEmFnUWk7Wx2UDZUK-jwfjz3MezovQUn3UVU-R7kWeSIZlXL2UGEI-5eJwoaGbxaA_l93rF4ESEvkwbGC6H9GKnVSBqJUzR1jmb_zjUUL8DSxUhvyk2TGCVvELMUJEcsL_eVsZfnL6_xu4ngn8T4pr5TkFR5ae3RY9NaccjcdftbD4t5nfdkHh4NXs0Fu0VuGrGYb0jpFs0s15oqtS0hVe4UGVuzuqJNFXC73CdtYyp88MWGADXTmH8vAfOMqeNQ4tQGaqQGTjGzAKcmDQ; xq_r_token=500b4e3d30d8b8237cdcf62998edbf723842f73a; xqat=db48cfe87b71562f38e03269b22f459d974aa8ae; acw_tc=2760824316049734804122452e827d19f916273b7a6a36908625432ede0e51; UM_distinctid=17591113fb865c-0c0bc6e8ea9c6e8-5c465d7b-1fa400-17591113fb9d53; device_id=df2de71e98cb84acedcb07542ad03de3' \
-H 'Cache-Control: no-cache' \
-H 'Accept-Language: zh-cn' \
-H 'Host: xueqiu.com' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15' \
-H 'Referer: https://xueqiu.com/?category=snb_article' \
-H 'Accept-Encoding: deflate, br' \
-H 'Connection: keep-alive' \
-H 'elastic-apm-traceparent: 00-eba5cbb0ffee1a2de98b1311efbb3149-a51a9e6d349eee5e-01'`
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		// want string
	}{
		// TODO: Add test cases.
		{
			name: "remove",
			args: args{
				url: cURLData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveSlash(tt.args.url)
			req, err := pcurl.ParseAndRequest(got)
			if err != nil {
				t.Fatal(err)
			}
			if len(req.Header) == 0 {
				t.Fatal("paese failed . header is empty")
			}
			t.Logf("header count %d", len(req.Header))
		})
	}
}
