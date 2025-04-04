import Alpine from "alpinejs";
import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";

Alpine.plugin(focus);
Alpine.plugin(anchor);

declare global {
  interface Window {
    Alpine: typeof Alpine
  }
}

window.Alpine = Alpine;


Alpine.data("dropdown", () => ({
  isOpen: false,
  isKeyboard: false,

  openDropdown() {
    this.isOpen = true;
    this.$dispatch("open");
  },
  closeDropdown(dismiss = false) {
    this.isOpen = false;
    this.$dispatch("close");
    if (dismiss) {
      this.$dispatch("dismiss");
    }
  },
}));

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

Alpine.data("sheet", () => ({
  isOpen: false,

  openSheet() {
    this.isOpen = true;
    this.$dispatch("open");
  },
  closeSheet(dismiss = false) {
    this.isOpen = false;
    this.$dispatch("close");
    if (dismiss) {
      this.$dispatch("dismiss");
    }
  },
}));

Alpine.data("clipboard", () => ({
  copyToClipboard(subject: string) {
    navigator.clipboard.writeText(subject);
  },
}));

Alpine.start();

window.addEventListener("load", () => {
  const slotItems = Array.from(document.querySelectorAll("[slot]"));
  for (const item of slotItems) {
    const id = item.getAttribute("slot");
    const wrapper = document.querySelector(`#${id}`);
    if (!wrapper) {
      return;
    }

    const slot = wrapper.querySelector(`slot[name='${id}']`);
    if (!slot) {
      return;
    }
    slot.replaceWith(item);
  }
});
