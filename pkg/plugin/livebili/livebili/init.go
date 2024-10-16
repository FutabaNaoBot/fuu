package livebili

import (
	"github.com/sirupsen/logrus"
	"time"
)

func (b *biliPlugin) init() error {
	conf := Config{}
	err := b.env.GetConf(&conf)
	if err != nil {
		return err
	}
	b.conf = conf
	for _, uid := range b.conf.Uids {
		b.liveState[uid] = false
	}
	go b.ticker()
	return nil
}

func (b *biliPlugin) ticker() {
	t := time.NewTicker(time.Second * time.Duration(b.conf.CheckDuration))
	for range t.C {
		err := b.doCheckLive()
		if err != nil {
			logrus.Errorf("check live error: %v", err)
		}
	}

}
