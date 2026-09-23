package repositories

import (
	"database/sql"
	"gorestapicrud/internal/models"

	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	// r.db - функция для получения подключения к базе данных
	db := r.db
	createdAt := time.Now()

	sqlStatement := `
		INSERT INTO users (
		full_name,
		age,
		habits,
		alive,
		created_at
		)

		VALUES($1, $2, $3, $4, $5)

		RETURNING id
	`

	// db.QueryRow() - функция для выполнения SQL-запроса >
	// возвращает sql.Rowобъект.
	// Затем мы используем эту sql.Row.Scan() - функцию для сканирования возвращенной строки и присвоения значения полю user.Id
	err := db.QueryRow(
		sqlStatement,
		user.FullName,
		user.Age,
		user.Habits,
		user.Alive,
		createdAt).Scan(&user.Id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user models.User, id int) (models.User, error) {
	db := r.db

	sqlStatement := `
	UPDATE users
	SET 
		full_name = $2,
		age = $3,
		habits = $4,
		alive = $5
	WHERE id = $1
	RETURNING id
	`

	err := db.QueryRow(sqlStatement, id, user.FullName, user.Age, user.Habits, user.Alive).Scan(&id)
	if err != nil {
		return models.User{}, err
	}

	user.Id = id

	return user, nil
}

func (r *UserRepository) GetUser(id int) (models.User, error) {
	db := r.db

	sqlStatement := `
	SELECT full_name, age, habits, alive, created_at FROM users WHERE id = $1
	`

	var user models.User

	err := db.QueryRow(sqlStatement, id).Scan(&user.FullName, &user.Age, &user.Habits, &user.Alive, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	user.Id = id

	return user, nil
}

func (r *UserRepository) GetAllUsers(page, limit int) ([]models.User, error) {
	db := r.db

	sqlStatement := `
	SELECT id, full_name, age, habits, alive, created_at FROM users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
	`

	offset := (page - 1) * limit
	// db.Query - возвращает массив строк
	rows, err := db.Query(sqlStatement, limit, offset)
	if err != nil {
		return nil, err
	}

	// rows нужно закрывать для того, чтобы освободить ресурсы сервера и не положить приложение
	defer rows.Close()

	// Создаем слайс, куда будем складывать пользователей
	var users []models.User

	for rows.Next() {
		// Создаем пустой объект пользователя для ТЕКУЩЕЙ строки
		var user models.User

		// Сканируем данные текущей строки в структуру
		err := rows.Scan(&user.Id, &user.FullName, &user.Age, &user.Habits, &user.Alive, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Добавляем пользователя в слайс
		users = append(users, user)
	}

	// Проверяем, не возникло ли ошибок во время итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	db := r.db

	sqlStatement := `
	DELETE FROM users WHERE id = $1
	RETURNING id
	`

	var deletedID int

	err := db.QueryRow(sqlStatement, id).Scan(&deletedID)
	if err != nil {
		return err
	}

	return nil
}
