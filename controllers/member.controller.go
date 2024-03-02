package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anangpermana/gin-gorm-base/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberController struct {
	DB *gorm.DB
}

func NewMemberController(DB *gorm.DB) MemberController {
	return MemberController{DB}
}

func Validation(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, f := range verr {
		err := f.ActualTag()
		fmt.Println(f.Field())
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}
	return errs
}

func (mc *MemberController) CreateMember(ctx *gin.Context) {
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Bad Request", "errors": Validation((verr))})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Bad Request"})
		return
	}

	now := time.Now()
	newMember := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		Handphone: payload.Handphone,
		Password:  payload.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := mc.DB.Create(&newMember)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{"success": false, "message": "Member with that email or handphone already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"success": false, "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Success add new data", "data": newMember})
}

func (mc *MemberController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Bad Request", "errors": Validation((verr))})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Bad Request"})
		return
	}

	var updateMember models.Member
	result := mc.DB.First(&updateMember, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No data with that id exists",
		})
		return
	}

	now := time.Now()

	memberToUpdate := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		Handphone: payload.Handphone,
		Password:  payload.Password,
		UpdatedAt: now,
	}

	// Update the member and check for errors
	if err := mc.DB.Model(&updateMember).Updates(memberToUpdate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update member",
		})
		return
	}

	// Respond with the updated data
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully updated member",
		"data":    updateMember,
	})
}

func (mc *MemberController) GetAll(ctx *gin.Context) {

	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	var search = ctx.Query("search")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var members []models.Member
	query := mc.DB.Order("created_at desc").Limit(intLimit).Offset(offset)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	results := query.Find(&members)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"success": false, "message": results.Error})
		return
	}

	var totalMembersCount int64
	resTotalMembers := mc.DB.Model(&models.Member{}).Count(&totalMembersCount)

	if resTotalMembers.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"success": false, "message": resTotalMembers.Error})
		return
	}

	if search != "" {
		totalMembersCount = int64(len(members))
	}

	totalPages := int(math.Ceil(float64(totalMembersCount) / float64(intLimit)))

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success get data",
		"data":    members,
		"meta": gin.H{
			"total":        totalMembersCount,
			"totalPages":   totalPages,
			"currentPages": intPage,
		},
	})
}

func (mc *MemberController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var member models.Member
	result := mc.DB.First(&member, "id = ?", id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "member not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "berhasil",
		"data":    member,
	})
}

func (mc *MemberController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Find the member by id
	var member models.Member
	result := mc.DB.First(&member, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "ID not found",
		})
		return
	}

	// Delete the member
	deleteResult := mc.DB.Delete(&member)
	if deleteResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Member deleted successfully"})
}

func (mc *MemberController) MultipleDelete(ctx *gin.Context) {
	var request struct {
		Ids []uuid.UUID `json:"ids" binding:"required,unique"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload"})
		return
	}

	result := mc.DB.Where("id IN (?)", request.Ids).Delete(&models.Member{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete members"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No members with the provided IDs found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Members deleted successfully"})
}
