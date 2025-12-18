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
  const hero = document.querySelector('.lp-hero');
  hero.style.setProperty('--hero-bg', `url('${brand.heroImage}')`);
  document.querySelector('.hero-badges').innerHTML = brand.badges.map(b => `<span class="badge">${b}</span>`).join('');
  document.querySelector('.hero-brand').textContent = brand.name;
  document.querySelector('.hero-tagline').textContent = brand.tagline;
  document.querySelector('.hero-story').textContent = brand.story;
  const ctas = document.querySelectorAll('.hero-cta');
  ctas.forEach(cta => { cta.href = brand.cta.link; cta.textContent = brand.cta.label; });
  const bigText = document.querySelector('.hero-big-text');
  bigText.innerHTML = (brand.heroWords || []).map(w => `<span>${w}</span>`).join('');
  applyHeroContrast(hero, brand.heroImage);
}

function renderStory(brand) {
  document.querySelector('.story-text').textContent = brand.story;
  const storyImg = document.querySelector('.story-img');
  storyImg.src = brand.heroImage;
  storyImg.alt = brand.name;
}

function renderProducts(brand) {
  const wrap = document.querySelector('.product-grid');
  wrap.innerHTML = '';
  brand.products.slice(0, 3).forEach(p => wrap.appendChild(createProductCard(p)));
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
  const faq = document.querySelector('#faq-list');
  if (!faq) return;
  faq.innerHTML = '';
  (brand.faq || []).forEach(item => {
    const div = document.createElement('div');
    div.className = 'faq-item';
    div.innerHTML = `<strong>Q：${item.q}</strong><div>A：${item.a}</div>`;
    faq.appendChild(div);
  });
}

function renderFooterCTA(brand) {
  const btns = document.querySelectorAll('.hero-cta');
  btns.forEach(btn => {
    btn.href = brand.cta.link;
    btn.textContent = brand.cta.label || '立即購買';
  });
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
  // 這裡可以加入亮度判斷，暫時設為 bright 提升視覺感
  el.classList.add('hero-bright');
  el.style.setProperty('--hero-ov-boost', '0.2');
}

async function initIndex() {
  const data = await loadBrands();
  const brand = data.brands[0];
  if (brand) {
    renderHero(brand);
    renderStory(brand);
    renderProducts(brand);
    renderFooterCTA(brand);
  }
}

async function initBrand() {
  const slug = getQuery('slug');
  const data = await loadBrands();
  const brand = data.brands.find(b => b.slug === slug);
  if (!brand) {
    const container = document.querySelector('.page-shell');
    if (container) container.innerHTML = '<p style="padding: 100px; text-align: center;">找不到品牌資料。</p>';
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
  if (page === 'landing') initIndex();
  if (page === 'brand') initBrand();
});
