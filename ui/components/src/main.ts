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

window.Alpine = Alpine;

// Checkbox

const checkboxStore = new Map<string, any>();

Alpine.magic("checkbox", () => (id: string) => checkboxStore.get(id));

Alpine.data("checkbox", (checked = false, id = null) => ({
  id: id,
  checked: checked,

  init() {
    if (this.id !== null) {
      checkboxStore.set(this.id, this);
    }
    this.$dispatch("init", { checked: this.checked });
  },
  async toggleCheckbox() {
    this.checked = !this.checked;
    await this.$nextTick();
    this.$dispatch("changed", { checked: this.checked });
  },
  isChecked(): boolean {
    return this.checked;
  },
}));

// Switch

const switchStore = new Map<string, any>();

Alpine.magic("switch_", () => (id: string) => switchStore.get(id));

Alpine.data("switch_", (checked = false, id = null) => ({
  id: id,
  checked: checked,

  init() {
    if (this.id !== null) {
      switchStore.set(this.id, this);
    }
    this.$dispatch("init", { checked: this.checked });
  },
  async toggleSwitch() {
    this.checked = !this.checked;
    await this.$nextTick();
    this.$dispatch("changed", { checked: this.checked });
  },
  isChecked(): boolean {
    return this.checked;
  },
}));

// Collapsible

const collapsibleStore = new Map<string, any>();

Alpine.magic("collapsible", () => (id: string) => collapsibleStore.get(id));

Alpine.data("collapsible", (open = false, id = null) => ({
  id: id,
  open: open,

  init() {
    if (this.id !== null) {
      collapsibleStore.set(this.id, this);
    }
    this.$dispatch("init", { open: this.open });
  },
  async toggleCollapsible() {
    this.open = !this.open;
    await this.$nextTick();
    this.$dispatch("changed", { open: this.open });
  },
  isOpen(): boolean {
    return this.open;
  },
}));

// Dialog
const dialogStore = new Map<string, any>();

Alpine.magic("dialog", () => (id: string) => dialogStore.get(id));

Alpine.data("dialog", (open = false, id = null) => ({
  id: id,
  open: open,

  init() {
    if (this.id !== null) {
      dialogStore.set(this.id, this);
    }
    this.$dispatch("init", { open: this.open });
  },
  async openDialog() {
    this.open = true;
    await this.$nextTick();
    this.$dispatch("opened");
  },
  async closeDialog(value: string) {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("closed", { value: value });
  },
  async dismissDialog() {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("dismissed");
  },
}));

// Dropdown

const dropdownStore = new Map<string, any>();

Alpine.magic("dropdown", () => (id: string) => dropdownStore.get(id));

Alpine.data("dropdown", (open = false, id = null) => ({
  id: id,
  open: open,
  usedKeyboard: false,

  init() {
    if (this.id !== null) {
      dropdownStore.set(this.id, this);
    }
    this.$dispatch("init", { open: this.open });
  },
  async openDropdown() {
    this.open = true;
    await this.$nextTick();
    this.$dispatch("opened");
  },
  async closeDropdown(value: string) {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("closed", { value: value });
  },
  async dismissDropdown() {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("dismissed");
  },
}));

// Select

const selectStore = new Map<string, any>();

Alpine.magic("select", () => (id: string) => selectStore.get(id));

Alpine.data("select", (value: string | null = null, clearable = false, id = null) => ({
  id: id,
  value: value?.length == 0 ? null : value,
  open: false,
  usedKeyboard: false,
  clearable: clearable,

  async openSelect() {
    this.open = true;
    await this.$nextTick();
    this.$dispatch("opened");
  },
  async closeSelect(value: string) {
    this.open = false;
    this.value = value;
    this.$dispatch("changed", { value: this.value });
    await this.$nextTick();
    this.$dispatch("closed", { value: this.value });
  },
  async dismissSelect() {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("dismissed");
  },
  async setSelect(value: string) {
    this.value = value;
    this.$dispatch("changed", { value: this.value });
  }
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

Alpine.data("tabs", () => ({
  value: null,

  setActiveTab(tab: string) {
    this.value = tab;
    this.$dispatch("change");
  },
}));

// Radio

const radioStore = new Map<string, any>();

Alpine.magic("radio", () => (id: string) => radioStore.get(id));

Alpine.data("radio", (name: string, value = null, id = null) => ({
  id: id,
  value: value,
  name: name,

  init() {
    if (this.id !== null) {
      radioStore.set(this.id, this);
    }
    this.$dispatch("init", { value: this.value });
  },
  async setRadioValue(value: string | null) {
    this.value = value;
    await this.$nextTick();
    this.$dispatch("changed", { value: this.value });
  },
  getCheckedValue(): string | null {
    return this.value;
  },
}));

Alpine.magic("clipboard", () => {
  return (data: string) => navigator.clipboard.writeText(data);
});

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
