package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/rajaghaneshan/go-fullstack/api/models"
)

var users = []models.User{
	{
		Nickname: "Alexander",
		Email:    "alexander@gmail.com",
		Password: "password",
	}, {
		Nickname: "Bob",
		Email:    "bobthebuilder@gmail.com",
		Password: "password",
	}, {
		Nickname: "Cathy",
		Email:    "cathy321@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello from Title 1",
	}, {
		Title:   "Title 2",
		Content: "Hello from Title 2",
	}, {
		Title:   "Title 3",
		Content: "Hello from Title 3",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("error attaching foreign key: %v:", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v:", err)
		}

		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
