package main // Определение пакета для теста. Тесты обычно находятся в том же пакете, что и тестируемый код.

import (
	"errors"  // Импорт пакета для создания ошибок.
	"fmt"     // Импорт пакета для форматированного вывода строк.
	"sync"    // Импорт пакета для синхронизации горутин.
	"testing" // Импорт пакета для написания тестов.
)

// MockCommandRunner имитирует реализацию CommandRunner для тестирования.
type MockCommandRunner struct {
	// Поля для хранения состояния мока или для настройки ответов могут быть добавлены здесь.
}

// RunCommand имитирует выполнение команды и возвращает заранее определённый результат.
func (m *MockCommandRunner) RunCommand(laravelRootPath string, command string) (string, error) {
	// Возвращает разные результаты в зависимости от входных параметров.
	if command == "fail" {
		// Для команды "fail" возвращает ошибку.
		return "", errors.New("command failed")
	}
	// Для любой другой команды возвращает строку, указывающую на успешное выполнение.
	return fmt.Sprintf("executed %s in %s", command, laravelRootPath), nil
}

// TestRunLaravelCommand проверяет функцию runLaravelCommand.
func TestRunLaravelCommand(t *testing.T) {
	// Создаём экземпляр MockCommandRunner для использования в тесте.
	runner := &MockCommandRunner{}
	// Создаём канал для результатов выполнения команд. Его размер соответствует количеству тестируемых команд.
	results := make(chan string, 2)
	// Инициализируем WaitGroup для ожидания завершения работы горутин.
	var wg sync.WaitGroup

	// Запланировать выполнение двух команд.
	wg.Add(2)

	// Запускаем горутину для выполнения успешной команды.
	go runLaravelCommand(&wg, results, runner, "/fake-path", "migrate")
	// Запускаем вторую горутину для выполнения команды, которая завершится с ошибкой.
	go runLaravelCommand(&wg, results, runner, "/fake-path", "fail")

	// Ожидаем завершения обеих горутин.
	wg.Wait()
	// Закрываем канал после того, как все результаты были отправлены.
	close(results)

	// Определяем ожидаемые результаты выполнения команд.
	successExpected := "Output of 'migrate': executed migrate in /fake-path\n"
	failExpected := "Error running 'fail': command failed\n"
	// Проверяем результаты, полученные из канала.
	for result := range results {
		if result != successExpected && result != failExpected {
			// Если результат не соответствует ни одному из ожидаемых, сообщаем об ошибке.
			t.Errorf("Unexpected result: %s", result)
		}
	}
}
