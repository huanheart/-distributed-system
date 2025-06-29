package mysql

import (
	"MyChat/config"
	"MyChat/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

func InitMysql() error {
	host := config.GetConfig().MysqlHost
	port := config.GetConfig().MysqlPort
	dbname := config.GetConfig().MysqlDatabaseName
	username := config.GetConfig().MysqlUser
	password := config.GetConfig().MysqlPassword
	charset := config.GetConfig().MysqlCharset

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)

	var log logger.Interface
	if gin.Mode() == "debug" {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	return migration()
}

func GetMusicfile(user_id int64, file_id string) (*model.MusicFile, error) {
	musicfile := new(model.MusicFile)
	err := DB.Where("user_id = ? AND uuid = ?", user_id, file_id).First(musicfile).Error
	return musicfile, err
}

func GetMusicfileByFileId(file_id string) (*model.MusicFile, error) {
	musicfile := new(model.MusicFile)
	err := DB.Where("uuid = ?", file_id).First(musicfile).Error
	return musicfile, err
}

// 查找对应userid 用户对音乐music_id 的raction行为
func GetMusicReaction(userID int64, musicUUID string) (*model.MusicReaction, error) {
	reaction := new(model.MusicReaction)
	err := DB.Where("user_id = ? AND music_uuid = ?", userID, musicUUID).First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //没有找到，返回 nil
		}
		return nil, err // 其他错误
	}
	return reaction, nil // 找到了
}

func GetUserByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("email = ?", email).First(user).Error
	return user, err
}

func GetUserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("username = ?", username).First(user).Error
	return user, err
}

func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	return user, err
}

func InsertMusicFile(file *model.MusicFile) (*model.MusicFile, error) {
	err := DB.Create(&file).Error
	return file, err
}

func MarkMusicFileUploaded(filePath string, value int64) error {
	return DB.Model(&model.MusicFile{}).
		Where("file_path = ?", filePath).
		Update("is_upload", value).Error
}

func SetCountDuration(filePath string, value float64) error {
	return DB.Model(&model.MusicFile{}).
		Where("file_path = ?", filePath).
		Update("duration", value).Error
}

func migration() error {
	return DB.AutoMigrate(
		new(model.User),
		new(model.MusicFile),
		new(model.MusicReaction),
	)
}
