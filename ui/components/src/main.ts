import hljs from "highlight.js/lib/core";
import go from "highlight.js/lib/languages/go";
import Alpine from "alpinejs";
import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";
import persist from "@alpinejs/persist";
import intersect from "@alpinejs/intersect";

// TODO: make this dynamic
hljs.registerLanguage("go", go);

document.addEventListener("DOMContentLoaded", () => {
  hljs.highlightAll();
});

Alpine.plugin(persist);
Alpine.plugin(focus);
Alpine.plugin(anchor);
Alpine.plugin(intersect);

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

export type ColorSchemeStore = {
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

export type LocaleStore = {
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

export type DeviceStore = {
  isMobile: boolean;
};

Alpine.store("device", {
  isMobile:
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    ),
});

type FetchFunction<T> = ({ signal }: { signal: AbortSignal }) => Promise<T>;
type Resolve<T> = (value: T | PromiseLike<T>) => void;
type Reject = (reason?: any) => void;

interface QueuedFetchItem<T> {
  fetchFn: FetchFunction<T>;
  resolve: Resolve<T>;
  reject: Reject;
}

class QueuedFetcher {
  private items: QueuedFetchItem<any>[] = [];
  private isProcessing = false;
  private abortControllers: AbortController[] = [];

  /**
   * Queues a fetch operation.
   * @param fetchFn A function that returns a Promise for the fetch operation.
   * @returns A Promise that will resolve with the result of the fetch.
   */
  queue<T>(fetchFn: FetchFunction<T>): Promise<T> {
    return new Promise<T>((resolve, reject) => {
      this.items.push({ fetchFn, resolve, reject });
      this.processQueue();
    });
  }

  /**
   * Processes the next item in the queue.
   * @private
   */
  private async processQueue(): Promise<void> {
    if (this.isProcessing || this.queue.length === 0) {
      return;
    }

    this.isProcessing = true;
    const item = this.items.shift();

    if (item) {
      const controller = new AbortController();
      this.abortControllers.push(controller);
      const signal = controller.signal;

      try {
        const result = await item.fetchFn({ signal });
        item.resolve(result);
      } catch (error: any) {
        // Check if the error is due to an abort
        if (error?.name === "AbortError") {
          item.reject(error);
        } else {
          item.reject(error);
        }
      } finally {
        // Remove the AbortController for the completed request
        this.abortControllers = this.abortControllers.filter(
          (c) => c !== controller
        );
        this.isProcessing = false;
        this.processQueue(); // Process the next item in the queue
      }
    } else {
      this.isProcessing = false; // Should not happen, but for safety
    }
  }

  /**
   * Aborts all currently pending fetch operations in the queue.
   */
  abort(): void {
    this.abortControllers.forEach((controller) => {
      controller.abort();
    });
    this.abortControllers = [];
    this.items.forEach((item) => {
      item.reject(new Error("Fetch aborted"));
    });
    this.items = [];
    this.isProcessing = false;
  }
}

const queuedFetcher = new QueuedFetcher();
const prefetchStore = new Map<string, { doc: Document }>();

function replaceDoc(doc: Document) {
  const head = document.head;
  const body = document.body;

  head.innerHTML = doc.head.innerHTML;
  body.innerHTML = doc.body.innerHTML;
  document.title = doc.title;

  Alpine.initTree(body);
  document.dispatchEvent(new Event("DOMContentLoaded"));
  window.scrollTo(0, 0);
}

window.addEventListener("popstate", () => {
  const data = prefetchStore.get(window.location.pathname);
  if (data) {
    replaceDoc(data.doc);
  } else {
    location.reload();
  }
});

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));
const prefetchDelay = 50;

const prefetch = {
  queue: async (url: string) => {
    if (prefetchStore.has(url)) {
      return;
    }

    const doc = await queuedFetcher.queue(async ({ signal }) => {
      const resp = await fetch(url, { signal });
      if (!resp.ok) {
        console.error("Prefetch failed: ", resp);
        return null;
      }
      const data = await resp.text();
      const doc = new DOMParser().parseFromString(data, "text/html");
      await sleep(prefetchDelay);
      return doc;
    });

    if (doc) {
      prefetchStore.set(url, { doc });
    }
  },
  load: async (url: string, e?: MouseEvent) => {
    const data = prefetchStore.get(url);
    // If no data is available, or if the meta key is pressed, do nothing
    // and let the default behavior of the link take over.
    if (!data || e?.metaKey) {
      return;
    }
    e?.preventDefault();

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

const stores = Array.from(document.querySelectorAll("[data-store-key]"));

for (const store of stores) {
  const key = store.getAttribute("data-store-key");
  if (!key) {
    throw new Error("Global store must have a data-store-key attribute");
  }
  if (
    !(
      store instanceof HTMLScriptElement &&
      store.getAttribute("type") === "application/json"
    )
  ) {
    throw new Error("Global store must be a JSON script tag");
  }
  const value = JSON.parse(store.textContent || "{}");
  Alpine.store(key, value);
}

Alpine.start();
