package dto

// TagListSearchMapping - define search query keys in role list page
var TagListSearchMapping = map[string]string{
	"n": "name",
}

// TagCreateDto - dto for role's creation
type TagCreateDto struct {
	Name string `form:"name" json:"name" binding:"required"`
	//Status                  string `form:"status" json:"status" binding:"required"`
}

// TagEditDto - dto for role's modification
type TagEditDto struct {
	ID   int    `form:"id" json:"id" binding:"required,gte=1"`
	Name string `form:"name" json:"name" binding:"required"`
	//Status                  string `form:"status" json:"status" binding:"required"`
}
