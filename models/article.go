package models

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"` // index in database
	Tag   Tag `json:"tag"`                 // 这里应该是belong to 关系, 任何Article应该属于一个Tag

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// // hook
// func (a *Article) BeforeCreate(scope *gorm.Scope) error {
// 	// 这里使用的是fieldname, 注意是结构体的名字CreatedOn
// 	return scope.SetColumn("CreatedOn", time.Now().Unix())
// }
//
// // hook
// func (a *Article) BeforeUpdate(scope *gorm.Scope) error {
// 	return scope.SetColumn("ModifiedOn", time.Now().Unix())
// }

func ExistArticleByID(id int) bool {
	var a Article
	db.Select("id").Where("id = ?", id).First(&a)
	return a.ID > 0
}

func GetArticle(id int) (a Article) {
	db.Where("id = ?", id).First(&a) // 首先找到具有当前id 的文章
	// 不过说实话, 还是手动指定比较清晰; 可能组员会不知道这样的约定俗成
	// 这里通过查找a相关的tag填充到a.Tag
	db.Model(&a).Related(&a.Tag)
	/*
			   GORM 支持通过类名（结构体名）和外键（通常是 {结构体名}ID 的形式）来自动建立关联。
		     GORM 的关联默认约定是遵循这种命名约定，但您也可以使用 GORM 标签手动指定外键和引用。

			   GORM 默认的关联规则
			   外键命名约定：

			   GORM 默认使用主表结构体的名称加上 ID 作为外键。
			   例如，如果有两个结构体 User 和 Profile，GORM 默认会认为 Profile 中的 UserID 字段是外键，指向 User 表的 ID 字段。
			   关联类型：

			   has one（一对一）：主表中包含一个外键指向相关表。
			   belongs to（多对一）：子表中包含一个外键指向主表。
			   has many（一对多）：主表中包含多个外键指向相关表。
			   many to many（多对多）：使用中间表来保存关联。
			   例子
			   假设我们有两个结构体：User 和 Profile，并且我们希望通过 UserID 字段来建立关联。

			   go
			   Copy code
			   type User struct {
			       ID      int
			       Name    string
			       Profile Profile // 一对一关系
			   }

			   type Profile struct {
			       ID     int
			       UserID int    // 外键
			       Bio    string
			   }
			   在这种情况下，GORM 会自动识别 Profile 中的 UserID 字段作为外键，指向 User 表中的 ID 字段。

			   使用 GORM 标签自定义关联
			   如果您希望自定义外键名称或者指定更多关联选项，可以使用 GORM 标签。例如：

			   go
			   Copy code
			   type User struct {
			       ID   int
			       Name string
			       Profile Profile `gorm:"foreignKey:CustomUserID"` // 指定外键为 CustomUserID
			   }

			   type Profile struct {
			       ID           int
			       CustomUserID int    // 自定义外键名称
			       Bio          string
			   }
			   在上述例子中，我们使用了 gorm:"foreignKey:CustomUserID" 明确指定了外键字段。

			   总结
			   GORM 可以自动通过类名+ID 的形式建立关联。
			   如果需要自定义外键名称或关联选项，可以使用 gorm 标签指定。
	*/
	return
}

func GetArticleTotal(maps any) int {
	var count int
	db.Model(&Article{}).Where(maps).Count(&count)
	return count
}

func GetArticles(pageNum int, pageSize int, maps any) (a []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&a)
	/*
				type User struct {
				  gorm.Model
				  Username string
				  Orders   []Order
		      一个user有多个order
				}

				type Order struct {
				  gorm.Model
				  UserID uint
				  Price  float64
				}

				// Preload Orders when find users
				db.Preload("Orders").Find(&users)
		    在
				// SELECT * FROM users; 这里假定查到的user_id = (1,2,3,4)
				// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

				db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)
				// SELECT * FROM users;
				// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
				// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
				// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
	*/
	return
}

func EditArticle(id int, data any) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return true
}

func AddArticle(data map[string]any) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}
