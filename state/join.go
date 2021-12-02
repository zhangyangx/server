package state

import (
	"bytes"
	"fmt"
	"github.com/ratel-online/server/consts"
	"github.com/ratel-online/server/database"
	"github.com/ratel-online/server/model"
	"strconv"
)

type join struct{}

func (s *join) Next(player *model.Player) (consts.StateID, error) {
	buf := bytes.Buffer{}
	rooms := database.GetRooms()
	buf.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\n", "ID", "Type", "Players", "State"))
	for _, room := range rooms {
		buf.WriteString(fmt.Sprintf("%d\t%s\t%d\t%s\n", room.ID, consts.GameTypes[room.Type], database.GetRoomPlayers(room.ID), consts.RoomStates[room.State]))
	}
	err := player.WriteString(buf.String())
	if err != nil {
		return 0, player.WriteError(err)
	}
	signal, err := player.AskForString(player.Terminal())
	if err != nil {
		return 0, player.WriteError(err)
	}
	if isExit(signal) {
		return s.Back(player), nil
	}
	if isLs(signal) {
		return consts.StateJoin, nil
	}
	roomId, err := strconv.ParseInt(signal, 10, 64)
	if err != nil {
		return 0, player.WriteError(consts.ErrorsRoomInvalid)
	}
	err = database.JoinRoom(roomId, player.ID)
	if err != nil {
		return 0, player.WriteError(err)
	}
	err = database.RoomBroadcast(roomId, fmt.Sprintf("\r\r%s joined room!\n", player.Name), player.ID)
	if err != nil {
		return 0, player.WriteError(err)
	}
	return consts.StateWaiting, nil
}

func (*join) Back(player *model.Player) consts.StateID {
	return consts.StateHome
}
