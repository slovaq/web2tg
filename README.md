# Web stack
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/sqlite
    go get -u github.com/go-chi/chi
    go get -u -v github.com/slovaq/web2tg/web

## Как же все это использовать?
### API
Если необходимо создать пользователя то используйте **это**
`localhost:1111/api/user_create?`
`localhost:1111/api/user_create?login=analwormx&name=123&password=123`
А если узнать, существуют ли пользователь, то **это**
/api/user_get?login=analwormx&name=123&password=123`

#### Поля
	login - логин пользователя ( пример: torvarlds) 
	name - Linus Torvalds (или какое-то другое? прим. при регистрации только) 
	password - пароль

#### Что возращает, если пользователь существует или создан

```json
{
  "Success": true,
  "Result": {
    "ID": 6,
    "Name": "123",
    "Login": "analwormx",
    "Password": "hidden"
  },
  "Error": null
}
```
#### Что возращает, когда юзера нет или пустые поля
```json
{
  "Success": false,
  "Result": null,
  "Error": "login zhopa not found"
}
```
# Работа с городами.
Создание города
`localhost:1111/api/city_create?name=Vladimir&chat_id=0110`
`localhost:1111/api/city_create`
#### Поля
    name - Название города
	chat_id - ID чата
#### Возращает
```json
{
    "Success": true,
    "Result": {
        "ID": 2,
        "Name": "Vladimir",
        "ChatID": 110
    },
    "Error": null
}
```
Получение всех городов

`localhost:1111/api/city_getAll`

```json 
{
    "Success": true,
    "Result": [
        {
            "ID": 1,
            "Name": "Moscow",
            "ChatID": 0
        },
        {
            "ID": 2,
            "Name": "Vladimir",
            "ChatID": 110
        }
    ],
    "Error": null
}
```
Получение города по имени
`localhost:1111/api/city_get?city_name=Moscow`
```json
{
    "Success": true,
    "Result": {
        "ID": 1,
        "Name": "Moscow",
        "ChatID": 0
    },
    "Error": null
}
```
## Работа с записями
** Требует аунтефикацию**
### Поля
	login - логин юзера
	password - пароль юзера
	cities - города куда отправлять
	text - текст записи
	date - не реализовано
### Что возращает
```json
{
    "Success": true,
    "Result": {
        "ID": 5,
        "Date": "2021-01-09T05:39:37.713505371+05:00",
        "Destination": null,
        "Text": "ungabunga"
    },
    "Error": null
}
```
