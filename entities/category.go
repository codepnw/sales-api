package entities

import "time"

func TableCategories() string {
	return "categories"
}

type Category struct {
	Id        int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	Id   int    `json:"categoryId"`
	Name string `json:"name"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Meta       *QueryMeta         `json:"meta"`
}

type CategoryUpdate struct {
	Id        int       `json:"category_id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}
