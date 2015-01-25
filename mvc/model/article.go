package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"strings"
	"time"
)

const (
	ARTICLE_STATUS_PUBLIC ArticlePublishStatus = "public"
	ARTICLE_STATUS_DRAFT  ArticlePublishStatus = "draft"

	ARTICLE_CONTENT_MARKDOWN ArticleContentType = "markdown"
	ARTICLE_CONTENT_HTML     ArticleContentType = "html"
)

type ArticlePublishStatus string
type ArticleContentType string

var (
	article_brief_seperator = "<!--more-->"
	article_id_key          = "article:id:%d"
	article_id_list_key     = "article:id-list"
	article_id_pub_list_key = "article:id-pub-list"
	article_short_key       = "article:short"
	article_slug_key        = "article:slug"

	article_content_key   = "article:content:%d"
	article_read_meta_key = "article:read:%d"
	article_hit_meta_key  = "article:hit:%d"
)

func generateArticleID(offset int64) int64 {
	diff := time.Now().Unix() - core.Config.InstallTime
	key := fmt.Sprintf(article_id_key, diff)
	if core.Db.Exist(key) {
		return generateArticleID(offset + 1)
	}
	return diff/1800 + 1 + offset
}

func generateArticleShort(title string) string {
	short := util.MD5(title, fmt.Sprint(time.Now().UnixNano()))[0:8]
	if core.Db.HExist(article_short_key, short) {
		return generateArticleShort(title)
	}
	return short
}

// article data struct
type Article struct {
	Id       int64
	Title    string
	Slug     string
	Short    string
	AuthorId int64 // article author user id

	CreateTime     time.Time
	LastUpdateTime time.Time

	PublishStatus ArticlePublishStatus

	BriefContent string
	ContentType  ArticleContentType

	contentTemp string // content temp variable

	EnableComment bool // enable comment

	readMeta *ArticleReadMeta // dynamic loading, not auto-load
	hitMeta  *ArticleHitMeta
	content  *ArticleContent
}

func NewArticle(title, slug, content string, contentType ArticleContentType, status ArticlePublishStatus, author int64) *Article {
	a := &Article{
		Id:             generateArticleID(0),
		Title:          title,
		Slug:           slug,
		Short:          generateArticleShort(title),
		AuthorId:       author,
		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
		PublishStatus:  status,
		BriefContent:   content, // use whole content as brief if first new
		ContentType:    contentType,
		contentTemp:    content,
		EnableComment:  true,
	}
	a.setBrief()
	return a
}

func (a *Article) setBrief() {
	if strings.Contains(a.contentTemp, article_brief_seperator) {
		tmp := strings.Split(a.contentTemp, article_brief_seperator)
		a.BriefContent = tmp[0]
	}
}

/*
===== save article
*/
func (a *Article) Save() error {
	var err error
	// save article data
	if err = a.saveArticle(); err != nil {
		return err
	}
	// save content
	if err = a.saveContent(); err != nil {
		return err
	}
	// save indexes
	if err = a.saveIndexes(); err != nil {
		return err
	}
	// save meta
	if err = a.saveHitMeta(); err != nil {
		return err
	}
	if err = a.saveReadMeta(); err != nil {
		return err
	}
	return nil
}

func (a *Article) saveArticle() error {
	key := fmt.Sprintf(article_id_key, a.Id)
	return core.Db.SetJson(key, a)
}

func (a *Article) saveIndexes() error {
	var err error
	// save short index
	if err = core.Db.HSet(article_short_key, a.Short, util.Int642Bytes(a.Id)); err != nil {
		return err
	}
	// save slug index
	if err = core.Db.HSet(article_slug_key, a.Slug, util.Int642Bytes(a.Id)); err != nil {
		return err
	}
	// save id list index with short
	if err = core.Db.ZSet(article_id_list_key, a.Id, []byte(a.Short)); err != nil {
		return err
	}
	if a.PublishStatus == ARTICLE_STATUS_PUBLIC {
		if err = core.Db.ZSet(article_id_pub_list_key, a.Id, []byte(a.Short)); err != nil {
			return err
		}
	} else {
		if err = core.Db.ZDel(article_id_pub_list_key, []byte(a.Short)); err != nil {
			return err
		}
	}
	return nil
}

