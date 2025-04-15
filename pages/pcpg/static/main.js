(() => {
  // www/app/assets/main.ts
  var observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.addedNodes.length > 0) {
        for (const node of mutation.addedNodes) {
          const elem = node;
          const slot = elem.getAttribute("slot");
          if (!slot) {
            throw new Error("Slotless element has been streamed to the DOM");
          }
          document.querySelectorAll(`slot[name=${slot}]`).forEach((slot2) => {
            elem.setAttribute("x-show", "true");
            slot2.replaceWith(elem);
          });
        }
      }
    }
  });
  observer.observe(document.body, { childList: true });
  window.addEventListener("DOMContentLoaded", () => {
    observer.disconnect();
  });
})();
