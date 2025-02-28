package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/utils"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PassageTagsRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (p *PassageTagsRepositoryImpl) AutoMigrate() {
	p.db.AutoMigrate(&entity.PassageTags{})
}

func (p *PassageTagsRepositoryImpl) SetDb(db *gorm.DB) {
	p.db = db
}

func (p *PassageTagsRepositoryImpl) SetRedis(rdb *cache.Redis) {
	p.rdb = rdb
}

func (p *PassageTagsRepositoryImpl) GetPassageTags(passageId string) ([]string, error) {
	var tags []string

	tagsJson, err := p.rdb.Get("blog:post:" + passageId + ":tags")

	if err != nil {
		err2 := p.db.Table("passage_tags").
			Select("tag_name").
			Joins("JOIN tags ON tags.tag_id = passage_tags.tag_id").
			Where("passage_tags.passage_id = ?", passageId).
			Pluck("tags.tag_name", &tags).Error

		if err2 != nil {
			return nil, err
		}

		if err == redis.Nil {
			go func() {
				tagsToJson, _ := json.Marshal(tags)
				p.rdb.Set("blog:post:"+passageId+":tags", string(tagsToJson), utils.RandExpTime(60, 70))
			}()
		}

		return tags, nil
	}

	err = json.Unmarshal([]byte(tagsJson), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (p *PassageTagsRepositoryImpl) AddPassageTags(passageId string, tags []string) error {
	var passageTags []entity.PassageTags
	for _, tag := range tags {
		var tagid int
		p.db.Model(&entity.Tag{}).
			Select("tag_id").
			Where("tag_name = ?", tag).
			Pluck("tag_id", &tagid)

		passageTags = append(passageTags, entity.PassageTags{
			PassageId: passageId,
			TagId:     tagid,
		})
	}

	return p.db.Create(&passageTags).Error
}
func (p *PassageTagsRepositoryImpl) GetTagPassages(tagName string) ([]string, error) {
	var passagesId []string

	passagesIdJson, err := p.rdb.Get("blog:tags:" + tagName)

	if err != nil {
		err2 := p.db.Table("passage_tags").
			Select("passage_id").
			Joins("JOIN tags ON tags.tag_id = passage_tags.tag_id").
			Where("tags.tag_name = ?", tagName).
			Pluck("passage_id", &passagesId).Error

		if err2 != nil {
			return nil, err
		}

		if err == redis.Nil {
			go func() {
				passagesIdToJson, _ := json.Marshal(passagesId)
				p.rdb.Set("blog:tags:"+tagName, string(passagesIdToJson), utils.RandExpTime(60, 70))
			}()
		}

		return passagesId, nil
	}

	err = json.Unmarshal([]byte(passagesIdJson), &passagesId)
	if err != nil {
		return nil, err
	}

	return passagesId, nil
}
