package domain

type Project struct {
	Id          uint64
	Title       string
	Description string
	CreatorId   uint64
}

type Projects struct {
	Projects    []Project
	Total       uint64
	CurrentPage int32
	LastPage    int32
}

func (p Project) GetOwnerId() uint64 {
	return p.CreatorId
}
