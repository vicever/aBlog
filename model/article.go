package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"strings"
	"time"
)

var (
	article_brief_separator = "<!--more-->"

	article_content_key    = "article:content:%d"
	article_meta_key       = "article:meta:%d"
	article_slug_key       = "article:slug:%d"
	article_short_link_key = "article:short-link:%d"
	article_id_key         = "article:id:%d"
)

type Article struct {
	Id           int64
	Title        string
	SlugLink     string // customize link url
	ShortLink    string // short link url
	AuthorUserId int64

	CreateTime     time.Time
	LastUpdateTime time.Time
	PublishStatus  string

	Brief         string
	EnableComment bool
	FormatType    string

	contentTemp string // maintain content if no saved
}

func NewArticle(title string, author int64, status string, format string, content string) *Article {
	article := &Article{
		Title:         title,
		AuthorUserId:  author,
		PublishStatus: status,
		FormatType:    format,
		contentTemp:   content,

		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
		ShortLink:      generateArticleShortLink(title),
		Id:             generateArticleID(0),
	}
	return article
}

func (a *Article) Save() error {
	// save article data
	var err error
	key := fmt.Sprintf(article_id_key, a.Id)
	if err = core.Db.SetJson(key, a); err != nil {
		return err
	}

	// save indexes
	if err = a.saveIndexes(); err != nil {
		return err
	}

	// save meta
	meta := newArticleMeta(a.Id, int64(len(a.contentTemp)))
	if err = meta.Save(); err != nil {
		return err
	}

	// save whole content
	cnt := newArticleContent(a.Id, a.contentTemp)
	if err = cnt.Save(); err != nil {
		return err
	}

	return nil
}

func (a *Article) saveIndexes() error {
	// short link index
	key := fmt.Sprintf(article_short_link_key, a.ShortLink)
	if err := core.Db.Set(key, util.Int642Bytes(a.Id)); err != nil {
		return err
	}

	// slug link index
	key = fmt.Sprintf(article_slug_key, a.SlugLink)
	if err := core.Db.Set(key, util.Int642Bytes(a.Id)); err != nil {
		return err
	}
	return nil
}

func (a *Article) removeIndexes() error {
	key := fmt.Sprintf(article_short_link_key, a.ShortLink)
	if err := core.Db.Del(key); err != nil {
		return err
	}

	key = fmt.Sprintf(article_slug_key, a.SlugLink)
	if err := core.Db.Del(key); err != nil {
		return err
	}
	return nil
}

func (a *Article) Remove() error {
	var err error
	// remove article
	key := fmt.Sprintf(article_id_key, a.Id)
	if err = core.Db.Del(key); err != nil {
		return err
	}

	// remove meta
	if meta, _ := GetArticleMeta(a.Id); meta != nil {
		if err = meta.Remove(); err != nil {
			return err
		}
	}

	// remove content
	if cnt, _ := GetArticleContent(a.Id); cnt != nil {
		if err = cnt.Remove(); err != nil {
			return err
		}
	}

	// remove indexes
	return a.removeIndexes()
}

// replace current article data into existing article by id
func (a *Article) Replace(articleId int64) error {
	a.Id = articleId
	// overwrite article data
	var err error
	key := fmt.Sprintf(article_id_key, a.Id)
	if err = core.Db.SetJson(key, a); err != nil {
		return err
	}

	// overwrite indexes
	if err = a.saveIndexes(); err != nil {
		return err
	}

	// when replacing, refresh meta instead of overwrite
	meta, err := GetArticleMeta(articleId)
	if err != nil {
		return err
	}
	meta.CalReadTime(int64(len(a.contentTemp)))
	if err = meta.Save(); err != nil {
		return err
	}

	// overwrite whole content
	cnt := newArticleContent(a.Id, a.contentTemp)
	if err = cnt.Save(); err != nil {
		return err
	}

	return nil
}

func generateArticleID(offset int64) int64 {
	diff := time.Now().Unix() - core.Config.InstallTime
	key := fmt.Sprintf(article_id_key, diff)
	if core.Db.Exist(key) {
		return generateArticleID(offset + 1)
	}
	return diff/1800 + 1 + offset
}

func generateArticleShortLink(title string) string {
	short := util.MD5(title, fmt.Sprint(time.Now().UnixNano()))[0:8]
	key := fmt.Sprintf(article_short_link_key, short)
	if core.Db.Exist(key) {
		return generateArticleShortLink(title)
	}
	return short
}

/*
===== article meta, including reads count, comments count and read time info
*/

type ArticleMeta struct {
	ArticleId int64
	Comments  int64
	Reads     int64
	Words     int64 // words count
	ReadTime  int64 // proper read time
}

func newArticleMeta(articleId int64, words int64) *ArticleMeta {
	meta := &ArticleMeta{
		ArticleId: articleId,
		Comments:  0,
		Reads:     1,
	}
	meta.CalReadTime(words)
	return meta
}

func GetArticleMeta(articleId int64) (*ArticleMeta, error) {
	key := fmt.Sprintf(article_meta_key, articleId)
	meta := &ArticleMeta{}
	if err := core.Db.GetJson(key, meta); err != nil {
		return nil, err
	}
	if articleId != meta.ArticleId {
		return nil, nil
	}
	return meta, nil
}

func (aMeta *ArticleMeta) Save() error {
	key := fmt.Sprintf(article_meta_key, aMeta.ArticleId)
	return core.Db.SetJson(key, aMeta)
}

func (aMeta *ArticleMeta) Remove() error {
	key := fmt.Sprintf(article_meta_key, aMeta.ArticleId)
	return core.Db.Del(key)
}

// increase comments and reads count
func (aMeta *ArticleMeta) Incr(comments, reads int64) {
	aMeta.Comments += comments
	aMeta.Reads += reads
}

// calculate read time by word counts
func (aMeta *ArticleMeta) CalReadTime(wordCount int64) {
	aMeta.Words = wordCount
	aMeta.ReadTime = wordCount / 10 // todo : how to calculate read time and how to display it
}

/*
===== article content, saving whole text for article
*/

type ArticleContent struct {
	ArticleId int64
	Content   string
}

func newArticleContent(articleId int64, content string) *ArticleContent {
	return &ArticleContent{
		ArticleId: articleId,
		Content:   content,
	}
}

func GetArticleContent(articleId int64) (*ArticleContent, error) {
	key := fmt.Sprintf(article_content_key, articleId)
	cnt := &ArticleContent{
		ArticleId: articleId,
	}
	bytes, err := core.Db.Get(key)
	if err != nil {
		return nil, err
	}
	cnt.Content = string(bytes)
	return cnt, nil
}

// save article's content
func (aCnt *ArticleContent) Save() error {
	key := fmt.Sprintf(article_content_key, aCnt.ArticleId)
	return core.Db.Set(key, []byte(aCnt.Content))
}

// get brief from whole content
func (aCnt *ArticleContent) GetBrief() string {
	contentSlice := strings.Split(aCnt.Content, article_brief_separator)
	if len(contentSlice) != 2 {
		return ""
	}
	return contentSlice[0]
}

// remove content
func (aCnt *ArticleContent) Remove() error {
	key := fmt.Sprintf(article_content_key, aCnt.ArticleId)
	return core.Db.Del(key)
}
