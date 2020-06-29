package dto

// BookListSearchMapping - define search query keys in role list page
var BookListSearchMapping = map[string]string{
	"n": "name",
	"s": "status",
}

// BookCreateDto - dto for role's creation
type BookCreateDto struct {
	GeneralAuthDto
	//Name                    string `form:"name" json:"name" binding:"required"`
	Content       string `form:"content" json:"content" binding:"max=1000"`
	AllowComments int    `form:"allow_comments" json:"allow_comments"`
	//CreateUserId  int  `binding:"required,gte=1"`
	Status string   `form:"status" json:"status" binding:"required,oneof=private publish"`
	Files  []string `form:"files" json:"files" binding:"required"`
}

// BookEditDto - dto for role's modification
type BookEditDto struct {
	GeneralAuthDto
	Id int `form:"id" json:"id" binding:"required,gte=1"`
	//Name                    string `form:"name" json:"name" binding:"required"`
	Content       string `form:"content" json:"content" binding:"max=1000"`
	AllowComments int    `form:"allow_comments" json:"allow_comments"`
	//CreateUserId  int
	Status string   `form:"status" json:"status" binding:"required,oneof=private publish"`
	Files  []string `form:"files" json:"files" binding:"required"`
}

type BookChangeDto struct {
	GeneralAuthDto
	Id int `form:"id" json:"id" binding:"required,gte=1"`
	//CreateUserId int
	Status string `form:"status" json:"status" binding:"required,oneof=private publish"`
}
