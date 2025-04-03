# InMemDB

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Made with Go](https://img.shields.io/badge/Made_with-Go-1f425f.svg)](https://go.dev/)
[![Made with Love](https://img.shields.io/badge/Made_with-♥-red.svg)](#)
[![Go Reference](https://img.shields.io/badge/Go-Reference-blue.svg)](https://go.dev/)
[![Tests](https://github.com/DaniilZ77/InMemDB/actions/workflows/ci.yml/badge.svg)](https://github.com/DaniilZ77/InMemDB/actions/workflows/ci.yml)


## Overview

**InMemDB** — это in-memory key-value база данных с асинхронной репликацией. Она спроектирована для высокопроизводительного хранения и быстрого доступа к данным благодаря шардированию и встроенному механизму Write-Ahead Log (WAL) с пакетной записью операций на диск для обеспечения надёжности.

**Основные возможности:**
- Хранение данных в памяти для мгновенного доступа.
- Асинхронная репликация по модели master-slave.
- Шардирование (распределение данных по нескольким shard'ам) для равномерной нагрузки.
- Write-Ahead Log (WAL) для сохранности операций в случае сбоя.
- Простые команды для работы с данными (`SET`, `GET`, `DEL`).

## Grammar

Взаимодействие с InMemDB строится на использовании трёх типов команд:
```ebnf
query          = set_command | get_command | del_command

set_command    = "SET" argument argument
get_command    = "GET" argument
del_command    = "DEL" argument

argument       = punctuation | letter | digit { punctuation | letter | digit }

punctuation    = "*" | "/" | "_" | ...
letter         = "a" | ... | "z" | "A" | ... | "Z"
digit          = "0" | ... | "9"
```

## Quick Start

1. **Склонировать репозиторий**:
    ```bash
    git clone https://github.com/DaniilZ77/InMemDB.git
    cd InMemDB
    ```
2. **Собрать Docker-образ**:
    ```bash
    task docker-build
    ```
3. **Запустить контейнер**:
    ```bash
    task docker-run
    ```
4. **Запустить клиента**:
    ```bash
    task client SERVER_ADDRESS=:3223
    task client SERVER_ADDRESS=:3224
    ```