window.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("a").forEach((a) => {
    console.log(a);
    if (a.host !== window.location.host && a.href.startsWith("http")) {
      a.setAttribute("data-umami-event", "outbound-link-click");
      a.setAttribute("data-umami-event-url", a.href);
    }
  });
});
