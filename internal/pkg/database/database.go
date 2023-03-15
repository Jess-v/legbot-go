package database

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type DatabaseConnection struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type table struct {
	userId   string
	location string
}

func (db DatabaseConnection) Connect() (*pgx.Conn, error) {
	log.Info("Connecting to database...")
	conn, err := pgx.Connect(context.Background(), "postgres://"+db.User+":"+db.Password+"@"+db.Host+":"+db.Port+"/"+db.Name)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (db DatabaseConnection) InitDb() error {
	log.Info("Initializing database...")
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS legbot (userid VARCHAR(255) PRIMARY KEY, location VARCHAR(255) NOT NULL)")
	if err != nil {
		return err
	}
	return nil
}

func (db DatabaseConnection) CheckUser(userId string) (bool, error) {
	log.Info("Checking if user in database...")
	conn, err := db.Connect()
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())
	exists := false
	err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM legbot WHERE userid = $1)", userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db DatabaseConnection) NewUser(userId string, location string) error {
	log.Info("Adding new user to database...")
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "INSERT INTO legbot (userid, location) VALUES ($1, $2)", userId, location)
	if err != nil {
		return err
	}
	return nil
}

func (db DatabaseConnection) UpdateUser(userId string, location string) error {
	log.Info("Updating user in database...")
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE legbot SET location = $1 WHERE userid = $2", location, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db DatabaseConnection) GetLocation(userId string) (string, error) {
	log.Info("Getting location from database...")
	conn, err := db.Connect()
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())
	var location string
	err = conn.QueryRow(context.Background(), "SELECT location FROM legbot WHERE userid = $1", userId).Scan(&location)
	if err != nil {
		return "", err
	}
	return location, nil
}
