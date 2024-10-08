async function handleSubmit(event) {
    event.preventDefault();

    const form = document.getElementById('loginForm');
    const formData = new URLSearchParams(new FormData(form));
    const actionUrl = form.action;
    try {
        const response = await fetch(actionUrl, {
            method: 'POST',
            body: formData,
        });
        if (response.ok) {
            window.location.href = '/';
        } else {
            const errorText = await response.text();
            displayError(errorText);
        }
    } catch (error) {
        displayError('An unexpected error occurred. Please try again.');
    }
}

function displayError(message) {
    const errorMessageContainer = document.getElementById('error-message');
    errorMessageContainer.textContent = message;
}
