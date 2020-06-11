package models

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
)

type Image struct {
	Model
	CreateUserId int    `json:"create_user_id"`
	Url          string `json:"url"`
}

func (Image) Create(dto dto.UploadDto) (Image, int) {
	image := Image{
		Url:          dto.Url,
		CreateUserId: dto.CreateUserId,
	}
	util.Log.Notice("ImageModel:", image)
	result := db.Create(&image)
	if result.Error == nil {
		return image, 0
	}
	util.Log.Error(result.Error.Error())
	return image, e.ERROR
}

func (Image) Delete(image *Image) bool {
	//db.Delete(&Image{}, "id = ?", id)
	if db.Delete(image).GetErrors() == nil {
		return true
	}
	return false
}
