package crypto_math

import (
	"math/big"
)

// AdvancedEuclidAlgorithm - обобщенный (расширенный) алгоритм Евклида.
//
// Вход: натуральные числа x и y отличные от нуля.
//
// Выход: m, a, b - наибольший общий делитель и его линейное представление.
func AdvancedEuclidAlgorithm(_x *big.Int, _y *big.Int) (m, a, b *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	x := new(big.Int)
	y := new(big.Int)

	x.Set(_x)
	y.Set(_y)

	flagSwap := false

	// x <= 0 или y <= 0
	if x.Sign() <= 0 || y.Sign() <= 0 {
		panic("x and y must be positive numbers other than zero")
	}

	// x < y
	if x.Cmp(y) == -1 {
		x, y = y, x
		flagSwap = true
	}

	a2 := big.NewInt(1)
	a1 := big.NewInt(0)
	b2 := big.NewInt(0)
	b1 := big.NewInt(1)

	for y.BitLen() > 0 {

		// q = x / y
		q := new(big.Int).Div(x, y)

		// r = x - q * y
		r := new(big.Int).Sub(
			x,
			new(big.Int).Mul(q, y),
		)

		// a = a2 - q * a1
		a = new(big.Int).Sub(
			a2,
			new(big.Int).Mul(q, a1),
		)

		// b = b2 - q * b1
		b = new(big.Int).Sub(
			b2,
			new(big.Int).Mul(q, b1),
		)

		x = y
		y = r
		a2 = a1
		a1 = a
		b2 = b1
		b1 = b

		m = x
		a = a2
		b = b2
	}
	if flagSwap {
		a, b = b, a
	}

	return m, a, b
}

// EuclidAlgorithm - алгоритм Евклида для целых чисел.
//
// Вход: целые числа x и y.
//
// Выход: m - наибольший общий делитель.
func EuclidAlgorithm(_x *big.Int, _y *big.Int) *big.Int {

	// Копируем значения, чтобы не менять значения по указателю
	x := new(big.Int)
	y := new(big.Int)

	x.Set(_x)
	y.Set(_y)

	// x < y
	if x.Cmp(y) == -1 {
		x, y = y, x
	}

	// Если числа отрицательные
	if x.Sign() < 0 {
		x = x.Neg(x)
	}

	if y.Sign() < 0 {
		y = y.Neg(y)
	}

	r := big.NewInt(0)

	for y.Sign() != 0 {
		r = r.Mod(x, y)
		x.Set(y)
		y.Set(r)
	}
	return x
}

// Pow - Алгоритм быстрого возведения в степень.
//
// Вход: a - основание (число), n - положительная степень (число).
//
// Выход: result - число a^n.
func Pow(_a *big.Int, _n *big.Int) (result *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	n := new(big.Int)

	a.Set(_a)
	n.Set(_n)

	// n < 0
	if n.Sign() == -1 {
		panic("n must be greater than or equal to zero")
	}

	result = big.NewInt(1)

	for i := 0; i < n.BitLen(); i++ {
		if n.Bit(i) == 1 {
			result = result.Mul(result, a)
		}
		a = a.Mul(a, a)
	}
	return result
}

// PowMod - Алгоритм быстрого возведения в степень по модулю.
//
// Вход: a - основание (число), n - положительная степень (число),
// mod - модуль (положительное число отличное от нуля).
//
// Выход: result - число a^n по модулю mod.
func PowMod(_a *big.Int, _n *big.Int, _mod *big.Int) (result *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	n := new(big.Int)
	mod := new(big.Int)

	a.Set(_a)
	n.Set(_n)
	mod.Set(_mod)

	// n < 0
	if n.Sign() < 0 {
		panic("n must be greater than or equal to zero")
	}

	// mod <= 0
	if mod.Sign() <= 0 {
		panic("mod must be a positive number other than zero")
	}

	result = big.NewInt(1)

	for i := 0; i < n.BitLen(); i++ {
		if n.Bit(i) == 1 {
			result = result.Mod(
				result.Mul(result, a),
				mod,
			)
		}
		a = a.Mod(
			a.Mul(a, a),
			mod,
		)
	}
	return result
}

