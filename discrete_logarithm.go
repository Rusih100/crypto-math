package crypto_math

import (
	"github.com/Rusih100/polynomial"
	"math/big"
)

// DiscreteLogarithm - Дискретное логарифмирование.
//
// Вход: a порядка r по модулю p, b
//
// Выход: x
func DiscreteLogarithm(_a *big.Int, _b *big.Int, _p *big.Int) *big.Int {

	// Копируем значения, чтобы не менять по указателю
	a := new(big.Int)
	b := new(big.Int)
	p := new(big.Int)

	a.Set(_a)
	b.Set(_b)
	p.Set(_p)

	// Проверка входных данных
	if !MillerRabinTest(p) {
		panic("p is a prime number")
	}

	b = b.Mod(b, p)
	if b.Sign() == 0 {
		panic("b != 0")
	}

	// Нахождение порядка r числа a
	r := new(big.Int)
	r = PrimeNumberOrder(a, p)

	if r == nil {
		return nil
	}

	// p / 2
	p2 := new(big.Int).Div(p, constNum2)

	// Полиномы для ветвящегося отображения
	x1Arr := []*big.Int{
		big.NewInt(1),
	}

	x2Arr := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
	}

	x1 := polynomial.NewPolynomial(x1Arr)
	x2 := polynomial.NewPolynomial(x2Arr)

	// Ветвящееся отображение
	fx := func(x *big.Int, logX *polynomial.Polynomial) (*big.Int, *polynomial.Polynomial) {

		y := big.NewInt(0)
		logYArr := []*big.Int{
			big.NewInt(0),
		}
		logY := polynomial.NewPolynomial(logYArr)

		if x.Cmp(p2) < 0 {
			y = y.Mod(new(big.Int).Mul(a, x), p)
			logY = logY.Add(logX, x1)
			return y, logY

		} else {
			y = y.Mod(new(big.Int).Mul(b, x), p)
			logY = logY.Add(logX, x2)
			return y, logY
		}
	}

	// 1. Случайны U и V (Полагаем равными 2)
	u := big.NewInt(2)
	v := big.NewInt(2)

	// Переменные
	c := new(big.Int)
	c = c.Mul(
		PowMod(a, u, p),
		PowMod(b, v, p),
	)
	c = c.Mod(c, p)

	d := new(big.Int)
	d.Set(c)

	// Логарифмы
	logArr := []*big.Int{
		new(big.Int).Set(u),
		new(big.Int).Set(v),
	}

	logC := polynomial.NewPolynomial(logArr)
	logD := polynomial.NewPolynomial(logArr)

	for {
		c, logC = fx(c, logC)

		d, logD = fx(d, logD)
		d, logD = fx(d, logD)

		condition1 := new(big.Int).Mod(c, p)
		condition2 := new(big.Int).Mod(d, p)

		if condition1.Cmp(condition2) == 0 {
			break
		}
	}

	logC = logC.Sub(logD, logC)
	logC = logC.Mod(logC, r)

	item0 := new(big.Int)
	item1 := new(big.Int)

	item0 = logC.Get(0)
	item1 = logC.Get(1)

	item0 = item0.Neg(item0)
	item0 = item0.Mod(item0, r)

	count := new(big.Int)
	x := new(big.Int)
	offset := new(big.Int)

	count, x, offset = ModuloComparisonFirst(item1, item0, r)

	if count.Sign() == 0 {
		return nil
	}

	res := new(big.Int)

	for i := big.NewInt(0); i.Cmp(count) < 0; i.Add(i, constNum1) {

		res = PowMod(a, x, p)
		res = res.Mod(res.Sub(res, b), p)

		if res.Sign() == 0 {
			return x
		}

		if offset != nil {
			x.Add(x, offset)
		}
	}

	return nil
}
