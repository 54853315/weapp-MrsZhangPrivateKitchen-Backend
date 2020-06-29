package dto

// CommentListSearchMapping - define search query keys in role list page
var CommentListSearchMapping = map[string]string{
	"book_id": "book_id",
}

// CommentListDto - General list request params
type CommentListDto struct {
	GeneralListDto
	BookId int `form:"book_id" json:"skip" binding:"required,gte=1"`
}

// CommentCreateDto - dto for role's creation
type CommentCreateDto struct {
	CreateUserId int    `binding:"required,gte=1"`
	BookId       int    `form:"book_id" json:"book_id" binding:"required,gte=1"`
	Content      string `form:"content" json:"content" binding:"required,max=1000"`
	//Status  string `form:"status" json:"status" binding:"required"`
}

// CommentEditDto - dto for role's modification
type CommentEditDto struct {
	CreateUserId int    `binding:"required,gte=1"`
	Id           int    `form:"id" json:"id" binding:"required,gte=1"`
	BookId       int    `form:"book_id" json:"book_id" binding:"required,gte=1"`
	Content      string `form:"content" json:"content" binding:"required,max=1000"`
	//Status  string `form:"status" json:"status" binding:"required"`
}