// InverseElement - Нахождение обратного элемента по модулю через расширенный алгоритм Евклида.
//
// Вход: a > 0, mod > 0.
//
// Выход: Обратный элемент к a по модулю mod.
func InverseElement(_a *big.Int, _mod *big.Int) (result *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	mod := new(big.Int)

	a.Set(_a)
	mod.Set(_mod)

	// Проверка входных данных
	// a <= 0
	if a.Sign() <= 0 {
		panic("a > 0")
	}

	// mod <= 0
	if mod.Sign() <= 0 {
		panic("mod > 0")
	}

	_, _, result = AdvancedEuclidAlgorithm(mod, a)

	result = result.Mod(result, mod)

	return result
}

// Jacobi - Алгоритм вычисления символа Якоби (Алгоритм взят с Википедии).
//
// Вход: a (a: 0 <= a < n) , n - натуральное нечетное больше 1 (n >= 3).
//
// Выход: Символ Якоби - 0, 1 или -1.
func Jacobi(_a *big.Int, _n *big.Int) int64 {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	n := new(big.Int)

	a.Set(_a)
	n.Set(_n)

	// Проверка входных данных
	if n.Bit(0) == 0 {
		panic("n must be odd")
	}

	// n < 3
	if n.Cmp(constNum3) < 0 {
		panic("n must be greater than or equal to 3")
	}

	// a < 0 или a >= n
	if a.Sign() < 0 || a.Cmp(n) >= 0 {
		panic("a: 0 <= a < n")
	}

	// a == 0
	if a.Sign() == 0 {
		return 0
	}

	// 1. Проверка взаимной простоты
	gcd := new(big.Int)
	gcd = EuclidAlgorithm(a, n)

	// gcd != 1
	if gcd.Cmp(constNum1) != 0 {
		return 0
	}

	// 2. Инициализация
	var result int64 = 1

	for {
		// 3. Избавление от четности
		k := big.NewInt(0)
		for a.Bit(0) == 0 {
			k = k.Add(k, constNum1)
			a = a.Rsh(a, 1)
		}

		// k - нечетное и (n (mod 8) = 3 или n (mod 8) = 5)
		if k.Bit(0) == 1 &&
			(new(big.Int).Mod(n, constNum8).Cmp(constNum3) == 0 ||
				new(big.Int).Mod(n, constNum8).Cmp(constNum5) == 0) {
			result = -result
		}

		// 4. Квадратичный закон взамности

		// a (mod 4) = 3 и n (mod 4) = 3
		if new(big.Int).Mod(a, constNum4).Cmp(constNum3) == 0 &&
			new(big.Int).Mod(n, constNum4).Cmp(constNum3) == 0 {
			result = -result
		}
		c := new(big.Int)
		c.Set(a)
		a = a.Mod(n, c)
		n.Set(c)

		// 5. Выход из алгоритма?
		if a.BitLen() == 0 {
			return result
		}
	}
}

// PrimeNumberOrder - Нахождения порядка числа a по простому модулю mod.
//
// Вход: Числа a и mod (Простое)
//
// Выход: Порядок l числа a по модулю mod
func PrimeNumberOrder(_a *big.Int, _mod *big.Int) *big.Int {

	// Копируем значения, чтобы не менять по указателю
	a := new(big.Int)
	mod := new(big.Int)

	a.Set(_a)
	mod.Set(_mod)

	// Проверка модуля
	if mod.Cmp(constNum1) <= 0 {
		panic("mod > 1")
	}

	if !MillerRabinTest(mod) {
		panic("mod is a prime number")
	}

	a = a.Mod(a, mod)

	if a.Sign() == 0 {
		panic("a != 0")
	}

	// Нахождение функции Эйлера
	phi := new(big.Int).Sub(mod, constNum1)

	// Поиск порядка
	result := big.NewInt(0)

	temp := new(big.Int)

	// Перебор всех делителей
	for factor := big.NewInt(1); factor.Cmp(phi) <= 0; factor.Add(factor, constNum1) {

		if new(big.Int).Mod(phi, factor).Sign() == 0 {

			temp = Pow(a, factor)
			temp = temp.Sub(temp, constNum1)
			temp = temp.Mod(temp, mod)

			if temp.Sign() == 0 {
				result.Set(factor)
				break
			}

		}
	}

	if result.Sign() == 0 {
		return nil
	}

	return result
}
