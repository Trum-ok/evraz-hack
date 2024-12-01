package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	ApiKey = os.Getenv("API_KEY")
	Url    = "http://84.201.152.196:8020/v1/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Response struct {
	RequestID  int      `json:"request_id"`
	ResponseID int      `json:"response_id"`
	Model      string   `json:"model"`
	Choices    []Choice `json:"choices"`
}

func СGetRequestToLlm(text string) (string, error) {
	// Формируем тело запроса
	system_promt := "Ответ должен быть исключительно на русском языке пример Анализа проекта \"project_name\" Общее количество ошибок: 3 Архитектурных нарушений: 1 Несоответствий стандартам: 1 ### Архитектурное нарушение > \"chat_service.py\" (номер строки:номер символа, при наличии) > Необходимо вынести в слой адаптеров, работать через репозитории и интерфейсы из сервисов ```python user = User.query.filter_by(username=token).first() location = Location.query.filter_by(name=name).first() ``` ### Краткое описание нарушения (Add braces to if statement) > \"LinkFragmentValidator.cs\" (номер строки:номер символа, при наличии) > Severity Code Description Project File Line Error (active) RCS1007 Add braces to if statement Eurofurence.App.Server.Services LinkFragmentValidator.cs 35 ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealerId)) return ValidationResult.Error(\"Target must be of typ Guid\"); ``` > Предложенное исправление ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealer"

	user_promt := fmt.Sprintf(`Провести подробный coderiew проекта по следующим критериям, также не забудь проверять стилистику кода со стандартом Microsoft:
NuGet-зависимости:
Все пакеты должны быть обновлены до актуальных версий.
Уязвимые транзитивные зависимости должны быть обновлены или включены в проект.
Отсутствие лишних зависимостей и абсолютных путей.
Кодовая база:
Удалите незавершённые TODO, закомментированный/неиспользуемый код, устаревший [Obsolete] код.
Комментарии должны быть у ключевых элементов.
Исключите: неиспользуемые переменные, дублирование логов исключений, некорректное объединение строк, избыточные null-проверки.
Архитектура и дизайн:
Корректная регистрация сервисов в IoC (e.g., HostedService как Singleton).
Исключите возврат null вместо пустой коллекции.
Проверяйте аргументы внутри методов.
Используйте интерфейсы вместо конкретных реализаций коллекций.
LINQ:
Применяйте Chunk() вместо Skip().Take().
Для пользовательских типов с Union(), Distinct() и т.п. переопределите Equals() и GetHashCode().
Избегайте лишних вызовов ToList() и Distinct() после Union().
Entity Framework:
Используйте асинхронные методы материализации данных.
Не удаляйте сущности в циклах, группируйте SaveChangesAsync().
Фильтрацию данных выполняйте на стороне БД.
Применяйте AddAsync()/AddRangeAsync() только с SequenceHiLo.
Вот код для проверки: %s.`, text)
	requestBody := map[string]interface{}{
		"model": "mistral-nemo-instruct-2407",
		"messages": []map[string]string{
			{"role": "system", "content": user_promt},
			{"role": "user", "content": system_promt},
		},
		"max_tokens":  1024,
		"temperature": 0.1,
	}

	// Кодируем тело запроса в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодировки JSON: %w", err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Authorization", ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 240 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка HTTP: %d, ответ: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON. Тело ответа: %s\n", string(respBody))
		return "", fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}

