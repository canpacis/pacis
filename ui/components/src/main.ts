import hljs from "highlight.js/lib/core";
import go from "highlight.js/lib/languages/go";
import Alpine from "alpinejs";
import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";
import persist from "@alpinejs/persist";

hljs.registerLanguage("go", go);

Alpine.plugin(persist);
Alpine.plugin(focus);
Alpine.plugin(anchor);

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}

const cookieStorage = {
  getItem(key: string) {
    let cookies = document.cookie.split(";");
    for (let i = 0; i < cookies.length; i++) {
      let cookie = cookies[i].split("=");
      if (key == cookie[0].trim()) {
        return JSON.stringify(decodeURIComponent(cookie[1]));
      }
    }
    return null;
  },
  setItem(key: string, value: string) {
    document.cookie =
      key + " = " + encodeURIComponent(JSON.parse(value)) + "; path=/";
  },
};

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

Alpine.data("select", () => ({
  value: null,
  isOpen: false,
  isKeyboard: false,
  clearable: false,

  openSelect() {
    this.isOpen = true;
    this.$dispatch("open");
  },
  closeSelect(value: string | null, dismiss = false) {
    this.isOpen = false;
    if (!dismiss) {
      this.value = value;
    }
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

Alpine.magic("clipboard", () => {
  return (data: string) => navigator.clipboard.writeText(data);
});

const scheme = document.querySelector("html")!.classList.contains("light")
  ? "light"
  : "dark";

Alpine.store("colorScheme", {
  value: Alpine.$persist(scheme).as("pacis_color_scheme").using(cookieStorage),

  toggle() {
    const html = document.querySelector("html");
    if (!html) {
      return;
    }
    if (this.value === "dark") {
      html.classList.remove("dark");
      html.classList.add("light");
      this.value = "light";
    } else {
      html.classList.remove("light");
      html.classList.add("dark");
      this.value = "dark";
    }
  },
});

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

  hljs.highlightAll();
});
