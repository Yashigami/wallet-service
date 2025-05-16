package utils

import "strings"

/*
IsRepeatedError Функция определяет, является ли ошибка транзакционной и поддающейся повторной попытке.
Возвращает true, если текст ошибки содержит признаки конфликтов
*/
func IsRepeatedError(err error) bool {
	if err != nil {
		return false
	}
	msg := err.Error()

	// Проверяем наличие в сообщении ключевых слов, связанных с транзакционными ошибками
	return strings.Contains(msg, "deadlock") ||
		strings.Contains(msg, "could not serialize access") ||
		strings.Contains(msg, "canceling statement due to conflict")
}
