const expressionInput = document.getElementById('expressionInput');
const submitButton = document.getElementById('submitButton');
const expressionsList = document.getElementById('expressionsList');
const errorMessage = document.getElementById('errorMessage');

let intervalId; // Переменная для хранения идентификатора интервала
let expressions = []; // Массив для хранения выражений на клиенте

// Функция для периодического обновления списка выражений
function updateExpressionsList() {
    fetchExpressions();
}

// Функция для отправки выражения на сервер
function submitExpression() {
    id = generateUUID()
    const expression = expressionInput.value.trim();
    if (expression) {
        const data = {
            id: id,
            expression: expression
        };

        fetch('/api/v1/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (response.ok) {
                expressionInput.value = '';
                // Добавляем выражение в массив
                expressions.push({ expression: expression, id: id, status: 'Pending', result: null });
                renderExpressions()
            } else {
                console.error('Error submitting expression:', response.status);
                response.text().then(errorText => {
                    showErrorMessage(`Error submitting expression: ${response.status}`, errorText);
                });            }
        })
        .catch(error => {
            console.error('Error submitting expression:', error);
            showErrorMessage('An error occurred while submitting the expression.');
        });
    }

}
// Функция для получения списка выражений с сервера
function fetchExpressions() {
    fetch('/api/v1/expressions')
        .then(response => response.json())
        .then(data => {
            expressions = data.expressions; // Обновляем массив выражений
            renderExpressions();
        })
        .catch(error => console.error('Error fetching expressions:', error));
}

// Функция для отрисовки списка выражений
function renderExpressions() {
    expressionsList.innerHTML = '';
    expressions.forEach(expression => {
        const listItem = document.createElement('li');
        listItem.textContent = `Expression: ${expression.expression}, ID: ${expression.id}, Status: ${expression.status}, Result: ${expression.result}`;

        // Проверяем статус выражения и присваиваем соответствующий класс CSS
        if (expression.status === 'completed') {
            listItem.classList.add('completed');
        }else{
            listItem.classList.add('notcompleted');

        }

        expressionsList.appendChild(listItem);
    });
}
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// Начальная загрузка списка выражений
fetchExpressions();

// Запуск периодического обновления списка выражений
intervalId = setInterval(updateExpressionsList, 5000); // Обновление каждые 5 секунд

// Обработчики событий
submitButton.addEventListener('click', submitExpression);

expressionInput.addEventListener('keydown', event => {
    if (event.key === 'Enter') {
        submitExpression();
    }
});

window.addEventListener('beforeunload', () => {
    clearInterval(intervalId);
});



function showErrorMessage(message, additionalText = '') {
    let errorMessageText = message;
    if (additionalText) {
        errorMessageText += `: ${additionalText}`;
    }
    errorMessage.textContent = errorMessageText;
    errorMessage.style.display = 'block';

    // Скрыть сообщение об ошибке через 5 секунд
    setTimeout(() => {
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';
    }, 5000);
}