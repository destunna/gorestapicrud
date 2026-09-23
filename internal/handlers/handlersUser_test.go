package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gorestapicrud/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Создаем универсальный Mock-репозиторий.
// Каждое поле — это функция, поведение которой мы можем задавать прямо внутри конкретного теста.
type MockUserRepository struct {
	MockGetAllUsers func(page, limit int) ([]models.User, error)
	MockGetUser     func(id int) (models.User, error)
	MockUpdateUser  func(user models.User, id int) (models.User, error)
	MockCreateUser  func(user models.User) (models.User, error)
	MockDeleteUser  func(id int) error
}

func (m *MockUserRepository) GetAllUsers(page, limit int) ([]models.User, error) {
	if m.MockGetAllUsers != nil {
		return m.MockGetAllUsers(page, limit)
	}

	return []models.User{}, errors.New("MockGetAllUsers не переопределен в тесте")
}

func (m *MockUserRepository) GetUser(id int) (models.User, error) {
	if m.MockGetUser != nil {
		return m.MockGetUser(id)
	}

	return models.User{}, errors.New("MockGetUser не переопределен в тесте")
}

func (m *MockUserRepository) UpdateUser(user models.User, id int) (models.User, error) {
	if m.MockUpdateUser != nil {
		return m.MockUpdateUser(user, id)
	}

	return models.User{}, errors.New("MockUpdateUser не переопределен в тесте")
}

func (m *MockUserRepository) CreateUser(user models.User) (models.User, error) {
	if m.MockCreateUser != nil {
		return m.MockCreateUser(user)
	}

	return models.User{}, nil
}

func (m *MockUserRepository) DeleteUser(id int) error {
	if m.MockDeleteUser != nil {
		return m.MockDeleteUser(id)
	}

	return errors.New("MockDeleteUser не переопределен в тесте")
}

