package core

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type M map[string]string

type Repository struct {
	Name        string `bson:"_id,omitempty" json:"name,omitempty" validate:"-"`
	Interval    string `bson:"interval,omitempty" json:"interval,omitempty" validate:"required,cron"`
	Image       string `bson:"image,omitempty" json:"image,omitempty" validate:"required,containsrune=:"`
	StorageDir  string `bson:"storageDir,omitempty" json:"storageDir,omitempty" validate:"required"`
	LogRotCycle int    `bson:"logRotCycle,omitempty" json:"logRotCycle,omitempty" validate:"omitempty,min=0,max=30"`
	Envs        M      `bson:"envs,omitempty" json:"envs,omitempty" validate:"omitempty,dive,keys,required,endkeys,required"`
	Volumes     M      `bson:"volumes,omitempty" json:"volumes,omitempty" validate:"omitempty,dive,keys,required,endkeys,required"`
	User        string `bson:"user,omitempty" json:"user,omitempty" validate:"-"`
	BindIp      string `bson:"bindIp,omitempty" json:"bindIp,omitempty" validate:"omitempty,ip"`
	Retry       int    `bson:"retry,omitempty" json:"retry,omitempty" validate:"omitempty,min=1,max=3"`
	CreatedAt   int64  `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
	UpdatedAt   int64  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"-"`
}

func (c *Core) GetRepository(name string) (*Repository, error) {
	r := new(Repository)
	if err := c.repoColl.FindId(name).One(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Core) AddRepository(repo *Repository) error {
	repo.CreatedAt = time.Now().Unix()
	repo.UpdatedAt = time.Now().Unix()
	return c.repoColl.Insert(*repo)
}

func (c *Core) UpdateRepository(name string, update bson.M) error {
	var set bson.M
	switch v := update["$set"].(type) {
	case map[string]interface{}:
		set = bson.M(v)
	case bson.M:
		set = v
	default:
		set = bson.M{}
	}
	set["updatedAt"] = time.Now().Unix()
	return c.repoColl.UpdateId(name, update)
}

func (c *Core) RemoveRepository(name string) error {
	return c.repoColl.RemoveId(name)
}

func (c *Core) ListRepositories(query, proj bson.M) []Repository {
	result := []Repository{}
	c.repoColl.Find(query).Select(proj).Sort("_id").All(&result)
	return result
}
