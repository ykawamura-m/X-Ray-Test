package repository

import (
	"X-Ray-Test/model"
	"context"
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
func RegisterUser(c context.Context, name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return registerUser_MySQL(c, name, email, tel)
	case model.DYNAMO:
		return registerUser_DynamoDB(c, name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー登録
// MySQLにユーザー情報を登録する
func registerUser_MySQL(c context.Context, name string, email string, tel string) error {
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

	return db.Create(user).WithContext(c).Error
}

// ユーザー登録
// DynamoDBにユーザー情報を登録する
func registerUser_DynamoDB(c context.Context, name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	user := model.User{
		ID:    newUserID(),
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).RunWithContext(c)
}

// ユーザー更新
// DBのユーザー情報を更新する
func UpdateUser(c context.Context, id string, name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return updateUser_MySQL(c, id, name, email, tel)
	case model.DYNAMO:
		return updateUser_DynamoDB(c, id, name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー更新
// MySQLのユーザー情報を更新する
func updateUser_MySQL(c context.Context, id string, name string, email string, tel string) error {
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

	return db.Select("*").Updates(user).WithContext(c).Error
}

// ユーザー更新
// DynamoDBのユーザー情報を更新する
func updateUser_DynamoDB(c context.Context, id string, name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	user := model.User{
		ID:    id,
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).RunWithContext(c)
}

// ユーザー削除
// DBのユーザー情報を削除する
func DeleteUser(c context.Context, id string, db int) error {
	switch db {
	case model.MYSQL:
		return deleteUser_MySQL(c, id)
	case model.DYNAMO:
		return deleteUser_DynamoDB(c, id)
	}

	return errors.New("invalid db")
}

// ユーザー削除
// MySQLのユーザー情報を削除する
func deleteUser_MySQL(c context.Context, id string) error {
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

	return db.Where("id = ?", user.ID).Delete(user).WithContext(c).Error
}

// ユーザー削除
// DynamoDBのユーザー情報を削除する
func deleteUser_DynamoDB(c context.Context, id string) error {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	return table.Delete("ID", id).RunWithContext(c)
}

// ユーザー取得
// DBからユーザー情報を取得する
func GetUser(c context.Context, id string, db int) (model.User, error) {
	switch db {
	case model.MYSQL:
		return getUser_MySQL(c, id)
	case model.DYNAMO:
		return getUser_DynamoDB(c, id)
	}

	return model.User{}, errors.New("invalid db")
}

// ユーザー取得
// MySQLからユーザー情報を取得する
func getUser_MySQL(c context.Context, id string) (model.User, error) {
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

	err = db.Where("id = ?", id).Take(&user).WithContext(c).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// ユーザー取得
// DynamoDBからユーザー情報を取得する
func getUser_DynamoDB(c context.Context, id string) (model.User, error) {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	var user model.User
	err := table.Get("ID", id).OneWithContext(c, &user)

	return user, err
}

// 全ユーザー取得
// DBから全てのユーザー情報を取得する
func GetAllUsers(c context.Context) ([]model.User, error) {
	var users []model.User

	tmp, err := getAllUsers_MySQL(c)
	if err != nil {
		return nil, err
	}
	users = append(users, tmp...)

	tmp, err = getAllUsers_DynamoDB(c)
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
func getAllUsers_MySQL(c context.Context) ([]model.User, error) {
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
	err = db.Find(&users).WithContext(c).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 全ユーザー取得
// DynamoDBから全てのユーザー情報を取得する
func getAllUsers_DynamoDB(c context.Context) ([]model.User, error) {
	db := openDynamoDB()
	table := db.Table(DYNAMO_TABLENAME)

	var users []model.User
	err := table.Scan().AllWithContext(c, &users)

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
