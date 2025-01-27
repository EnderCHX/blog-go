package controller

import (
	"blog-go/blog"
	"blog-go/log"
	"errors"
	"strconv"

	"fmt"
	"net/http"
	"time"

	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewPost(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"code":    "Unauthorized",
			"data":    nil,
		})
		return
	}

	var passageReqBody blog.PassageReqBody
	ctx.BindJSON(&passageReqBody)

	passage := blog.NewPassage(
		blog.WithAuthorUsername(claims.(*auth.JWTPayload).Username),
		blog.WithContent(passageReqBody.Content),
		blog.WithTitle(passageReqBody.Title),
		blog.WithCreatedAt(time.Now()),
		blog.WithUpdatedAt(time.Now()),
		blog.WithDeleted(false),
	).GenPassageId()

	err := passage.Create()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "文章创建失败",
			"code":    "InternalServerError",
			"data":    err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"code":    "Success",
		"data":    passage,
	})

	// 创建标签
	var tags []blog.Tag
	for _, tag := range passageReqBody.Tags {
		tags = append(tags, blog.Tag{
			TagName: tag,
		})
	}

	err = blog.CreateTags(tags)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Logger.Info("[MySQL]标签重复", zap.String("error", err.Error()))
		}
	}

	//获取所有标签
	tagsall, err := blog.GetTags()
	if err != nil {
		log.Logger.Error("[MySQL]获取标签失败", zap.String("error", err.Error()))
	}

	//更新标签缓存
	err = blog.UpdateTagsCache(tagsall)

	if err != nil {
		log.Logger.Error("[Redis]缓存失败", zap.String("error", err.Error()))
	}

	tagsmap, _ := blog.GetTagsMap()

	var passageTags []blog.PassageTags

	// 创建文章标签映射
	for _, tag := range tags {
		passageTags = append(passageTags, blog.PassageTags{
			PassageId: passage.PassageId,
			TagId: func() int {
				r, _ := strconv.Atoi(tagsmap[tag.TagName])
				return r
			}(),
		})
		blog.RemoveTagCache(tag.TagName) //删除缓存
	}
	err = blog.CreatePassageTags(passageTags)
	if err != nil {
		log.Logger.Error("[MySQL]创建 passage_tags 失败", zap.String("error", err.Error()))
	}
}

func GetPassages(ctx *gin.Context) {
	passages, err := blog.GetPassagesCache()

	if err != nil {
		passages, resultErr := blog.GetPassages()
		if resultErr != nil {
			log.Logger.Error("[MySQL]获取文章列表失败", zap.String("error", resultErr.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "InternalServerError",
				"message": "获取文章列表失败",
				"data":    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "Success",
			"message": fmt.Sprintf("获取文章列表成功, 共有 %d 篇文章", len(passages)),
			"data":    passages,
		})
		if err == redis.Nil {
			err = blog.UpdatePassagesCache(passages)
			if err != nil {
				log.Logger.Error("[Redis]设置redis数据失败", zap.String("error", err.Error()))
			}
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": fmt.Sprintf("获取文章成功，共 %d 篇文章", len(passages)),
		"data":    passages,
	})
}

func GetPassagesByTag(ctx *gin.Context) {
	tagName := ctx.Param("tagname")

	passages, err := blog.GetPassagesByIdCache(tagName)

	if err != nil {
		passages, resultErr := blog.GetPassagesByTag(tagName)

		if resultErr != nil && resultErr.Error() == "passage not found" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "NotFound",
				"message": "未找到",
				"data":    nil,
			})
			return
		}

		if resultErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "InternalServerError",
				"message": "InternalServerError",
				"data":    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "Success",
			"message": fmt.Sprintf("获取文章列表成功, 共有 %d 篇文章", len(passages)),
			"data":    passages,
		})
		if err == redis.Nil {
			err = blog.UpdatePassagesByIdCache(passages, tagName)
			if err != nil {
				log.Logger.Error("[Redis]缓存更新失败", zap.String("error", err.Error()))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "获取成功",
		"data":    passages,
	})
}

func GetPassageById(ctx *gin.Context) {
	var passage blog.Passage
	passage.PassageId = ctx.Param("id")
	err := passage.GetCache()

	if err != nil {
		err2 := passage.Get()
		if err2 != nil {
			if errors.Is(err2, gorm.ErrRecordNotFound) {
				log.Logger.Error("[MySQL]文章 "+passage.PassageId+" 不存在", zap.String("error", err2.Error()))
				ctx.JSON(http.StatusNotFound, gin.H{
					"code":    "PassageNotFound",
					"message": "文章不存在",
					"data":    nil,
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    "Success",
			"message": "获取成功",
			"data":    passage,
		})

		if err == redis.Nil {
			err = passage.UpdateCache()
			if err != nil {
				log.Logger.Error("[Redis]缓存失败", zap.String("error", err.Error()))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "获取成功",
		"data":    passage,
	})
}

func GetPassageTags(ctx *gin.Context) {
	passageId := ctx.Param("id")
	tags, err := blog.GetPassageTagsCache(passageId)

	if err != nil {
		tags, err2 := blog.GetPassageTags(passageId)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "InternalServerError",
				"message": "服务器内部错误",
				"data":    nil,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    "Success",
			"message": "获取标签成功",
			"data":    tags,
		})

		if err.Error() == "no tags in cache" && len(tags) > 0 {
			err = blog.UpdatePassageTagsCache(passageId, tags)
			if err != nil {
				log.Logger.Error("[Redis]缓存失败", zap.String("error", err.Error()))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "获取成功",
		"data":    tags,
	})
}

func GetTags(ctx *gin.Context) {
	tags, err := blog.GetTagsCache()
	log.Logger.Debug("[Redis]GetTags")
	if err != nil {
		tags, err2 := blog.GetTags()
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "InternalServerError",
				"message": "服务器内部错误",
				"data":    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "Success",
			"message": "获取标签成功",
			"data":    tags,
		})

		if err.Error() == "no tags in cache" {
			log.Logger.Debug("[Redis]Set Cache")
			err = blog.UpdateTagsCache(tags)
			if err != nil {
				log.Logger.Error("[Redis]缓存失败", zap.String("error", err.Error()))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "获取标签成功",
		"data":    tags,
	})
}
