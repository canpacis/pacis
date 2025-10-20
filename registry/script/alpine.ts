import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";
import Alpine from "alpinejs";
import { Notyf, type NotyfNotification } from "notyf";

Alpine.plugin(anchor);
Alpine.plugin(focus);

Alpine.data("data", (id: string) => {
  const raw = document.querySelector(`script[type="application/json"]#${id}`)
    ?.textContent ?? "{}";
  return JSON.parse(raw);
});

Alpine.data("accordion", (defaultValue: string = "") => ({
  active: defaultValue,
  select(value: string, root: HTMLElement | null) {
    if (this.active === value) {
      this.active = "";
    } else {
      this.active = value;
    }
    if (root) {
      root.dispatchEvent(
        new CustomEvent("changed", { detail: { value: this.active } }),
      );
    }
  },
}));

Alpine.data("dialog", () => ({
  opened: false,
  open(root: HTMLElement | null) {
    this.opened = true;
    if (root) {
      root.dispatchEvent(new CustomEvent("open"));
    }
  },
  close(root: HTMLElement | null) {
    this.opened = false;
    if (root) {
      root.dispatchEvent(new CustomEvent("closed"));
    }
  },
}));

Alpine.data("dropdown", () => ({
  keyboard: false,
  mouse: false,
  get opened(): boolean {
    return this.keyboard || this.mouse;
  },
  open(w: "mouse" | "keyboard" = "mouse", root: HTMLElement | null) {
    if (w === "mouse") {
      this.mouse = true;
    } else {
      this.keyboard = true;
    }
    if (root) {
      root.dispatchEvent(new CustomEvent("open"));
    }
  },
  close(root: HTMLElement | null) {
    this.mouse = false;
    this.keyboard = false;
    if (root) {
      root.dispatchEvent(new CustomEvent("closed"));
    }
  },
}));

Alpine.data("select", (defaultValue: string) => ({
  value: defaultValue,
  keyboard: false,
  mouse: false,
  get opened(): boolean {
    return this.keyboard || this.mouse;
  },
  label(root: HTMLElement, placeholder: string): string {
    if (this.value.length === 0) {
      return placeholder;
    }
    const element = Array.from(root.querySelectorAll("[role='option']")).find(
      (option) => option.getAttribute("data-value") === this.value,
    );
    return element?.textContent ?? placeholder;
  },
  select(value: string, root: HTMLElement | null) {
    this.value = value;
    if (root) {
      root.dispatchEvent(new CustomEvent("changed", { detail: { value } }));
    }
  },
  open(w: "mouse" | "keyboard" = "mouse", root: HTMLElement | null) {
    if (w === "mouse") {
      this.mouse = true;
    } else {
      this.keyboard = true;
    }
    if (root) {
      root.dispatchEvent(new CustomEvent("open"));
    }
  },
  close(root: HTMLElement | null) {
    this.mouse = false;
    this.keyboard = false;
    if (root) {
      root.dispatchEvent(new CustomEvent("closed"));
    }
  },
}));

Alpine.data("sheet", () => ({
  opened: false,
  open(root: HTMLElement | null) {
    this.opened = true;
    if (root) {
      root.dispatchEvent(new CustomEvent("open"));
    }
  },
  close(root: HTMLElement | null) {
    this.opened = false;
    if (root) {
      root.dispatchEvent(new CustomEvent("closed"));
    }
  },
}));

Alpine.data("tabs", (defaultValue: string) => ({
  active: defaultValue,
  select(value: string, root: HTMLElement | null = null) {
    this.active = value;
    if (root) {
      root.dispatchEvent(new CustomEvent("changed", { detail: { value } }));
    }
  },
}));

const toast = new Notyf({ ripple: false });

const toastObserver = new MutationObserver((mutations) => {
  for (const mutation of mutations) {
    const target = mutation.target as HTMLElement;
    if (target.classList.contains("notyf__toast--disappear")) {
      target.classList.remove("notyf__toast--disappear");
      target.classList.add("animate-out");
      target.classList.add("fade-out");
      target.addEventListener("animationend", () => {
        target.remove();
      });
    }
  }
});

new MutationObserver((mutations) => {
  for (const mutation of mutations) {
    for (const node of mutation.addedNodes) {
      toastObserver.observe(node, { attributes: true });
    }
  }
}).observe(document.body.querySelector(".notyf") as HTMLElement, {
  attributes: false,
  childList: true,
  subtree: false,
});

Alpine.data("toast", (message: string = "", duration: number = 2000) => ({
  instance: null as unknown as NotyfNotification,
  show() {
    this.instance = toast.open({
      message: message,
      duration: duration,
      className:
        "border rounded-md shadow-lg animate-in fade-in w-64 p-4 text-sm cursor-default",
    });
  },
}));

Alpine.data("tooltip", (delay: number = 0) => ({
  opened: false,
  timout: 0,
  open(root: HTMLElement | null) {
    this.timout = setTimeout(() => {
      this.opened = true;
    }, delay);

    if (root) {
      root.dispatchEvent(new CustomEvent("open"));
    }
  },
  close(root: HTMLElement | null) {
    clearTimeout(this.timout);
    this.opened = false;
    if (root) {
      root.dispatchEvent(new CustomEvent("closed"));
    }
  },
}));

Alpine.start();
