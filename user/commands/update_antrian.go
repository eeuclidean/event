package commands

import (
	"strings"
)

type UpdateAntrianCommand struct {
	PoliID  string `json:"poli_id,omitempty"`
	Tanggal string `json:"tanggal,omitempty"`
	Kuota   int    `json:"kuota,omitempty"`
}

func (command UpdateAntrianCommand) GetAntrainPoliID() string {
	return command.PoliID + strings.ReplaceAll(command.Tanggal, "-", "")
}