func PythonGetRequestToLlm(text string) (string, error) {
	system_promt := "Ответ должен быть исключительно на русском языке пример анализа проекта \"project_name\" Общее количество ошибок: 3 Архитектурных нарушений: 1 Несоответствий стандартам: 1 ### Архитектурное нарушение > \"chat_service.py\" (номер строки:номер символа, при наличии) > Необходимо вынести в слой адаптеров, работать через репозитории и интерфейсы из сервисов ```python user = User.query.filter_by(username=token).first() location = Location.query.filter_by(name=name).first() ``` ### Краткое описание нарушения (Add braces to if statement) > \"LinkFragmentValidator.cs\" (номер строки:номер символа, при наличии) > Severity Code Description Project File Line Error (active) RCS1007 Add braces to if statement Eurofurence.App.Server.Services LinkFragmentValidator.cs 35 ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealerId)) return ValidationResult.Error(\"Target must be of typ Guid\"); ``` > Предложенное исправление ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealer"
	user_promt := fmt.Sprintf(`Проведи ревью кода %s с учётом следующих стандартов:

Архитектура: Строгое разделение на три слоя:

Приложение (App Layer): Содержит бизнес-логику, DTO, сервисы, ошибки. Использует DI, не зависит от адаптеров.
Адаптеры (Adapters): Содержат интеграции, контроллеры, репозитории. SQLAlchemy таблицы описываются в snake_case с naming_convention.
Композиты (Composites): Управляют сборкой приложения, настройками, зависимостями.
Кодирование: Соблюдай PEP8. Для форматирования используй isort и yapf. Строки переносятся при достижении 80 символов (максимум — 100).

Данные:

Хранение дат в БД — naive UTC. Для сериализации используем ISO 8601.
Форматы: Decimal → строка, UUID → строка, Enum → name/value.
Логирование: Используй модуль logging. Формат: JSON, строки через для оптимизации.

Тесты: Пиши юнит-тесты с мокацией адаптеров. Для интеграционных тестов репозиториев используй SQLite в памяти.

Валидация данных: Все входные данные валидируются с помощью DTO (Pydantic).

SQLAlchemy:

Указывай naming_convention.
Избегай диалекто-зависимых конструкций (для MSSQL — отмечай как # TODO: dialect dependent).
Асинхронность: Разрешена только при наличии обоснования. Используй gevent и патчинг, включая драйвер MSSQL.

Прочие моменты:

Авторизация: через JWT (PyJWT). Данные пользователя извлекаются из токена.
Конфигурация: параметры передаются через env и описываются в классах-наследниках BaseSettings.
Транзакции: управление через паттерн "Единица работы". Вложенные транзакции запрещены.
Документирование:

Поддерживай читаемость кода.
Докстринги должны быть оформлены по PEP 256/257.
Убедись, что код соответствует этим стандартам, легко читаем, поддерживаем и следует проектным принципам.
Авторизация через JWT с pyjwt. Оформление кода по PEP8 с инструментами isort и yapf. Строки кода ограничены 80-100 символами`, text)
	requestBody := map[string]interface{}{
		"model": "mistral-nemo-instruct-2407",
		"messages": []map[string]string{
			{"role": "system", "content": user_promt},
			{"role": "user", "content": system_promt},
		},
		"max_tokens":  1024,
		"temperature": 0.1,
	}

	// Кодируем тело запроса в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодировки JSON: %w", err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Authorization", ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 240 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка HTTP: %d, ответ: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON. Тело ответа: %s\n", string(respBody))
		return "", fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil

}

func CArchitectureGetRequestToLlm(text string) (string, error) {
	system_promt := `Ответ должен быть исключительно на русском языке ввида Пример Code Review
1. Общие замечания
Структура проекта:

Структура проекта логична и разделена по функциональности.
Папка utils/ хорошо подходит для общих функций. Однако стоит убедиться, что здесь находятся только функции, которые переиспользуются в нескольких местах.
Тесты:

Хорошо, что тесты вынесены в отдельную папку.
Убедитесь, что все критически важные функции покрыты тестами (например, чтение больших файлов, обработка пустых файлов).
Документация:

README содержит описание проекта, но желательно добавить инструкции по сборке и запуску приложения.`
	user_promt := fmt.Sprintf("Надо провести отчет этой архитектуры %s, как самый крутой программист на планете", text)
	// Кодируем тело запроса в JSON
	requestBody := map[string]interface{}{
		"model": "mistral-nemo-instruct-2407",
		"messages": []map[string]string{
			{"role": "system", "content": user_promt},
			{"role": "user", "content": system_promt},
		},
		"max_tokens":  1024,
		"temperature": 0.1,
	}
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодировки JSON: %w", err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Authorization", ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 240 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка HTTP: %d, ответ: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON. Тело ответа: %s\n", string(respBody))
		return "", fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}

