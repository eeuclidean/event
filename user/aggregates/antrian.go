package aggregates

import (
	"strings"
	"time"
)

const (
	ANTRIAN_CREATED = "antrian_created"
	ANTRIAN_UPDATED = "antrian_updated"
)

const (
	ANTRIAN_BRANCH = 1
	ANTRIAN_POLI   = 2
)

func NewEmptyAntrianPoli(poli Poli, tanggal string) Antrian {
	antrian := newAntrian(poli, tanggal)
	antrian.ID = poli.ID + strings.ReplaceAll(tanggal, "-", "")
	antrian.Kuota = poli.Max
	antrian.PoliID = poli.ID
	antrian.Type = ANTRIAN_POLI
	antrian.Terisi = 0
	return antrian
}

func NewAntrianPoli(poli Poli, tanggal string) Antrian {
	antrian := newAntrian(poli, tanggal)
	antrian.ID = poli.ID + strings.ReplaceAll(tanggal, "-", "")
	antrian.Kuota = poli.Max
	antrian.PoliID = poli.ID
	antrian.Type = ANTRIAN_POLI
	return antrian
}

func NewAntrianBranch(poli Poli, tanggal string) Antrian {
	antrian := newAntrian(poli, tanggal)
	antrian.ID = poli.BranchID + strings.ReplaceAll(tanggal, "-", "")
	antrian.Kuota = -1
	antrian.Type = ANTRIAN_BRANCH
	return antrian
}

func newAntrian(poli Poli, tanggal string) Antrian {
	return Antrian{
		BranchID: poli.BranchID,
		Tanggal:  tanggal,
		Checkin:  0,
		Selesai:  0,
		Terisi:   1,
		Status:   ANTRIAN_CREATED,
		Created:  time.Now(),
	}
}

type Antrian struct {
	ID       string    `json:"id,omitempty" bson:"_id,omitempty"`
	BranchID string    `json:"branch_id,omitempty" bson:"branch_id,omitempty"`
	PoliID   string    `json:"poli_id,omitempty" bson:"poli_id,omitempty"`
	Tanggal  string    `json:"tanggal,omitempty" bson:"tanggal,omitempty"`
	Kuota    int       `json:"kuota,omitempty" bson:"kuota,omitempty"`
	Type     int       `json:"type,omitempty" bson:"type,omitempty"`
	Checkin  int       `json:"check_in,omitempty" bson:"check_in,omitempty"`
	Terisi   int       `json:"terisi,omitempty" bson:"terisi,omitempty"`
	Selesai  int       `json:"selesai,omitempty" bson:"selesai,omitempty"`
	Status   string    `json:"status,omitempty" bson:"status,omitempty"`
	Created  time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}
