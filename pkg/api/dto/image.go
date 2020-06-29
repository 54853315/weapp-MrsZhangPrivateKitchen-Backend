package dto

import (
	"mime/multipart"
)

type UploadDto struct {
	GeneralAuthDto
	File *multipart.FileHeader `form:"file,default=0" binding:"required"`
	Url  string
}
