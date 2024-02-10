package main // Определение пакета, в который входит файл.

import (
	"fmt"     // Импорт пакета fmt для форматированного ввода/вывода.
	"os"      // Импорт пакета os для работы с аргументами командной строки и переменными окружения.
	"os/exec" // Импорт пакета exec для запуска внешних команд.
	"sync"    // Импорт пакета sync для синхронизации горутин.
)

// CommandRunner определяет интерфейс для запуска команд.
type CommandRunner interface {
	RunCommand(laravelRootPath, command string) (string, error) // Описание метода для запуска команды.
}

// RealCommandRunner реализует интерфейс CommandRunner для выполнения настоящих команд.
type RealCommandRunner struct{}

// RunCommand выполняет реальную команду Laravel используя CLI.
func (r *RealCommandRunner) RunCommand(laravelRootPath, command string) (string, error) {
	cmd := exec.Command("php", "artisan", command) // Составление команды для запуска.
	cmd.Dir = laravelRootPath                      // Установка директории, в которой будет выполняться команда.
	output, err := cmd.CombinedOutput()            // Запуск команды и сбор вывода.
	return string(output), err                     // Возвращение результата выполнения и ошибки, если таковая имеется.
}

// runLaravelCommand выполняет команду Laravel и отправляет результаты в канал.
func runLaravelCommand(wg *sync.WaitGroup, results chan<- string, runner CommandRunner, laravelRootPath, command string) {
	defer wg.Done()                                            // Отложенный вызов для уменьшения счетчика WaitGroup после завершения функции.
	output, err := runner.RunCommand(laravelRootPath, command) // Запуск команды.
	if err != nil {
		// Отправка сообщения об ошибке в канал, если команда завершится с ошибкой.
		results <- fmt.Sprintf("Error running '%s': %s\n", command, err)
		return
	}
	// Отправка успешного вывода команды в канал.
	results <- fmt.Sprintf("Output of '%s': %s\n", command, output)
}

// main является точкой входа в программу.
func main() {
	if len(os.Args) < 3 { // Проверка наличия достаточного количества аргументов.
		fmt.Println("Usage: <path to laravel root> <command1> <command2> ...") // Инструкция по использованию.
		os.Exit(1)                                                             // Выход из программы с кодом ошибки.
	}

	laravelRootPath := os.Args[1] // Путь к корневой папке Laravel.
	commands := os.Args[2:]       // Команды Laravel для выполнения.

	var wg sync.WaitGroup                       // Инициализация WaitGroup для управления горутинами.
	results := make(chan string, len(commands)) // Создание канала для результатов с буфером для хранения всех результатов.

	runner := &RealCommandRunner{} // Создание экземпляра RealCommandRunner для выполнения команд.
	for _, command := range commands {
		wg.Add(1) // Увеличение счетчика WaitGroup перед запуском горутины.
		// Запуск горутины для выполнения команды Laravel.
		go runLaravelCommand(&wg, results, runner, laravelRootPath, command)
	}

	wg.Wait()      // Ожидание завершения всех горутин.
	close(results) // Закрыти

	for result := range results { // Перебираем результаты из канала.
		fmt.Print(result) // Выводим каждый результат в стандартный вывод.
	}
}
