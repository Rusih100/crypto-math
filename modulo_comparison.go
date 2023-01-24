package crypto_math

import "math/big"

// ModuloComparisonFirst - Решение сравнения первой степени.
//
// Вход: Сравнение вида ax = b (по модулю mod): числа a, b и mod. (a > 0, mod > 0).
//
// Выход: Количество решений, первое решение, сдвиг для получения следущего решения.
func ModuloComparisonFirst(_a *big.Int, _b *big.Int, _mod *big.Int) (countSolutions, x1, offset *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	b := new(big.Int)
	mod := new(big.Int)

	a.Set(_a)
	b.Set(_b)
	mod.Set(_mod)

	// Проверка входных данных
	// a == 0
	if a.Sign() == 0 {
		panic("a != 0")
	}

	// mod <= 0
	if mod.Sign() <= 0 {
		panic("mod > 0")
	}

	// Переход к положительным числам
	// a (mod) + mod
	a = new(big.Int).Add(new(big.Int).Mod(a, mod), mod)
	b = new(big.Int).Add(new(big.Int).Mod(b, mod), mod)

	// Проверяем разрешимость сравнения
	gcd := new(big.Int)
	gcd = EuclidAlgorithm(a, mod)

	// Если неразрешимо
	// b (mod gcd) != 0
	if new(big.Int).Mod(b, gcd).Sign() != 0 {
		return big.NewInt(0), nil, nil
	}

	// Единственное решение
	// gcd == 1
	if gcd.Cmp(constNum1) == 0 {

		x := new(big.Int)

		// Записываем в x обратный к а элемент, далее умножаем на b
		x = InverseElement(a, mod)
		x = x.Mod(new(big.Int).Mul(x, b), mod)

		return big.NewInt(1), x, nil
	}

	// Множество решений

	// Переход к новому сравнению
	a1 := new(big.Int).Div(a, gcd)
	b1 := new(big.Int).Div(b, gcd)
	mod1 := new(big.Int).Div(mod, gcd)

	x := new(big.Int)
	// Записываем в x обратный к а элемент, далее умножаем на b
	x = InverseElement(a1, mod1)
	x = x.Mod(new(big.Int).Mul(x, b1), mod1)

	return gcd, x, mod1
}

