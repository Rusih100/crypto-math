package crypto_math

import (
	"github.com/Rusih100/polynomial"
	"math/big"
	"os"
)

// Реализация базового конечного поля
// Операции выполняются по модулю p

type BaseGaloisField struct {
	p *big.Int
}

// Set - Задает начальное значение p
func (f *BaseGaloisField) Set(p *big.Int) *BaseGaloisField {

	f.p = big.NewInt(0)

	if p.Cmp(constNum2) < 0 {
		panic("p >= 2")
	}

	if !MillerRabinTest(p) {
		panic("p is a prime number")
	}

	f.p.Set(p)

	return f
}

// Add - Складывает два элемента в поле
func (f *BaseGaloisField) Add(a, b *big.Int) *big.Int {

	// Проверки
	if a.Sign() < 0 || a.Cmp(f.p) >= 0 {
		panic("The element a does not belong to the field")
	}

	if b.Sign() < 0 || b.Cmp(f.p) >= 0 {
		panic("The element b does not belong to the field")
	}

	return new(big.Int).Mod(
		new(big.Int).Add(a, b),
		f.p,
	)

}

// Mul - Умножает два элемента в поле
func (f *BaseGaloisField) Mul(a, b *big.Int) *big.Int {

	// Проверки
	if a.Sign() < 0 || a.Cmp(f.p) >= 0 {
		panic("The element a does not belong to the field")
	}

	if b.Sign() < 0 || b.Cmp(f.p) >= 0 {
		panic("The element b does not belong to the field")
	}

	return new(big.Int).Mod(
		new(big.Int).Mul(a, b),
		f.p,
	)
}

// Строковое представление
func (f *BaseGaloisField) String() string {
	return "GF(" + f.p.String() + ")"
}

// CayleyTableAdd - Таблица Кэли для сложения
//
// Файл сохраняется по абсолютному пути
func (f *BaseGaloisField) CayleyTableAdd(path string) {

	name := f.p.String() + "_add"

	// Создание файла
	file, err := os.Create(path + "/" + name + ".csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Первая строка
	result := "\t"

	temp := new(big.Int)

	for i := big.NewInt(0); i.Cmp(f.p) < 0; i.Add(i, constNum1) {
		result = result + i.String()

		if i.Cmp(new(big.Int).Sub(f.p, constNum1)) != 0 {
			result = result + "\t"
		}
	}
	result = result + "\n"

	_, err = file.WriteString(result)
	if err != nil {
		panic(err)
	}

	result = ""

	for i := big.NewInt(0); i.Cmp(f.p) < 0; i.Add(i, constNum1) {
		result = i.String() + "\t"

		for j := big.NewInt(0); j.Cmp(f.p) < 0; j.Add(j, constNum1) {
			temp = f.Add(i, j)
			result = result + temp.String()

			if j.Cmp(new(big.Int).Sub(f.p, constNum1)) != 0 {
				result = result + "\t"
			}
		}
		result = result + "\n"
		_, err = file.WriteString(result)
		if err != nil {
			panic(err)
		}
	}
}

// CayleyTableMul - Таблица Кэли для умножения
//
// Файл сохраняется по абсолютному пути
func (f *BaseGaloisField) CayleyTableMul(path string) {

	name := f.p.String() + "_mul"

	// Создание файла
	file, err := os.Create(path + "/" + name + ".csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Первая строка
	result := "\t"

	temp := new(big.Int)

	for i := big.NewInt(0); i.Cmp(f.p) < 0; i.Add(i, constNum1) {
		result = result + i.String()

		if i.Cmp(new(big.Int).Sub(f.p, constNum1)) != 0 {
			result = result + "\t"
		}
	}
	result = result + "\n"

	_, err = file.WriteString(result)
	if err != nil {
		panic(err)
	}

	result = ""

	for i := big.NewInt(0); i.Cmp(f.p) < 0; i.Add(i, constNum1) {
		result = i.String() + "\t"

		for j := big.NewInt(0); j.Cmp(f.p) < 0; j.Add(j, constNum1) {
			temp = f.Mul(i, j)
			temp.Mod(temp, f.p)

			result = result + temp.String()

			if j.Cmp(new(big.Int).Sub(f.p, constNum1)) != 0 {
				result = result + "\t"
			}
		}
		result = result + "\n"
		_, err = file.WriteString(result)
		if err != nil {
			panic(err)
		}
	}
}

// NewBaseGaloisField - Создает BaseGaloisField и задает ему начальное значение p
func NewBaseGaloisField(p *big.Int) *BaseGaloisField {
	return new(BaseGaloisField).Set(p)
}

// Реализация расширения базового конечного поля

type GaloisField struct {
	p   *big.Int
	n   *big.Int
	mod *polynomial.Polynomial
}

// Set - Задает начальные значения и проверяет может ли поле быть построено
func (g *GaloisField) Set(p *big.Int, n *big.Int, poly *polynomial.Polynomial) *GaloisField {

	// Проверка входных данных

	// Копируем и проверяем p
	g.p = big.NewInt(0)

	if p.Cmp(constNum2) < 0 {
		panic("p >= 2")
	}

	if !MillerRabinTest(p) {
		panic("p is a prime number")
	}

	g.p.Set(p)

	// Копируем и проверяем n

	if n.Cmp(constNum1) <= 0 {
		panic("n > 1")
	}

	g.n = big.NewInt(0)
	g.n.Set(n)

	// Копируем Полином
	g.mod = new(polynomial.Polynomial)
	g.mod.SetPolynomial(poly)

	// Проверка на количество элементов
	offsetArr := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
	}
	offset := polynomial.NewPolynomial(offsetArr)

	// Счетчик элементов
	counter := big.NewInt(2)

	// Начальный полином
	xArr := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
	}

	x := polynomial.NewPolynomial(xArr)

	for x.Value(g.p).Cmp(constNum1) != 0 {
		x = x.Mul(offset, x)
		_, x = new(polynomial.Polynomial).QuoRem(x, g.mod)
		x.Mod(x, g.p)

		counter.Add(counter, constNum1)
	}

	want := new(big.Int)
	want = Pow(g.p, g.n)

	if want.Cmp(counter) != 0 {
		panic("The field cannot be created")
	}

	return g
}

