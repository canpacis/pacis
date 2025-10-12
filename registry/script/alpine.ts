import Alpine from "alpinejs";

// Alpine.data("data", (id: string) => {
//   const raw = document.querySelector(`script[type="application/json"]#${id}`)
//     ?.textContent ?? "{}";
//   return JSON.parse(raw);
// });

// const dropdownStore = new Map<string, any>();

// Alpine.magic("dropdown", () => (id: string) => dropdownStore.get(id));

// Alpine.data("dropdown", (open = false, id = null) => ({
//   id: id,
//   open: open,
//   usedKeyboard: false,

//   init() {
//     if (this.id !== null) {
//       dropdownStore.set(this.id, this);
//     }
//     this.$dispatch("init", { open: this.open });
//   },
//   async openDropdown() {
//     this.open = true;
//     await this.$nextTick();
//     this.$dispatch("opened");
//   },
//   async closeDropdown(value: string) {
//     this.open = false;
//     await this.$nextTick();
//     this.$dispatch("closed", { value: value });
//   },
//   async dismissDropdown() {
//     this.open = false;
//     await this.$nextTick();
//     this.$dispatch("dismissed");
//   },
// }));

Alpine.start();
