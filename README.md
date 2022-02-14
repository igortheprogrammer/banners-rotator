# Banners rotator - [сервис ротации баннеров](https://github.com/OtusGolang/final_project/blob/master/02-banners-rotation.md)

![](https://goreportcard.com/badge/github.com/igortheprogrammer/banners-rotator)

Сервис "Ротатор баннеров" предназначен для выбора наиболее эффективных (кликабельных) баннеров, в условиях меняющихся
предпочтений пользователей и набора баннеров.

---

## Команды

0. `make generate` выполняет необходимую кодогенерацию.
1. `make build` собирает бинарник проекта.
2. `make run` разворачивает сервис в докере.
3. `make stop` останавливает сервис, гасит контейнеры.
4. `make test` запускает юнит-тесты.
5. `make integration-tests` запускает интеграционные тесты.
6. `make lint` запускает линтер golangci-lint.

## API (gRPC) эндпоинты

1. Создание нового слота

```
CreateSlot {"description": string} -> {"id": string, "description": string}
```

2. Создание нового баннера

```
CreateBanner {"description": string} -> {"id": string, "description": string}
```

3. Создание новой группы

```
CreateGroup {"description": string} -> {"id": string, "description": string}
```

4. Создание ротации

```
CreateRotation {"slot_id": int64, "banner_id": int64} -> {"message": string}
```

5. Удаление ротации

```
DeleteRotation {"slot_id": int64, "banner_id": int64} -> {"message": string}
```

6. Создание события клика

```
CreateClickEvent {"slot_id": int64, "banner_id": int64, "group_id": int64} -> {"message": string}
```

7. Получение баннера для отображения в слоте

```
BannerForSlot {"slot_id": int64, "group_id": int64} -> {"id": string, "description": string}
```
