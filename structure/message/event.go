package message

import "time"

// Photo is a message
// 車の制限速度をチェックしている写真を表現
type Photo struct {
	// License is the license plate number of the car
	// 車のナンバープレート
	License string
	// Speed is the speed of the car
	// 車の速度
	Speed int
}

// NoLicense is a Photo
func (p *Photo) NoLicense() bool {
	return p.License == ""
}

type PhotoMessage struct {
	ID           string
	Photo        string
	CreationTime *time.Time
	Speed        *int
}

type TimeoutMessage struct {
	PhotoMessage
}
