package categories

type Category struct {
	CategoryId int    `db:"category_id" json:"categoryId"`
	Title      string `db:"title" json:"title"`
	Desc       string `db:"desc" json:"desc"`
}
