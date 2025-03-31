package main

import (
	stack2 "github.com/golang-collections/collections/stack"
	"strconv"
)

func FibonacciIterative(n int) int {
	// Функція вираховує і повертає n-не число фібоначчі
	// Імплементація без використання рекурсії
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	}
	res := 1
	prev := 1
	acc := 0
	for i := 2; i < n; i++ {
		acc = res + prev
		prev = res
		res = acc
	}
	return acc
}

func FibonacciRecursive(n int) int {
	// Функція вираховує і повертає n-не число фібоначчі
	// Імплементація з використанням рекурсії
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func IsPrime(n int) bool {
	// Функція повертає `true` якщо число `n` - просте.
	// Інакше функція повертає `false`
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsBinaryPalindrome(n int) bool {
	// Функція повертає `true` якщо число `n` у бінарному вигляді є паліндромом
	// Інакше функція повертає `false`
	//
	// Приклади:
	// Число 7 (111) - паліндром, повертаємо `true`
	// Число 5 (101) - паліндром, повертаємо `true`
	// Число 6 (110) - не є паліндромом, повертаємо `false`
	direct := strconv.FormatInt(int64(n), 2)
	for i := 0; i < len(direct)/2; i++ {
		if direct[i] != direct[len(direct)-1-i] {
			return false
		}
	}
	return true
}

func ValidParentheses(s string) bool {
	// Функція повертає `true` якщо у вхідній стрічці дотримані усі правила високристання дужок
	// Правила:
	// 1. Допустимі дужки `(`, `[`, `{`, `)`, `]`, `}`
	// 2. У кожної відкритої дужки є відповідна закриваюча дужка того ж типу
	// 3. Закриваючі дужки стоять у правильному порядку
	//    "[{}]" - правильно
	//    "[{]}" - не правильно
	// 4. Кожна закриваюча дужка має відповідку відкриваючу дужку
	stack := stack2.Stack{}
	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack.Push(char)
		case ')':
			top := stack.Pop()
			if top != '(' {
				return false
			}
		case ']':
			top := stack.Pop()
			if top != '[' {
				return false
			}
		case '}':
			top := stack.Pop()
			if top != '{' {
				return false
			}
		}
	}
	return true
}

func Increment(num string) int {
	// Функція на вхід отримує стрічку яка складається лише з символів `0` та `1`
	// Тобто стрічка містить певне число у бінарному вигляді
	// Потрібно повернути число на один більше
	val, _ := strconv.ParseInt(num, 2, 32)
	return int(val) + 1
}
