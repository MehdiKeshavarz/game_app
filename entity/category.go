package entity

type Category string

const (
	FootballCategory Category = "football"
	HistoryCategory  Category = "history"
	TechCategory     Category = "tech"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory, HistoryCategory, TechCategory:
		return true
	default:
		return false
	}
}
