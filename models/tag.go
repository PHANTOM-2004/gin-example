package models

import log "github.com/sirupsen/logrus"

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps any) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps any) (count int64) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	})

	return true
}

func EditTag(id int, maps any) bool {
	// 注意这里使用的是Updates, 另外还有一个Update不会更新零值
	db.Model(&Tag{}).Where("id = ?", id).Updates(maps)
	return true
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

func CleanAllTag() bool {
	log.Info("Running tag cleaning")
	defer log.Info("all soft deleted tags cleaned")

	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})
	return true
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	return tag.ID > 0
}

func ExistTag(name string) bool {
	var tag Tag
	// 这里的语法具体看一下文档
	db.Select("id").Where("name = ?", name).First(&tag)

	/*NOTE When query with struct, GORM will only query with those fields has non-zero value,
	* that means if your field’s value is 0, '', false or other zero values,
	* it won’t be used to build query conditions, for example:

		db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
		SELECT * FROM users WHERE name = "jinzhu";
	*/
	return tag.ID > 0
}

/*
* 这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。
gorm所支持的回调方法：

创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
删除：BeforeDelete、AfterDelete
查询：AfterFind

实际上可以猜测到这里可以利用反射, 检测是否实现了这几个方法
如果实现的花就会在框架之中自动调用
* */

// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
// 	err := scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return err
// }
//
// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
// 	err := scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return err
// }
