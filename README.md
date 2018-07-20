# utils [![Build Status](https://www.travis-ci.org/kodix/utils.svg?branch=master)](https://www.travis-ci.org/kodix/utils) [![Go Report Card](https://goreportcard.com/badge/github.com/kodix/utils)](https://goreportcard.com/report/github.com/kodix/utils)

## Health
Пакет с http обработчиком для эндпоинта health (показывает количество текущих соединений и лимит) и 
посредником, управляющим счетчиком соединений.

## Fields
Пакет с типом "Массив строк" для удобной работы с БД и json

## Must
Будь проще, пиши меньше

## MW
Пакет с набором http-посредников:
    
    Auth - проверяет http-заголовки запроса 
    Metrics - добавляет метрики Prometheus'a
    Logs - копирует стандартный логгер и передает его в реквест с уникальным префиксом