func (a *Article) saveContent() error {
	a.content = &ArticleContent{
		ArticleId:   a.Id,
		ContentType: a.ContentType,
		Content:     a.contentTemp,
	}
	key := fmt.Sprintf(article_content_key, a.Id)
	return core.Db.SetJson(key, a.content)
}

func (a *Article) saveReadMeta() error {
	a.readMeta = &ArticleReadMeta{
		Words: util.WordCount(a.contentTemp),
	}
	a.readMeta.ReadingTime = util.ReadingTimeCount(a.readMeta.Words)
	key := fmt.Sprintf(article_read_meta_key, a.Id)
	return core.Db.SetJson(key, a.readMeta)
}

func (a *Article) saveHitMeta() error {
	a.hitMeta = &ArticleHitMeta{1, 0}
	key := fmt.Sprintf(article_hit_meta_key, a.Id)
	return core.Db.SetJson(key, a.hitMeta)
}

/*
===== remove article
*/

func (a *Article) Remove() error {
	var err error
	// remove article data
	if err = a.removeArticle(); err != nil {
		return err
	}
	// remove content
	if err = a.removeContent(); err != nil {
		return err
	}
	// remove indexes
	if err = a.removeContent(); err != nil {
		return err
	}
	// remove meta
	if err = a.removeContent(); err != nil {
		return err
	}
	if err = a.removeContent(); err != nil {
		return err
	}
	return nil
}

func (a *Article) removeArticle() error {
	key := fmt.Sprintf(article_id_key, a.Id)
	return core.Db.Del(key)
}

func (a *Article) removeIndexes() error {
	var err error
	// save short index
	if err = core.Db.HDel(article_short_key, a.Short); err != nil {
		return err
	}
	// save slug index
	if err = core.Db.HDel(article_slug_key, a.Slug); err != nil {
		return err
	}
	// save id list index with short
	if err = core.Db.ZDel(article_id_list_key, []byte(a.Short)); err != nil {
		return err
	}
	if a.PublishStatus == ARTICLE_STATUS_PUBLIC {
		if err = core.Db.ZDel(article_id_pub_list_key, []byte(a.Short)); err != nil {
			return err
		}
	}
	return nil
}

func (a *Article) removeContent() error {
	key := fmt.Sprintf(article_content_key, a.Id)
	return core.Db.Del(key)
}

func (a *Article) removeReadMeta() error {
	key := fmt.Sprintf(article_read_meta_key, a.Id)
	return core.Db.Del(key)
}

func (a *Article) removeHitMeta() error {
	key := fmt.Sprintf(article_hit_meta_key, a.Id)
	return core.Db.Del(key)
}

/*
===== update article
*/

func (a *Article) UpdateTo(id int64) error {
	oldArticle, err := GetArticleById(id)
	if err != nil || oldArticle == nil {
		return err
	}
	// overwrite some data by old article
	a.Id = oldArticle.Id
	a.Short = oldArticle.Short

	// clean indexes
	if err = a.removeIndexes(); err != nil {
		return err
	}

	// save article data
	if err = a.saveArticle(); err != nil {
		return err
	}
	// save content
	if err = a.saveContent(); err != nil {
		return err
	}
	// save indexes
	if err = a.saveIndexes(); err != nil {
		return err
	}
	if err = a.saveReadMeta(); err != nil {
		return err
	}

	return nil
}

type ArticleReadMeta struct {
	Words       int64
	ReadingTime int64
}

type ArticleHitMeta struct {
	Hits     int64
	Comments int64
}

type ArticleContent struct {
	ArticleId   int64
	Content     string
	ContentType ArticleContentType
}

/*
===== get article
*/

func GetArticleById(id int64) (*Article, error) {
	key := fmt.Sprintf(article_id_key, id)
	a := new(Article)
	if err := core.Db.GetJson(key, a); err != nil {
		return nil, err
	}
	if a.Id != id {
		return nil, nil
	}
	return a, nil

}
