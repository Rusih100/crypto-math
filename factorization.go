package crypto_math

import (
	"crypto/rand"
	bigFloat "github.com/ALTree/bigfloat"
	"github.com/Rusih100/polynomial"
	"math"
	"math/big"
	rnd "math/rand"
	"sort"
	"time"
)

// RoPollardFactor - Ро-метод Полларда.
//
// Вход: Число n, начальное значение c, функция fx (полином)
//
// Выход: Нетривиальный делитель числа p числа n
func RoPollardFactor(_n *big.Int, _c *big.Int, _fx *polynomial.Polynomial) *big.Int {

	// Копируем значения чтобы не менять значения по указателю
	n := new(big.Int)
	n.Set(_n)

	c := new(big.Int)
	c.Set(_c)

	fx := new(polynomial.Polynomial)
	fx.SetPolynomial(_fx)

	// Проверка входных данных
	if n.Sign() <= 0 {
		panic("n > 0")
	}

	if c.Sign() <= 0 || c.Cmp(n) >= 0 {
		panic("0 < c < n")
	}

	a := new(big.Int)
	b := new(big.Int)

	a.Set(c)
	b.Set(c)

	// Засекание времени для отсечки
	startTime := time.Now()
	stopper := 0.0

	// Размер отсчеки взависмости от размера числа
	if n.BitLen() <= 64 {
		stopper = 1
	} else {
		stopper = float64(n.BitLen() / 20)
	}

	for {
		a = fx.Value(a)
		a = a.Mod(a, n)

		b = fx.Value(b)
		b = b.Mod(b, n)
		b = fx.Value(b)
		b = b.Mod(b, n)

		d := new(big.Int)
		d = EuclidAlgorithm(new(big.Int).Sub(a, b), n)

		if d.Cmp(constNum1) > 0 && d.Cmp(n) < 0 {
			return d
		}

		if d.Cmp(n) == 0 {
			return nil
		}

		if time.Since(startTime).Seconds() > stopper {
			return nil
		}
	}
}

