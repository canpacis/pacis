import hljs from "highlight.js/lib/core";
import go from "highlight.js/lib/languages/go";
import Alpine from "alpinejs";
import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";
import persist from "@alpinejs/persist";

// TODO: make this dynamic
hljs.registerLanguage("go", go);

document.addEventListener("DOMContentLoaded", () => {
  hljs.highlightAll();
});

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

Alpine.magic("appswitch", () => (id: string) => switchStore.get(id));

Alpine.data("appswitch", (checked = false, id = null) => ({
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

Alpine.data(
  "select",
  (value: string | null = null, clearable = false, id = null) => ({
    id: id,
    value: value?.length == 0 ? null : value,
    open: false,
    usedKeyboard: false,
    clearable: clearable,

    init() {
      if (this.id !== null) {
        selectStore.set(this.id, this);
      }
      this.$dispatch("init", { value: this.value });
    },
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
    },
  })
);

// Sheet

const sheetStore = new Map<string, any>();

Alpine.magic("sheet", () => (id: string) => sheetStore.get(id));

Alpine.data("sheet", (open = false, id = null) => ({
  id: id,
  open: open,

  init() {
    if (this.id !== null) {
      sheetStore.set(this.id, this);
    }
    this.$dispatch("init", { open: this.open });
  },
  async openSheet() {
    this.open = true;
    await this.$nextTick();
    this.$dispatch("opened");
  },
  async closeSheet() {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("closed");
  },
  isOpen(): boolean {
    return this.open;
  },
}));

// Tabs

const tabsStore = new Map<string, any>();

Alpine.magic("tabs", () => (id: string) => tabsStore.get(id));

Alpine.data("tabs", (value = null, id = null) => ({
  id: id,
  value: value,

  init() {
    if (this.id !== null) {
      tabsStore.set(this.id, this);
    }
    this.$dispatch("init", { value: this.value });
  },
  async setTab(tab: string) {
    this.value = tab;
    await this.$nextTick();
    this.$dispatch("changed", { value: this.value });
  },
  getTab(): string | null {
    return this.value;
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

// Tooltip

const tooltipStore = new Map<string, any>();

Alpine.magic("tooltip", () => (id: string) => tooltipStore.get(id));

Alpine.data("tooltip", (open = false, id = null) => ({
  id: id,
  open: open,
  cancel: new AbortController(),

  init() {
    if (this.id !== null) {
      tooltipStore.set(this.id, this);
    }
    this.$dispatch("init", { open: this.open });
  },
  queueOpenTooltip(delay: number = 0) {
    if (this.open) {
      return;
    }
    setTimeout(async () => {
      if (this.open) {
        return;
      }
      if (this.cancel.signal.aborted) {
        this.cancel = new AbortController();
        return;
      }
      await this.openTooltip();
      this.cancel = new AbortController();
    }, delay);
  },
  abortTooltip() {
    this.cancel.abort();
    if (this.open) {
      this.open = false;
      this.cancel = new AbortController();
    }
    this.$dispatch("aborted");
  },
  async openTooltip() {
    this.open = true;
    await this.$nextTick();
    this.$dispatch("opened");
  },
  async closeTooltip() {
    this.open = false;
    await this.$nextTick();
    this.$dispatch("closed");
  },
  isOpen(): boolean {
    return this.open;
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

type ColorSchemeStore = {
  value: "dark" | "light";
  toggle: () => void;
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

type LocaleStore = {
  value: string;
  set: (value: string) => void;
};

const locale = document.querySelector("html")!.getAttribute("lang") || "en";

Alpine.store("locale", {
  value: Alpine.$persist(locale).as("pacis_locale").using(cookieStorage),

  set(value: string) {
    this.value = value;
    const html = document.querySelector("html");
    if (!html) {
      return;
    }
    html.setAttribute("lang", value);
  },
});

const prefetchStore = new Map<string, { doc: Document }>();

function replaceDoc(doc: Document) {
  const head = document.head;
  const body = document.body;

  head.innerHTML = doc.head.innerHTML;
  body.innerHTML = doc.body.innerHTML;
  document.title = doc.title;

  Alpine.initTree(body);
  document.dispatchEvent(new Event("DOMContentLoaded"));
  document.dispatchEvent(new Event("load"));
  window.scrollTo(0, 0);
}

window.addEventListener("popstate", (e) => {
  if (e.state && e.state.page) {
    const data = prefetchStore.get(e.state.page);
    if (data) {
      replaceDoc(data.doc);
    } else {
      location.reload();
    }
  }
});

const prefetch = {
  get: async (url: string) => {
    if (prefetchStore.has(url)) {
      return;
    }
    const resp = await fetch(url);

    if (!resp.ok) {
      console.error("Prefetch failed: ", resp);
      return;
    }
    const data = await resp.text();
    const doc = new DOMParser().parseFromString(data, "text/html");
    prefetchStore.set(url, { doc });
  },
  set: async (url: string, e: MouseEvent) => {
    const data = prefetchStore.get(url);
    if (!data) {
      return;
    }
    e.preventDefault();

    window.history.pushState({ page: url }, "", url);
    replaceDoc(data.doc);
  },
  clear: async (url?: string) => {
    if (url) {
      prefetchStore.delete(url);
      return;
    }
    for (const key of prefetchStore.keys()) {
      prefetchStore.delete(key);
    }
  },
};

Alpine.effect(() => {
  (Alpine.store("colorScheme") as ColorSchemeStore).value;
  (Alpine.store("locale") as LocaleStore).value;
  prefetch.clear();
});

Alpine.magic("prefetch", () => prefetch);

Alpine.start();
