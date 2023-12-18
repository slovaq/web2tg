package repository

import (
	"fmt"

	"github.com/slovaq/web2tg/internal/model"
	"github.com/slovaq/web2tg/internal/repository/sqlite3"
	"github.com/slovaq/web2tg/internal/repository/sqlite3/user"
)

type User interface {
	GetUserByLoginOrName(login string) (*model.User, error)
}
type Repository struct {
	User User
}

func NewRepository(driver iDb) (*Repository, error) {
	switch dr := driver.GetDriverImplementation().(type) {
	// case *mongolib.IMongo:
	// 	return &Repository{
	// 		driverType: dr.DriverType(),
	// 		User:       user.Init(dr),
	// 	}
	case *sqlite3.ISqlite3:
		userSchema, err := user.Init(dr)
		if err != nil {
			return nil, err
		}
		return &Repository{

			User: userSchema,
		}, nil
	default:
		fmt.Printf("dr: %v\n", dr)
	}

	return nil, fmt.Errorf("repository: driver type not supported or config is broken")
}
