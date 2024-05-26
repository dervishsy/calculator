// Initialization of variables
const expressionInput = document.getElementById('expressionInput');
const submitButton = document.getElementById('submitButton');
const expressionsList = document.getElementById('expressionsList');
const errorMessage = document.getElementById('errorMessage');

let intervalId; // Interval variable for periodic updates
let expressions = []; // Expression array for rendering

// Periodic updates of the expression list
function updateExpressionsList() {
    fetchExpressions();
}

// Submit an expression to the server
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
                // Add the expression to the array
                expressions.push({ expression: expression, id: id, status: 'Pending', result: null });
                renderExpressions()
            } else {
                // Show error message if the request was not successful
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
// Fetch the list of expressions from the server
function fetchExpressions() {
    fetch('/api/v1/expressions')
        .then(response => response.json())
        .then(data => {
            expressions = data.expressions; // Update the expression array with the fetched data
            renderExpressions();
        })
        .catch(error => console.error('Error fetching expressions:', error));
}

// Render the list of expressions on the page
function renderExpressions() {
    expressionsList.innerHTML = '';
    expressions.forEach(expression => {
        const listItem = document.createElement('li');
        listItem.textContent = `Expression: ${expression.expression}, ID: ${expression.id}, Status: ${expression.status}, Result: ${expression.result}`;

        // Check the status of the expression and assign the appropriate CSS class
        if (expression.status === 'completed') {
            listItem.classList.add('completed');
        }else{
            listItem.classList.add('notcompleted');

        }

        expressionsList.appendChild(listItem);
    });
}

// Generate a UUID string for the expression id
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// Show an error message to the user with a timeout
function showErrorMessage(message, additionalText = '') {
    let errorMessageText = message;
    if (additionalText) {
        errorMessageText += `: ${additionalText}`;
    }
    errorMessage.textContent = errorMessageText;
    errorMessage.style.display = 'block';

    // Hide the error message after 5 seconds
    setTimeout(() => {
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';
    }, 5000);
}
// Fetch the list of expressions from the server
fetchExpressions();

// Start the periodic expression list update
intervalId = setInterval(updateExpressionsList, 5000); // Update every 5 seconds

// Event listeners
submitButton.addEventListener('click', submitExpression);

expressionInput.addEventListener('keydown', event => {
    if (event.key === 'Enter') {
        submitExpression();
    }
});

window.addEventListener('beforeunload', () => {
    clearInterval(intervalId);
});
