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
	/*
	   type User struct {
	     gorm.Model
	     Username string
	     Orders   []Order
	   }

	   type Order struct {
	     gorm.Model
	     UserID uint
	     Price  float64
	   }

	   // Preload Orders when find users
	   db.Preload("Orders").Find(&users)
	   // SELECT * FROM users;
	   // SELECT * FROM orders WHERE user_id IN (1,2,3,4);

	   db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)
	   // SELECT * FROM users;
	   // SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
	   // SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
	   // SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
	*/
	db.Where("id = ?", id).Preload("Tag").First(&a) // 首先找到具有当前id 的文章
	return
}

func GetArticleTotal(maps any) (count int64) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
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
