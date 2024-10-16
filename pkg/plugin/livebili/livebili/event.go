package livebili

import (
	"bytes"
	"encoding/json"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"log"
	"net/http"
)

func (b *biliPlugin) doCheckLive() error {
	var uids []int64
	for k, _ := range b.liveState {
		uids = append(uids, k)
	}
	data := map[string]interface{}{
		"uids": uids,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	live := LiveResp{}
	err = json.Unmarshal(body, &live)
	if err != nil {
		return err
	}
	if live.Code != 0 {
		return fmt.Errorf("code: %d,msg: %s", live.Code, live.Msg)
	}

	for _, info := range live.Data {
		err = b.sendRoomInfo(&info)
		if err != nil {
			return err
		}
	}
	return nil

}

func (b *biliPlugin) sendRoomInfo(info *RoomInfo) error {
	lastStatus := b.liveState[int64(info.Uid)]
	living := IsLiving(info.LiveStatus)
	if lastStatus != living {
		b.liveState[int64(info.Uid)] = living
	}
	if !living {
		return nil
	}
	b.env.RangeBot(func(ctx *zero.Ctx) bool {
		msgChan := []message.MessageSegment{
			message.AtAll(),
			message.Text("开播啦！"),
			message.Text(info.Title),
			message.Image(info.CoverFromUser),
		}
		for _, group := range b.conf.Groups {
			ctx.SendGroupMessage(group, msgChan)
		}
		return true
	})

	return nil
}
