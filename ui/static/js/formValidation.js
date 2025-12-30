document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("configForm");
    const button = form.querySelector("button");
    const inputs = form.querySelectorAll("input[required]");

    function validateForm() {
        let allValid = true;

        inputs.forEach(input => {
            const pattern = input.getAttribute("pattern");
            let isValid = true;

            if (!input.value.trim()) {
                isValid = false;
            }

            if (pattern) {
                const regex = new RegExp("^" + pattern + "$");
                if (!regex.test(input.value)) {
                    isValid = false;
                }
            }

            if (!isValid) {
                allValid = false;
                input.classList.add("invalid");
                input.classList.remove("valid");
            } else {
                input.classList.add("valid");
                input.classList.remove("invalid");
            }
        });

        button.disabled = !allValid;
    }

    inputs.forEach(input => {
        input.addEventListener("input", validateForm);
        input.addEventListener("change", validateForm);
        input.addEventListener("blur", validateForm);
    });

    validateForm();

    form.addEventListener("submit", (e) => {
        if (button.disabled) {
            e.preventDefault();
            alert("Please fix invalid fields before generating config.");
        }
    });
});
