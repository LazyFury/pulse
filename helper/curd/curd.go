package curd

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lazyfury/pulse/framework"
	"github.com/lazyfury/pulse/framework/driver"
	"github.com/lazyfury/pulse/helper/execl"
	"github.com/lazyfury/pulse/helper/response"
	"github.com/lazyfury/pulse/helper/utils"
	"gorm.io/gorm"
)

var searchAct = []string{"in", "not_in", "like", "eq", "gt", "gte", "lte", "lt", "is_null", "is_not_null", "asc", "desc"}

type Pagination struct {
	Page  int
	Limit int
	Query *gorm.DB
}

type CRUD struct {
	DB     *driver.DB
	Model  framework.IModel
	Export execl.IExport
}

func New(db *driver.DB, model framework.IModel) *CRUD {
	return &CRUD{
		DB:    db,
		Model: model,
	}
}

func (c *CRUD) WithExportAttrs(e execl.IExport) *CRUD {
	c.Export = e
	return c
}

func (c *CRUD) GetDB() *driver.DB {
	return c.DB
}

func (c *CRUD) GetTableName() string {
	model := c.Model
	if model == nil {
		return ""
	}

	return model.TableName()
}

func (c *CRUD) GetList(ctx *gin.Context) *Pagination {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	var params map[string]string = make(map[string]string)
	ctx.ShouldBindQuery(&params)
	delete(params, "page")
	delete(params, "limit")
	return &Pagination{
		Page:  page,
		Limit: limit,
		Query: c.Where(utils.StringMapToInterfaceMap(params)),
	}
}

func (c *CRUD) superWhere(action string, tx *gorm.DB, key string, v interface{}) *gorm.DB {
	if key == "" || v == nil {
		return tx
	}
	if !c.isLegalKey(key) {
		log.Warn("illegal key: " + key)
		return tx
	}

	switch action {
	case "in":
		tx = tx.Where(key+" IN (?)", utils.TryInterfaceToStringToArray(v))
	case "not_in":
		tx = tx.Where(key+" NOT IN (?)", utils.TryInterfaceToStringToArray(v))
	case "like":
		tx = tx.Where(key+" LIKE ?", "%"+v.(string)+"%")
	case "eq":
		tx = tx.Where(key+" = ?", v)
	case "gt":
		tx = tx.Where(key+" > ?", v)
	case "lt":
		tx = tx.Where(key+" < ?", v)
	case "gte":
		tx = tx.Where(key+" >= ?", v)
	case "lte":
		tx = tx.Where(key+" <= ?", v)
	case "is_null":
		tx = tx.Where(key + " IS NULL")
	case "is_not_null":
		tx = tx.Where(key + " IS NOT NULL")
	// desc
	case "desc":
		tx = tx.Order(key + " DESC")
	case "asc":
		tx = tx.Order(key + " ASC")
	}

	return tx
}

func (c *CRUD) Where(params map[string]interface{}) *gorm.DB {
	tx := c.DB.Model(c.Model)
	for k, v := range params {

		// 外键查询
		fk_reg := fmt.Sprintf(`^(\S+)__(\S+)__fk__(%s)$`, strings.Join(searchAct, "|"))
		if find := regexp.MustCompile(fk_reg).FindStringSubmatch(k); len(find) > 0 {
			delete(params, k)

			preload := find[1]
			key := find[2]

			opration := find[3]

			if preload != "" {
				tx = tx.Joins(preload)
				tx = c.superWhere(opration, tx, fmt.Sprintf("%s.%s", preload, key), v)
			}
			continue
		}

		where_reg := fmt.Sprintf(`^(\S+)__(%s)$`, strings.Join(searchAct, "|"))
		if find := regexp.MustCompile(where_reg).FindStringSubmatch(k); len(find) > 0 {
			delete(params, k)
			key := find[1]

			if !c.isField(key) {
				continue
			}

			action := find[2]
			// log.Info("find pattern", "key", key, "action", action)
			tx = c.superWhere(action, tx, fmt.Sprintf("%s.%s", c.GetTableName(), key), v)
			continue
		}

		if v == "" {
			delete(params, k)
		}

	}

	params = c.filterModelFields(params)

	// tx = tx.Where(params)

	tx.Order("id DESC")

	for k, v := range params {
		tx = c.superWhere("eq", tx, fmt.Sprintf("%s.%s", c.GetTableName(), k), v)
	}

	return tx
}

