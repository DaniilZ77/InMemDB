# InMemDB

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Made with Go](https://img.shields.io/badge/Made_with-Go-1f425f.svg)](https://go.dev/)
[![Made with Love](https://img.shields.io/badge/Made_with-♥-red.svg)](#)
[![Go Reference](https://img.shields.io/badge/Go-Reference-blue.svg)](https://go.dev/)
[![Tests](https://github.com/DaniilZ77/InMemDB/actions/workflows/ci.yml/badge.svg)](https://github.com/DaniilZ77/InMemDB/actions/workflows/ci.yml)

## Overview

**InMemDB** — in-memory key-value база данных с поддержкой транзакций, MVCC и асинхронной репликацией. Оптимизирована для высоконагруженных приложений благодаря шардированию и продвинутым механизмам управления конкурентным доступом.

**Основные возможности:**
- 🚀 Хранение данных в памяти для мгновенного доступа
- ⚡ Поддержка ACID-транзакций (BEGIN, COMMIT, ROLLBACK)
- 🔄 Multi-Version Concurrency Control (MVCC) для неблокирующих операций
- 🌐 Асинхронная репликация по модели master-slave
- 🧩 Шардирование данных с автоматическим балансированием нагрузки
- 📝 Write-Ahead Log (WAL) с гарантированной записью операций
- 🛠️ Расширенные команды для работы с данными и транзакциями

## Grammar

Обновлённая грамматика с поддержкой транзакций:
```ebnf
query          = set_command | get_command | del_command | begin_command | commit_command | rollback_command

set_command    = "SET" argument argument
get_command    = "GET" argument
del_command    = "DEL" argument
begin_command  = "BEGIN"
commit_command = "COMMIT"
rollback_command = "ROLLBACK"

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
