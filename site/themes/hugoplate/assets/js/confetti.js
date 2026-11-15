// Confetti on the signup success page: a burst of particles at the pointer, at most every 200 ms.
(() => {
  if (matchMedia("(prefers-reduced-motion: reduce)").matches) return;

  const COUNT = 75;
  const POWER = 25;
  const GRAVITY = 0.5;
  const COLORS = ["#f44336", "#e91e63", "#9c27b0", "#3f51b5", "#03a9f4", "#4caf50", "#ffeb3b", "#ff9800"];

  const canvas = document.createElement("canvas");
  canvas.style.cssText = "position:fixed;inset:0;pointer-events:none;z-index:50";
  document.body.append(canvas);
  const ctx = canvas.getContext("2d");

  const resize = () => {
    canvas.width = innerWidth;
    canvas.height = innerHeight;
  };
  addEventListener("resize", resize);
  resize();

  let particles = [];
  let lastBurst = -Infinity;

  addEventListener("pointermove", (e) => {
    if (e.timeStamp - lastBurst < 200) return;
    lastBurst = e.timeStamp;
    const idle = particles.length === 0;
    for (let i = 0; i < COUNT; i++) {
      const angle = Math.random() * 2 * Math.PI;
      const speed = Math.random() * POWER;
      particles.push({
        x: e.clientX,
        y: e.clientY,
        vx: Math.cos(angle) * speed,
        vy: Math.sin(angle) * speed,
        size: 4 + Math.random() * 4,
        rotation: Math.random() * Math.PI,
        color: COLORS[i % COLORS.length],
      });
    }
    if (idle) requestAnimationFrame(draw);
  });

  function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    for (const p of particles) {
      p.vx *= 0.96;
      p.vy = p.vy * 0.96 + GRAVITY;
      p.x += p.vx;
      p.y += p.vy;
      p.rotation += p.vx / 10;
      ctx.save();
      ctx.translate(p.x, p.y);
      ctx.rotate(p.rotation);
      ctx.fillStyle = p.color;
      ctx.fillRect(-p.size / 2, -p.size / 4, p.size, p.size / 2);
      ctx.restore();
    }
    particles = particles.filter((p) => p.y < canvas.height + 20);
    if (particles.length) requestAnimationFrame(draw);
  }
})();
