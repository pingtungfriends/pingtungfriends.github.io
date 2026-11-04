/* The generated HTML remains complete and navigable without JavaScript. */
(() => {
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)');
  const towns = [...document.querySelectorAll('.has-brand')];
  const mapLabels = [...document.querySelectorAll('.town-label')];
  const mapSVG = document.querySelector('.map-stage svg');
  if (mapSVG) {
    const sizeMapLabels = () => {
      const scale = mapSVG.getScreenCTM()?.a;
      if (scale > 0) mapSVG.style.setProperty('--map-label-size', `${13.5 / scale}px`);
    };
    sizeMapLabels();
    if ('ResizeObserver' in window) new ResizeObserver(sizeMapLabels).observe(mapSVG);
    else window.addEventListener('resize', sizeMapLabels);
  }
  const brands = [...document.querySelectorAll('.brand-row')];
  const caption = document.querySelector('.map-caption');
  function activate(brand) {
    towns.forEach(town => {
      const active = town.dataset.slugs.split(' ').includes(brand.dataset.slug);
      town.classList.toggle('is-active', active);
      if (active) town.setAttribute('aria-current', 'location');
      else town.removeAttribute('aria-current');
    });
    mapLabels.forEach(label => label.classList.toggle('is-active', label.dataset.slugs.split(' ').includes(brand.dataset.slug)));
    if (caption) caption.textContent = `${brand.dataset.name} · ${brand.dataset.town}`;
  }
  if (brands.length) {
    activate(brands.find(brand => `#${brand.id}` === location.hash) || brands[0]);
    if ('IntersectionObserver' in window) {
      const inMiddle = new Set();
      const syncCurrentBrand = () => {
        const middle = innerHeight * .5;
        const distance = brand => {
          const rect = brand.getBoundingClientRect();
          return Math.abs(rect.top + rect.height / 2 - middle);
        };
        const nearest = [...inMiddle].sort((a, b) => distance(a) - distance(b))[0];
        if (nearest) activate(nearest);
      };
      const brandObserver = new IntersectionObserver(entries => {
        entries.forEach(entry => entry.isIntersecting ? inMiddle.add(entry.target) : inMiddle.delete(entry.target));
        syncCurrentBrand();
      }, { rootMargin: '-20% 0px -20% 0px', threshold: 0 });
      // Two groups can remain inside the observer band throughout a smooth scroll.
      // Re-evaluate their centres as they move, even without a new intersection.
      let scrollFrame = 0;
      window.addEventListener('scroll', () => {
        if (scrollFrame) return;
        scrollFrame = requestAnimationFrame(() => { scrollFrame = 0; syncCurrentBrand(); });
      }, { passive: true });
      brands.forEach(brand => brandObserver.observe(brand));
    }
    towns.forEach(town => town.addEventListener('click', event => {
      const target = document.querySelector(town.getAttribute('href'));
      if (!target) return;
      event.preventDefault();
      history.pushState(null, '', town.getAttribute('href'));
      activate(target);
      target.scrollIntoView({ behavior: reduced.matches ? 'instant' : 'smooth', block: 'start' });
      // Continue keyboard reading in the brand list after following a map link.
      target.setAttribute('tabindex', '-1');
      target.focus({ preventScroll: true });
    }));
  }

  const tabs = [...document.querySelectorAll('[role="tab"]')];
  const panels = tabs.map(tab => document.getElementById(tab.getAttribute('aria-controls')));
  function selectTab(index, focus = false) {
    tabs.forEach((tab, i) => {
      tab.setAttribute('aria-selected', String(i === index));
      tab.tabIndex = i === index ? 0 : -1;
      panels[i].hidden = i !== index;
    });
    if (focus) tabs[index].focus();
  }
  tabs.forEach((tab, i) => {
    tab.addEventListener('click', () => selectTab(i));
    tab.addEventListener('keydown', event => {
      let next;
      if (event.key === 'ArrowRight') next = (i + 1) % tabs.length;
      if (event.key === 'ArrowLeft') next = (i - 1 + tabs.length) % tabs.length;
      if (event.key === 'Home') next = 0;
      if (event.key === 'End') next = tabs.length - 1;
      if (next !== undefined) { event.preventDefault(); selectTab(next, true); }
    });
  });
  if (tabs.length) selectTab(0);

  document.querySelectorAll('[data-video]').forEach(frame => {
    frame.querySelector('.video-play').addEventListener('click', event => {
      event.preventDefault();
      const iframe = document.createElement('iframe');
      iframe.src = `https://www.youtube-nocookie.com/embed/${frame.dataset.video}?autoplay=1`;
      iframe.title = frame.dataset.videoTitle;
      iframe.allow = 'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share';
      iframe.allowFullscreen = true;
      iframe.referrerPolicy = 'strict-origin-when-cross-origin';
      frame.replaceChildren(iframe);
      iframe.focus();
    });
  });

  if ('IntersectionObserver' in window && !reduced.matches) {
    const reveals = [...document.querySelectorAll('[data-reveal]')];
    const observer = new IntersectionObserver(entries => {
      entries.forEach(entry => {
        if (!entry.isIntersecting) return;
        entry.target.classList.add('is-visible');
        observer.unobserve(entry.target);
      });
    }, { threshold: .08 });
    reveals.forEach((element, i) => {
      element.classList.add('reveal-ready');
      element.style.setProperty('--reveal-delay', `${(i % 3) * 70}ms`);
      observer.observe(element);
    });
    const shutterObserver = new IntersectionObserver(entries => {
      entries.forEach(entry => {
        if (!entry.isIntersecting) return;
        entry.target.classList.add('is-open');
        shutterObserver.unobserve(entry.target);
      });
    }, { threshold: .12 });
    document.querySelectorAll('[data-shutter]').forEach(element => {
      element.classList.add('shutter-ready');
      shutterObserver.observe(element);
    });
  }
  document.documentElement.classList.add('js');
})();
