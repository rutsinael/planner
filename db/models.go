package db

type Task struct {
	ID      string `db:"id"      json:"id,omitempty"`
	Date    string `db:"date"    json:"date"`
	Title   string `db:"title"   json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat"  json:"repeat"`
}
