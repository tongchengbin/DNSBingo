package store

import "time"

// 间隔扫描map  进行删除过期key

func ExpireMap() {
	for {
		expireTime := time.Now().Unix()
		for k, v := range DnsMap {
			expire := true
			for _, d := range v {
				if expireTime-d.Time > 3600 {
					expire = false
					break
				}
			}
			if expire {
				// remove
				delete(DnsMap, k)
			}
		}
	}
}