// ModuloComparisonSecond - Решение сравнения второй степени.
//
// Вход: Сравнение вида x^2 = a (по модулю p): числа a и p. (p - простое и p > 2).
//
// Выход: Решение сравнения второй степени.
func ModuloComparisonSecond(_a *big.Int, _p *big.Int) (xPos, xNeg *big.Int) {

	// Копируем значения, чтобы не менять значения по указателю
	a := new(big.Int)
	p := new(big.Int)

	a.Set(_a)
	p.Set(_p)

	// Проверка входных данных
	// p <= 2
	if p.Cmp(constNum2) <= 0 {
		panic("p > 2")
	}

	if !MillerRabinTest(p) {
		panic("p is a prime number")
	}

	// Переход к положительным числам
	a = a.Mod(a, p)

	// a == 0
	if a.Sign() == 0 {
		panic("a is not divisible by p")
	}

	// Проверяем квадратичный вычет a
	if Jacobi(a, p) != 1 {
		return nil, nil
	}

	// Перебором ищем квадратичный невычет N
	N := big.NewInt(1)

	// Пока N < p; N++
	for ; N.Cmp(p) < 0; N = N.Add(N, constNum1) {
		if Jacobi(N, p) == -1 {
			break
		}
	}

	// 1. Представление p в виде p = 2^k * h + 1
	h := new(big.Int).Sub(p, constNum1)

	k := big.NewInt(0)
	for h.Bit(0) == 0 {
		k = k.Add(k, constNum1)
		h = h.Rsh(h, 1)
	}

	// 2.
	a1 := new(big.Int)
	a1 = PowMod(
		a,
		new(big.Int).Div(new(big.Int).Add(h, constNum1), constNum2),
		p,
	)

	a2 := new(big.Int)
	a2 = InverseElement(a, p)

	N1 := new(big.Int)
	N1 = PowMod(N, h, p)

	N2 := big.NewInt(1)

	j := big.NewInt(0)

	// 3.
	// i = 0; i <= k - 2; i++
	for i := big.NewInt(0); i.Cmp(new(big.Int).Sub(k, constNum2)) <= 0; i.Add(i, constNum1) {

		// 3.1
		b := new(big.Int).Mod(
			new(big.Int).Mul(a1, N2),
			p,
		)

		// 3.2
		bPow2 := new(big.Int) // Квадрат b
		bPow2 = PowMod(b, constNum2, p)

		c := new(big.Int).Mod(
			new(big.Int).Mul(a2, bPow2),
			p,
		)

		// 3.3
		dPower := new(big.Int)
		dPower = Pow(
			constNum2,
			new(big.Int).Sub(new(big.Int).Sub(k, constNum2), i),
		)

		d := new(big.Int)
		d = PowMod(c, dPower, p)

		// d == 1
		if d.Cmp(constNum1) == 0 {
			j = big.NewInt(0)
		}

		// d == -1
		if d.Cmp(new(big.Int).Add(p, big.NewInt(-1))) == 0 {
			j = big.NewInt(1)
		}

		// 3.4
		N1Power := new(big.Int)
		N1Power = Pow(constNum2, i)
		N1Power = N1Power.Mul(N1Power, j)

		temp := new(big.Int)
		temp = PowMod(N1, N1Power, p)

		N2 = new(big.Int).Mod(
			new(big.Int).Mul(N2, temp),
			p,
		)
	}

	xPos = new(big.Int).Mod(
		new(big.Int).Mul(a1, N2),
		p,
	)

	xNeg = new(big.Int).Mod(
		new(big.Int).Mul(a1, N2),
		p,
	)

	xNeg.Neg(xNeg)
	xNeg = xNeg.Mod(xNeg, p)

	return xPos, xNeg
}

// ModuloComparisonSystem - Решение системы сравнений.
//
// Вход: Массив коэфицентов bArray и массив модулей mArray.
//
// Выход: Решение системы сравнений, если все модули взаимопросты.
func ModuloComparisonSystem(bArray []*big.Int, mArray []*big.Int) (result *big.Int) {

	// Длины массивов
	bArrayLen := len(bArray)
	mArrayLen := len(mArray)

	// Проверка входных данных
	if bArrayLen == 0 {
		panic("bArray: An empty array was passed")
	}

	if mArrayLen == 0 {
		panic("mArray: An empty array was passed")
	}

	if bArrayLen != mArrayLen {
		panic("Arrays of various lengths were transmitted")
	}

	// Проверка взаимопростоты модулей
	testGCD := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	for i := 0; i < mArrayLen; i++ {
		for j := i + 1; j < mArrayLen; j++ {
			x = mArray[i]
			y = mArray[j]

			testGCD = EuclidAlgorithm(x, y)

			if testGCD.Cmp(constNum1) != 0 {
				return nil
			}
		}
	}

	// Ищем произведение модулей
	M := big.NewInt(1)

	for i := 0; i < mArrayLen; i++ {
		x = mArray[i]
		M = M.Mul(M, x)
	}

	// Ищем решение
	result = big.NewInt(0)
	Mj := new(big.Int)
	Nj := new(big.Int)

	for j := 0; j < mArrayLen; j++ {
		bj := bArray[j]
		mj := mArray[j]

		Mj = Mj.Div(M, mj)
		Nj = InverseElement(Mj, mj)

		sub := new(big.Int).Mul(bj, Mj)
		sub = sub.Mod(sub, M)
		sub = sub.Mul(sub, Nj)
		sub = sub.Mod(sub, M)

		result.Add(result, sub)
		result.Mod(result, M)
	}

	return result
}
