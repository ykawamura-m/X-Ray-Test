package repository

import (
	"X-Ray-Test/model"
	"context"
	"crypto/rand"
	"errors"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-xray-sdk-go/xray"
	driver "github.com/go-sql-driver/mysql"
	"github.com/guregu/dynamo"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ユーザー登録
// DBにユーザー情報を登録する
func RegisterUser(ctx context.Context, name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return registerUser_MySQL(ctx, name, email, tel)
	case model.DYNAMO:
		return registerUser_DynamoDB(ctx, name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー登録
// MySQLにユーザー情報を登録する
func registerUser_MySQL(ctx context.Context, name string, email string, tel string) error {
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

	return db.WithContext(ctx).Create(user).Error
}

// ユーザー登録
// DynamoDBにユーザー情報を登録する
func registerUser_DynamoDB(ctx context.Context, name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(os.Getenv("DYNAMO_TABLENAME"))

	user := model.User{
		ID:    newUserID(),
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).RunWithContext(ctx)
}

// ユーザー更新
// DBのユーザー情報を更新する
func UpdateUser(ctx context.Context, id string, name string, email string, tel string, db int) error {
	switch db {
	case model.MYSQL:
		return updateUser_MySQL(ctx, id, name, email, tel)
	case model.DYNAMO:
		return updateUser_DynamoDB(ctx, id, name, email, tel)
	}

	return errors.New("invalid db")
}

// ユーザー更新
// MySQLのユーザー情報を更新する
func updateUser_MySQL(ctx context.Context, id string, name string, email string, tel string) error {
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

	return db.WithContext(ctx).Select("*").Updates(user).Error
}

// ユーザー更新
// DynamoDBのユーザー情報を更新する
func updateUser_DynamoDB(ctx context.Context, id string, name string, email string, tel string) error {
	db := openDynamoDB()
	table := db.Table(os.Getenv("DYNAMO_TABLENAME"))

	user := model.User{
		ID:    id,
		Name:  name,
		Email: email,
		Tel:   tel,
		DB:    model.DYNAMO,
	}

	return table.Put(user).RunWithContext(ctx)
}

// ユーザー削除
// DBのユーザー情報を削除する
func DeleteUser(ctx context.Context, id string, db int) error {
	switch db {
	case model.MYSQL:
		return deleteUser_MySQL(ctx, id)
	case model.DYNAMO:
		return deleteUser_DynamoDB(ctx, id)
	}

	return errors.New("invalid db")
}

// ユーザー削除
// MySQLのユーザー情報を削除する
func deleteUser_MySQL(ctx context.Context, id string) error {
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

	return db.WithContext(ctx).Where("id = ?", user.ID).Delete(user).Error
}

// ユーザー削除
// DynamoDBのユーザー情報を削除する
func deleteUser_DynamoDB(ctx context.Context, id string) error {
	db := openDynamoDB()
	table := db.Table(os.Getenv("DYNAMO_TABLENAME"))

	return table.Delete("ID", id).RunWithContext(ctx)
}

// ユーザー取得
// DBからユーザー情報を取得する
func GetUser(ctx context.Context, id string, db int) (model.User, error) {
	switch db {
	case model.MYSQL:
		return getUser_MySQL(ctx, id)
	case model.DYNAMO:
		return getUser_DynamoDB(ctx, id)
	}

	return model.User{}, errors.New("invalid db")
}

// ユーザー取得
// MySQLからユーザー情報を取得する
func getUser_MySQL(ctx context.Context, id string) (model.User, error) {
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

	err = db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// ユーザー取得
// DynamoDBからユーザー情報を取得する
func getUser_DynamoDB(ctx context.Context, id string) (model.User, error) {
	db := openDynamoDB()
	table := db.Table(os.Getenv("DYNAMO_TABLENAME"))

	var user model.User
	err := table.Get("ID", id).OneWithContext(ctx, &user)

	return user, err
}

// 全ユーザー取得
// DBから全てのユーザー情報を取得する
func GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User

	tmp, err := getAllUsers_MySQL(ctx)
	if err != nil {
		return nil, err
	}
	users = append(users, tmp...)

	tmp, err = getAllUsers_DynamoDB(ctx)
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
func getAllUsers_MySQL(ctx context.Context) ([]model.User, error) {
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
	err = db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 全ユーザー取得
// DynamoDBから全てのユーザー情報を取得する
func getAllUsers_DynamoDB(ctx context.Context) ([]model.User, error) {
	db := openDynamoDB()
	table := db.Table(os.Getenv("DYNAMO_TABLENAME"))

	var users []model.User
	err := table.Scan().AllWithContext(ctx, &users)

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
	config := driver.Config{
		Net:                  "tcp",
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.UTC,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	instrumentedDB, err := xray.SQLContext("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	dialector := mysql.New(mysql.Config{
		Conn: instrumentedDB,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, nil
}

// DynamoDB接続
func openDynamoDB() *dynamo.DB {
	baseSess := session.Must(session.NewSession())
	sess := xray.AWSSession(baseSess)
	return dynamo.New(sess, &aws.Config{Region: aws.String(os.Getenv("DYNAMO_REGION"))})
}
