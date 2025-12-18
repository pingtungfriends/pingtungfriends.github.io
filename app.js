async function loadBrands() {
  const res = await fetch('data/brands.json');
  if (!res.ok) throw new Error('無法載入品牌資料');
  return res.json();
}

function getQuery(name) {
  const params = new URLSearchParams(window.location.search);
  return params.get(name);
}

function createProductCard(product) {
  const div = document.createElement('div');
  div.className = 'card';
  div.innerHTML = `
    <img src="${product.image}" alt="${product.name}" loading="lazy" />
    <h3>${product.name}</h3>
    <p>${product.summary}</p>
    <p class="price">$${product.price} · ${product.spec}</p>
    <a class="btn-primary" href="${product.link}" target="_blank" rel="noopener">立即購買</a>
  `;
  return div;
}

function createBrandCard(brand) {
  const div = document.createElement('div');
  div.className = 'brand-card';
  div.innerHTML = `
    <img src="${brand.heroImage}" alt="${brand.name}" loading="lazy" />
    <div>
      <h3>${brand.name}</h3>
      <p class="small">${brand.tagline}</p>
      <p class="small">${brand.origin}</p>
    </div>
    <a class="btn-primary" href="brand.html?slug=${brand.slug}">查看一頁式</a>
  `;
  return div;
}

function renderHero(brand) {
  const hero = document.querySelector('.hero');
  hero.style.setProperty('--hero-bg', `url('${brand.heroImage}')`);
  hero.innerHTML = `
    <img src="${brand.heroImage}" alt="${brand.name}" />
    <div class="content hero-grid">
      <div>
        <div class="badges">${brand.badges.map(b => `<span class="badge">${b}</span>`).join('')}</div>
        <h1>${brand.name}</h1>
        <p>${brand.tagline}</p>
        <p>${brand.story}</p>
        <a class="btn-primary" href="${brand.cta.link}" target="_blank" rel="noopener">${brand.cta.label}</a>
      </div>
      <div class="glass-panel">
        <div class="big-text">${(brand.heroWords || []).map(w => `<span>${w}</span>`).join('')}</div>
      </div>
    </div>
  `;
  applyHeroContrast(hero, brand.heroImage);
}

function renderStory(brand) {
  document.querySelector('#story').innerHTML = `
    <h2>品牌故事｜${brand.origin}</h2>
    <p>${brand.story}</p>
    <div class="badges">${brand.badges.map(b => `<span class="badge" style="color:var(--leaf-dark);background:rgba(47,122,77,0.12);border:1px solid rgba(47,122,77,0.25);">${b}</span>`).join('')}</div>
  `;
}

function renderProducts(brand) {
  const wrap = document.querySelector('#products .grid');
  brand.products.forEach(p => wrap.appendChild(createProductCard(p)));
}

function renderVideos(brand) {
  const list = document.querySelector('#videos .video-list');
  brand.videos.forEach(v => {
    const item = document.createElement('div');
    item.className = 'video-item';
    item.innerHTML = `<div>${v.title}</div><a href="${v.url}" target="_blank" rel="noopener">立即播放</a>`;
    list.appendChild(item);
  });
}

function renderFAQ(brand) {
  const faq = document.querySelector('#faq');
  brand.faq.forEach(item => {
    const div = document.createElement('div');
    div.className = 'faq-item';
    div.innerHTML = `<strong>Q：${item.q}</strong><div>A：${item.a}</div>`;
    faq.appendChild(div);
  });
}

function renderFooterCTA(brand) {
  const btn = document.querySelector('#footer-cta');
  btn.href = brand.cta.link;
  btn.textContent = brand.cta.label;
}

function renderHeroShowcase(brands) {
  const el = document.querySelector('.hero-showcase');
  if (!el || !brands.length) return;
  const brand = brands[0];
  el.style.setProperty('--hero-bg', `url('${brand.heroImage}')`);
  el.querySelector('.hero-brand').textContent = brand.name;
  el.querySelector('.hero-tagline').textContent = brand.tagline;
  el.querySelector('.hero-story').textContent = brand.story;
  const cta = el.querySelector('.hero-cta');
  cta.href = brand.cta.link;
  cta.textContent = brand.cta.label;
  const words = brand.heroWords?.length ? brand.heroWords : (brand.tagline || '').split(' ').slice(0, 3);
  const bigText = el.querySelector('.hero-big-text');
  bigText.innerHTML = words.map(w => `<span>${w}</span>`).join('');
  const pager = el.querySelector('.hero-pager');
  pager.innerHTML = `01 <span>/</span> ${String(brands.length).padStart(2, '0')}`;
  applyHeroContrast(el, brand.heroImage);
}

function sampleImageLightness(src, defaultValue = null) {
  return new Promise(resolve => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.src = src;
    img.onload = () => {
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      const size = 32;
      canvas.width = size;
      canvas.height = size;
      ctx.drawImage(img, 0, 0, size, size);
      const data = ctx.getImageData(0, 0, size, size).data;
      let sum = 0;
      for (let i = 0; i < data.length; i += 4) {
        const r = data[i], g = data[i + 1], b = data[i + 2];
        sum += 0.2126 * r + 0.7152 * g + 0.0722 * b; // 相對亮度
      }
      const avg = sum / (data.length / 4);
      resolve(avg);
    };
    img.onerror = () => resolve(defaultValue);
  });
}

async function applyHeroContrast(el, imgUrl) {
  if (!el || !imgUrl) return;
  el.classList.remove('hero-light', 'hero-dark', 'hero-bright');
  const lightness = await sampleImageLightness(imgUrl, null);
  // >180 表示背景偏亮：提高覆蓋強度
  if (lightness !== null && lightness > 180) {
    el.classList.add('hero-bright');
    el.style.setProperty('--hero-ov-boost', '0.25');
  } else {
    el.style.setProperty('--hero-ov-boost', '0');
  }
}

async function initIndex() {
  const data = await loadBrands();
  renderHeroShowcase(data.brands);
  const wrap = document.querySelector('.brand-list');
  data.brands.forEach(brand => wrap.appendChild(createBrandCard(brand)));
}

async function initBrand() {
  const slug = getQuery('slug');
  const data = await loadBrands();
  const brand = data.brands.find(b => b.slug === slug);
  if (!brand) {
    document.querySelector('.container').innerHTML = '<p>找不到品牌資料。</p>';
    return;
  }
  document.title = brand.seo?.title || brand.name;
  renderHero(brand);
  renderStory(brand);
  renderProducts(brand);
  renderVideos(brand);
  renderFAQ(brand);
  renderFooterCTA(brand);
}

window.addEventListener('DOMContentLoaded', () => {
  const page = document.body.dataset.page;
  if (page === 'index') initIndex();
  if (page === 'brand') initBrand();
});
