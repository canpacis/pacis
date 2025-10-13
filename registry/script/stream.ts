const observer = new MutationObserver((mutations) => {
  for (const mutation of mutations) {
    if (mutation.addedNodes.length > 0) {
      for (const node of mutation.addedNodes) {
        const elem = node;
        if (!(elem instanceof HTMLElement)) {
          return;
        }
        const slot = elem.getAttribute("slot");
        if (!slot) {
          return;
        }
        document.querySelectorAll(`slot[name=${slot}]`).forEach((target) => {
          target.replaceWith(elem);
        });
      }
    }
  }
});

observer.observe(document.body, { childList: true });
window.addEventListener("DOMContentLoaded", () => {
  observer.disconnect();
});
