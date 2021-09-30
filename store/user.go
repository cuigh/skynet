package store

import (
	"context"
	"time"

	"github.com/cuigh/auxo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserStatusBlocked = 0
	UserStatusNormal  = 1
)

type User struct {
	Id         string    `json:"id" bson:"_id"`
	Name       string    `json:"name" bson:"name" valid:"required"`
	LoginName  string    `json:"login_name" bson:"login_name" valid:"required"`
	Roles      []string  `json:"roles" bson:"roles"`
	Email      string    `json:"email,omitempty" bson:"email" valid:"email"`
	Phone      string    `json:"phone,omitempty" bson:"phone"`
	Wecom      string    `json:"wecom,omitempty" bson:"wecom"`
	Admin      bool      `json:"admin" bson:"admin"`
	Status     int32     `json:"status" bson:"status"` // 0-禁用, 1-正常
	Password   string    `json:"-" bson:"password"`
	Salt       string    `json:"-" bson:"salt" mongo:"create:true,update:false,deep:false"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	ModifyTime time.Time `json:"modify_time" bson:"modify_time"`
}

type UserStore interface {
	Create(u *User) error
	Modify(t *User) error
	ModifyProfile(t *User) error
	Find(id string) (u *User, err error)
	FindByName(loginName string) (*User, error)
	Search(name, loginName, filter string, pageIndex, pageSize int64) (users []*User, total int64, err error)
	Fetch(ids []string) (users []*User, err error)
	SetPassword(id string, password, salt string) error
	SetStatus(id string, status int32) error
	Count(ctx context.Context) (int64, error)
	CreateIndexes(ctx context.Context) error
}

type userStore struct {
	c *mongo.Collection
}

func NewUserStore(db *mongo.Database) UserStore {
	return &userStore{db.Collection("user")}
}

func (s *userStore) Find(id string) (u *User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := s.c.FindOne(ctx, bson.M{"_id": id})
	t := &User{}
	if err := r.Decode(&t); err == nil {
		return t, nil
	} else if err == mongo.ErrNoDocuments {
		return nil, nil
	} else {
		return nil, err
	}
}

func (s *userStore) Fetch(ids []string) (users []*User, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.D{
		{"password", 0},
		{"salt", 0},
	})
	cur, err := s.c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userStore) FindByName(loginName string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := s.c.FindOne(ctx, bson.M{"login_name": loginName})
	t := &User{}
	if err := r.Decode(&t); err == nil {
		return t, nil
	} else if err == mongo.ErrNoDocuments {
		return nil, nil
	} else {
		return nil, err
	}
}

func (s *userStore) Search(name, loginName, filter string, pageIndex, pageSize int64) (users []*User, total int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := bson.M{}
	if name != "" {
		args["name"] = name
	}
	if loginName != "" {
		args["login_name"] = loginName
	}
	switch filter {
	case "admin":
		args["admin"] = true
	case "active":
		args["status"] = UserStatusNormal
	case "blocked":
		args["status"] = UserStatusBlocked
	}

	// fetch total count
	total, err = s.c.CountDocuments(ctx, args)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetProjection(bson.D{
		{"password", 0},
		{"salt", 0},
	}).SetSkip(pageSize * (pageIndex - 1)).SetLimit(pageSize).SetSort(bson.M{"login_name": 1})
	cur, err := s.c.Find(ctx, args, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *userStore) SetPassword(id string, pwd, salt string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"password": pwd,
		"salt":     salt,
	}
	_, err = s.c.UpdateByID(ctx, id, bson.M{"$set": update})
	return
}

func (s *userStore) SetStatus(id string, status int32) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.c.UpdateByID(ctx, id, bson.M{"$set": bson.M{"status": status}})
	return
}

func (s *userStore) Create(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.Id = createId()
	u.CreateTime = time.Now()
	u.ModifyTime = u.CreateTime
	_, err := s.c.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errors.Format("登录名 %s 已经存在", u.LoginName)
	}
	return err
}

func (s *userStore) Modify(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"name":        u.Name,
		"login_name":  u.LoginName,
		"roles":       u.Roles,
		"email":       u.Email,
		"phone":       u.Phone,
		"wecom":       u.Wecom,
		"admin":       u.Admin,
		"status":      u.Status,
		"modify_time": time.Now(),
	}
	r, err := s.c.UpdateByID(ctx, u.Id, bson.M{"$set": update})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find user '%d'", u.Id)
	}
	return nil
}

func (s *userStore) ModifyProfile(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"name":        u.Name,
		"login_name":  u.LoginName,
		"email":       u.Email,
		"phone":       u.Phone,
		"modify_time": time.Now(),
	}
	r, err := s.c.UpdateByID(ctx, u.Id, bson.M{"$set": update})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find user '%d'", u.Id)
	}
	return nil
}

func (s *userStore) Count(ctx context.Context) (int64, error) {
	filter := bson.M{}
	return s.c.CountDocuments(ctx, filter)
}

func (s *userStore) CreateIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys:    bson.D{{"login_name", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := s.c.Indexes().CreateOne(ctx, index)
	return err
}
