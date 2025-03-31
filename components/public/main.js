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

  Alpine.magic("clipboard", () => {
    return (subject) => navigator.clipboard.writeText(subject);
  });
});

window.addEventListener("load", () => {
  const slotItems = Array.from(document.querySelectorAll("[slot]"))
  for (const item of slotItems) {
    const id = item.getAttribute("slot")
    const wrapper = document.querySelector(`#${id}`)
    if (!wrapper) {
      return
    }

    const slot = wrapper.shadowRoot.querySelector(`slot[name='${id}']`)
    if (!slot) {
      return
    }
    slot.replaceWith(item)
  }
})