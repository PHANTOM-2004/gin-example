package models

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
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
