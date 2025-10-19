import anchor from "@alpinejs/anchor";
import focus from "@alpinejs/focus";
import Alpine from "alpinejs";

Alpine.plugin(anchor);
Alpine.plugin(focus);

Alpine.data("data", (id: string) => {
  const raw = document.querySelector(`script[type="application/json"]#${id}`)
    ?.textContent ?? "{}";
  return JSON.parse(raw);
});

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

Alpine.data("tabs", (defaultValue: string) => ({
  active: defaultValue,
  select(value: string, root: HTMLElement | null = null) {
    this.active = value;
    if (root) {
      root.dispatchEvent(new CustomEvent("changed", { detail: { value } }));
    }
  },
}));

Alpine.start();