// Add - Складывает два элемента в поле
func (g *GaloisField) Add(a, b *polynomial.Polynomial) *polynomial.Polynomial {

	maxValue := new(big.Int)
	maxValue = Pow(g.p, g.n)

	// Проверка на принадлежность полю
	if a.Value(g.p).Sign() < 0 || a.Value(g.p).Cmp(maxValue) >= 0 {
		panic("The element a does not belong to the field")
	}

	if b.Value(g.p).Sign() < 0 || b.Value(g.p).Cmp(maxValue) >= 0 {
		panic("The element b does not belong to the field")
	}

	result := new(polynomial.Polynomial)
	result = result.Add(a, b)
	result = result.Mod(result, g.p)

	return result
}

// Mul - Умножает два элемента в поле
func (g *GaloisField) Mul(a, b *polynomial.Polynomial) *polynomial.Polynomial {

	maxValue := new(big.Int)
	maxValue = Pow(g.p, g.n)

	// Проверка на принадлежность полю
	if a.Value(g.p).Sign() < 0 || a.Value(g.p).Cmp(maxValue) >= 0 {
		panic("The element a does not belong to the field")
	}

	if b.Value(g.p).Sign() < 0 || b.Value(g.p).Cmp(maxValue) >= 0 {
		panic("The element b does not belong to the field")
	}

	result := new(polynomial.Polynomial)
	result = result.Mul(a, b)
	result.Mod(result, g.p)

	_, result = result.QuoRem(result, g.mod)
	result.Mod(result, g.p)

	return result
}

// Строковое представление
func (g *GaloisField) String() string {
	return "GF(" + g.p.String() + "^" + g.n.String() + ")"
}

