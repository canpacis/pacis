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
