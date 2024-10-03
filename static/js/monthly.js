document.querySelectorAll(".toggle-transactions").forEach(button => {
    button.addEventListener("click", function() {
        const index = this.getAttribute("data-payout-index");
        const transactions = document.querySelector(`#${index} > .transactions`);
        if (transactions) {
            transactions.classList.toggle("hidden");
            this.classList.toggle("before:rotate-180");
        }
    });
});
