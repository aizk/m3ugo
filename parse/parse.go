package parse

import (
	"bytes"
	"errors"
	"github.com/grafov/m3u8"
	"github.com/liunian1004/m3ugo/network"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
)

func ParseM3uUrl(m3uUrl string) (list *m3u8.MediaPlaylist, err error) {
	u, err := url.Parse(m3uUrl)
	if err != nil {
		logrus.Errorln("parse m3uUrl error")
		return
	}

	res, err := network.Client.R().
		SetHeader("Host", u.Host).
		Get(m3uUrl)
	if err != nil {
		return
	}
	data := res.Body()

	playlist, listType, err := m3u8.Decode(*bytes.NewBuffer(data), true)
	if err != nil {
		return
	}

	if listType == m3u8.MEDIA {
		list = playlist.(*m3u8.MediaPlaylist)
		return
	}

	err = errors.New("parse m3u fail: unsupport type")
	return
}

func ParesM3uFile(path string) (list *m3u8.MediaPlaylist, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	playlist, listType, err := m3u8.Decode(*bytes.NewBuffer(data), true)
	if err != nil {
		return
	}

	if listType == m3u8.MEDIA {
		list = playlist.(*m3u8.MediaPlaylist)
		return
	}

	err = errors.New("parse m3u fail: unsupport type")


	return
}