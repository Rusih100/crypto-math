package crypto_math

import (
	"math/big"
	rnd "math/rand"
	"time"
)

// RandNumber - Генерация k-битного случайного нечетного числа.
//
// Вход: Разрядность k генерируемого числа.
//
// Выход: Случайное k-битное нечетное число.
func RandNumber(k int) (result *big.Int) {

	// Проверка входных данных
	if k <= 1 {
		panic("k > 1")
	}

	randNumber := new(big.Int)

	// Получаем k битное число из 1: 1000...001
	randNumber = randNumber.SetBit(randNumber, k-1, 1)
	randNumber = randNumber.SetBit(randNumber, 0, 1)

	// Случайные числа
	rnd.Seed(time.Now().UnixNano())

	// Побитовая догенерация случайного числа с помощью OR
	for i := 1; i < randNumber.BitLen()-1; i++ {
		bit := rnd.Int31n(2)
		if bit == 1 {
			randNumber = randNumber.SetBit(randNumber, i, uint(bit))
		}
	}
	return randNumber
}

// SimpleNumber - Генерация k-битного простого числа.
//
// Вход: Разрядность k искомого простого числа, параметр t >= 1.
//
// Выход: Число, простое с вероятностью 1 - 1 / (4**t).
func SimpleNumber(k int, t int) (result *big.Int) {

	// Проверка входных данных
	if k <= 1 {
		panic("k > 1")
	}

	if t < 1 {
		panic("t >= 1")
	}

	// Генерируем случайное нечетное k-битное число
	randNumber := new(big.Int)
	randNumber = RandNumber(k)

	for i := 0; i < t; i++ {
		if !MillerRabinTest(randNumber) {
			randNumber = randNumber.Add(randNumber, constNum1)
			i = 0
		}
	}
	return randNumber
}
