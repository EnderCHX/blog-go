package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TagsRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (t *TagsRepositoryImpl) AutoMigrate() {
	t.db.AutoMigrate(&entity.Tag{})
}

func (p *TagsRepositoryImpl) SetDb(db *gorm.DB) {
	p.db = db
}

func (p *TagsRepositoryImpl) SetRedis(rdb *cache.Redis) {
	p.rdb = rdb
}

func (t *TagsRepositoryImpl) GetTags() ([]entity.Tag, error) {
	var tags []entity.Tag

	tagsJson, err := t.rdb.Get("blog:tags")

	if err != nil {
		err2 := t.db.Find(&tags).Error

		if err2 != nil {
			return nil, err
		}

		if err == redis.Nil {
			go func() {
				tagsToJson, _ := json.Marshal(tags)
				t.rdb.Set("blog:tags", string(tagsToJson), utils.RandExpTime(60, 70))
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

func (t *TagsRepositoryImpl) GetTagName(id int) (string, error) {
	var tagname string
	tagname, err := t.rdb.Get(fmt.Sprintf("blog:tag:info:name:%d", id))

	if err != nil {

		err2 := t.db.Model(&entity.Tag{}).Select("tag_name").Where("tag_id = ?", id).Pluck("tag_name", &tagname).Error
		if err2 != nil {
			return "", err2
		}

		if err == redis.Nil {
			go func() {
				t.rdb.Set(fmt.Sprintf("blog:tag:info:name:%d", id), tagname, utils.RandExpTime(60, 70))
			}()
		}

		return tagname, nil
	}

	return tagname, nil
}

func (t *TagsRepositoryImpl) GetTagId(tagname string) (int, error) {
	var tagid int

	id, err := t.rdb.Get("blog:tag:info:id:" + tagname)

	if err != nil {

		err2 := t.db.Model(&entity.Tag{}).Select("tag_id").Where("tag_name = ?", tagname).Pluck("tag_id", &tagid).Error
		if err2 != nil {
			return -1, err2
		}

		if err == redis.Nil {
			go func() {
				t.rdb.Set("blog:tag:info:id:"+tagname, strconv.Itoa(tagid), utils.RandExpTime(60, 70))
			}()
		}

		return tagid, nil
	}

	tagid, err = strconv.Atoi(id)
	if err != nil {
		return -1, err

	}

	return tagid, nil
}

func (t *TagsRepositoryImpl) AddTags(tags []entity.Tag) error {
	var err1 error

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err1 = t.rdb.Del("blog:tags")
		wg.Done()
	}()

	err2 := t.db.Create(&tags).Error
	wg.Wait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("%v\n%v", err1, err2)
	}
}

func (t *TagsRepositoryImpl) DeleteTags(tags []entity.Tag) error {
	var err1 error

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		errStr := ""
		err11 := t.rdb.Del("blog:tags")
		if err11 != nil {
			errStr += err11.Error() + "\n"
		}

		for _, tag := range tags {
			err12 := t.rdb.Del(fmt.Sprintf("blog:tag:info:id:%s", tag.TagName))
			err13 := t.rdb.Del(fmt.Sprintf("blog:tag:info:name:%d", tag.TagId))

			if err12 != nil {
				errStr += err12.Error() + "\n"
			}
			if err13 != nil {
				errStr += err13.Error() + "\n"
			}
		}

		wg.Done()
	}()

	err2 := t.db.Delete(&tags).Error
	wg.Wait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("%v\n%v", err1, err2)
	}
}
