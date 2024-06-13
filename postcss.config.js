const purgecss = {
  content: ['./**/*.html'],
  defaultExtractor: (content) => {
    const defaultSelectors = content.match(/[A-Za-z0-9_-]+/g) || [];
    const extendedSelectors = content.match(/[^<>"=\s]+/g) || [];
    return defaultSelectors.concat(extendedSelectors);
  },
  safelist: [
    /^swiper-/,
    /^lb-/,
    /^gl/,
    /^go/,
    /^gc/,
    /^gs/,
    /^gi/,
    /^gz/,
    /^gprev/,
    /^gnext/,
    /^desc/,
    /^zoom/,
    /^search/,
    /^:is/,
    /dark/,
    /show/,
    /dragging/,
    /fullscreen/,
    /loaded/,
    /visible/,
    /current/,
    /active/,
  ],
};

module.exports = {
  plugins: {
    tailwindcss: {},
    "@fullhuman/postcss-purgecss":
      process.env.HUGO_ENVIRONMENT === "production" ? purgecss : false,
    autoprefixer: process.env.HUGO_ENVIRONMENT === "production" ? {} : false,
  },
};
