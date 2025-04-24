window.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("a").forEach((a) => {
    if (a.host !== window.location.host && a.href.startsWith("http")) {
      a.setAttribute("data-umami-event", "outbound-link-click");
      a.setAttribute("data-umami-event-url", a.href);
    }
  });
});

window.addEventListener("alpine:init", () => {
  window.addEventListener("DOMContentLoaded", () => {
    const user = Alpine.store("user");
    if (user.logged_in)
      window.umami.identify({
        email: user.email,
        name: user.name,
        id: user.id,
      });
  });
});