// RoOnePollardFactor - (Po - 1)-метод Полларда.
//
// Вход: Число n
//
// Выход: Нетривиальный делитель числа p числа n
func RoOnePollardFactor(_n *big.Int) *big.Int {

	// Копируем значения чтобы не менять значения по указателю
	n := new(big.Int)
	n.Set(_n)

	// Проверка входных данных
	if n.Sign() <= 0 {
		panic("n > 0")
	}

	if n.Cmp(constNum3) <= 0 {
		return nil
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

	d := new(big.Int)
	d = EuclidAlgorithm(a, n)

	if d.Cmp(constNum2) >= 0 {
		return d
	}

	// Для логарифма
	lnN := new(big.Float).SetInt(n)
	lnN = bigFloat.Log(lnN)

	lnN64, _ := lnN.Float64()

	var lnI float64

	for i := 0; i < len(simpleNumbersBase); i++ {
		lnI = math.Log(
			float64(simpleNumbersBase[i]),
		)

		l := int64(lnN64 / lnI)

		power := new(big.Int)
		power = Pow(
			big.NewInt(simpleNumbersBase[i]),
			big.NewInt(l),
		)

		a = PowMod(a, power, n)
	}

	d = EuclidAlgorithm(new(big.Int).Sub(a, constNum1), n)

	if d.Cmp(constNum1) == 0 || d.Cmp(n) == 0 {
		return nil
	}

	return d
}

// BruteForceFactorization - Факторизация числа перебором.
//
// Вход: Число n
//
// Выход: Массив различных делителей числа n
func BruteForceFactorization(_n *big.Int) (result []*big.Int) {

	// Копируем значения, что бы не ихменять их по указателю
	n := new(big.Int)
	n.Set(_n)

	factor := big.NewInt(2)

	// Проверка на 2

	if new(big.Int).Mod(n, factor).Sign() == 0 {
		result = append(result, new(big.Int).Set(factor))
	}

	for new(big.Int).Mod(n, factor).Sign() == 0 {
		n = n.Div(n, factor)
	}

	temp := new(big.Int)

	// Цикл до корня из n, пропуская четные
	for factor = big.NewInt(3); temp.Mul(factor, factor).Cmp(n) <= 0; factor = factor.Add(factor, constNum2) {

		if new(big.Int).Mod(n, factor).Sign() == 0 {
			result = append(result, new(big.Int).Set(factor))
		}

		for new(big.Int).Mod(n, factor).Sign() == 0 {
			n.Div(n, factor)
		}
	}

	if n.Cmp(constNum1) != 0 {
		result = append(result, new(big.Int).Set(n))
	}

	return result
}

// Factorization - Факторизация числа смешанными методами.
//
// Вход: Число n
//
// Выход: Массив различных делителей числа n
//
// Примечание: Может не факторизовать слишком большой делитель
func Factorization(_n *big.Int) (result []*big.Int) {

	// Копируем значения, что бы не ихменять их по указателю
	n := new(big.Int)
	n.Set(_n)

	// Если число меньше 32 бит, перебираем все делители
	if n.BitLen() <= 32 {
		return BruteForceFactorization(n)
	}

	// Иначе методы Полларда

	// Случайные числа
	rnd.Seed(time.Now().UnixNano())

	// Константа шагов
	steps := 10

	if n.BitLen() > 120 {
		steps = n.BitLen() / 20
	}

	// Добавляем факторизумое число в массив
	result = append(result, new(big.Int).Set(n))

	// Переменная для записи делителя
	factor := new(big.Int)

	polyArr := []*big.Int{
		big.NewInt(1),
		big.NewInt(0),
		big.NewInt(1),
	}

	poly := polynomial.NewPolynomial(polyArr)

	for s := 0; s < steps; s++ {

		// Итерация по элементам
		for i := 0; i < len(result); i++ {

			// Пропускаем числа <= 3
			if result[i].Cmp(constNum3) <= 0 {
				continue
			}

			// Случайно выбираем метод факторизации
			choice := rnd.Int31n(2)

			switch choice {
			case 0:
				// Ро метод

				// Генерируем случайное число c
				c, err := rand.Int(
					rand.Reader,
					new(big.Int).Sub(result[i], constNum2),
				)

				if err != nil {
					panic(err)
				}
				c = c.Add(c, constNum1)
				//

				factor = RoPollardFactor(result[i], c, poly)

			case 1:
				// (Ро - 1) метод
				factor = RoOnePollardFactor(result[i])
			}

			// Если делитель найден
			if factor != nil {

				// Проходимся по всем элементам и убираем повторы
				for k := 0; k < len(result); k++ {
					for new(big.Int).Mod(result[k], factor).Sign() == 0 {
						result[k].Div(result[k], factor)
					}
				}

				// Добавляем делитель в конец
				result = append(result, new(big.Int).Set(factor))
			}
		}
	}

	// Сортируем массив
	sort.Slice(result, func(i, j int) bool {
		return result[i].Cmp(result[j]) < 0
	})

	// Удаление 1
	for result[0].Cmp(constNum1) == 0 {
		result = result[1:]
	}

	// Проверка тестом простоты
	for m := 0; m < len(result); m++ {
		if !MillerRabinTest(result[m]) && result[m].BitLen() <= 128 {

			// Уход в рекурсию
			temp := []*big.Int{}
			temp = Factorization(result[m])
			result = append(result, temp...)

			// Сортируем массив
			sort.Slice(result, func(i, j int) bool {
				return result[i].Cmp(result[j]) < 0
			})

			// Убираем повторы
			for i := 1; i < len(result); i++ {
				if result[i-1].Cmp(result[i]) == 0 {
					result = append(result[:i], result[i+1:]...)
					i--
				}
			}
		}
	}
	return result
}
