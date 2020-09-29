package object

import "time"

type TimeImpl struct {
	bTime time.Time
	MTime time.Time
	CTime time.Time
}

func NewTimeImpl() TimeImpl {
	now := time.Now()
	return TimeImpl{
		bTime: now,
		MTime: now,
		CTime: now,
	}
}

func (t *TimeImpl) BirthTime() time.Time {
	return t.bTime
}

func (t *TimeImpl) ModTime() time.Time {
	return t.MTime
}

func (t *TimeImpl) ChangeTime() time.Time {
	return t.CTime
}