// CayleyTableAdd - Таблица Кэли для сложения
//
// Файл сохраняется по абсолютному пути
func (g *GaloisField) CayleyTableAdd(path string) {

	// Максимальное количество элементов
	maxValue := new(big.Int)
	maxValue = Pow(g.p, g.n)

	name := g.p.String() + "^" + g.n.String() + "_add"

	// Создание файла
	file, err := os.Create(path + "/" + name + ".csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Создаем массив для итераций по многочленам
	var iArr []*big.Int
	for i := big.NewInt(0); i.Cmp(g.n) <= 0; i.Add(i, constNum1) {
		iArr = append(iArr, big.NewInt(0))
	}

	// Первая строка
	result := "\t"

	temp := new(polynomial.Polynomial)

	for O := big.NewInt(1); O.Cmp(maxValue) <= 0; O.Add(O, constNum1) {

		temp.Set(iArr)
		result = result + temp.String()

		if O.Cmp(maxValue) != 0 {
			result = result + "\t"
		}

		iArr[0] = iArr[0].Add(iArr[0], constNum1)

		for i := 0; i < len(iArr); i++ {

			if iArr[i].Cmp(g.p) == 0 {
				iArr[i].Mod(iArr[i], g.p)
				iArr[i+1].Add(iArr[i+1], constNum1)
			}
		}

	}
	result = result + "\n"

	_, err = file.WriteString(result)
	if err != nil {
		panic(err)
	}

	// Массив строк

	// Создаем массивы для итерации
	iArr = []*big.Int{}
	for i := big.NewInt(0); i.Cmp(g.n) <= 0; i.Add(i, constNum1) {
		iArr = append(iArr, big.NewInt(0))
	}

	a := new(polynomial.Polynomial)
	b := new(polynomial.Polynomial)

	for I := big.NewInt(1); I.Cmp(maxValue) <= 0; I.Add(I, constNum1) {

		jArr := []*big.Int{}
		for j := big.NewInt(0); j.Cmp(g.n) <= 0; j.Add(j, constNum1) {
			jArr = append(jArr, big.NewInt(0))
		}

		a.Set(iArr)
		result = a.String() + "\t"

		for J := big.NewInt(1); J.Cmp(maxValue) <= 0; J.Add(J, constNum1) {

			b.Set(jArr)

			temp = temp.Add(a, b)
			temp = temp.Mod(temp, g.p)

			result = result + temp.String()

			if J.Cmp(maxValue) != 0 {
				result = result + "\t"
			}

			jArr[0] = jArr[0].Add(jArr[0], constNum1)

			for j := 0; j < len(jArr); j++ {

				if jArr[j].Cmp(g.p) == 0 {
					jArr[j].Mod(jArr[j], g.p)
					jArr[j+1].Add(jArr[j+1], constNum1)
				}
			}

		}

		result = result + "\n"

		_, err = file.WriteString(result)
		if err != nil {
			panic(err)
		}

		// Увеличение полинома a

		if I.Cmp(new(big.Int).Sub(maxValue, constNum1)) != 0 {
			result = result + "\t"
		}

		iArr[0] = iArr[0].Add(iArr[0], constNum1)

		for i := 0; i < len(iArr); i++ {

			if iArr[i].Cmp(g.p) == 0 {
				iArr[i].Mod(iArr[i], g.p)
				iArr[i+1].Add(iArr[i+1], constNum1)
			}
		}

	}

}

// CayleyTableMul - Таблица Кэли для умножения
//
// Файл сохраняется по абсолютному пути
func (g *GaloisField) CayleyTableMul(path string) {

	// Максимальное количество элементов
	maxValue := new(big.Int)
	maxValue = Pow(g.p, g.n)

	name := g.p.String() + "^" + g.n.String() + "_mul"

	// Создание файла
	file, err := os.Create(path + "/" + name + ".csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Создаем массив для итераций по многочленам
	var iArr []*big.Int
	for i := big.NewInt(0); i.Cmp(g.n) <= 0; i.Add(i, constNum1) {
		iArr = append(iArr, big.NewInt(0))
	}

	// Первая строка
	result := "\t"

	temp := new(polynomial.Polynomial)

	for O := big.NewInt(1); O.Cmp(maxValue) <= 0; O.Add(O, constNum1) {

		temp.Set(iArr)
		result = result + temp.String()

		if O.Cmp(maxValue) != 0 {
			result = result + "\t"
		}

		iArr[0] = iArr[0].Add(iArr[0], constNum1)

		for i := 0; i < len(iArr); i++ {

			if iArr[i].Cmp(g.p) == 0 {
				iArr[i].Mod(iArr[i], g.p)
				iArr[i+1].Add(iArr[i+1], constNum1)
			}
		}

	}
	result = result + "\n"

	_, err = file.WriteString(result)
	if err != nil {
		panic(err)
	}

	// Массив строк

	// Создаем массивы для итерации
	iArr = []*big.Int{}
	for i := big.NewInt(0); i.Cmp(g.n) <= 0; i.Add(i, constNum1) {
		iArr = append(iArr, big.NewInt(0))
	}

	a := new(polynomial.Polynomial)
	b := new(polynomial.Polynomial)

	for I := big.NewInt(1); I.Cmp(maxValue) <= 0; I.Add(I, constNum1) {

		jArr := []*big.Int{}
		for j := big.NewInt(0); j.Cmp(g.n) <= 0; j.Add(j, constNum1) {
			jArr = append(jArr, big.NewInt(0))
		}

		a.Set(iArr)
		result = a.String() + "\t"

		for J := big.NewInt(1); J.Cmp(maxValue) <= 0; J.Add(J, constNum1) {

			b.Set(jArr)

			temp = temp.Mul(a, b)
			temp = temp.Mod(temp, g.p)
			_, temp = temp.QuoRem(temp, g.mod)
			temp = temp.Mod(temp, g.p)

			result = result + temp.String()

			if J.Cmp(maxValue) != 0 {
				result = result + "\t"
			}

			jArr[0] = jArr[0].Add(jArr[0], constNum1)

			for j := 0; j < len(jArr); j++ {

				if jArr[j].Cmp(g.p) == 0 {
					jArr[j].Mod(jArr[j], g.p)
					jArr[j+1].Add(jArr[j+1], constNum1)
				}
			}

		}

		result = result + "\n"

		_, err = file.WriteString(result)
		if err != nil {
			panic(err)
		}

		// Увеличение полинома a

		if I.Cmp(new(big.Int).Sub(maxValue, constNum1)) != 0 {
			result = result + "\t"
		}

		iArr[0] = iArr[0].Add(iArr[0], constNum1)

		for i := 0; i < len(iArr); i++ {

			if iArr[i].Cmp(g.p) == 0 {
				iArr[i].Mod(iArr[i], g.p)
				iArr[i+1].Add(iArr[i+1], constNum1)
			}
		}

	}

}

// NewGaloisField - Создает GaloisField и задает ему начальные значения
func NewGaloisField(p *big.Int, n *big.Int, poly *polynomial.Polynomial) *GaloisField {
	return new(GaloisField).Set(p, n, poly)
}
