// Shows the "Er is iets fout gegaan" message (#failed) on the signup and contact pages when a
// request fails without an answer from our backend: no connection, a timeout, or an error page
// from the proxy in front of it. The backend's own errors arrive as an <hx-partial> and show
// their own message instead. See https://htmx.org/essays/hypermedia-friendly-scripting/
(function () {
  "use strict";
  const showFailed = (show) => {
    const failed = document.getElementById("failed");
    if (failed) failed.hidden = !show;
  };
  document.addEventListener("htmx:before:request", () => showFailed(false));
  document.addEventListener("htmx:error", () => showFailed(true));
  document.addEventListener("htmx:before:swap", (event) => {
    const { ctx, tasks } = event.detail;
    if (ctx.response.status >= 400) {
      showFailed(!tasks.some((task) => task.type === "partial"));
    }
  });
})();
