package repository

import (
	"X-Ray-Test/model"
	"crypto/rand"
	"errors"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MYSQL_USER       = "admin"
	MYSQL_PASS       = "password"
	MYSQL_HOST       = "x-ray-test.cdeh7vyviaaz.ap-northeast-1.rds.amazonaws.com"
	MYSQL_PROTOCOL   = "tcp(" + MYSQL_HOST + ":3306)"
	MYSQL_DBNAME     = "x_ray_test"
	DYNAMO_REGION    = "ap-northeast-1"
	DYNAMO_TABLENAME = "x-ray-test"
)

// ユーザー登録
// DBにユーザー情報を登録する
func RegisterUser(name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return registerUser_MySQL(name, email, tel)
	case model.DYNAMO:
		return registerUser_DynamoDB(name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー登録
// MySQLにユーザー情報を登録する
func registerUser_MySQL(name string, email string, tel string) error {
	db, err := openMySQL()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	user := model.User{
		ID:    newUserID(),
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.MYSQL,
	}

	return db.Create(user).Error
}

// ユーザー登録
// DynamoDBにユーザー情報を登録する
func registerUser_DynamoDB(name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	user := model.User{
		ID:    newUserID(),
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).Run()
}

// ユーザー更新
// DBのユーザー情報を更新する
func UpdateUser(id string, name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return updateUser_MySQL(id, name, email, tel)
	case model.DYNAMO:
		return updateUser_DynamoDB(id, name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー更新
// MySQLのユーザー情報を更新する
func updateUser_MySQL(id string, name string, email string, tel string) error {
	db, err := openMySQL()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	user := model.User{
		ID:    id,
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.MYSQL,
	}

	return db.Select("*").Updates(user).Error
}

// ユーザー更新
// DynamoDBのユーザー情報を更新する
func updateUser_DynamoDB(id string, name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	user := model.User{
		ID:    id,
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).Run()
}

// ユーザー削除
// DBのユーザー情報を削除する
func DeleteUser(id string, db int) error {
	switch db {
	case model.MYSQL:
		return deleteUser_MySQL(id)
	case model.DYNAMO:
		return deleteUser_DynamoDB(id)
	}

	return errors.New("invalid db")
}

// ユーザー削除
// MySQLのユーザー情報を削除する
func deleteUser_MySQL(id string) error {
	db, err := openMySQL()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	user := model.User{
		ID: id,
	}

	return db.Where("id = ?", user.ID).Delete(user).Error
}

// ユーザー削除
// DynamoDBのユーザー情報を削除する
func deleteUser_DynamoDB(id string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	return table.Delete("ID", id).Run()
}

// ユーザー取得
// DBからユーザー情報を取得する
func GetUser(id string, db int) (model.User, error) {
	switch db {
	case model.MYSQL:
		return getUser_MySQL(id)
	case model.DYNAMO:
		return getUser_DynamoDB(id)
	}

	return model.User{}, errors.New("invalid db")
}

// ユーザー取得
// MySQLからユーザー情報を取得する
func getUser_MySQL(id string) (model.User, error) {
	var user model.User

	db, err := openMySQL()
	if err != nil {
		return user, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return user, err
	}
	defer sqlDB.Close()

	err = db.Where("id = ?", id).Take(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// ユーザー取得
// DynamoDBからユーザー情報を取得する
func getUser_DynamoDB(id string) (model.User, error) {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	var user model.User
	err := table.Get("ID", id).One(&user)

	return user, err
}

// 全ユーザー取得
// DBから全てのユーザー情報を取得する
func GetAllUsers() ([]model.User, error) {
	var users []model.User

	tmp, err := getAllUsers_MySQL()
	if err != nil {
		return nil, err
	}
	users = append(users, tmp...)

	tmp, err = getAllUsers_DynamoDB()
	if err != nil {
		return nil, err
	}
	users = append(users, tmp...)

	sort.SliceStable(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, nil
}

// 全ユーザー取得
// MySQLから全てのユーザー情報を取得する
func getAllUsers_MySQL() ([]model.User, error) {
	db, err := openMySQL()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	var users []model.User
	err = db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 全ユーザー取得
// DynamoDBから全てのユーザー情報を取得する
func getAllUsers_DynamoDB() ([]model.User, error) {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	var users []model.User
	err := table.Scan().All(&users)

	return users, err
}

// ID採番
func newUserID() string {
	ms := ulid.Timestamp(time.Now().In(time.UTC))
	entropy := rand.Reader
	return ulid.MustNew(ms, entropy).String()
}

// MySQL接続
func openMySQL() (*gorm.DB, error) {
	dsn := MYSQL_USER + ":" + MYSQL_PASS + "@" + MYSQL_PROTOCOL + "/" + MYSQL_DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// DynamoDB接続
func openDynamoDB() *dynamo.DB {
	sess := session.Must(session.NewSession())
	return dynamo.New(sess, &aws.Config{Region: aws.String(DYNAMO_REGION)})
}