func TestHandleGetAllUsers(t *testing.T) {
	e := echo.New()

	// Настраиваем mock: возвращаем слайс из двух пользователей
	mockRepo := &MockUserRepository{
		MockGetAllUsers: func(page, limit int) ([]models.User, error) {
			// Проверяем, что хэндлер передал правильные параметры пагинации
			assert.Equal(t, 1, page)
			assert.Equal(t, 2, limit) // Лимит жестко задан в вашем хэндлере равным 2

			return []models.User{
				{FullName: "Alex", Age: "15", Habits: "Smoke", Alive: true},
				{FullName: "Mary", Age: "20", Habits: "None", Alive: true},
			}, nil
		},
	}

	// Инициализируем хэндлер с нашей заглушкой
	h := NewUserHandler(mockRepo)

	// Имитируем запрос: GET /users?page=1
	req := httptest.NewRequest(http.MethodGet, "/users?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем тестируемый метод хэндлера
	err := h.HandleGetAllUsers(c)

	// Проверяем результаты (Asserts)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code) // Ожидаем статус 200

	// Распаковываем полученный JSON-массив
	var resUsers []models.User
	err = json.Unmarshal(rec.Body.Bytes(), &resUsers)
	assert.NoError(t, err)

	// Проверяем, что вернулось ровно 2 пользователя и данные корректны
	assert.Len(t, resUsers, 2)
	assert.Equal(t, "Alex", resUsers[0].FullName)
	assert.Equal(t, "Mary", resUsers[1].FullName)
}

func TestHandleGetUser(t *testing.T) {
	e := echo.New()
	mockRepo := &MockUserRepository{
		MockGetUser: func(id int) (models.User, error) {
			return models.User{
				FullName: "Alex",
				Age:      "15",
				Habits:   "Smoke",
				Alive:    true,
			}, nil
			// Описываем жестко заданный ответ: функция сразу возвращает заполненную модель пользователя и пустую ошибку (nil), имитируя, что пользователь найден в БД.
		},
	}

	h := NewUserHandler(mockRepo) // Создаем хэндлер с нашей заглушкой
	// Инициализируется структура хэндлеров `UserHandler`. Вместо реальной БД мы передаем туда нашу заглушку `mockRepo`.

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	// Имитируется входящий HTTP-запрос методом `GET` на URL `/users/1`. Тело запроса пустое (`nil`), так как это получение данных.

	rec := httptest.NewRecorder()
	// Создается виртуальный «диктофон» (ResponseRecorder), который перехватит и запишет в себя ответ сервера (статус-код и JSON), не отправляя его в сеть.

	c := e.NewContext(req, rec)
	// Создается контекст `echo.Context`, связывающий наш фейковый запрос (`req`) и устройство для записи ответа (`rec`).

	c.SetParamNames("id")
	// Имитируется работа роутера Echo: регистрируется имя параметра пути, который мы ожидаем вытащить из URL (в данном случае это `:id`).

	c.SetParamValues("1")
	// Этому параметру `:id` присваивается конкретное текстовое значение `"1"`. Хэндлер внутри прочитает его через `c.Param("id")`.

	err := h.HandleGetUser(c)
	// Вызывается сам тестируемый метод. Внутри него код прочитает ID "1", вызовет `mockRepo.GetUser` и запишет JSON-ответ в наш `rec`.

	assert.NoError(t, err)
	// Проверяется, что сам метод `HandleGetUser` не вернул системную ошибку Go (ожидается `nil`).

	assert.Equal(t, http.StatusOK, rec.Code)
	// Проверяется, что в «диктофон» записался HTTP статус-код `200 OK`. Если код другой, тест упадет.

	var resUser models.User
	// Объявляется пустая переменная структуры `UserUser`, в которую мы попытаемся распаковать ответ от хэндлера.

	err = json.Unmarshal(rec.Body.Bytes(), &resUser)
	// Считываются байты тела ответа из «диктофона» (там лежит JSON-строка) и декодируются обратно в структуру `resUser`.

	assert.NoError(t, err) // Убедимся, что JSON успешно распарсился
	// Проверяется, что хэндлер вернул валидный JSON и перевод из строки в структуру прошел без ошибок.

	assert.Equal(t, "Alex", resUser.FullName)
	// Проверяется, совпадает ли имя пользователя в ответе хэндлера с тем, что нам выдал Mock репозитория (`"Alex"`).

	assert.Equal(t, "15", resUser.Age)
	// Проверяется, что возраст в JSON-ответе равен `"15"`.

	assert.Equal(t, "Smoke", resUser.Habits)
	// Проверяется, что привычки в JSON-ответе равны `"Smoke"`.

	assert.True(t, resUser.Alive)
	// Проверяется финальное булево поле: статус жизни пользователя должен быть равен `true`.
}

func TestHandleUpdateUser(t *testing.T) {
	e := echo.New()
	mockRepo := &MockUserRepository{
		MockUpdateUser: func(user models.User, id int) (models.User, error) {
			return models.User{
				FullName: "Alexis",
				Age:      "18",
				Habits:   "Smoke a lot",
				Alive:    true,
			}, nil
		},
	}

	h := NewUserHandler(mockRepo)

	// Создаем фейковый JSON для отправки в теле PUT-запроса
	inputUser := models.User{FullName: "Alexis", Age: "18"}
	bodyBytes, _ := json.Marshal(inputUser)

	// Передаем bodyBytes вместо nil
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(bodyBytes))
	// Обязательно указываем заголовок, чтобы c.Bind() понял, что это JSON
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Добавляем параметры URL, чтобы strconv.Atoi(id) успешно отработал
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.HandleUpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resUser models.User
	err = json.Unmarshal(rec.Body.Bytes(), &resUser)
	assert.NoError(t, err)

	assert.Equal(t, "Alexis", resUser.FullName)
	assert.Equal(t, "18", resUser.Age)
	assert.Equal(t, "Smoke a lot", resUser.Habits)
	assert.True(t, resUser.Alive)
}

func TestHandleCreateUser(t *testing.T) {
	e := echo.New()
	mockRepo := &MockUserRepository{
		MockCreateUser: func(user models.User) (models.User, error) {
			return models.User{
				FullName: "Alexis",
				Age:      "18",
				Habits:   "Smoke a lot",
				Alive:    true,
			}, nil
		},
	}

	h := NewUserHandler(mockRepo)

	inputUser := models.User{
		FullName: "Alexis",
		Age:      "18",
		Habits:   "Smoke a lot",
		Alive:    true,
	}

	bodyBytes, _ := json.Marshal(inputUser)

	// Передаем bodyBytes вместо nil
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	// Обязательно указываем заголовок, чтобы c.Bind() понял, что это JSON
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleCreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resUser models.User
	err = json.Unmarshal(rec.Body.Bytes(), &resUser)
	assert.NoError(t, err)

	assert.Equal(t, "Alexis", resUser.FullName)
	assert.Equal(t, "18", resUser.Age)
	assert.Equal(t, "Smoke a lot", resUser.Habits)
	assert.True(t, resUser.Alive)
}

func TestHandleDeleteUser(t *testing.T) {
	e := echo.New()
	mockRepo := &MockUserRepository{
		MockDeleteUser: func(id int) error {
			return nil
		},
	}

	h := NewUserHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPut, "/users/1", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.HandleDeleteUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resUser models.User
	err = json.Unmarshal(rec.Body.Bytes(), &resUser)
	assert.NoError(t, err)
}