func (c *CRUD) isLegalKey(key string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`).MatchString(key)
}

// func get all keys
func (c *CRUD) getModelKeys() []string {
	value := reflect.ValueOf(c.Model).Elem()
	keys := value.NumField()
	keysArr := make([]string, 0)
	for i := range keys {
		if value.Field(i).Kind() == reflect.Ptr {
			// 嵌入的 struct 比如 gin.Model 的 id,craeted_at... 字段
			_val := value.Type().Field(i).Type.Elem()
			if _val.Kind() == reflect.Struct {
				_keys := _val.NumField()
				for j := range _keys {
					keysArr = append(keysArr, _val.Field(j).Tag.Get("json"))
				}
			}

		}
		key := value.Type().Field(i).Tag.Get("json")
		if key == "" {
			continue
		}
		keysArr = append(keysArr, key)
	}

	return keysArr
}

func (c *CRUD) isField(key string) bool {
	return utils.InArray[string](c.getModelKeys(), key)
}

func (c *CRUD) filterModelFields(data map[string]interface{}) map[string]interface{} {
	keysArr := c.getModelKeys()
	for k := range data {
		if !utils.InArray[string](keysArr, k) {
			delete(data, k)
		}
	}
	return data
}

func (c *CRUD) GetListHandler(list any) func(ctx *gin.Context) {
	return c.GetListHandlerWithWhere(list, nil)
}

// get list handler
func (c *CRUD) GetListHandlerWithWhere(list any, extraWhere func(tx *gorm.DB) *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		start := time.Now()
		pagination := c.GetList(ctx)
		var query *gorm.DB = pagination.Query
		page := pagination.Page
		limit := pagination.Limit

		// query
		if extraWhere != nil {
			query = extraWhere(query)
		}
		var count int64 = 0
		query.Count(&count)

		err := query.Offset((page - 1) * limit).Limit(limit).Find(list).Error

		if err != nil {
			response.Error(response.ErrInternalError, err.Error()).Done(ctx)
			return
		}
		end := time.Now()
		response.Success(map[string]any{
			"list":  list,
			"total": count,
			"page":  page,
			"limit": limit,
			"pages": math.Ceil(float64(count) / float64(limit)),
			"time":  strconv.FormatInt(end.Sub(start).Milliseconds(), 10) + "ms",
		}).Done(ctx)
	}
}

func (c *CRUD) GetById(id int64) (interface{}, error) {
	var model interface{}
	err := c.DB.Model(c.Model).Where("id = ?", id).First(&model).Error
	return model, err
}

func (c *CRUD) CreateHandler(ctx *gin.Context, model interface{}, check func(model interface{}) error) {
	if err := ctx.ShouldBindBodyWith(model, binding.JSON); err != nil {
		response.Error(response.ErrBadRequest, "获取参数错误:"+err.Error()).Done(ctx)
		return
	}

	if check != nil {
		if err := check(model); err != nil {
			response.Error(response.ErrBadRequest, err.Error()).Done(ctx)
			return
		}
	}

	if err := c.Create(model); err != nil {
		response.Error(response.ErrInternalError, err.Error()).Done(ctx)
		return
	}
	response.Success(model).Done(ctx)
}

func (c *CRUD) UpdateHandler(ctx *gin.Context, model framework.IModel, check func(model interface{}) error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			response.Error(response.ErrBadRequest, "内部错误 Painc").Done(ctx)
		}
	}()

	if err := ctx.ShouldBindBodyWith(model, binding.JSON); err != nil {
		response.Error(response.ErrBadRequest, "获取参数错误:"+err.Error()).Done(ctx)
		return
	}

	if check != nil {
		if err := check(model); err != nil {
			response.Error(response.ErrBadRequest, err.Error()).Done(ctx)
			return
		}
	}

	if err := c.Update(model); err != nil {
		response.Error(response.ErrInternalError, err.Error()).Done(ctx)
		return
	}
	response.Success(model).Done(ctx)
}

func (c *CRUD) DeleteHandler(ctx *gin.Context, check func(model interface{}) error) {
	type data struct {
		Ids []uint `json:"ids"`
		Id  uint   `json:"id"`
	}

	var d = &data{}

	if err := ctx.ShouldBindBodyWithJSON(d); err != nil {
		response.Error(response.ErrBadRequest, "参数错误:"+err.Error()).Done(ctx)
		return
	}

	if check != nil {
		if err := check(c.Model); err != nil {
			response.Error(response.ErrBadRequest, err.Error()).Done(ctx)
			return
		}
	}

	if d.Id > 0 && len(d.Ids) == 0 {
		d.Ids = []uint{d.Id}
	}

	if len(d.Ids) == 0 {
		response.Error(response.ErrBadRequest, "请选择要删除的记录").Done(ctx)
		return
	}

	if err := c.Delete(d.Ids...); err != nil {
		response.Error(response.ErrInternalError, err.Error()).Done(ctx)
		return
	}

	response.Success("删除成功").Done(ctx)
}

func (c *CRUD) Create(model interface{}) error {
	return c.DB.Model(c.Model).Create(model).Error
}

func (c *CRUD) Update(model interface {
	GetId() uint
}) error {
	return c.DB.Model(c.Model).Where("id = ?", model.GetId()).Save(model).Error
}

func (c *CRUD) Delete(ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}

	return c.DB.Model(c.Model).Where("id in ?", ids).Delete(c.Model).Error
}

func (c *CRUD) ExportHanderWithWhere(list any, extraWhere func(tx *gorm.DB) *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params map[string]string = make(map[string]string)
		ctx.ShouldBindQuery(&params)
		delete(params, "page")
		delete(params, "limit")

		var query *gorm.DB = c.Where(utils.StringMapToInterfaceMap(params))

		// query
		if extraWhere != nil {
			query = extraWhere(query)
		}
		var count int64 = 0
		query.Count(&count)

		err := query.Find(list).Error

		if err != nil {
			response.Error(response.ErrInternalError, err.Error()).Done(ctx)
			return
		}

		if c.Export == nil {
			response.Error(response.ErrInternalError, "导出数据为空").Done(ctx)
			return
		}

		workbook, err := c.Export.Export(list)

		if err != nil {
			response.Error(response.ErrInternalError, err.Error()).Done(ctx)
			return
		}

		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		ctx.Header("Content-Disposition", "attachment; filename="+"list.xlsx")
		err = workbook.Write(ctx.Writer)

		if err != nil {
			response.Error(response.ErrInternalError, err.Error()).Done(ctx)
			return
		}

		ctx.Writer.Flush()
	}
}