func JsGetRequestToLlm(text string) (string, error) {
	system_promt := "Ответ должен быть исключительно на русском языке пример Анализа проекта \"project_name\" Общее количество ошибок: 3 Архитектурных нарушений: 1 Несоответствий стандартам: 1 ### Архитектурное нарушение > \"chat_service.py\" (номер строки:номер символа, при наличии) > Необходимо вынести в слой адаптеров, работать через репозитории и интерфейсы из сервисов ```python user = User.query.filter_by(username=token).first() location = Location.query.filter_by(name=name).first() ``` ### Краткое описание нарушения (Add braces to if statement) > \"LinkFragmentValidator.cs\" (номер строки:номер символа, при наличии) > Severity Code Description Project File Line Error (active) RCS1007 Add braces to if statement Eurofurence.App.Server.Services LinkFragmentValidator.cs 35 ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealerId)) return ValidationResult.Error(\"Target must be of typ Guid\"); ``` > Предложенное исправление ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealer"
	user_promt := fmt.Sprintf(`Пожалуйста, проведите ревью этого кода %s с учётом следующих стандартов:
1. **Функции**:
   - Используйте обычные функции вместо стрелочных (за исключением простых callback-функций).
   - Стрелочные функции с телом должны быть заменены на обычные функции.
2. **Компоненты**:
   - Все компоненты должны использовать стиль именования PascalCase.
   - Компоненты должны быть разделены на UI Kit, Components, Containers, и Pages в соответствии с проектной структурой.
   - Обязательно используйте CSS Modules с именами классов в snake_case.
   - Для каждого компонента должен быть файл index.ts для экспорта.
3. **Типы и интерфейсы**:
   - Типы компонентов должны быть описаны в отдельном файле types.ts.
   - Все props компонентов должны следовать определенному шаблону интерфейса (например, className?, style?).
4. **Конвенции именования**:
   - Переменные, константы и функции должны использовать camelCase.
   - Классы, типы, интерфейсы и enums должны использовать PascalCase.
   - Селекторы CSS и типы ответов API должны быть в snake_case.
   - Константы должны быть в UPPER_SNAKE_CASE.
   - Названия функций должны отражать выполняемое действие (например, getAccounts, sortItems).
   - Булевы переменные и функции должны начинаться с is или has.
5. **Обработчики событий**:
   - Функции-обработчики событий должны начинаться с handle (например, handleClick, handleSelect).
   - Пропсы для обработчиков событий должны начинаться с on (например, onClick, onChange).
6. **Читаемость кода**:
   - Используйте описательные и осмысленные имена для переменных, функций и компонентов.
   - Избегайте аббревиатур или неясных имен (например, d, list1 должно быть более описательным).
7. **Запрещено использование стрелочных функций как методов объектов**:
   - Не используйте стрелочные функции в качестве методов объектов, так как это может привести к проблемам с привязкой this.
Пожалуйста, убедитесь, что код соответствует этим стандартам, легко читаем и поддерживаем, а также правильно типизирован.`, text)
	requestBody := map[string]interface{}{
		"model": "mistral-nemo-instruct-2407",
		"messages": []map[string]string{
			{"role": "system", "content": user_promt},
			{"role": "user", "content": system_promt},
		},
		"max_tokens":  1024,
		"temperature": 0.1,
	}
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодировки JSON: %w", err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Authorization", ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 240 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка HTTP: %d, ответ: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON. Тело ответа: %s\n", string(respBody))
		return "", fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}

func GoGetRequestToLlm(text string) (string, error) {
	system_promt := "Ответ должен быть исключительно на русском языке пример Анализа проекта \"project_name\" Общее количество ошибок: 3 Архитектурных нарушений: 1 Несоответствий стандартам: 1 ### Архитектурное нарушение > \"chat_service.py\" (номер строки:номер символа, при наличии) > Необходимо вынести в слой адаптеров, работать через репозитории и интерфейсы из сервисов ```python user = User.query.filter_by(username=token).first() location = Location.query.filter_by(name=name).first() ``` ### Краткое описание нарушения (Add braces to if statement) > \"LinkFragmentValidator.cs\" (номер строки:номер символа, при наличии) > Severity Code Description Project File Line Error (active) RCS1007 Add braces to if statement Eurofurence.App.Server.Services LinkFragmentValidator.cs 35 ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealerId)) return ValidationResult.Error(\"Target must be of typ Guid\"); ``` > Предложенное исправление ```csharp if (!Guid.TryParse(fragment.Target, out Guid dealer"
	user_promt := fmt.Sprintf(`Проведите ревью этого Go-кода %s с учётом следующих стандартов:
1. Структура проекта:
Разделение на логические модули (handlers, services, repositories).
Пакеты и файлы — в snake_case.
Публичные функции должны быть документированы.
2. Функции:
Короткие, выполняющие одну задачу.
Имена отражают действия (например, GetUserByID).
Более 3 параметров — использовать структуры.
3. Переменные и константы:
camelCase для переменных, UPPER_SNAKE_CASE для констант.
Осмысленные имена (например, userID, ErrInvalidInput).
4. Типы и интерфейсы:
PascalCase для структур, интерфейсов.
Поля структур — camelCase.
Интерфейсы названы по назначению (например, UserRepository).
5. Обработка ошибок:
Использовать errors или fmt.Errorf с контекстом.
Проверка ошибок обязательна, без подавления.
6. Тестирование:
Покрытие ≥80%, тесты структурированы (*_test.go).
Осмысленные имена (например, TestGetUserByID_Success).
7. Производительность и безопасность:
Использовать context.Context.
Безопасные запросы к БД, управление ресурсами через defer.
Избегать утечек памяти в горутинах.
Код должен быть читабельным, поддерживаемым, форматированным через go fmt и соответствовать рекомендациям Go.`, text)
	requestBody := map[string]interface{}{
		"model": "mistral-nemo-instruct-2407",
		"messages": []map[string]string{
			{"role": "system", "content": user_promt},
			{"role": "user", "content": system_promt},
		},
		"max_tokens":  1024,
		"temperature": 0.1,
	}
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодировки JSON: %w", err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Authorization", ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 240 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка HTTP: %d, ответ: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON. Тело ответа: %s\n", string(respBody))
		return "", fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}
