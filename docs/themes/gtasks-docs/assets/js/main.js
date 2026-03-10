/* gtasks-docs main.js */

// ── Theme (runs before paint to prevent flash) ────────────────
(function () {
  var stored = localStorage.getItem('gtasks-theme');
  var theme = stored || 'dark';
  document.documentElement.setAttribute('data-theme', theme);
})();

// ── Micro-animation keyframes injected via JS ─────────────────
(function () {
  var style = document.createElement('style');
  style.textContent = [
    '@keyframes gt-fade-up {',
    '  from { opacity: 0; transform: translateY(14px); }',
    '  to   { opacity: 1; transform: translateY(0); }',
    '}',
    '@keyframes gt-slide-in {',
    '  from { opacity: 0; transform: translateX(-18px); }',
    '  to   { opacity: 1; transform: translateX(0); }',
    '}',
    '@keyframes gt-slide-in-right {',
    '  from { opacity: 0; transform: translateX(18px); }',
    '  to   { opacity: 1; transform: translateX(0); }',
    '}',
    '@keyframes gt-stagger-in {',
    '  from { opacity: 0; transform: translateY(10px); }',
    '  to   { opacity: 1; transform: translateY(0); }',
    '}',
    '[data-fade] { opacity: 0; }',
    '[data-slide-in] { opacity: 0; }',
    '[data-slide-in-right] { opacity: 0; }',
    '[data-stagger] { opacity: 0; }',
  ].join('\n');
  document.head.appendChild(style);
})();

