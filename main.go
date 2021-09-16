package main

import (
	"github.com/grafov/m3u8"
	"github.com/guonaihong/clop"
	"github.com/liunian1004/m3ugo/fetch/bd"
	"github.com/liunian1004/m3ugo/network"
	"github.com/liunian1004/m3ugo/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Link           string   `clop:"-l; --link" usage:"link of video address"`
	URL            string   `clop:"-u; --m3u8-url" usage:"url of m3u8 file"`
	File           string   `clop:"-f; --m3u8-file" usage:"local m3u8 file"`
	OutputFileName string   `clop:"-o; --out-file-name" usage:"out file name"`
	Headers        []string `clop:"-H; --header; greedy" usage:"http header. Example: Referer:http://www.example.com"`
}

func main() {
	var err error
	conf := &Config{}
	err = clop.Bind(conf)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	if conf.Link == "" && conf.URL == "" && conf.File == "" {
		logrus.Errorln("must have link or m3u8 url or file path")
		return
	}

	if conf.URL != "" || conf.File != "" {
		if conf.OutputFileName == "" {
			logrus.Errorln("use m3u8 url or file must set filename with -o")
		}
		return
	}

	if len(conf.Headers) > 0 {
		for _, header := range conf.Headers {
			s := strings.SplitN(header, ":", 2)
			key := strings.TrimRight(s[0], " ")
			if len(s) == 2 {
				network.Client.SetHeader(key, strings.TrimLeft(s[1], " "))
			}
		}
	}

	viper.SetConfigFile("config.yml")
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalln("read config: ",err)
	}

	if len(viper.GetStringMapString("header")) != 0 {
		for k, v := range viper.GetStringMapString("header") {
			logrus.Infof("set header %s: %s", k, v)
			network.Client.SetHeader(k, v)
		}
	}

	// m3u8 & filename resolve
	var fileName string
	var m3uUrl string

	// get m3u8 URL
	if conf.URL != "" {
		m3uUrl = conf.URL
		if conf.OutputFileName != "" {
			fileName = conf.OutputFileName
		} else {
			logrus.Fatalln("download filename is empty please set with -o argument")
		}
	} else {
		f := &bd.BDFetch{}
		m3uUrl, fileName, err = f.Fetch(conf.Link)
		if err != nil {
			logrus.Errorln(err)
			return
		}
	}

	logrus.Infoln("m3u url: ", m3uUrl)
	// parse m3u8 file
	var playList *m3u8.MediaPlaylist
	if conf.File != "" {
		playList, err = parse.ParesM3uFile(conf.File)
		if err != nil {
			logrus.Errorln(err)
			return
		}
	} else {
		playList, err = parse.ParseM3uUrl(m3uUrl)
		if err != nil {
			logrus.Errorln(err)
			return
		}
	}


	// download file
	err = network.Download(playList, m3uUrl, fileName)
	if err != nil {
		logrus.Errorln("download error", err)
		return
	}

	// concurrent - merge？

	// trans mp4?
}
