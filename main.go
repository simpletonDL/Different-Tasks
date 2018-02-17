package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
Структура, описывающая ТЕСТ. Заглавные буквы нужны, чтобы отличать
тесты, которые тестируют написанный мной алгоритм, от ТЕСТОВ, которые
являются входными данными в задаче.
 */
type Test struct {
	probability float64
	time int
}

/*
Собственно функция, реализующая алгоритм, решающий задачау для
набора ТЕСТОВ tests и ограничения времени maxTests. Возвращает
искомую максимальную вероятность падения и соответствующий ей
набор ТЕСТОВ (срез из индексов).
 */
func solve_bag(tests []Test, maxTime int) (float64, []int) {
	/*
	Описание алгоритма:
	Задача на самом деле является интерпретацией задчи о рюкзаке.
	Почему? Мы хотим максимизировать вероятность, что программа упадёт,
	то есть что упадёт хотябы один тест. Это вероятность равна
	1 - вероятность того, что все тесты пройдут. Значит мы хотим
	минимизировать вероятность того, что все тесты пройдут. Таким
	образом получаем задачу о рюкзаке для минимизации. Правда,
	здесь вместо суммы произведение, но на самом деле это не важно.
	dp[i][j] = минимальная вероятность, которую можно получить
	использовав j первых ТЕСТОВ, суммарное время которых не
	превышает i. Далее итеративно считаем динамику и искомая
	вероятность хранится в dp[maxTime][countTest], где countTest -
	колличество ТЕСТОВ.
	 */
	countTest := len(tests)

	dp := make([][]float64, maxTime + 1)
	path := make([][]bool, maxTime + 1)

	for i := range dp {
		dp[i] = make([]float64, countTest + 1)
		path[i] = make([]bool, countTest + 1)

		for j := range dp[i] {
			dp[i][j] = 1.0
		}
	}

	for j := 1; j <= countTest; j++ {
		currentTest := tests[j-1]

		for i := 1; i <= maxTime; i++ {
			dp[i][j] = dp[i][j-1];
			if (i - currentTest.time) >= 0 && dp[i - currentTest.time][j-1] * (1.0 - currentTest.probability) < dp[i][j]{
				dp[i][j] = dp[i - currentTest.time][j-1] * (1.0 - currentTest.probability)
				path[i][j] = true
			}
		}
	}

	/*
	Продолжение описание алгоритма. Теперь про восстановление ответа.
	Чтобы восстановить ответ, в процессе заполнения массива dp будем
	заполнять массив path. path[i][j] = 1, если ТЕСТ номер j-1 (при индексации с 0)
	находится в оптимальном наборе ТЕСТОВ, если можно использовать только j первых
	тестов и с суммарным ограничением по времени i. Теперь в двумерной таблице
	dp мы знаем, из какой ячейки в какую мы пришли и можно легко восстановить ответ,
	начиная с ячейки dp[maxIime][countTest]
	 */
	var optimal_subset []int
	i, j := maxTime, countTest

	for i >=0 && j >= 0 {
		if (path[i][j]) {
			optimal_subset = append(optimal_subset, j-1)
			i -= tests[j-1].time
		}
		j--
	}

	for i, j := 0, len(optimal_subset)-1; i < j; i, j = i+1, j-1 {
		optimal_subset[i], optimal_subset[j] = optimal_subset[j], optimal_subset[i]
	}

	return 1.0 - dp[maxTime][countTest], optimal_subset
}

/* Все остальные функции для тестирования задания, я решил их написать,
чтобы затестить рюкзак :D
*/

/*
Наивный алгоритм, который решает задачу полным перебором по всем возможным
подмножествам ТЕСТОВ. Вход и выход такие же. Работает эта штука экспоненциально,
поэтому больше, чем 30 ТЕСТОВ брать не стоит.
 */
func solve_bag_by_search(tests []Test, maxTime int) (float64, []int) {
	answer := 1.0
	var optimal_subset []int

	for i := 0; i < 1 << uint(len(tests)); i++ {
		probability := 1.0
		var subset []int

		time := 0
		for j := 0; j < len(tests); j++ {
			if(i >> uint(j)) % 2 == 1 {
				probability *= 1 - tests[j].probability
				time += tests[j].time
				subset = append(subset, j)
			}
		}
		if time <= maxTime && probability < answer {
			answer = probability
			optimal_subset = subset
		}
	}

	return 1.0 - answer, optimal_subset
}

/*
Эта функция генерирует случайным образом набор из countTest ТЕСТОВ, время исполнения
которых меньше MaxTestTime и возвращает решение задачи для этого набора тестов и ограничения
по времени maxTime, полулченное рюкзаком и перебором, а так же возвращает true/false, если
искомые вероятности совпадают/ не совпадают.
 */
func runTest(countTest int, maxTime int, maxTestTime int) (bool, float64, float64, []int, []int) {
	tests := make([]Test, countTest)
	for i := range tests {
		tests[i].probability = rand.Float64()
		tests[i].time = rand.Intn(maxTestTime)
	}

	time_get, subset_1 := solve_bag(tests, maxTime)
	time_expected, subset_2 := solve_bag_by_search(tests, maxTime)

	return math.Abs(time_get - time_expected) < 0.0000000001, time_get, time_expected, subset_1, subset_2
}

/*
Функция запускает мультитест, для каждого теста выписывает: прошёл он или нет (то есть совпадают ли
вероятности, полученные рюкзаком и перебором), вероятность полученную рюкзаком и оптимальные поднаборы
тестов, полученные рюкзаком и перебором. Вообще последние не обязаны совпадать, так как вполне могут
существовать различные поднаборы, на которых достигается максимальная вероятнсть.
 */
func runMultiTest()  {
	const count = 10
	const countTest = 10
	const maxTime = 200
	const maxTestTime = 200

	done := true

	for i := 1; i <= count; i++ {
		accepted, time_get, time_expected, subset_get, subset_expected := runTest(countTest, maxTime, maxTestTime)
		fmt.Print("Test", i, ": ", accepted, ", probability: ", time_get)
		if (!accepted) {
			fmt.Print(", difference ", math.Abs(time_get - time_expected))
			done = false
		}
		fmt.Println(", subset get: ", subset_get, ", subset expected: ", subset_expected)
	}

	if (done) {
		fmt.Println("accepted")
	} else {
		print("fail")
	}
}

/*
Пример ввода данных: Сначала количество тестов и ограничение времени, потом
задаются тесты (вероятность и время)
2 2
0.5 1
0.7 1
 */

func main() {
	rand.Seed(time.Now().UnixNano())
	//runMultiTest() // Можно запустить мультест, если Вам понадобится

	var countTest, maxTime int
	fmt.Scan(&countTest, &maxTime)

	tests := make([]Test, countTest)
	for i := range tests {
		fmt.Scanf("%f %d/n", &(tests[i].probability), &(tests[i].time))
	}

	time, subset := solve_bag(tests, maxTime)
	fmt.Println(time, subset)
}