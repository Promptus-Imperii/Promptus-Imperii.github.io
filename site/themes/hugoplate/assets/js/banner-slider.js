// Banner slider on the home page. The home page template (layouts/index.html) loads this
// together with Swiper.
// See: https://swiperjs.com/
(function () {
  "use strict";

  // Banner Slider
  // ----------------------------------------
  new Swiper(".banner-slider", {
    spaceBetween: 24,
    loop: true,
    pagination: {
      el: ".banner-slider-pagination",
      type: "bullets",
      clickable: true,
    },
    autoplay: {
      delay: 3000,
    },
    speed: 1000,
    effect: 'fade',
    fadeEffect: {
      crossFade: true
    },
    keyboard: {
      enabled: true,
      onlyInViewport: true,
    },
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
  });
})();
