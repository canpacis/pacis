document.addEventListener("alpine:init", () => {
  Alpine.data("dialog", () => ({
    isOpen: false,

    openDialog() {
      this.isOpen = true;
      this.$dispatch("open");
    },
    closeDialog(dismiss = false) {
      this.isOpen = false;
      this.$dispatch("close");
      if (dismiss) {
        this.$dispatch("dismiss");
      }
    },
  }));
});
