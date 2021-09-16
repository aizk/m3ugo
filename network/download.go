package network

import (
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/grafov/m3u8"
	"github.com/liunian1004/m3ugo/decode"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"strings"
)

func Download(playList *m3u8.MediaPlaylist, m3uUrl, fileName string) (err error) {
	f, err := os.OpenFile(fmt.Sprintf("./%s.ts", fileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(m3uUrl)
	if err != nil {
		logrus.Errorln("parse m3uUrl error")
		return
	}
	var key []byte
	var iv []byte
	var res *resty.Response
	if playList.Key != nil {
		res, err = Client.R().Get(playList.Key.URI)
		if err != nil {
			logrus.Errorln("download key file error", err)
			return
		}
		iv, err = hex.DecodeString(strings.TrimPrefix(playList.Key.IV, "0x"))
		if err != nil {
			logrus.Errorln("hex decode IV error", err)
			return
		}
		key = res.Body()
	}

	for i, segment := range playList.Segments {
		if segment == nil {
			break
		}

		// resolve url
		newUrl, err1 := u.Parse(segment.URI)
		if err1 != nil {
			logrus.Errorln(fmt.Sprintf("parse segment URI:%s error: %s", segment.URI, err1.Error()))
			return
		}

		// get data
		res, err = Client.R().SetHeader("Host", newUrl.Host).Get(newUrl.String())
		if err != nil {
			logrus.Errorln("fetch segment URI data error" + err.Error())
			return
		}

		data := res.Body()

		// need decrypt?
		if playList.Key != nil {
			data, err = decode.AESCBCDecrypt(data, key, iv)
			if err != nil {
				logrus.Errorln("decode AESCBCDecrypt error", err)
				return
			}
		}

		_, err = f.Write(data)
		if err != nil {
			logrus.Errorln("write data to fail error")
			return
		}

		logrus.Printf("dowload index %d success.\n", i)
	}

	return
}