package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/utils"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PassageRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (p *PassageRepositoryImpl) AutoMigrate() {
	p.db.AutoMigrate(&entity.Passage{})
}

func (p *PassageRepositoryImpl) SetDb(db *gorm.DB) {
	p.db = db
}

func (p *PassageRepositoryImpl) SetRedis(rdb *cache.Redis) {
	p.rdb = rdb
}

func (p *PassageRepositoryImpl) GetPassage(id string) (entity.Passage, error) {
	var passage entity.Passage
	passageStr, err := p.rdb.Get("blog:post:" + id)
	if err != nil {
		err2 := p.db.Where("passage_id = ?", id).Where("deleted = ?", 0).First(&passage).Error
		if err2 != nil {
			return passage, err
		}

		if err == redis.Nil {
			go func() {
				p.rdb.Set("blog:post:"+id, passage.ToJson(), utils.RandExpTime(60, 70))
			}()
		}

		return passage, nil
	}

	err = json.Unmarshal([]byte(passageStr), &passage)
	if err != nil {
		return passage, err
	}

	return passage, nil
}

func (p *PassageRepositoryImpl) GetPassages() ([]entity.Passage, error) {
	var passages []entity.Passage
	passagesJson, err := p.rdb.Get("blog:posts")
	if err != nil {
		err2 := p.db.
			Select("passage_id",
				"title",
				"author_username",
				"created_at",
				"updated_at").
			Where("deleted = ?", 0).
			Order("created_at DESC").
			Find(&passages).Error

		if err2 != nil {
			return nil, err
		}

		if err == redis.Nil {
			go func() {
				passagesToJson, _ := json.Marshal(passages)
				p.rdb.Set("blog:posts", string(passagesToJson), utils.RandExpTime(60, 70))
			}()
		}

		return passages, nil
	}

	err = json.Unmarshal([]byte(passagesJson), &passages)
	if err != nil {
		return nil, err
	}

	return passages, nil
}

func (p *PassageRepositoryImpl) GetPassagesByDate(start, end time.Time) ([]string, error) {
	var passages []string
	err := p.db.
		Model(&entity.Passage{}).
		Select("passage_id").
		Where("deleted = ?", 0).
		Where("created_at >= ?", start).
		Where("created_at <= ?", end).
		Order("created_at DESC").
		Find(&passages).Error

	if err != nil {
		return nil, err
	}

	return passages, err
}

func (p *PassageRepositoryImpl) AddPassage(passage entity.Passage) error {
	var err1 error

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err1 = p.rdb.Del("blog:posts")
		wg.Done()
	}()

	err2 := p.db.Create(&passage).Error

	wg.Wait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("%v\n%v", err1, err2)
	}
}

func (p *PassageRepositoryImpl) UpdatePassage(passage entity.Passage) error {
	var err1 error

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go p.DelPassageCache(passage.PassageId, wg, &err1)

	err2 := p.db.Model(&passage).Where("passage_id = ?", passage.PassageId).Updates(&passage).Error

	wg.Wait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("%v\n%v", err1, err2)
	}
}

func (p *PassageRepositoryImpl) DeletePassage(id string) error {
	var err1 error

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go p.DelPassageCache(id, wg, &err1)

	passage := entity.NewPassage(entity.WithPassageId(id), entity.WithDeletedAt(time.Now()))
	err2 := p.db.UpdateColumns(&passage).Error

	wg.Wait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("%v\n%v", err1, err2)
	}
}

func (p *PassageRepositoryImpl) DelPassageCache(passageid string, wg *sync.WaitGroup, err *error) {
	err11 := p.rdb.Del("blog:post:" + passageid)
	err12 := p.rdb.Del("blog:posts")

	if err11 == nil && err12 == nil {
		err = nil
	} else {
		*err = fmt.Errorf("%v\n%v", err11, err12)
	}

	wg.Done()
}
