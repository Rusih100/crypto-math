package crypto_math

import (
	"crypto/rand"
	"math/big"
)

// FermatTest - Тест Ферма.
//
// Вход: n - целое число, n > 1.
//
// Выход: true - "Число n, вероятно, простое" или false - "Число n составное".
func FermatTest(_n *big.Int) bool {

	// Копируем значения, чтобы не менять значения по указателю
	n := new(big.Int)
	n.Set(_n)

	if n.Bit(0) == 0 && n.Cmp(constNum2) != 0 {
		return false
	}

	// n > 0 и n < 5
	if n.Cmp(constNum1) > 0 && n.Cmp(constNum5) < 0 {
		return true
	}

	// n <= 1
	if n.Cmp(constNum1) <= 0 {
		panic("n > 1")
	}

	// Генерируем случайное число 2 ≤ a < n - 1:
	// 0 ≤ a < n - 3	| +2
	a, err := rand.Int(
		rand.Reader,
		new(big.Int).Sub(n, constNum3),
	)
	if err != nil {
		panic(err)
	}
	a = a.Add(a, constNum2)
	//

	r := new(big.Int)
	r = PowMod(
		a,
		new(big.Int).Sub(n, constNum1),
		n,
	)
	if r.Cmp(constNum1) == 0 {
		return true
	}
	return false
}

// SolovayStrassenTest - Тест Соловэя-Штрассена.
//
// Вход: n - целое число, n > 1.
//
// Выход: true - "Число n, вероятно, простое" или false - "Число n составное".
func SolovayStrassenTest(_n *big.Int) bool {

	// Копируем значения, чтобы не менять значения по указателю
	n := new(big.Int)
	n.Set(_n)

	if n.Bit(0) == 0 && n.Cmp(constNum2) != 0 {
		return false
	}

	// n > 0 и n < 5
	if n.Cmp(constNum1) > 0 && n.Cmp(constNum5) < 0 {
		return true
	}

	// n <= 1
	if n.Cmp(constNum1) <= 0 {
		panic("n > 1")
	}

	// Генерируем случайное число 2 ≤ a < n - 1:
	// 0 ≤ a < n - 3	| +2
	a, err := rand.Int(
		rand.Reader,
		new(big.Int).Sub(n, constNum3),
	)
	if err != nil {
		panic(err)
	}
	a = a.Add(a, constNum2)
	//

	r := new(big.Int)

	r = PowMod(
		a,
		new(big.Int).Div(
			new(big.Int).Sub(n, constNum1),
			constNum2),
		n,
	)

	// r != 1 и r != n - 1
	if r.Cmp(constNum1) != 0 && r.Cmp(new(big.Int).Sub(n, constNum1)) != 0 {
		return false
	}

	s := Jacobi(a, n)

	// (r - s) (mod n) == 0
	if new(big.Int).Mod(new(big.Int).Sub(r, big.NewInt(s)), n).Sign() == 0 {
		return true
	}
	return false
}

// MillerRabinTest - Тест Миллера-Рабина.
//
// Вход: n - целое число, n > 1.
//
// Выход: true - "Число n, вероятно, простое" или false - "Число n составное".
func MillerRabinTest(_n *big.Int) bool {

	// Копируем значения, чтобы не менять значения по указателю
	n := new(big.Int)
	n.Set(_n)

	// Проверка на четность
	if n.Bit(0) == 0 && n.Cmp(constNum2) != 0 {
		return false
	}

	// n > 0 и n < 5
	if n.Cmp(constNum1) > 0 && n.Cmp(constNum5) < 0 {
		return true
	}

	// n <= 1
	if n.Cmp(constNum1) <= 0 {
		panic("n > 1")
	}

	// Представления числа

	// n - 1
	t := new(big.Int).Sub(n, constNum1)

	s := big.NewInt(0)
	for t.Bit(0) == 0 {
		s = s.Add(s, constNum1)
		t = t.Rsh(t, 1)
	}

	// Количество раундов теста
	const k = 20
	x := new(big.Int)

	for i := 0; i < k; i++ {

		nextIter := false

		// Генерируем случайное число 2 ≤ a < n - 1:
		// 0 ≤ a < n - 3	| +2
		a, err := rand.Int(
			rand.Reader,
			new(big.Int).Sub(n, constNum3),
		)
		if err != nil {
			panic(err)
		}
		a = a.Add(a, constNum2)
		//

		x = PowMod(a, t, n)

		if x.Cmp(constNum1) == 0 || x.Cmp(new(big.Int).Sub(n, constNum1)) == 0 {
			continue
		}

		for j := big.NewInt(0); j.Cmp(new(big.Int).Sub(s, constNum1)) < 0; j.Add(j, constNum1) {

			x = PowMod(x, constNum2, n)

			if x.Cmp(constNum1) == 0 {
				return false
			}
			if x.Cmp(new(big.Int).Sub(n, constNum1)) == 0 {
				nextIter = true
				break
			}
		}

		if nextIter {
			continue
		}

		return false
	}

	return true
}
