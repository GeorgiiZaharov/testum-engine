![Go](https://img.shields.io/badge/Language-Go-00ADD8?logo=go&logoColor=white)
[![Go Report Card](https://goreportcard.com/badge/github.com/GeorgiiZaharov/testum-engine)](https://goreportcard.com/report/github.com/GeorgiiZaharov/testum-engine)
[![Go Tests](https://github.com/GeorgiiZaharov/testum-engine/actions/workflows/tests.yml/badge.svg)](https://github.com/GeorgiiZaharov/testum-engine/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/GeorgiiZaharov/testum-engine/graph/badge.svg)](https://codecov.io/gh/GeorgiiZaharov/testum-engine)
[![CodeFactor](https://www.codefactor.io/repository/github/georgiizaharov/testum-engine/badge/main)](https://www.codefactor.io/repository/github/georgiizaharov/testum-engine/overview/main)
![GitHub last commit](https://img.shields.io/github/last-commit/GeorgiiZaharov/testum-engine)
[![Release](https://img.shields.io/github/v/release/GeorgiiZaharov/testum-engine?style=flat-square)](https://github.com/GeorgiiZaharov/testum-engine/releases)

# Testum

![Testum Logo](./docs/imgs/logo.jpg)

**Testum** — это веб-ориентированная клиент-серверная система, предназначенная для автоматизации процесса создания, проведения и анализа тестирования студентов.

## Основные возможности

Система позволяет:

- Генерировать тесты из текстового файла специального формата.  
- Управлять доступом к тестам.  
- Проводить тестирование студентов.  
- Реализовывать условную логику отображения блоков вопросов.  
- Хранить и анализировать результаты тестирования.  
- Разграничивать права доступа между ролями:
  - **Студент** — прохождение тестов и просмотр своих результатов.  
  - **Лектор** — создание, управление и анализ тестов.  
  - **Администратор** — управление ролями пользователей.

## Сборка и установка

1. Клонируем репозиторий:
```bash
git clone git@github.com:GeorgiiZaharov/testum-engine.git
cd testum-engine
````

2. Создаём `.env` на основе `.env.example`.

3. Сборка и запуск:

```bash
make
```

4. Запуск тестов с покрытием:

```bash
make test
```

5. Очистка локальной установки Go (опционально):

```bash
make clear
```