document.addEventListener('DOMContentLoaded', function () {

  // ── Theme toggle ────────────────────────────────────────────
  var btn = document.getElementById('theme-toggle');
  if (btn) {
    var cur = document.documentElement.getAttribute('data-theme');
    btn.textContent = cur === 'dark' ? '○' : '●';
    btn.addEventListener('click', function () {
      var next = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
      document.documentElement.setAttribute('data-theme', next);
      localStorage.setItem('gtasks-theme', next);
      btn.textContent = next === 'dark' ? '○' : '●';
    });
  }

  // ── Hero staggered fade-up ──────────────────────────────────
  document.querySelectorAll('[data-fade]').forEach(function (el) {
    var delay = parseInt(el.getAttribute('data-fade'), 10) * 90;
    el.style.animation = 'gt-fade-up 0.5s cubic-bezier(0.22,1,0.36,1) ' + (80 + delay) + 'ms both';
  });

  // ── Intersection observer for scroll animations ─────────────
  if ('IntersectionObserver' in window) {

    var slideObs = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        if (e.isIntersecting) {
          e.target.style.animation = 'gt-slide-in 0.55s cubic-bezier(0.22,1,0.36,1) both';
          slideObs.unobserve(e.target);
        }
      });
    }, { threshold: 0.12 });

    document.querySelectorAll('[data-slide-in]').forEach(function (el) {
      slideObs.observe(el);
    });

    var slideRightObs = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        if (e.isIntersecting) {
          e.target.style.animation = 'gt-slide-in-right 0.55s cubic-bezier(0.22,1,0.36,1) 80ms both';
          slideRightObs.unobserve(e.target);
        }
      });
    }, { threshold: 0.12 });

    document.querySelectorAll('[data-slide-in-right]').forEach(function (el) {
      slideRightObs.observe(el);
    });

    var staggerObs = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        if (e.isIntersecting) {
          var cells = e.target.querySelectorAll('[data-stagger]');
          cells.forEach(function (cell, i) {
            cell.style.animation = 'gt-stagger-in 0.45s cubic-bezier(0.22,1,0.36,1) ' + (i * 55) + 'ms both';
          });
          staggerObs.unobserve(e.target);
        }
      });
    }, { threshold: 0.08 });

    document.querySelectorAll('.features-grid').forEach(function (grid) {
      staggerObs.observe(grid);
    });

    var faqObs = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        if (e.isIntersecting) {
          var items = e.target.querySelectorAll('.faq-item');
          items.forEach(function (item, i) {
            item.style.opacity = '0';
            item.style.animation = 'gt-fade-up 0.4s cubic-bezier(0.22,1,0.36,1) ' + (i * 45) + 'ms both';
          });
          faqObs.unobserve(e.target);
        }
      });
    }, { threshold: 0.08 });

    document.querySelectorAll('.faq-list').forEach(function (list) {
      faqObs.observe(list);
    });

  } else {
    document.querySelectorAll('[data-fade],[data-slide-in],[data-slide-in-right],[data-stagger]').forEach(function (el) {
      el.style.opacity = '1';
    });
  }

  // ── Install tabs ────────────────────────────────────────────
  document.querySelectorAll('.tab-btn').forEach(function (tabBtn) {
    tabBtn.addEventListener('click', function () {
      var target = tabBtn.getAttribute('data-tab');
      document.querySelectorAll('.tab-btn').forEach(function (b) { b.classList.remove('active'); });
      document.querySelectorAll('.tab-panel').forEach(function (p) { p.classList.remove('active'); });
      tabBtn.classList.add('active');
      var panel = document.getElementById('tab-' + target);
      if (panel) {
        panel.classList.add('active');
        panel.style.opacity = '0';
        panel.style.transition = 'opacity 0.18s';
        requestAnimationFrame(function () { panel.style.opacity = '1'; });
      }
    });
  });

  // ── Copy buttons ────────────────────────────────────────────
  document.querySelectorAll('.copy-btn').forEach(function (copyBtn) {
    copyBtn.addEventListener('click', function () {
      var targetId = copyBtn.getAttribute('data-copy');
      var text = targetId ? (document.getElementById(targetId) || {}).textContent : null;
      if (!text) return;
      navigator.clipboard.writeText(text.trim()).then(function () {
        var orig = copyBtn.textContent;
        copyBtn.textContent = 'copied!';
        copyBtn.style.color = '#50E3C2';
        copyBtn.style.borderColor = '#50E3C2';
        setTimeout(function () {
          copyBtn.textContent = orig;
          copyBtn.style.color = '';
          copyBtn.style.borderColor = '';
        }, 1600);
      });
    });
  });

  // ── FAQ smooth open/close ────────────────────────────────────
  document.querySelectorAll('details.faq-item').forEach(function (det) {
    var body = det.querySelector('.faq-body');
    if (!body) return;
    det.addEventListener('toggle', function () {
      if (det.open) {
        body.style.maxHeight = '0';
        body.style.overflow = 'hidden';
        body.style.transition = 'max-height 0.25s cubic-bezier(0.22,1,0.36,1)';
        requestAnimationFrame(function () { body.style.maxHeight = body.scrollHeight + 'px'; });
      } else {
        body.style.maxHeight = body.scrollHeight + 'px';
        requestAnimationFrame(function () {
          body.style.maxHeight = '0';
          body.addEventListener('transitionend', function h() {
            body.style.maxHeight = '';
            body.style.overflow = '';
            body.style.transition = '';
            body.removeEventListener('transitionend', h);
          });
        });
      }
    });
  });

  // ── Active sidebar link ──────────────────────────────────────
  var path = window.location.pathname;
  document.querySelectorAll('.sidebar-link').forEach(function (link) {
    var href = link.getAttribute('href');
    if (href === path || href === path.replace(/\/$/, '') || href + '/' === path) {
      link.classList.add('active');
    }
  });

  // ── Nav underline micro-interaction ─────────────────────────
  var navStyle = document.createElement('style');
  navStyle.textContent =
    '.nav-links li a:not(.nav-gh) { position: relative; }' +
    '.nav-links li a:not(.nav-gh)::after { content:""; position:absolute; bottom:-3px; left:0; width:0; height:1px;' +
    '  background:#50E3C2; transition:width 0.18s cubic-bezier(0.22,1,0.36,1); }' +
    '.nav-links li a:not(.nav-gh):hover::after { width:100%; }';
  document.head.appendChild(navStyle);

  // ── Feature cell hover: teal left border ────────────────────
  document.querySelectorAll('.feature-cell').forEach(function (cell) {
    cell.addEventListener('mouseenter', function () {
      cell.style.transition = 'background 0.15s, border-left-color 0.12s';
      cell.style.borderLeft = '2px solid #50E3C2';
    });
    cell.addEventListener('mouseleave', function () {
      cell.style.borderLeft = '';
    });
  });

  // ── Step rows: left indicator on hover ──────────────────────
  document.querySelectorAll('.step-row').forEach(function (row) {
    row.style.transition = 'background 0.12s, border-left 0.12s';
    row.addEventListener('mouseenter', function () {
      row.style.borderLeft = '2px solid #50E3C2';
    });
    row.addEventListener('mouseleave', function () {
      row.style.borderLeft = '';
    });
  });

  // ── Btn: ripple on click ─────────────────────────────────────
  document.querySelectorAll('.btn').forEach(function (b) {
    b.addEventListener('click', function (e) {
      var rect = b.getBoundingClientRect();
      var rip = document.createElement('span');
      rip.style.cssText = [
        'position:absolute',
        'border-radius:50%',
        'width:8px',
        'height:8px',
        'background:rgba(255,255,255,0.35)',
        'pointer-events:none',
        'transform:scale(0)',
        'animation:gt-ripple 0.5s ease-out both',
        'left:' + (e.clientX - rect.left - 4) + 'px',
        'top:' + (e.clientY - rect.top - 4) + 'px',
      ].join(';');
      if (getComputedStyle(b).position === 'static') b.style.position = 'relative';
      b.style.overflow = 'hidden';
      b.appendChild(rip);
      setTimeout(function () { rip.remove(); }, 500);
    });
  });

  var ripStyle = document.createElement('style');
  ripStyle.textContent = '@keyframes gt-ripple { to { transform:scale(24); opacity:0; } }';
  document.head.appendChild(ripStyle);

});
