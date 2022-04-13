package main

import (
	t "wb/subs"
)

// Запуск сервера из отдельного пакета
func main() {

	// Звгружаем в кеш все данные из бд
	t.Load_cash()
	// Запускаем http-server
	go Start_http()
	// Подписываемся на nats-server
	t.Create_Subscriber()
}
