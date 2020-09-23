package database

import (
	"context"

	"github.com/arefiev1997/goapi/internal/config"
	"github.com/go-pg/pg"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database interface {
	GetClasses(ctx context.Context) (result []Class, err error)
	CreateClass(ctx context.Context, c Class) error
	DeleteClass(ctx context.Context, id int) error
	GetStudentsByClass(ctx context.Context, classID int) (result []Student, err error)
	CreateStudent(ctx context.Context, student Student) error
	Close()
}

type DB struct {
	conn *sqlx.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn: conn,
	}, nil
}

type Class struct {
	tableName struct{} `sql:"class"`
	ID        int
	Number    int
	Letter    string
}

type Student struct {
	tableName  struct{} `sql:"student"`
	ID         int
	Name       string
	Surname    string
	Patronymic string
	Age        int
	ClassID    int `db:"class" pg:"class" sql:"class"`
}

func (d *DB) GetClasses(ctx context.Context) (result []Class, err error) {
	q := "SELECT id, number, letter FROM class;"
	if err = d.conn.SelectContext(ctx, &result, q); err != nil {
		return nil, err
	}
	return result, err
}

func (d *DB) CreateClass(ctx context.Context, c Class) error {
	q := "INSERT INTO class (number, letter) VALUES ($1, $2);"
	_, err := d.conn.ExecContext(ctx, q, c.Number, c.Letter)
	return err
}

func (d *DB) DeleteClass(ctx context.Context, id int) error {
	q := "DELETE FROM class WHERE id = $1;"
	_, err := d.conn.ExecContext(ctx, q, id)
	return err
}

func (d *DB) GetStudentsByClass(ctx context.Context, classID int) (result []Student, err error) {
	q := "SELECT id, name, surname, patronymic, age, class FROM student WHERE class = $1;"
	if err := d.conn.SelectContext(ctx, &result, q, classID); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) CreateStudent(ctx context.Context, student Student) error {
	q := "INSERT INTO student (name, surname, patronymic, age, class) VALUES ($1, $2, $3, $4, $5);"
	_, err := d.conn.ExecContext(ctx, q, student.Name, student.Surname, student.Patronymic, student.Age, student.ClassID)
	return err
}

func (d *DB) Close() {
	d.conn.Close()
}

type DBorm struct {
	db *pg.DB
}

func NewDBorm(path string) (*DBorm, error) {
	opt, err := pg.ParseURL(path)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)
	return &DBorm{
		db: db,
	}, nil
}

func (d *DBorm) GetClasses(ctx context.Context) (result []Class, err error) {
	if err := d.db.Model(&result).Select(); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBorm) CreateClass(ctx context.Context, c Class) error {
	_, err := d.db.Model(&c).Insert()
	return err
}

func (d *DBorm) DeleteClass(ctx context.Context, id int) error {
	_, err := d.db.Model(&Class{}).Where("id = ?", id).Delete()
	return err
}

func (d *DBorm) GetStudentsByClass(ctx context.Context, classID int) (result []Student, err error) {
	err = d.db.Model(&result).Where("class = ?", classID).Select()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBorm) CreateStudent(ctx context.Context, student Student) error {
	_, err := d.db.Model(&student).Insert()
	return err
}

func (d *DBorm) Close() {
	d.db.Close()
}
