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

// PhotoMessage is a message
type PhotoMessage struct {
	ID           string
	Photo        string
	CreationTime *time.Time
	Speed        *int
}

// IllegalStatePanicMessage is a message
// 意図的にパニックを起こすメッセージ
type IllegalStatePanicMessage struct{}

// HasCreationTime is a PhotoMessage
func (p *PhotoMessage) HasCreationTime() bool {
	return p.CreationTime != nil
}

// HasSpeed is a PhotoMessage
func (p *PhotoMessage) HasSpeed() bool {
	return p.Speed != nil
}

// UpdatePhotoMessage is a PhotoMessage
func (p *PhotoMessage) UpdatePhotoMessage(this *PhotoMessage) {
	if this.HasCreationTime() {
		p.CreationTime = this.CreationTime
	}
	if this.HasSpeed() {
		p.Speed = this.Speed
	}
}

type TimeoutMessage struct {
	PhotoMessage
}
